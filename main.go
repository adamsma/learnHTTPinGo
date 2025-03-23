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

	var line string

	for {

		data := make([]byte, 8)
		_, err := file.Read(data)
		parts := strings.Split(string(data), "\n")

		if err != nil {

			if errors.Is(err, io.EOF) {
				break
			}

			log.Fatalf("error: %s\n", err.Error())
			break
		}

		for i, part := range parts {

			if i == len(parts)-1 {
				line += parts[len(parts)-1]
				break
			}

			line += part
			fmt.Printf("read: %s\n", line)
			line = ""
		}

	}

	if line != "" {
		fmt.Printf("read: %s\n", line)
	}

}
