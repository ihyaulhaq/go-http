package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/ihyaulhaq/go-http/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("can't listen to %s:  %w", addr, err)
	}

	s := &Server{
		listener: l,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()

}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}
		go s.handle(conn)

	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		if tc, ok := conn.(*net.TCPConn); ok {
			tc.CloseWrite()
		}
		conn.Close()
	}()
	body := []byte("Hello World!\n")
	h := response.GetDefaultHeaders(len(body))
	if err := response.WriteStatusLine(conn, response.StatusOk); err != nil {
		return
	}
	if err := response.WriteHeaders(conn, h); err != nil {
		return
	}
	// write body so the client receives the full HTTP response
	if _, err := conn.Write(body); err != nil {
		return
	}
}
