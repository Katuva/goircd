package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type UserCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("USER", func(s *server.Server) server.Command {
		return &UserCommand{server: s}
	})
}

func (c *UserCommand) Name() string {
	return "USER"
}

func (c *UserCommand) Execute(client *server.Client, params string) {
	if client.IsRegistered() {
		client.SendNumeric(utils.ERR_ALREADYREGISTRED, ":Unauthorized command (already registered)")
		return
	}

	parts := strings.SplitN(params, " ", 4)
	if len(parts) < 4 {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "USER :Not enough parameters")
		return
	}

	username := parts[0]
	// parts[1] and parts[2] are mode and unused in modern IRC
	realname := parts[3]

	if strings.HasPrefix(realname, ":") {
		realname = realname[1:]
	}

	client.SetUser(username, realname)

	if client.Nick != "" && !client.IsRegistered() {
		client.SetRegistered()
		server.SendWelcomeMessages(client)
	}
}

func (c *UserCommand) Help() string {
	return "USER <username> <mode> <unused> :<realname> - Specifies username and real name"
}
