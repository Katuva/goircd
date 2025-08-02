package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Check if the server address is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client/test_client.go <server:port>")
		os.Exit(1)
	}

	// Get the server address
	serverAddr := os.Args[1]

	// Connect to the server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to", serverAddr)
	fmt.Println("Enter IRC commands (e.g., NICK test, USER test 0 * :Test User)")
	fmt.Println("Type 'quit' to exit")

	// Create a channel to signal when to exit
	done := make(chan struct{})

	// Start a goroutine to read from the server
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading from server: %v\n", err)
				close(done)
				return
			}
			fmt.Print("< ", message)
		}
	}()

	// Read from stdin and send to the server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.ToLower(line) == "quit" {
			break
		}

		// Send the command to the server
		fmt.Fprintf(conn, "%s\r\n", line)
	}

	// Signal that we're done
	close(done)
}
