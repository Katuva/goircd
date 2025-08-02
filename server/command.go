package server

type Command interface {
	Name() string

	Execute(client *Client, params string)

	Help() string
}

type CommandInitFunc func(server *Server) Command

var CommandRegistry = make(map[string]CommandInitFunc)

func RegisterCommandInit(name string, initFunc CommandInitFunc) {
	CommandRegistry[name] = initFunc
}
