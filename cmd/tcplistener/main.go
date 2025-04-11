package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

// const inputPath = "messages.txt"
const port = ":42069"

func main() {

	// file, err := os.Open(inputPath)
	// if err != nil {
	// 	log.Fatalf("could not open %s: %s", inputPath, err)
	// }
	// defer file.Close()

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not open listening port: %s", err.Error())
	}
	defer l.Close()

	fmt.Println("Listening for TCP traffic on", port)

	for {

		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("could not accept connection: %s", err.Error())
		}

		fmt.Println("Reading data from connection")
		fmt.Println("=====================================")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error processing request: %v", err)
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Connection closed")

	}

}

// func getLinesChannel(f io.ReadCloser) <-chan string {

// 	ch := make(chan string)

// 	go func() {

// 		line := ""
// 		defer close(ch)

// 		for {

// 			data := make([]byte, 8)
// 			n, err := f.Read(data)

// 			if err != nil {

// 				if errors.Is(err, io.EOF) {
// 					break
// 				}

// 				log.Fatalf("error: %s\n", err.Error())

// 			}

// 			parts := strings.Split(string(data[:n]), "\n")

// 			for i, part := range parts {

// 				if i == len(parts)-1 {
// 					line += parts[len(parts)-1]
// 					break
// 				}

// 				line += part
// 				ch <- line
// 				line = ""
// 			}

// 		}

// 		if line != "" {
// 			ch <- line
// 		}

// 	}()

// 	return ch

// }
