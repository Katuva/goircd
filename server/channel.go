package server

import (
	"goircd/utils"
	"strings"
	"sync"
	"time"
)

type ChannelMode int

const (
	ModePrivate ChannelMode = iota
	ModeSecret
	ModeInviteOnly
	ModeTopicSettableByOpsOnly
	ModeNoExternalMessages
	ModeModerated
	ModeLimit
	ModeKey
	ModeOp
	ModeVoice
	ModeBan
)

type Channel struct {
	Name       string
	Topic      string
	TopicSetBy string
	TopicSetAt time.Time
	Key        string
	Limit      int
	Modes      map[ChannelMode]bool
	clients    map[*Client]bool
	operators  map[*Client]bool
	voiced     map[*Client]bool
	banned     map[string]bool
	inviteList map[string]bool
	mu         sync.RWMutex
	createdAt  time.Time
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name:       name,
		Modes:      make(map[ChannelMode]bool),
		clients:    make(map[*Client]bool),
		operators:  make(map[*Client]bool),
		voiced:     make(map[*Client]bool),
		banned:     make(map[string]bool),
		inviteList: make(map[string]bool),
		createdAt:  time.Now(),
	}
}

func (ch *Channel) ApplyDefaultModes(modeString string) {
	if modeString == "" {
		return
	}

	adding := true
	for _, char := range modeString {
		switch char {
		case '+':
			adding = true
		case '-':
			adding = false
		case 'n':
			ch.SetMode(ModeNoExternalMessages, adding)
		case 't':
			ch.SetMode(ModeTopicSettableByOpsOnly, adding)
		case 's':
			ch.SetMode(ModeSecret, adding)
		case 'p':
			ch.SetMode(ModePrivate, adding)
		case 'i':
			ch.SetMode(ModeInviteOnly, adding)
		case 'm':
			ch.SetMode(ModeModerated, adding)
		}
	}
}

func (ch *Channel) AddClient(client *Client) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.clients[client] = true
	client.JoinChannel(ch)
}

func (ch *Channel) RemoveClient(client *Client) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	delete(ch.clients, client)
	delete(ch.operators, client)
	delete(ch.voiced, client)

	client.LeaveChannel(ch.Name)
}

func (ch *Channel) GetClients() []*Client {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	clients := make([]*Client, 0, len(ch.clients))
	for client := range ch.clients {
		clients = append(clients, client)
	}

	return clients
}

func (ch *Channel) HasClient(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	_, exists := ch.clients[client]
	return exists
}

func (ch *Channel) SetTopic(topic, setBy string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.Topic = topic
	ch.TopicSetBy = setBy
	ch.TopicSetAt = time.Now()
}

func (ch *Channel) GetTopic() (string, string, time.Time) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	return ch.Topic, ch.TopicSetBy, ch.TopicSetAt
}

func (ch *Channel) SetMode(mode ChannelMode, enabled bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.Modes[mode] = enabled
}

func (ch *Channel) HasMode(mode ChannelMode) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	return ch.Modes[mode]
}

func (ch *Channel) SetOperator(client *Client, isOp bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if isOp {
		ch.operators[client] = true
	} else {
		delete(ch.operators, client)
	}
}

func (ch *Channel) IsOperator(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	return ch.operators[client]
}

func (ch *Channel) SetVoiced(client *Client, hasVoice bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if hasVoice {
		ch.voiced[client] = true
	} else {
		delete(ch.voiced, client)
	}
}

func (ch *Channel) IsVoiced(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	return ch.voiced[client]
}

func (ch *Channel) SetBanned(mask string, isBanned bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if isBanned {
		ch.banned[mask] = true
	} else {
		delete(ch.banned, mask)
	}
}

func (ch *Channel) IsBanned(client *Client) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	clientMask := utils.FormatMask(client.Nick, client.User, client.Host)

	for mask := range ch.banned {
		if utils.MatchesBanMask(clientMask, mask) {
			return true
		}
	}

	return false

}

func (ch *Channel) GetBanList() []string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	banList := make([]string, 0, len(ch.banned))
	for mask := range ch.banned {
		banList = append(banList, mask)
	}

	return banList
}

func (ch *Channel) SetInvited(nick string, isInvited bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if isInvited {
		ch.inviteList[strings.ToLower(nick)] = true
	} else {
		delete(ch.inviteList, strings.ToLower(nick))
	}
}

func (ch *Channel) IsInvited(nick string) bool {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	return ch.inviteList[strings.ToLower(nick)]
}

func (ch *Channel) Broadcast(message string) {
	clients := ch.GetClients()

	for _, client := range clients {
		client.Send(message)
	}
}

func (ch *Channel) BroadcastFrom(sender *Client, message string) {
	clients := ch.GetClients()

	for _, client := range clients {
		if client != sender {
			client.Send(message)
		}
	}
}
