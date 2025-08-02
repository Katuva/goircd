package commands

import (
	"goircd/logger"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type OpCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("OP", func(s *server.Server) server.Command {
		return &OpCommand{server: s}
	})
}

func (c *OpCommand) Name() string {
	return "OP"
}

func (c *OpCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	if len(parts) < 2 {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "OP :Not enough parameters")
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

	channel.SetOperator(targetClient, true)

	logger.ChannelOp(client.Nick, "gave operator status to", targetNick, channelName)

	modeMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " MODE " + channelName + " +o " + targetNick
	channel.Broadcast(modeMsg)
}

func (c *OpCommand) Help() string {
	return "OP <channel> <nickname> - Gives channel operator status to a user"
}
