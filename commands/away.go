package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type AwayCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("AWAY", func(s *server.Server) server.Command {
		return &AwayCommand{server: s}
	})
}

func (c *AwayCommand) Name() string {
	return "AWAY"
}

func (c *AwayCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if params == "" {
		client.SetAway(false)
		client.SendNumeric(utils.RPL_UNAWAY, ":You are no longer marked as being away")
		return
	}

	awayMessage := strings.TrimSpace(params)
	if strings.HasPrefix(awayMessage, ":") {
		awayMessage = awayMessage[1:]
	}

	client.SetAway(true)
	client.AwayMessage = awayMessage
	client.SendNumeric(utils.RPL_NOWAWAY, ":You have been marked as being away")
}

func (c *AwayCommand) Help() string {
	return "AWAY [:<message>] - Sets or removes away status"
}
