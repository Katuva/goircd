package commands

import (
	"goircd/server"
)

type PongCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("PONG", func(s *server.Server) server.Command {
		return &PongCommand{server: s}
	})
}

func (c *PongCommand) Name() string {
	return "PONG"
}

func (c *PongCommand) Execute(client *server.Client, params string) {
	client.UpdateLastPing()
}

func (c *PongCommand) Help() string {
	return "PONG <server> [<server2>] - Response to a PING message"
}
