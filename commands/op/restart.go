package op

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"goircd/server"
	"goircd/utils"
)

type RestartCommand struct {
	server *server.Server
}

func init() {
	server.RegisterCommandInit("RESTART", func(s *server.Server) server.Command {
		return &RestartCommand{server: s}
	})
}

func (c *RestartCommand) Name() string {
	return "RESTART"
}

func (c *RestartCommand) Execute(client *server.Client, params string) {
	if !client.IsRegistered() {
		client.SendNumeric(utils.ERR_NOTREGISTERED, ":You have not registered")
		return
	}

	if !client.IsOperator() {
		client.SendNumeric(utils.ERR_NOPRIVILEGES, ":Permission Denied- You're not an IRC operator")
		return
	}

	reason := "Server restart"
	if params != "" {
		if strings.HasPrefix(params, ":") {
			reason = params[1:]
		} else {
			reason = params
		}
	}

	restartMsg := ":" + utils.SERVER_NAME + " NOTICE * :Server is restarting: " + reason

	notifiedClients := make(map[*server.Client]bool)
	for _, channel := range c.server.GetAllChannels() {
		for _, client := range channel.GetClients() {
			if !notifiedClients[client] {
				client.Send(restartMsg)
				notifiedClients[client] = true
			}
		}
	}

	time.Sleep(1 * time.Second)

	c.restartServer()
}

func (c *RestartCommand) restartServer() {
	execPath, err := os.Executable()
	if err != nil {
		os.Exit(0)
		return
	}

	workDir, err := os.Getwd()
	if err != nil {
		workDir = filepath.Dir(execPath)
	}

	cmd := exec.Command(execPath)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Start()
	if err != nil {
		os.Exit(0)
		return
	}

	os.Exit(0)
}

func (c *RestartCommand) Help() string {
	return "RESTART [reason] - Restarts the server"
}
