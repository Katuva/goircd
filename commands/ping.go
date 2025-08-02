package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type PingCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("PING", func(s *server.Server) server.Command {
		return &PingCommand{server: s}
	})
}

func (c *PingCommand) Name() string {
	return "PING"
}

func (c *PingCommand) Execute(client *server.Client, params string) {
	if params == "" {
		client.SendNumeric(utils.ERR_NOORIGIN, ":No origin specified")
		return
	}

	token := params
	if strings.Contains(params, " ") {
		token = strings.Split(params, " ")[0]
	}

	client.UpdateLastPing()

	client.Send(":" + utils.SERVER_NAME + " PONG " + utils.SERVER_NAME + " :" + token)
}

func (c *PingCommand) Help() string {
	return "PING <server> - Pings the server to check if it's still connected"
}
