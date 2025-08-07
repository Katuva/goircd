package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
	Security SecurityConfig `yaml:"security"`
	Channels ChannelsConfig `yaml:"channels"`
	Users    UsersConfig    `yaml:"users"`
	MOTD     string         `yaml:"-"`
}

type ServerConfig struct {
	Name        string `yaml:"name"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Description string `yaml:"description"`
	MOTD        string `yaml:"motd"`
	MaxClients  int    `yaml:"max_clients"`
	Timeout     int    `yaml:"timeout"`
}

type LoggingConfig struct {
	Level    string `yaml:"level"`
	File     string `yaml:"file"`
	MaxSize  int    `yaml:"max_size"`
	MaxAge   int    `yaml:"max_age"`
	Compress bool   `yaml:"compress"`
	Console  bool   `yaml:"console"`
}

type SecurityConfig struct {
	RequireAuth    bool       `yaml:"require_auth"`
	MaxNickLength  int        `yaml:"max_nick_length"`
	MaxChannelName int        `yaml:"max_channel_name"`
	MaskHosts      bool       `yaml:"mask_hosts"`
	Secret         string     `yaml:"secret"`
	BannedNicks    []string   `yaml:"banned_nicks"`
	AllowedHosts   []string   `yaml:"allowed_hosts"`
	Operators      []Operator `yaml:"operators"`
}

type Operator struct {
	Nick     string `yaml:"nick"`
	Whois    string `yaml:"whois"`
	Vhost    string `yaml:"vhost"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type ChannelsConfig struct {
	DefaultModes       string   `yaml:"default_modes"`
	MaxChannels        int      `yaml:"max_channels"`
	MaxUsersPerChannel int      `yaml:"max_users_per_channel"`
	DefaultChannels    []string `yaml:"default_channels"`
	AutoJoin           bool     `yaml:"auto_join"`
}

type UsersConfig struct {
	MaxIdleTime      int `yaml:"max_idle_time"`
	PingInterval     int `yaml:"ping_interval"`
	MaxMessageLength int `yaml:"max_message_length"`
}

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

func Load(configPath string) error {
	var err error

	once.Do(func() {
		instance, err = loadConfig(configPath)
	})

	return err
}

func Get() *Config {
	mu.RLock()
	defer mu.RUnlock()
	if instance == nil {
		panic("config not loaded: call config.Load() first")
	}
	return instance
}

func Reload(configPath string) error {
	mu.Lock()
	defer mu.Unlock()

	newConfig, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	instance = newConfig
	return nil
}

func loadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	setDefaults(&config)

	if config.Server.MOTD != "" {
		data, err := os.ReadFile(config.Server.MOTD)

		if err != nil {
			return nil, fmt.Errorf("failed to read MOTD file: %w", err)
		}

		config.MOTD = string(data)
	}

	return &config, nil
}

func setDefaults(config *Config) {
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 6667
	}
	if config.Server.Name == "" {
		config.Server.Name = "GoIRCd"
	}
	if config.Server.MaxClients == 0 {
		config.Server.MaxClients = 100
	}
	if config.Server.Timeout == 0 {
		config.Server.Timeout = 300 // 5 minutes
	}

	// Logging defaults
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.MaxSize == 0 {
		config.Logging.MaxSize = 10 // 10MB
	}
	if config.Logging.MaxAge == 0 {
		config.Logging.MaxAge = 30 // 30 days
	}
	if config.Logging.Console == false {
		config.Logging.Console = true
	}

	// Security defaults
	if config.Security.MaxNickLength == 0 {
		config.Security.MaxNickLength = 16
	}
	if config.Security.MaxChannelName == 0 {
		config.Security.MaxChannelName = 50
	}

	// Channels defaults
	if config.Channels.MaxChannels == 0 {
		config.Channels.MaxChannels = 20
	}
	if config.Channels.MaxUsersPerChannel == 0 {
		config.Channels.MaxUsersPerChannel = 100
	}

	// Users defaults
	if config.Users.MaxIdleTime == 0 {
		config.Users.MaxIdleTime = 600 // 10 minutes
	}
	if config.Users.PingInterval == 0 {
		config.Users.PingInterval = 60 // 1 minute
	}
	if config.Users.MaxMessageLength == 0 {
		config.Users.MaxMessageLength = 512
	}
}
