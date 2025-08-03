package op

import (
	"goircd/config"
	"goircd/hash"
	"goircd/logger"
	"strings"

	"goircd/server"
	"goircd/utils"
)

type OperCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("OPER", func(s *server.Server) server.Command {
		return &OperCommand{server: s}
	})
}

func (c *OperCommand) Name() string {
	return "OPER"
}

func (c *OperCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	parts := strings.Split(params, " ")
	if len(parts) < 2 {
		client.SendNumeric(utils.ERR_NEEDMOREPARAMS, "OPER :Not enough parameters")
		return
	}

	username := parts[0]
	password := parts[1]

	cfg := config.Get()

	for _, operator := range cfg.Security.Operators {
		if operator.Nick == username {
			if hash.Verify(password, operator.Password) {
				if utils.MatchesHost(client.Host, operator.Host) {
					client.SetOperator(true)

					if operator.Whois != "" {
						client.Whois = operator.Whois
					}

					if operator.Vhost != "" {
						client.Vhost = operator.Vhost
					}

					client.SendNumeric(utils.RPL_YOUREOPER, ":You are now an IRC operator")
					client.SendMessage(utils.FormatUserMask(client.Nick, client.User, client.Host), "MODE", client.Nick, "+o")
					logger.Info("IRCOP: User %s has logged in using %s", client.Nick, username)
					return
				} else {
					client.SendNumeric(utils.ERR_NOOPERHOST, ":No O-lines for your host")
					logger.Warning("IRCOP: User %s attempted to login with %s from an invalid host %s", client.Nick, username, client.Host)
					return
				}
			}
		}
	}

	client.SendNumeric(utils.ERR_PASSWDMISMATCH, ":Password incorrect")
	logger.Warning("IRCOP: User %s attempted to login to %s with an incorrect password", client.Nick, username)
}

func (c *OperCommand) Help() string {
	return "OPER <username> <password> - Authenticates as an IRC operator"
}
