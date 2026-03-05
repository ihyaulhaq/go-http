package request

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ihyaulhaq/go-http/internal/headers"
)

type parseState int

const (
	stateInitialized parseState = iota
	requestStateDone
	requestStateParsingHeaders
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const bufferSize = 8
const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	req := &Request{
		state:   stateInitialized,
		Headers: headers.Headers{},
	}

	for req.state != requestStateDone {

		// sizing buffer
		if len(buf) == readToIndex {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		nRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		readToIndex += nRead

		nParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if nParsed > 0 {
			copy(buf, buf[nParsed:readToIndex])
			readToIndex -= nParsed
		}

	}

	if req.state != requestStateDone {
		return nil, fmt.Errorf("incomplete request")
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalbytes := 0
	for r.state != requestStateDone {
		n, err := r.ParseSingle(data[totalbytes:])
		if err != nil {
			return 0, err
		}

		if n == 0 {
			break
		}
		totalbytes += n
	}
	return totalbytes, nil
}

func (r *Request) ParseSingle(data []byte) (int, error) {
	switch r.state {
	case stateInitialized:
		reqLine, consumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if consumed == 0 {
			return 0, nil
		}

		r.state = requestStateParsingHeaders
		r.RequestLine = reqLine

		return consumed, nil
	case requestStateParsingHeaders:
		consumed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {

			r.state = requestStateDone
		}
		return consumed, nil

	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}

}

func parseRequestLine(raw []byte) (RequestLine, int, error) {

	str := string(raw)

	idx := strings.Index(str, crlf)
	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	line := str[:idx]
	consumed := idx + 2

	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("invalid number of parts in request line")
	}

	method := parts[0]
	target := parts[1]
	versionLiteral := parts[2]

	for _, r := range method {
		if r < 'A' || r > 'Z' {
			return RequestLine{}, 0, fmt.Errorf("invalid method: %s", method)
		}

	}

	if !strings.HasPrefix(versionLiteral, "HTTP/") {
		return RequestLine{}, 0, fmt.Errorf("invalid http version")
	}

	version := strings.TrimPrefix(versionLiteral, "HTTP/")
	if version != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("unsupported http version: %s", versionLiteral)
	}
	return RequestLine{
			HttpVersion:   "1.1",
			RequestTarget: target,
			Method:        method,
		},
		consumed, nil

}
