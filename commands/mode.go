package commands

import (
	"fmt"
	"strconv"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type ModeCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("MODE", func(s *server.Server) server.Command {
		return &ModeCommand{server: s}
	})
}

func (c *ModeCommand) Name() string {
	return "MODE"
}

func (c *ModeCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if params == "" {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "MODE :Not enough parameters")
		return
	}

	parts := strings.SplitN(params, " ", 2)
	target := parts[0]

	if strings.HasPrefix(target, "#") || strings.HasPrefix(target, "&") {
		c.handleChannelMode(client, target, parts)
	} else {
		c.handleUserMode(client, target, parts)
	}
}

func (c *ModeCommand) handleChannelMode(client *server.Client, channelName string, parts []string) {
	channel := c.server.GetChannel(channelName)
	if channel == nil {
		client.SendNumeric(utils.ERR_NOSUCHCHANNEL, channelName+" :No such channel")
		return
	}

	if len(parts) == 1 {
		c.showChannelModes(client, channel)
		return
	}

	modeString := parts[1]
	modeParams := ""
	if strings.Contains(modeString, " ") {
		modeParts := strings.SplitN(modeString, " ", 2)
		modeString = modeParts[0]
		modeParams = modeParts[1]
	}

	if (modeString == "+b" || modeString == "b") && modeParams == "" {

		if !channel.HasClient(client) && !client.IsOperator() {
			client.SendNumeric(utils.ERR_NOTONCHANNEL, channelName+" :You're not on that channel")
			return
		}
		c.showBanList(client, channel)
		return
	}

	if !channel.IsOperator(client) {
		client.SendNumeric(utils.ERR_CHANOPRIVSNEEDED, channelName+" :You're not channel operator")
		return
	}

	adding := true
	paramIndex := 0
	paramList := strings.Split(modeParams, " ")
	var modeChanges strings.Builder
	var appliedModes strings.Builder
	var appliedParams []string

	for _, char := range modeString {
		if char == '+' {
			adding = true
			continue
		}
		if char == '-' {
			adding = false
			continue
		}

		switch char {
		case 'o':
			if paramIndex >= len(paramList) {
				continue
			}
			targetNick := paramList[paramIndex]
			paramIndex++

			targetClient := c.server.GetClient(targetNick)
			if targetClient == nil {
				client.SendNumeric(utils.ERR_NOSUCHNICK, targetNick+" :No such nick/channel")
				continue
			}

			if !channel.HasClient(targetClient) {
				client.SendNumeric(utils.ERR_USERNOTINCHANNEL, targetNick+" "+channelName+" :They aren't on that channel")
				continue
			}

			channel.SetOperator(targetClient, adding)
			if adding {
				modeChanges.WriteString(fmt.Sprintf("+%c", char))
			} else {
				modeChanges.WriteString(fmt.Sprintf("-%c", char))
			}
			appliedModes.WriteString(string(char))
			appliedParams = append(appliedParams, targetNick)
		case 'v':
			if paramIndex >= len(paramList) {
				continue
			}
			targetNick := paramList[paramIndex]
			paramIndex++

			targetClient := c.server.GetClient(targetNick)
			if targetClient == nil {
				client.SendNumeric(utils.ERR_NOSUCHNICK, targetNick+" :No such nick/channel")
				continue
			}

			if !channel.HasClient(targetClient) {
				client.SendNumeric(utils.ERR_USERNOTINCHANNEL, targetNick+" "+channelName+" :They aren't on that channel")
				continue
			}

			channel.SetVoiced(targetClient, adding)
			if adding {
				modeChanges.WriteString(fmt.Sprintf("+%c", char))
			} else {
				modeChanges.WriteString(fmt.Sprintf("-%c", char))
			}
			appliedModes.WriteString(string(char))
			appliedParams = append(appliedParams, targetNick)
		case 'b':
			if paramIndex >= len(paramList) {

				if adding {
					c.showBanList(client, channel)
				}
				continue
			}
			mask := paramList[paramIndex]
			paramIndex++

			channel.SetBanned(mask, adding)
			if adding {
				modeChanges.WriteString(fmt.Sprintf("+%c", char))
			} else {
				modeChanges.WriteString(fmt.Sprintf("-%c", char))
			}
			appliedModes.WriteString(string(char))
			appliedParams = append(appliedParams, mask)
		case 'k':
			if adding {
				if paramIndex >= len(paramList) {
					client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "MODE :Not enough parameters")
					continue
				}
				key := paramList[paramIndex]
				paramIndex++

				if channel.Key != "" {
					client.SendNumeric(utils.ERR_KEYSET, channelName+" :Channel key already set")
					continue
				}

				channel.Key = key
				modeChanges.WriteString(fmt.Sprintf("+%c", char))
				appliedModes.WriteString(string(char))
				appliedParams = append(appliedParams, key)
			} else {
				channel.Key = ""
				modeChanges.WriteString(fmt.Sprintf("-%c", char))
				appliedModes.WriteString(string(char))
			}
		case 'l':
			if adding {
				if paramIndex >= len(paramList) {
					client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "MODE :Not enough parameters")
					continue
				}
				limitStr := paramList[paramIndex]
				paramIndex++

				limit, err := strconv.Atoi(limitStr)
				if err != nil || limit < 0 {
					client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "MODE :Invalid limit")
					continue
				}

				channel.Limit = limit
				modeChanges.WriteString(fmt.Sprintf("+%c", char))
				appliedModes.WriteString(string(char))
				appliedParams = append(appliedParams, limitStr)
			} else {
				channel.Limit = 0
				modeChanges.WriteString(fmt.Sprintf("-%c", char))
				appliedModes.WriteString(string(char))
			}

		case 'i', 'm', 'n', 'p', 's', 't':
			var mode server.ChannelMode

			switch char {
			case 'i':
				mode = server.ModeInviteOnly
			case 'm':
				mode = server.ModeModerated
			case 'n':
				mode = server.ModeNoExternalMessages
			case 'p':
				mode = server.ModePrivate
			case 's':
				mode = server.ModeSecret
			case 't':
				mode = server.ModeTopicSettableByOpsOnly
			}

			channel.SetMode(mode, adding)
			if adding {
				modeChanges.WriteString(fmt.Sprintf("+%c", char))
			} else {
				modeChanges.WriteString(fmt.Sprintf("-%c", char))
			}
			appliedModes.WriteString(string(char))

		default:
			client.SendNumeric(utils.ERR_UNKNOWNMODE, string(char)+" :is unknown mode char to me for "+channelName)
		}
	}

	if modeChanges.Len() > 0 {
		modeMsg := ":" + utils.FormatUserMask(client.Nick, client.User, client.Host) + " MODE " + channelName + " " + modeChanges.String()

		if len(appliedParams) > 0 {
			modeMsg += " " + strings.Join(appliedParams, " ")
		}

		channel.Broadcast(modeMsg)
	}
}

func (c *ModeCommand) handleUserMode(client *server.Client, target string, parts []string) {
	if target != client.Nick {
		client.SendNumeric(utils.ERR_USERSDONTMATCH, ":Cannot change mode for other users")
		return
	}

	if len(parts) == 1 {
		c.showUserModes(client)
		return
	}

	modeString := parts[1]

	adding := true
	var modeChanges strings.Builder

	for _, char := range modeString {
		if char == '+' {
			adding = true
			continue
		}
		if char == '-' {
			adding = false
			continue
		}

		switch char {
		case 'i':
			if adding {
				modeChanges.WriteString(fmt.Sprintf("+%c", char))
			} else {
				modeChanges.WriteString(fmt.Sprintf("-%c", char))
			}
		case 'w':
			if adding {
				modeChanges.WriteString(fmt.Sprintf("+%c", char))
			} else {
				modeChanges.WriteString(fmt.Sprintf("-%c", char))
			}
		case 'o':
			if adding {
				client.SendNumeric(utils.ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
				continue
			}

			client.SetOperator(false)
			modeChanges.WriteString(fmt.Sprintf("-%c", char))
		default:
			client.SendNumeric(utils.ERR_UMODEUNKNOWNFLAG, ":Unknown MODE flag")
		}
	}

	if modeChanges.Len() > 0 {
		client.SendMessage(utils.FormatUserMask(client.Nick, client.User, client.Host), "MODE", client.Nick, modeChanges.String())
	}
}

func (c *ModeCommand) showChannelModes(client *server.Client, channel *server.Channel) {
	var modes strings.Builder
	var params []string

	if channel.HasMode(server.ModeInviteOnly) {
		modes.WriteString("i")
	}
	if channel.HasMode(server.ModeModerated) {
		modes.WriteString("m")
	}
	if channel.HasMode(server.ModeNoExternalMessages) {
		modes.WriteString("n")
	}
	if channel.HasMode(server.ModePrivate) {
		modes.WriteString("p")
	}
	if channel.HasMode(server.ModeSecret) {
		modes.WriteString("s")
	}
	if channel.HasMode(server.ModeTopicSettableByOpsOnly) {
		modes.WriteString("t")
	}

	if channel.Key != "" {
		modes.WriteString("k")
		params = append(params, channel.Key)
	}
	if channel.Limit > 0 {
		modes.WriteString("l")
		params = append(params, strconv.Itoa(channel.Limit))
	}

	modeString := "+"
	if modes.Len() > 0 {
		modeString += modes.String()
	}

	if len(params) > 0 {
		client.SendNumeric(utils.RPL_CHANNELMODEIS, channel.Name+" "+modeString+" "+strings.Join(params, " "))
	} else {
		client.SendNumeric(utils.RPL_CHANNELMODEIS, channel.Name+" "+modeString)
	}
}

func (c *ModeCommand) showBanList(client *server.Client, channel *server.Channel) {
	banList := channel.GetBanList()

	for _, banMask := range banList {
		client.SendNumeric(utils.RPL_BANLIST, channel.Name+" "+banMask)
	}

	client.SendNumeric(utils.RPL_ENDOFBANLIST, channel.Name+" :End of channel ban list")
}

func (c *ModeCommand) showUserModes(client *server.Client) {
	var modes strings.Builder

	if client.IsOperator() {
		modes.WriteString("o")
	}

	modeString := "+"
	if modes.Len() > 0 {
		modeString += modes.String()
	}

	client.SendNumeric(utils.RPL_UMODEIS, modeString)
}

func (c *ModeCommand) Help() string {
	return "MODE <channel> [<mode> [<mode parameters>]] - Changes or displays channel modes\n" +
		"MODE <nickname> [<mode>] - Changes or displays user modes"
}
