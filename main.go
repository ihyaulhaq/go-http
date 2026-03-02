package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	// "os"
	"strings"
)

// const inputFilePath = "messages.txt"
const port = ":42069"

func main() {
	// file, err := os.Open(inputFilePath)

	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("cant lisen to port %s : %s\n", port, err)
	}
	defer listener.Close()
	defer fmt.Println("Connention has been close")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connention has been made ")
		str := getLinesChannel(conn)

		fmt.Println("=====================================")

		for line := range str {
			fmt.Println(line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		curentLine := ""

		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)

			if err != nil {
				if curentLine != "" {
					lines <- curentLine
				}

				if errors.Is(err, io.EOF) {
					break
				}

				fmt.Printf("error: %s\n", err.Error())
				return
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")

			for l := 0; l < len(parts)-1; l++ {

				line := curentLine + parts[l]
				lines <- line
				curentLine = ""
			}
			curentLine += parts[len(parts)-1]
		}

	}()

	return lines
}
