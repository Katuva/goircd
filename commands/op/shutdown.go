package op

import (
	"goircd/logger"
	"goircd/server"
	"goircd/utils"
)

type ShutdownCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("SHUTDOWN", func(s *server.Server) server.Command {
		return &ShutdownCommand{server: s}
	})
}

func (c *ShutdownCommand) Name() string {
	return "SHUTDOWN"
}

func (c *ShutdownCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if !client.IsOperator() {
		client.SendNumeric(utils.ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	logger.Info("IRCOP: User %s has issued a shutdown command", client.Nick)

	if params != "" {
		c.server.ShutdownAsync(params)
	} else {
		c.server.ShutdownAsync()
	}
}

func (c *ShutdownCommand) Help() string {
	return "RESTART [reason] - Restarts the server"
}
