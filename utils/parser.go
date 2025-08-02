package utils

import (
	"strings"
)

func ParseMessage(message string) (string, string, []string) {
	var prefix string
	var command string
	var params []string

	message = strings.TrimSuffix(message, "\r\n")

	if strings.HasPrefix(message, ":") {
		parts := strings.SplitN(message, " ", 2)
		prefix = parts[0][1:]
		if len(parts) > 1 {
			message = parts[1]
		} else {
			message = ""
		}
	}

	parts := strings.SplitN(message, " ", 2)
	if len(parts) > 0 {
		command = strings.ToUpper(parts[0])
	}

	if len(parts) > 1 {
		paramStr := parts[1]

		if idx := strings.Index(paramStr, " :"); idx != -1 {
			if idx > 0 {
				params = strings.Split(paramStr[:idx], " ")
			}
			params = append(params, paramStr[idx+2:])
		} else {
			params = strings.Split(paramStr, " ")
		}
	}

	return prefix, command, params
}

func FormatMessage(prefix, command string, params []string) string {
	var message strings.Builder

	if prefix != "" {
		message.WriteString(":")
		message.WriteString(prefix)
		message.WriteString(" ")
	}

	message.WriteString(command)

	for i, param := range params {
		message.WriteString(" ")

		if i == len(params)-1 && (strings.Contains(param, " ") || strings.HasPrefix(param, ":")) {
			message.WriteString(":")
			message.WriteString(param)
		} else {
			message.WriteString(param)
		}
	}

	message.WriteString("\r\n")

	return message.String()
}

func ParseUserMask(mask string) (string, string, string) {
	var nick, user, host string

	parts := strings.SplitN(mask, "!", 2)
	if len(parts) > 0 {
		nick = parts[0]
	}

	if len(parts) > 1 {
		userHost := parts[1]
		parts = strings.SplitN(userHost, "@", 2)

		if len(parts) > 0 {
			user = parts[0]
		}

		if len(parts) > 1 {
			host = parts[1]
		}
	}

	return nick, user, host
}

func FormatUserMask(nick, user, host string) string {
	if user != "" && host != "" {
		return nick + "!" + user + "@" + host
	}
	return nick
}

func IsValidChannelName(name string) bool {
	if name == "" {
		return false
	}

	if !strings.HasPrefix(name, "#") && !strings.HasPrefix(name, "&") {
		return false
	}

	if strings.ContainsAny(name, " ,\x00\x07\x0D\x0A") {
		return false
	}

	return true
}

func IsValidNickname(nick string) bool {
	if nick == "" {
		return false
	}

	if strings.HasPrefix(nick, "#") || strings.HasPrefix(nick, "&") ||
		strings.HasPrefix(nick, ":") || strings.HasPrefix(nick, "$") {
		return false
	}

	if strings.ContainsAny(nick, " ,\x00\x07\x0D\x0A") {
		return false
	}

	return true
}
