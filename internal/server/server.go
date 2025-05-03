package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Port     int
	listener net.Listener
	closed   atomic.Bool
	Handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	s := Server{Port: port, listener: l, Handler: handler}
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	sc := response.StatusCodeSuccess

	var b bytes.Buffer
	hErr := s.Handler(&b, req)
	if hErr != nil {
		hErr.Write(&b)
		sc = hErr.StatusCode
	}

	h := response.GetDefaultHeaders(b.Len())
	response.WriteStatusLine(conn, sc)
	response.WriteHeaders(conn, h)
	_, err = conn.Write(b.Bytes())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

}
