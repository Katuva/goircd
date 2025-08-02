package commands

import (
	"goircd/logger"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type BanCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("BAN", func(s *server.Server) server.Command {
		return &BanCommand{server: s}
	})
}

func (c *BanCommand) Name() string {
	return "BAN"
}

func (c *BanCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	if len(parts) < 2 {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "BAN :Not enough parameters")
		return
	}

	channelName := parts[0]
	banMask := parts[1]

	channel := c.server.GetChannel(channelName)
	if channel == nil {
		client.SendNumeric(utils.ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	if !channel.HasClient(client) {
		client.SendNumeric(utils.ERR_NOTONCHANNEL, channelName+" :You're not on that channel")
		return
	}

	if !channel.IsOperator(client) {
		client.SendNumeric(utils.ERR_CHANOPRIVSNEEDED, channelName+" :You're not channel operator")
		return
	}

	if !strings.ContainsAny(banMask, "*?") && !strings.Contains(banMask, "!") && !strings.Contains(banMask, "@") {
		targetClient := c.server.GetClient(banMask)
		if targetClient != nil {
			banMask = utils.FormatMask(targetClient.Nick, targetClient.User, targetClient.Host)
		} else {
			banMask = banMask + "!*@*"
		}
	}

	if !utils.ValidateHostPattern(banMask) {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "BAN :Invalid ban mask")
		return
	}

	channel.SetBanned(banMask, true)

	logger.Ban(client.Nick, banMask, channelName, "")

	modeMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " MODE " + channelName + " +b " + banMask
	channel.Broadcast(modeMsg)
}

func (c *BanCommand) Help() string {
	return "BAN <channel> <mask> - Bans users matching the mask from the channel"
}
