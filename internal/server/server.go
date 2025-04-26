package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Port     int
	listener *net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {

	var l *net.Listener
	go func(l *net.Listener) {

		lstr, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		l = &lstr
		if err != nil {
			log.Fatalf("could not open listening port: %s", err.Error())
		}

	}(l)

	s := Server{Port: port, listener: l}
	s.closed.Store(false)

	return &s, nil

}

func (s *Server) Close() error {
	return nil
}

func (s *Server) listen() {

}

func (s *Server) handle(conn net.Conn) {

}
