package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type UserhostCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("USERHOST", func(s *server.Server) server.Command {
		return &UserhostCommand{server: s}
	})
}

func (c *UserhostCommand) Name() string {
	return "USERHOST"
}

func (c *UserhostCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if params == "" {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "USERHOST :Not enough parameters")
		return
	}

	nicknames := strings.Fields(params)
	if len(nicknames) > 5 {
		nicknames = nicknames[:5]
	}

	var responses []string

	for _, nick := range nicknames {
		targetClient := c.server.GetClient(nick)
		if targetClient != nil {
			operatorMark := ""
			if targetClient.IsOperator() {
				operatorMark = "*"
			}

			awayStatus := "+"

			if targetClient.IsAway() {
				awayStatus = "-"
			}

			response := nick + operatorMark + "=" + awayStatus + targetClient.User + "@" + targetClient.Host
			responses = append(responses, response)
		}
	}

	if len(responses) > 0 {
		client.SendNumeric(utils.RPL_USERHOST, ":"+strings.Join(responses, " "))
	} else {
		client.SendNumeric(utils.RPL_USERHOST, ":")
	}
}

func (c *UserhostCommand) Help() string {
	return "USERHOST <nickname> [<nickname> ...] - Shows host information for specified users"
}
