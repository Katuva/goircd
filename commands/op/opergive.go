package op

import (
	"goircd/logger"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type OperGiveCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("OPERGIVE", func(s *server.Server) server.Command {
		return &OperGiveCommand{server: s}
	})
}

func (c *OperGiveCommand) Name() string {
	return "OPERGIVE"
}

func (c *OperGiveCommand) Execute(client *server.Client, params string) {
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
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "OPERGIVE :Not enough parameters")
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

	channel.SetOperator(targetClient, true)

	logger.IRCOp(client.Nick, "gave operator status to", targetNick, channelName)

	modeMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " MODE " + channelName + " +o " + targetNick
	channel.Broadcast(modeMsg)

	targetClient.Send(":" + utils.SERVER_NAME + " NOTICE " + targetNick + " :You have been given channel operator status in " + channelName + " by IRC operator " + client.Nick)
}

func (c *OperGiveCommand) Help() string {
	return "OPERGIVE <channel> <nickname> - IRC operator command to give channel operator status to a user"
}
