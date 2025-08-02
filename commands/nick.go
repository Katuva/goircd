package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type NickCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("NICK", func(s *server.Server) server.Command {
		return &NickCommand{server: s}
	})
}

func (c *NickCommand) Name() string {
	return "NICK"
}

func (c *NickCommand) Execute(client *server.Client, params string) {
	if params == "" {
		client.SendNumeric(utils.ERR_NONICKNAMEGIVEN, ":No nickname given")
		return
	}

	nickname := strings.Split(params, " ")[0]

	if !utils.IsValidNickname(nickname) {
		client.SendNumeric(utils.ERR_ERRONEUSNICKNAME, nickname+" :Erroneous nickname")
		return
	}

	if existingClient := c.server.GetClient(nickname); existingClient != nil && existingClient != client {
		client.SendNumeric(utils.ERR_NICKNAMEINUSE, nickname+" :Nickname is already in use")
		return
	}

	oldNick := client.Nick

	client.SetNick(nickname)

	c.server.AddClient(client)

	if oldNick == "" && client.User != "" && !client.IsRegistered() {
		client.SetRegistered()
		server.SendWelcomeMessages(client)
	}
}

func (c *NickCommand) Help() string {
	return "NICK <nickname> - Sets your nickname"
}
