package commands

import (
	"goircd/logger"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type DeOpCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("DEOP", func(s *server.Server) server.Command {
		return &DeOpCommand{server: s}
	})
}

func (c *DeOpCommand) Name() string {
	return "DEOP"
}

func (c *DeOpCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	if len(parts) < 2 {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "DEOP :Not enough parameters")
		return
	}

	channelName := parts[0]
	targetNick := parts[1]

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

	if !channel.IsOperator(targetClient) {
		client.SendNumeric(utils.ERR_USERNOTINCHANNEL, targetNick+" "+channelName+" :They aren't a channel operator")
		return
	}

	channel.SetOperator(targetClient, false)

	logger.ChannelOp(client.Nick, "removed operator status from", targetNick, channelName)

	modeMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " MODE " + channelName + " -o " + targetNick
	channel.Broadcast(modeMsg)
}

func (c *DeOpCommand) Help() string {
	return "DEOP <channel> <nickname> - Removes channel operator status from a user"
}
