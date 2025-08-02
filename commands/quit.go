package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type QuitCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("QUIT", func(s *server.Server) server.Command {
		return &QuitCommand{server: s}
	})
}

func (c *QuitCommand) Name() string {
	return "QUIT"
}

func (c *QuitCommand) Execute(client *server.Client, params string) {
	var quitMessage string

	if params != "" {
		quitMessage = params
		if strings.HasPrefix(quitMessage, ":") {
			quitMessage = quitMessage[1:]
		}
	} else {
		quitMessage = client.Nick
	}

	quitMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " QUIT"
	if quitMessage != "" {
		quitMsg += " :" + quitMessage
	}

	for _, channel := range client.GetChannels() {
		channel.BroadcastFrom(client, quitMsg)
		channel.RemoveClient(client)
	}

	client.Close()
}

func (c *QuitCommand) Help() string {
	return "QUIT [<message>] - Quit the IRC server with an optional message"
}
