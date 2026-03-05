package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ihyaulhaq/go-http/internal/request"
	"github.com/ihyaulhaq/go-http/internal/response"
	"github.com/ihyaulhaq/go-http/internal/server"
)

const PORT = 3000

func main() {
	server, err := server.Serve(PORT, handler)
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

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	default:
		fmt.Fprint(w, "All good frfr\n")
		return nil
	}

}
