package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "DISCORD_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := NewDiscordServer(token)
	if err := server.Connect(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to Discord: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Discord plugin started, listening for messages...")

	// Handle messages
	go func() {
		for msg := range server.Messages() {
			fmt.Printf("Received message from %s in channel %s: %s\n",
				msg.AccountID, msg.ChannelID, msg.Content)
			// TODO: forward to gateway via gRPC/HTTP
		}
	}()

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Discord plugin shutting down...")
	cancel()
}
