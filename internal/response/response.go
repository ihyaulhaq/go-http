package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/ihyaulhaq/go-http/internal/headers"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type writerState int

const (
	WriteStatusLine writerState = iota
	WriterHeaders
	WriteBody
	WriteDone
)

type Writer struct {
	respose io.Writer
	state   writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{respose: w, state: WriteStatusLine}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriteStatusLine {
		return fmt.Errorf("Write status line called out of order")
	}

	_, err := fmt.Fprintf(w.respose, "HTTP/1.1 %d %s\r\n", statusCode, statusCode)
	if err != nil {
		return err
	}

	w.state = WriterHeaders

	return nil
}

func (s StatusCode) String() string {
	switch s {
	case StatusOk:
		return "OK"
	case StatusBadRequest:
		return "Bad Request"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.state != WriterHeaders {
		return fmt.Errorf("WriteHeaders called out of order")
	}

	for key, value := range h.Data {
		_, err := fmt.Fprintf(w.respose, "%s: %s\r\n", key, value)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w.respose, "\r\n")
	if err != nil {
		return err
	}
	w.state = WriteBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != WriteBody {
		return 0, fmt.Errorf("WriteBody called out of order")
	}

	n, err := w.respose.Write(p)
	if err != nil {
		return n, err
	}

	w.state = WriteDone
	return n, nil
}
