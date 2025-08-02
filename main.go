package main

import (
	"flag"
	"fmt"
	"goircd/config"
	"goircd/hash"
	"goircd/logger"
	"golang.org/x/term"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"goircd/server"

	_ "goircd/commands"
	_ "goircd/commands/op"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	genPasswordHash := flag.Bool("gen-password-hash", false, "Generate a password hash for the given password and exit")
	flag.Parse()

	if err := config.Load(*configPath); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if *genPasswordHash {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			logger.Fatal("Failed to read password: %v", err)
		}
		fmt.Println()

		fmt.Print("Confirm password: ")
		confirmPasswordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			logger.Fatal("Failed to read password: %v", err)
		}
		fmt.Println()

		password := strings.TrimSpace(string(passwordBytes))
		confirmPassword := strings.TrimSpace(string(confirmPasswordBytes))

		if password == "" {
			logger.Fatal("Password cannot be empty")
		}

		if password != confirmPassword {
			logger.Fatal("Passwords do not match")
		}

		logger.Info("Generating password hash...")

		hashedPassword, err := hash.Make(password)
		if err != nil {
			logger.Fatal("Failed to generate password hash: %v", err)
		}

		fmt.Printf("Password hash: %s\n", hashedPassword)
		os.Exit(0)
	}

	cfg := config.Get()

	logger.Info("Starting goircd server on %s:%d", cfg.Server.Host, cfg.Server.Port)

	ircServer, err := server.NewServer(cfg.Server.Host, cfg.Server.Port)
	if err != nil {
		logger.Fatal("Failed to create server: %v", err)
	}

	go func() {
		if err := ircServer.Start(); err != nil {
			logger.Fatal("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ircServer.Shutdown()
}
