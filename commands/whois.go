package commands

import (
	"strconv"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type WhoisCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("WHOIS", func(s *server.Server) server.Command {
		return &WhoisCommand{server: s}
	})
}

func (c *WhoisCommand) Name() string {
	return "WHOIS"
}

func (c *WhoisCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if params == "" {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "WHOIS :Not enough parameters")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	targetNick := parts[0]

	targetClient := c.server.GetClient(targetNick)
	if targetClient == nil {
		client.SendNumeric(utils.ERR_NOSUCHNICK, targetNick+" :No such nick/channel")
		return
	}

	client.SendNumeric(utils.RPL_WHOISUSER, targetNick+" "+targetClient.User+" "+targetClient.Host+" * :"+targetClient.RealName)

	client.SendNumeric(utils.RPL_WHOISSERVER, targetNick+" "+utils.SERVER_NAME+" :"+utils.SERVER_VERSION)

	if targetClient.IsOperator() {
		client.SendNumeric(utils.RPL_WHOISOPERATOR, targetNick+" :is an IRC operator")
	}

	idleTime := int(targetClient.GetIdleTime().Seconds())
	client.SendNumeric(utils.RPL_WHOISIDLE, targetNick+" "+strconv.Itoa(idleTime)+" :seconds idle")

	channels := targetClient.GetChannels()
	if len(channels) > 0 {
		var channelList strings.Builder
		for i, channel := range channels {
			if channel.IsOperator(targetClient) {
				channelList.WriteString("@")
			} else if channel.IsVoiced(targetClient) {
				channelList.WriteString("+")
			}
			channelList.WriteString(channel.Name)

			if i < len(channels)-1 {
				channelList.WriteString(" ")
			}
		}
		client.SendNumeric(utils.RPL_WHOISCHANNELS, targetNick+" :"+channelList.String())
	}

	client.SendNumeric(utils.RPL_ENDOFWHOIS, targetNick+" :End of WHOIS list")
}

func (c *WhoisCommand) Help() string {
	return "WHOIS <nickname> - Shows information about the specified user"
}
