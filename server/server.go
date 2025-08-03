package server

import (
	"bufio"
	"errors"
	"fmt"
	"goircd/config"
	"goircd/logger"
	"goircd/utils"
	"net"
	"os"
	"strings"
	"sync"
)

type Server struct {
	host           string
	port           int
	listener       net.Listener
	clients        map[string]*Client
	channels       map[string]*Channel
	commands       map[string]Command
	mu             sync.RWMutex
	shutdown       chan struct{}
	waitGroup      sync.WaitGroup
	shutdownCmd    chan string
	isShuttingDown bool
	shutdownOnce   sync.Once
}

func NewServer(host string, port int) (*Server, error) {
	server := &Server{
		host:        host,
		port:        port,
		clients:     make(map[string]*Client),
		channels:    make(map[string]*Channel),
		commands:    make(map[string]Command),
		shutdown:    make(chan struct{}),
		shutdownCmd: make(chan string, 1),
	}

	if err := server.loadCommands(); err != nil {
		logger.Fatal("failed to load commands: %w", err)
	}

	go server.handleShutdown()

	return server, nil
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener
	logger.Info("Server listening on %s", addr)

	go s.acceptConnections()

	return nil
}

func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.shutdown:
					return
				default:
					logger.Error("Error accepting connection: %v", err)
					continue
				}
			}

			s.waitGroup.Add(1)
			go s.handleClient(conn)
		}
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer s.waitGroup.Done()
	defer conn.Close()

	client := NewClient(conn, s)

	logger.Info("Client connecting from %s", client.Host)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		s.processCommand(client, line)
	}

	if err := scanner.Err(); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			logger.Error("Error reading from client: %v", err)
		}
	}

	s.RemoveClient(client)
}

func (s *Server) processCommand(client *Client, line string) {
	if line == "" {
		return
	}

	parts := strings.SplitN(line, " ", 2)
	cmdName := strings.ToUpper(parts[0])
	var params string
	if len(parts) > 1 {
		params = parts[1]
	}

	s.mu.RLock()
	cmd, exists := s.commands[cmdName]
	s.mu.RUnlock()

	if exists {
		cmd.Execute(client, params)
	} else {
		client.SendNumeric(421, fmt.Sprintf("%s :Unknown command", cmdName))
	}
}

func (s *Server) RegisterCommand(name string, cmd Command) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commands[name] = cmd
	logger.Debug("Registered command: %s", name)
}

func (s *Server) loadCommands() error {
	for _, initFunc := range CommandRegistry {
		cmd := initFunc(s)
		s.RegisterCommand(strings.ToUpper(cmd.Name()), cmd)
	}

	return nil
}

func (s *Server) AddClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client.Nick] = client
}

func (s *Server) RemoveClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, channel := range s.channels {
		channel.RemoveClient(client)
	}

	if client.Nick != "" {
		delete(s.clients, client.Nick)
	}

	logger.Info("Client disconnected: %s", client.Nick)
}

func (s *Server) GetClient(nick string) *Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clients[nick]
}

func (s *Server) GetChannel(name string) *Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.channels[name]
}

func (s *Server) GetAllChannels() []*Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	channels := make([]*Channel, 0, len(s.channels))
	for _, channel := range s.channels {
		channels = append(channels, channel)
	}

	return channels
}

func (s *Server) CreateChannel(name string) *Channel {
	s.mu.Lock()
	defer s.mu.Unlock()

	channel := NewChannel(name)
	s.channels[name] = channel
	return channel
}

func (s *Server) handleShutdown() {
	reason := <-s.shutdownCmd
	s.doShutdown(reason)
}

func (s *Server) doShutdown(reason ...string) {
	s.shutdownOnce.Do(func() {
		logger.Info("Server shutting down...")

		var shutdownReason string
		if len(reason) > 0 {
			shutdownReason = "Server is shutting down: " + reason[0]
		} else {
			shutdownReason = "Server is shutting down"
		}

		close(s.shutdown)

		if s.listener != nil {
			s.listener.Close()
		}

		shutdownMsg := ":" + utils.SERVER_NAME + " NOTICE * :" + shutdownReason

		s.mu.RLock()
		for _, client := range s.clients {
			client.Send(shutdownMsg)
			client.Close()
		}
		s.mu.RUnlock()

		s.waitGroup.Wait()

		logger.Info("Server shutdown complete")

		os.Exit(0)
	})
}

func (s *Server) ShutdownAsync(reason ...string) {
	var shutdownReason string
	if len(reason) > 0 {
		shutdownReason = reason[0]
	}

	select {
	case s.shutdownCmd <- shutdownReason:
		// Shutdown initiated
	default:
		// Shutdown already in progress
	}
}

// Shutdown initiates a synchronous shutdown (for signal handlers)
func (s *Server) Shutdown() {
	s.doShutdown()
}

func SendWelcomeMessages(client *Client) {
	client.SendNumeric(utils.RPL_WELCOME, "Welcome to the Internet Relay Network "+
		utils.FormatUserMask(client.Nick, client.User, client.Host))

	client.SendNumeric(utils.RPL_YOURHOST, "Your host is "+utils.SERVER_NAME+
		", running version "+utils.SERVER_VERSION)

	client.SendNumeric(utils.RPL_CREATED, "This server was created "+utils.SERVER_CREATED)

	client.SendNumeric(utils.RPL_MYINFO, utils.SERVER_NAME+" "+utils.SERVER_VERSION+
		" iwr oiw")

	cfg := config.Get()
	if cfg.MOTD != "" {
		client.SendNumeric(utils.RPL_MOTDSTART, ":- "+utils.SERVER_NAME+" Message of the day - ")

		formattedMOTD := utils.FormatMOTD(cfg.MOTD)

		motdLines := strings.Split(strings.TrimSpace(formattedMOTD), "\n")
		for _, line := range motdLines {
			client.SendNumeric(utils.RPL_MOTD, ":- "+line)
		}

		client.SendNumeric(utils.RPL_ENDOFMOTD, ":End of /MOTD command")
	}

	logger.Info("Client %s from %s registered", client.Nick, client.Host)
}
