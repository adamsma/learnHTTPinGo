package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputPath = "messages.txt"

func main() {

	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("could not open %s: %s", inputPath, err)
	}
	defer file.Close()

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Println("read:", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string)

	go func() {

		line := ""
		defer close(ch)

		for {

			data := make([]byte, 8)
			n, err := f.Read(data)

			if err != nil {

				if errors.Is(err, io.EOF) {
					break
				}

				log.Fatalf("error: %s\n", err.Error())

			}

			parts := strings.Split(string(data[:n]), "\n")

			for i, part := range parts {

				if i == len(parts)-1 {
					line += parts[len(parts)-1]
					break
				}

				line += part
				ch <- line
				line = ""
			}

		}

		if line != "" {
			ch <- line
		}

	}()

	return ch

}
