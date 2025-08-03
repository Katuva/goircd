package commands

import (
	"fmt"
	"goircd/config"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type JoinCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("JOIN", func(s *server.Server) server.Command {
		return &JoinCommand{server: s}
	})
}

func (c *JoinCommand) Name() string {
	return "JOIN"
}

func (c *JoinCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if params == "" {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "JOIN :Not enough parameters")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	channelNames := strings.Split(parts[0], ",")

	var keys []string
	if len(parts) > 1 {
		keys = strings.Split(parts[1], ",")
	}

	cfg := config.Get()

	for i, channelName := range channelNames {
		if !utils.IsValidChannelName(channelName) {
			client.SendNumeric(utils.ERR_NOSUCHCHANNEL, channelName+" :No such channel")
			continue
		}

		if len(client.GetChannels()) > cfg.Channels.MaxChannels {
			client.SendNumeric(utils.ERR_TOOMANYCHANNELS, channelName+" :You have joined too many channels")
			continue
		}

		var key string
		if i < len(keys) {
			key = keys[i]
		}

		channel := c.server.GetChannel(channelName)
		isNewChannel := false

		if channel == nil {
			if len(channelName) > cfg.Security.MaxChannelName {
				client.SendNumeric(utils.ERR_BADCHANMASK, channelName+" :Channel name too long")
				return
			}

			channel = c.server.CreateChannel(channelName)
			isNewChannel = true
		} else {
			if channel.HasMode(server.ModeInviteOnly) && !channel.IsInvited(client.Nick) && !channel.IsOperator(client) {
				client.SendNumeric(utils.ERR_INVITEONLYCHAN, channelName+" :Cannot join channel (+i)")
				continue
			}

			if channel.IsBanned(client) {
				client.SendNumeric(utils.ERR_BANNEDFROMCHAN, channelName+" :Cannot join channel (+b)")
				continue
			}

			if channel.Key != "" && key != channel.Key && !channel.IsOperator(client) {
				client.SendNumeric(utils.ERR_BADCHANNELKEY, channelName+" :Cannot join channel (+k)")
				continue
			}

			if len(channel.GetClients()) >= cfg.Channels.MaxUsersPerChannel {
				client.SendNumeric(utils.ERR_CHANNELISFULL, channelName+" :Cannot join channel (+l)")
				continue
			}
		}

		channel.AddClient(client)

		if isNewChannel {
			channel.SetOperator(client, true)

			if cfg.Channels.DefaultModes != "" {
				channel.ApplyDefaultModes(cfg.Channels.DefaultModes)

				channel.Broadcast(fmt.Sprintf(":%s MODE %s %s", client.Nick, channelName, cfg.Channels.DefaultModes))
			}
		}

		joinMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " JOIN " + channelName
		channel.Broadcast(joinMsg)

		topic, setBy, setAt := channel.GetTopic()
		if topic != "" {
			client.SendNumeric(utils.RPL_TOPIC, channelName+" :"+topic)
			client.SendNumeric(utils.RPL_TOPICWHOTIME, channelName+" "+setBy+" "+fmt.Sprintf("%d", setAt.Unix()))
		} else {
			client.SendNumeric(utils.RPL_NOTOPIC, channelName+" :No topic is set")
		}

		sendNameReply(client, channel)
	}
}

func (c *JoinCommand) Help() string {
	return "JOIN <channel>{,<channel>} [<key>{,<key>}] - Joins the specified channels"
}

func sendNameReply(client *server.Client, channel *server.Channel) {
	var nickList strings.Builder

	for _, member := range channel.GetClients() {
		if nickList.Len() > 0 {
			nickList.WriteString(" ")
		}

		if channel.IsOperator(member) {
			nickList.WriteString("@")
		} else if channel.IsVoiced(member) {
			nickList.WriteString("+")
		}

		nickList.WriteString(member.Nick)
	}

	chanPrefix := "="
	if channel.HasMode(server.ModeSecret) {
		chanPrefix = "@"
	} else if channel.HasMode(server.ModePrivate) {
		chanPrefix = "*"
	}

	client.SendNumeric(utils.RPL_NAMEREPLY, chanPrefix+" "+channel.Name+" :"+nickList.String())
	client.SendNumeric(utils.RPL_ENDOFNAMES, channel.Name+" :End of NAMES list")
}
