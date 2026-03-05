package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/ihyaulhaq/go-http/internal/request"
	"github.com/ihyaulhaq/go-http/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("can't listen to %s:  %w", addr, err)
	}

	s := &Server{
		listener: l,
		handler:  handler,
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		WriteHandlerError(conn, &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		})

	}

	var buf bytes.Buffer

	if handlerErr := s.handler(&buf, req); handlerErr != nil {
		WriteHandlerError(conn, handlerErr)
		return
	}

	body := buf.Bytes()
	h := response.GetDefaultHeaders(len(body))
	if err := response.WriteStatusLine(conn, response.StatusOk); err != nil {
		return
	}
	if err := response.WriteHeaders(conn, h); err != nil {
		return
	}
	conn.Write(body)
}

func WriteHandlerError(w io.Writer, h *HandlerError) {
	body := []byte(h.Message)
	headers := response.GetDefaultHeaders(len(body))
	response.WriteStatusLine(w, h.StatusCode)
	response.WriteHeaders(w, headers)
	w.Write(body)
}
