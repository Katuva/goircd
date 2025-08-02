package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type WhoCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("WHO", func(s *server.Server) server.Command {
		return &WhoCommand{server: s}
	})
}

func (c *WhoCommand) Name() string {
	return "WHO"
}

func (c *WhoCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	var mask string
	if params != "" {
		parts := strings.SplitN(params, " ", 2)
		mask = parts[0]
	}

	if strings.HasPrefix(mask, "#") || strings.HasPrefix(mask, "&") {
		channel := c.server.GetChannel(mask)
		if channel != nil {
			for _, member := range channel.GetClients() {
				sendWhoReply(client, member, mask, channel)
			}
		}
	} else {
		for _, channel := range c.server.GetAllChannels() {
			for _, member := range channel.GetClients() {
				if mask == "" || strings.Contains(member.Nick, mask) ||
					strings.Contains(member.User, mask) ||
					strings.Contains(member.Host, mask) {
					sendWhoReply(client, member, channel.Name, channel)
					break
				}
			}
		}
	}

	client.SendNumeric(utils.RPL_ENDOFWHO, mask+" :End of WHO list")
}

func sendWhoReply(client *server.Client, target *server.Client, channelName string, channel *server.Channel) {
	// Format: <channel> <user> <host> <server> <nick> <H|G>[*][@|+] :<hopcount> <real name>
	// H = Here, G = Gone (away), * = IRC operator, @ = channel operator, + = voiced

	status := "H" // Assume user is here (not away)

	// Add operator/voice status
	if target.IsOperator() {
		status += "*"
	}

	if channel != nil {
		if channel.IsOperator(target) {
			status += "@"
		} else if channel.IsVoiced(target) {
			status += "+"
		}
	}

	reply := channelName + " " +
		target.User + " " +
		target.Host + " " +
		utils.SERVER_NAME + " " +
		target.Nick + " " +
		status + " :0 " +
		target.RealName

	client.SendNumeric(utils.RPL_WHOREPLY, reply)
}

func (c *WhoCommand) Help() string {
	return "WHO [<mask>] - Lists users who match the given mask, or all users if no mask is given"
}
