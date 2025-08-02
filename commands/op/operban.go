package op

import (
	"goircd/logger"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type OperBanCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("OPERBAN", func(s *server.Server) server.Command {
		return &OperBanCommand{server: s}
	})
}

func (c *OperBanCommand) Name() string {
	return "OPERBAN"
}

func (c *OperBanCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if !client.IsOperator() {
		client.SendNumeric(utils.ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	parts := strings.SplitN(params, " ", 3)
	if len(parts) < 2 {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "OPERBAN :Not enough parameters")
		return
	}

	channelName := parts[0]
	banMask := parts[1]
	var reason string
	if len(parts) > 2 {
		reason = parts[2]
		if strings.HasPrefix(reason, ":") {
			reason = reason[1:]
		}
	} else {
		reason = "Banned by IRC operator"
	}

	channel := c.server.GetChannel(channelName)
	if channel == nil {
		client.SendNumeric(utils.ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	if !strings.ContainsAny(banMask, "*?") && !strings.Contains(banMask, "!") && !strings.Contains(banMask, "@") {
		targetClient := c.server.GetClient(banMask)
		if targetClient != nil {
			banMask = targetClient.Nick + "!" + targetClient.User + "@" + targetClient.Host
		} else {
			banMask = banMask + "!*@*"
		}
	}

	channel.SetBanned(banMask, true)

	logger.Ban(client.Nick, banMask, channelName, reason)

	modeMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " MODE " + channelName + " +b " + banMask
	channel.Broadcast(modeMsg)

	for _, chanClient := range channel.GetClients() {
		clientMask := chanClient.Nick + "!" + chanClient.User + "@" + chanClient.Host
		if matchesBanMask(clientMask, banMask) && chanClient != client && !chanClient.IsOperator() {
			logger.Kick(client.Nick, chanClient.Nick, channelName, reason)

			kickMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " KICK " + channelName + " " + chanClient.Nick + " :" + reason
			channel.Broadcast(kickMsg)

			channel.RemoveClient(chanClient)
		}
	}
}

func matchesBanMask(clientMask, banMask string) bool {
	return strings.Contains(clientMask, banMask) || banMask == clientMask
}

func (c *OperBanCommand) Help() string {
	return "OPERBAN <channel> <mask> [<reason>] - IRC operator command to ban users matching the mask from the channel"
}
