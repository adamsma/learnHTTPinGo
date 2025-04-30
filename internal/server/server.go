package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Port     int
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	s := Server{Port: port, listener: l}
	s.closed.Store(false)

	go s.listen()

	return &s, nil

}

func (s *Server) Close() error {

	s.closed.Store(true)

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {

	for {

		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("could not accept connection: %s", err.Error())
			continue
		}

		go s.handle(conn)

	}
}

func (s *Server) handle(conn net.Conn) {

	defer conn.Close()

	msg := "HTTP/1.1 200 OK\r\n" + // Status line
		"Content-Type: text/plain\r\n" + // Example header
		"\r\n" + // Blank line to separate headers from the body
		"Hello World!\n" // Body
	conn.Write([]byte(msg))

}
