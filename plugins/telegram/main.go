package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "TELEGRAM_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := NewTelegramServer(token)
	if err := server.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start telegram server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Telegram plugin started, polling for updates...")

	// Handle messages
	go func() {
		for msg := range server.Messages() {
			fmt.Printf("Received message from %s: %s\n", msg.AccountID, msg.Content)
			// TODO: forward to gateway via gRPC/HTTP
		}
	}()

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Telegram plugin shutting down...")
	cancel()
}
