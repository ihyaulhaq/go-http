package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)

	if err != nil {
		log.Fatalf("cant open file %s : %s\n", inputFilePath, err)
	}
	str := getLinesChannel(file)

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	for line := range str {

		fmt.Printf("read: %s\n", line)
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
					curentLine = ""
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
