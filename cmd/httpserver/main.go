package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {

	h := func(w io.Writer, req *request.Request) *server.HandlerError {

		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandlerError{
				StatusCode: response.StatusCodeBadRequest,
				Message:    "Your problem is not my problem\n",
			}
		case "/myproblem":
			return &server.HandlerError{
				StatusCode: response.StatusCodeInternalServerError,
				Message:    "Woopsie, my bad\n",
			}
		default:
			w.Write([]byte("All good, frfr\n"))
		}

		return nil

	}

	server, err := server.Serve(port, h)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
