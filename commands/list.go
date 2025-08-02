package commands

import (
	"strings"

	"goircd/server"
	"goircd/utils"
)

type ListCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("LIST", func(s *server.Server) server.Command {
		return &ListCommand{server: s}
	})
}

func (c *ListCommand) Name() string {
	return "LIST"
}

func (c *ListCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	client.SendNumeric(utils.RPL_LISTSTART, "Channel :Users  Name")

	channels := c.server.GetAllChannels()

	if params != "" {
		channelNames := strings.Split(params, ",")
		filteredChannels := make([]*server.Channel, 0)

		for _, channelName := range channelNames {
			channelName = strings.TrimSpace(channelName)
			if channel := c.server.GetChannel(channelName); channel != nil {
				filteredChannels = append(filteredChannels, channel)
			}
		}
		channels = filteredChannels
	}

	for _, channel := range channels {
		if channel.HasMode(server.ModeSecret) && !channel.HasClient(client) {
			continue
		}

		userCount := len(channel.GetClients())
		topic, _, _ := channel.GetTopic()

		listLine := channel.Name + " " + string(rune(userCount+'0'))
		if topic != "" {
			listLine += " :" + topic
		}

		client.SendNumeric(utils.RPL_LIST, listLine)
	}

	client.SendNumeric(utils.RPL_LISTEND, ":End of LIST")
}

func (c *ListCommand) Help() string {
	return "LIST [<channels>] - List channels on the server"
}
