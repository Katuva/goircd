package commands

import (
	"strings"
	"time"

	"goircd/server"
	"goircd/utils"
)

type TopicCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("TOPIC", func(s *server.Server) server.Command {
		return &TopicCommand{server: s}
	})
}

func (c *TopicCommand) Name() string {
	return "TOPIC"
}

func (c *TopicCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if params == "" {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "TOPIC :Not enough parameters")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	channelName := parts[0]

	channel := c.server.GetChannel(channelName)
	if channel == nil {
		client.SendNumeric(utils.ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	if !channel.HasClient(client) {
		client.SendNumeric(utils.ERR_NOTONCHANNEL, channelName+" :You're not on that channel")
		return
	}

	if len(parts) == 1 {
		topic, setBy, setAt := channel.GetTopic()
		if topic == "" {
			client.SendNumeric(utils.RPL_NOTOPIC, channelName+" :No topic is set")
		} else {
			client.SendNumeric(utils.RPL_TOPIC, channelName+" :"+topic)
			client.SendNumeric(utils.RPL_TOPICWHOTIME, channelName+" "+setBy+" "+time.Unix(setAt.Unix(), 0).Format(time.RFC3339))
		}
		return
	}

	if channel.HasMode(server.ModeTopicSettableByOpsOnly) && !channel.IsOperator(client) {
		client.SendNumeric(utils.ERR_CHANOPRIVSNEEDED, channelName+" :You're not channel operator")
		return
	}

	newTopic := parts[1]
	if strings.HasPrefix(newTopic, ":") {
		newTopic = newTopic[1:]
	}

	if len(newTopic) > utils.MAX_TOPIC_LENGTH {
		newTopic = newTopic[:utils.MAX_TOPIC_LENGTH]
	}

	channel.SetTopic(newTopic, client.Nick)

	topicMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " TOPIC " + channelName + " :" + newTopic
	channel.Broadcast(topicMsg)
}

func (c *TopicCommand) Help() string {
	return "TOPIC <channel> [<topic>] - Changes or displays the topic of a channel"
}
