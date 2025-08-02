package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type PartCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("PART", func(s *server.Server) server.Command {
		return &PartCommand{server: s}
	})
}

func (c *PartCommand) Name() string {
	return "PART"
}

func (c *PartCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if params == "" {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "PART :Not enough parameters")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	channelNames := strings.Split(parts[0], ",")

	var partMessage string
	if len(parts) > 1 {
		partMessage = parts[1]
		if strings.HasPrefix(partMessage, ":") {
			partMessage = partMessage[1:]
		}
	} else {
		partMessage = client.Nick
	}

	for _, channelName := range channelNames {
		channel := c.server.GetChannel(channelName)
		if channel == nil {
			client.SendNumeric(utils.ERR_NOSUCHCHANNEL, channelName+" :No such channel")
			continue
		}

		if !channel.HasClient(client) {
			client.SendNumeric(utils.ERR_NOTONCHANNEL, channelName+" :You're not on that channel")
			continue
		}

		partMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " PART " + channelName
		if partMessage != "" {
			partMsg += " :" + partMessage
		}
		channel.Broadcast(partMsg)

		channel.RemoveClient(client)
	}
}

func (c *PartCommand) Help() string {
	return "PART <channel>{,<channel>} [<reason>] - Leaves the specified channels"
}
