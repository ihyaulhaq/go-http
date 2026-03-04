package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const address = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		log.Fatalf("resolve error: %v", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("dial error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", address)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("read error: %v", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatalf("write error: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Message sent: %s", input)

	}
}
