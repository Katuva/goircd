package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type PrivmsgCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("PRIVMSG", func(s *server.Server) server.Command {
		return &PrivmsgCommand{server: s}
	})
}

func (c *PrivmsgCommand) Name() string {
	return "PRIVMSG"
}

func (c *PrivmsgCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if params == "" {
		client.SendNumeric(utils.ERR_NORECIPIENT, ":No recipient given (PRIVMSG)")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	if len(parts) < 2 {
		client.SendNumeric(utils.ERR_NOTEXTTOSEND, ":No text to send")
		return
	}

	target := parts[0]
	message := parts[1]

	if strings.HasPrefix(message, ":") {
		message = message[1:]
	}

	if message == "" {
		client.SendNumeric(utils.ERR_NOTEXTTOSEND, ":No text to send")
		return
	}

	formattedMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " PRIVMSG " + target + " :" + message

	if strings.HasPrefix(target, "#") || strings.HasPrefix(target, "&") {
		channel := c.server.GetChannel(target)
		if channel == nil {
			client.SendNumeric(utils.ERR_NOSUCHCHANNEL, target+" :No such channel")
			return
		}

		if !channel.HasClient(client) && channel.HasMode(server.ModeNoExternalMessages) {
			client.SendNumeric(utils.ERR_CANNOTSENDTOCHAN, target+" :Cannot send to channel")
			return
		}

		if channel.HasMode(server.ModeModerated) && !channel.IsVoiced(client) && !channel.IsOperator(client) {
			client.SendNumeric(utils.ERR_CANNOTSENDTOCHAN, target+" :Cannot send to channel (+m)")
			return
		}

		channel.BroadcastFrom(client, formattedMsg)
	} else {
		targetClient := c.server.GetClient(target)
		if targetClient == nil {
			client.SendNumeric(utils.ERR_NOSUCHNICK, target+" :No such nick/channel")
			return
		}

		targetClient.Send(formattedMsg)
	}
}

func (c *PrivmsgCommand) Help() string {
	return "PRIVMSG <target> :<message> - Sends a message to a user or channel"
}
