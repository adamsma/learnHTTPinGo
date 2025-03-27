package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	serverAddr := "localhost:42069"

	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("Unable to resolve UDP address: %s", err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Unable to initialize UDP connection: %s", err.Error())
	}
	defer conn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", serverAddr)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Unable to read input")
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatalf("Error sending message: %v\n", err)
		}

		fmt.Printf("Message sent: %s", line)
	}

}
