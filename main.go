package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const inputPath = "messages.txt"

func main() {

	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("could not open %s: %s", inputPath, err)
	}
	defer file.Close()

	data := make([]byte, 8)

	for {
		count, err := file.Read(data)
		if err != nil {

			if errors.Is(err, io.EOF) {
				break
			}

			log.Fatalf("error: %s\n", err.Error())
			break
		}

		fmt.Printf("read: %s\n", data[:count])
	}

}
