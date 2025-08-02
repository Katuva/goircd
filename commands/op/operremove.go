package op

import (
	"goircd/logger"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type OperRemoveCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("OPERREMOVE", func(s *server.Server) server.Command {
		return &OperRemoveCommand{server: s}
	})
}

func (c *OperRemoveCommand) Name() string {
	return "OPERREMOVE"
}

func (c *OperRemoveCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if !client.IsOperator() {
		client.SendNumeric(utils.ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	if len(parts) < 2 {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "OPERREMOVE :Not enough parameters")
		return
	}

	channelName := parts[0]
	targetNick := parts[1]

	channel := c.server.GetChannel(channelName)
	if channel == nil {
		client.SendNumeric(utils.ERR_NOSUCHCHANNEL, channelName+" :No such channel")
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

	logger.IRCOp(client.Nick, "removed operator status from", targetNick, channelName)

	modeMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " MODE " + channelName + " -o " + targetNick
	channel.Broadcast(modeMsg)

	targetClient.Send(":" + utils.SERVER_NAME + " NOTICE " + targetNick + " :Your channel operator status in " + channelName + " has been removed by IRC operator " + client.Nick)
}

func (c *OperRemoveCommand) Help() string {
	return "OPERREMOVE <channel> <nickname> - IRC operator command to remove channel operator status from a user"
}
