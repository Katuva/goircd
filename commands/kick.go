package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type KickCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("KICK", func(s *server.Server) server.Command {
		return &KickCommand{server: s}
	})
}

func (c *KickCommand) Name() string {
	return "KICK"
}

func (c *KickCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	parts := strings.SplitN(params, " ", 3)
	if len(parts) < 2 {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "KICK :Not enough parameters")
		return
	}

	channelName := parts[0]
	targetNick := parts[1]
	var reason string
	if len(parts) > 2 {
		reason = parts[2]
		if strings.HasPrefix(reason, ":") {
			reason = reason[1:]
		}
	} else {
		reason = targetNick
	}

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

	targetClient := c.server.GetClient(targetNick)
	if targetClient == nil {
		client.SendNumeric(utils.ERR_NOSUCHNICK, targetNick+" :No such nick/channel")
		return
	}

	if !channel.HasClient(targetClient) {
		client.SendNumeric(utils.ERR_USERNOTINCHANNEL, targetNick+" "+channelName+" :They aren't on that channel")
		return
	}

	kickMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " KICK " + channelName + " " + targetNick + " :" + reason
	channel.Broadcast(kickMsg)

	channel.RemoveClient(targetClient)
}

func (c *KickCommand) Help() string {
	return "KICK <channel> <nickname> [<reason>] - Kicks a user from a channel"
}
