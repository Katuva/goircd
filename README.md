# GoIRCd - Go IRC Server

GoIRCd is a simple IRC server implemented in Go, following the RFC 1459 specification. It supports the basic IRC functionality including channels, private messaging, and operator commands.

## Features

- Follows RFC 1459 specification
- Concurrent client handling using goroutines
- Modular command system with dynamic loading
- Support for channels and private messaging
- Operator commands for server administration
- Simple and clean codebase

## Supported Commands

### Basic Commands
- `NICK` - Set or change your nickname
- `USER` - Set your username and real name
- `JOIN` - Join a channel
- `PART` - Leave a channel
- `PRIVMSG` - Send a message to a user or channel
- `MODE` - Change channel or user modes
- `TOPIC` - View or change a channel's topic
- `PING`/`PONG` - Server ping/pong for connection maintenance

### Channel Operator Commands
- `OP` - Give channel operator status to a user
- `DEOP` - Remove channel operator status from a user
- `KICK` - Kick a user from a channel
- `BAN` - Ban a user from a channel

### IRC Operator Commands
- `OPER` - Authenticate as an IRC operator
- `KILL` - Disconnect a user from the server
- `RESTART` - Restart the server
- `OPERGIVE` - Give channel operator status to a user (without being a channel operator)
- `OPERREMOVE` - Remove channel operator status from a user (without being a channel operator)
- `OPERBAN` - Ban a user from a channel (without being a channel operator)

## Building and Running

### Prerequisites
- Go 1.24 or later

### Building
```bash
go build -o goircd
```

### Running
```bash
./goircd -port 6667 -host 0.0.0.0
```

Command-line options:
- `-port` - Port to listen on (default: 6667)
- `-host` - Host to bind to (default: 0.0.0.0)

## Testing

A simple test client is included to test the server functionality:

```bash
go run client/test_client.go localhost:6667
```

## Operator Authentication

To authenticate as an operator, use the following command:
```
OPER admin password
```

## Logging System

The server includes a privacy-focused logging system that logs important actions without including message content. The logging system:

- Logs command execution, operator actions, bans, kicks, and connection events
- Ensures privacy by not logging message content
- Provides different log levels (Debug, Info, Warning, Error, Fatal)
- Uses a centralized logger package for consistent logging throughout the codebase

## Project Structure

- `main.go` - Entry point for the server
- `server/` - Core server implementation
  - `server.go` - Server struct and methods
  - `client.go` - Client struct and methods
  - `channel.go` - Channel struct and methods
  - `command.go` - Command interface and registration
- `commands/` - IRC command implementations
  - `nick.go`, `user.go`, etc. - Individual command implementations
  - `op/` - Operator command implementations
- `logger/` - Logging system
  - `logger.go` - Logger implementation with privacy features
- `client/` - Test client for testing the server
  - `test_client.go` - Simple IRC client implementation
- `utils/` - Utility functions and constants
  - `parser.go` - IRC message parsing utilities
  - `constants.go` - IRC numeric replies and constants

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.