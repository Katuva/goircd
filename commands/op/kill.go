package op

import (
	"goircd/logger"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type KillCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("KILL", func(s *server.Server) server.Command {
		return &KillCommand{server: s}
	})
}

func (c *KillCommand) Name() string {
	return "KILL"
}

func (c *KillCommand) Execute(client *server.Client, params string) {
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
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "KILL :Not enough parameters")
		return
	}

	targetNick := parts[0]
	reason := parts[1]
	if strings.HasPrefix(reason, ":") {
		reason = reason[1:]
	}

	targetClient := c.server.GetClient(targetNick)
	if targetClient == nil {
		client.SendNumeric(utils.ERR_NOSUCHNICK, targetNick+" :No such nick/channel")
		return
	}

	killMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " KILL " + targetNick + " :" + reason
	targetClient.Send(killMsg)

	quitMsg := ":" + utils.FormatUserMask(targetClient.Nick, targetClient.User, targetClient.Host) + " QUIT :Killed by " + client.Nick + " (" + reason + ")"
	for _, channel := range targetClient.GetChannels() {
		channel.BroadcastFrom(targetClient, quitMsg)
	}

	for _, channel := range targetClient.GetChannels() {
		channel.RemoveClient(targetClient)
	}

	targetClient.Close()

	logger.IRCOp(client.Nick, "killed", targetClient.Nick, "Reason: "+reason)
}

func (c *KillCommand) Help() string {
	return "KILL <nickname> <reason> - Disconnects a client from the server"
}
