package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

const badRequestHTML = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`

const internalErrorHTML = `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

const successHTML = `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

func main() {

	h := func(w *response.Writer, req *request.Request) {

		switch t := req.RequestLine.RequestTarget; {
		case t == "/yourproblem":

			w.WriteStatusLine(response.StatusCodeBadRequest)

			h := response.GetDefaultHeaders(len(badRequestHTML))
			h.Override("content-type", "text/html")
			w.WriteHeaders(h)

			w.WriteBody([]byte(badRequestHTML))

		case t == "/myproblem":

			w.WriteStatusLine(response.StatusCodeInternalServerError)

			h := response.GetDefaultHeaders(len(internalErrorHTML))
			h.Override("content-type", "text/html")
			w.WriteHeaders(h)

			w.WriteHeaders(h)
			w.WriteBody([]byte(internalErrorHTML))

		case strings.HasPrefix(t, "/httpbin"):

			proxyHandler(w, req)

		case t == "/video":

			w.WriteStatusLine(response.StatusCodeSuccess)

			vid, err := os.ReadFile("assests/vim.mp4")
			if err != nil {
				fmt.Println("error loading video:", err)
			}

			h := response.GetDefaultHeaders(len(vid))
			h.Override("content-type", "video/mp4")

			w.WriteHeaders(h)

			w.WriteBody(vid)

		default:

			w.WriteStatusLine(response.StatusCodeSuccess)

			h := response.GetDefaultHeaders(len(successHTML))
			h.Set("content-type", "text/html")
			w.WriteHeaders(h)

			w.WriteBody([]byte(successHTML))

		}

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

func proxyHandler(w *response.Writer, req *request.Request) {

	w.WriteStatusLine(response.StatusCodeSuccess)

	h := response.GetDefaultHeaders(0)
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-SHA256")
	h.Set("Trailer", "X-Content-Length")
	h.Override("content-type", "text/html")
	h.Delete("Content-Length")
	w.WriteHeaders(h)

	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	const maxChunkSize = 1024
	buf := make([]byte, maxChunkSize)
	msg := make([]byte, 0)

	for {

		n, err := res.Body.Read(buf)

		if n > 0 {
			_, err = w.WriteChunkedBody(buf[:n])
			msg = append(msg, buf[:n]...)
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}

	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("error writing chunked body done:", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(msg))

	log.Printf("Hash: %x, Length: %d", hash, len(msg))

	trailers := headers.NewHeaders()
	trailers.Override("X-Content-SHA256", hash)
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(msg)))

	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Println("error writing trailers:", err)
	}

}
