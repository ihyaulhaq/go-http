package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ihyaulhaq/go-http/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("cant lisen to port %s : %s\n", port, err)
	}
	defer listener.Close()

	fmt.Println("=====================================")
	fmt.Println("Listening for TCP traffic on port", port)
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error:%s", err.Error())
		}
		fmt.Println("Request Line: ")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		fmt.Println("Headers: ")
		for key, value := range req.Headers.All() {
			fmt.Printf("- %s: %s \n", key, value)
		}

		fmt.Println("Body: ")
		fmt.Printf("%s\n", string(req.Body))
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
		fmt.Println("=====================================")
	}
}
