package server

import (
	"fmt"
	"log"
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

type Handler func(w *response.Writer, req *request.Request) *HandlerError

type Router struct {
	routes map[string]Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
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

			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	for {

		req, err := request.RequestFromReader(conn)
		if err != nil {
			hErr := &HandlerError{
				StatusCode: response.StatusBadRequest,
				Message:    err.Error(),
			}

			hErr.WriteHandlerError(conn)
			return
		}

		w := response.NewWriter(conn)
		if handlerErr := s.handler(w, req); handlerErr != nil {
			handlerErr.WriteHandlerError(conn)
			return
		}

		if req.Headers.Get("Connection") == "close" {
			break
		}
	}
}

func (h *HandlerError) WriteHandlerError(conn net.Conn) {
	body := []byte(h.Message)
	headers := response.GetDefaultHeaders(len(body))
	w := response.NewWriter(conn)
	w.WriteStatusLine(h.StatusCode)
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Handler),
	}
}

func (r *Router) Handle(method, path string, handler Handler) {
	key := method + " " + path
	r.routes[key] = handler
}

func (r *Router) GET(path string, h Handler)    { r.Handle("GET", path, h) }
func (r *Router) POST(path string, h Handler)   { r.Handle("POST", path, h) }
func (r *Router) DELETE(path string, h Handler) { r.Handle("DELETE", path, h) }

func (r *Router) ServeHTTP(w *response.Writer, req *request.Request) *HandlerError {
	key := req.RequestLine.Method + " " + req.RequestLine.RequestTarget

	handler, ok := r.routes[key]
	if !ok {
		return &HandlerError{
			StatusCode: response.StatusNotFound,
			Message:    "404 Not Found",
		}
	}
	return handler(w, req)
}
