package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn        net.Conn
	server      *Server
	Nick        string
	User        string
	RealName    string
	Host        string
	Vhost       string
	channels    map[string]*Channel
	registered  bool
	isOperator  bool
	isAway      bool
	Whois       string
	AwayMessage string
	lastPing    time.Time
	mu          sync.RWMutex
	sendMu      sync.Mutex
}

func NewClient(conn net.Conn, server *Server) *Client {
	host, _, _ := net.SplitHostPort(conn.RemoteAddr().String())

	return &Client{
		conn:     conn,
		server:   server,
		Host:     host,
		channels: make(map[string]*Channel),
		lastPing: time.Now(),
	}
}

func (c *Client) Send(message string) {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	if !strings.HasSuffix(message, "\r\n") {
		message = message + "\r\n"
	}

	_, err := c.conn.Write([]byte(message))
	if err != nil {
		// Handle error, possibly disconnect client
	}
}

func (c *Client) SendNumeric(numeric int, message string) {
	prefix := fmt.Sprintf(":%s %03d %s ", c.server.host, numeric, c.Nick)
	if c.Nick == "" {
		prefix = fmt.Sprintf(":%s %03d * ", c.server.host, numeric)
	}
	c.Send(prefix + message)
}

func (c *Client) SendMessage(source, command, target, message string) {
	c.Send(fmt.Sprintf(":%s %s %s :%s", source, command, target, message))
}

func (c *Client) JoinChannel(channel *Channel) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.channels[channel.Name] = channel
}

func (c *Client) LeaveChannel(channelName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.channels, channelName)
}

func (c *Client) IsInChannel(channelName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.channels[channelName]
	return exists
}

func (c *Client) GetChannels() []*Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channels := make([]*Channel, 0, len(c.channels))
	for _, channel := range c.channels {
		channels = append(channels, channel)
	}

	return channels
}

func (c *Client) GetHost() string {
	if c.isOperator && c.Vhost != "" {
		return c.Vhost
	}

	return c.Host
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SetNick(nick string) {
	c.Nick = nick
}

func (c *Client) SetUser(user, realName string) {
	c.User = user
	c.RealName = realName
}

func (c *Client) SetRegistered() {
	c.registered = true
}

func (c *Client) IsRegistered() bool {
	return c.registered
}

func (c *Client) SetOperator(isOp bool) {
	c.isOperator = isOp
}

func (c *Client) IsOperator() bool {
	return c.isOperator
}

func (c *Client) SetAway(isAway bool) {
	c.isAway = isAway
}

func (c *Client) IsAway() bool {
	return c.isAway
}

func (c *Client) UpdateLastPing() {
	c.lastPing = time.Now()
}

func (c *Client) GetIdleTime() time.Duration {
	return time.Since(c.lastPing)
}
