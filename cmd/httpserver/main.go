package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ihyaulhaq/go-http/internal/server"
)

const PORT = 3000

func main() {
	server, err := server.Serve(PORT)
	if err != nil {
		log.Fatalf("Error startting server : %v", err)
	}
	defer server.Close()
	log.Println("Server Started on port", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Serveer gracefully stopped")
}
