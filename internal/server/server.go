package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Port     int
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("could not open listening port: %s", err.Error())
	}

	s := Server{Port: port, listener: l}
	s.closed.Store(false)

	go s.listen()

	return &s, nil

}

func (s *Server) Close() error {

	s.listener.Close()
	s.closed.Store(true)

	return nil
}

func (s *Server) listen() {

	for !s.closed.Load() {

		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection: %s", err.Error())
		}

		go s.handle(conn)

	}
}

func (s *Server) handle(conn net.Conn) {

	msg := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!"
	conn.Write([]byte(msg))
	conn.Close()

}
