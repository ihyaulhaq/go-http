package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/ihyaulhaq/go-http/internal/headers"
	"github.com/ihyaulhaq/go-http/internal/request"
	"github.com/ihyaulhaq/go-http/internal/response"
	"github.com/ihyaulhaq/go-http/internal/server"
)

const PORT = 42069

func main() {
	router := server.NewRouter()

	router.GET("/video", youtubeHandler)
	router.GET("/", defaultHandler)

	server, err := server.Serve(PORT, router.ServeHTTP)
	if err != nil {
		log.Fatalf("Error startting server : %v", err)
	}
	defer server.Close()
	log.Println("Server Started on port", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Serveer gracefully stopped")
}

func handler(res *response.Writer, req *request.Request) *server.HandlerError {

	url := req.RequestLine.RequestTarget

	if strings.HasPrefix(url, "/httpbin") {
		return proxyhandler(res, req)
	}

	switch url {
	case "/yourproblem":
		body := []byte(`<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
	</html>`)
		h := response.GetDefaultHeaders(len(body))
		h.Set("Content-Type", "text/html")
		res.WriteStatusLine(response.StatusOk)
		res.WriteHeaders(h)
		res.WriteBody(body)
		return nil

	case "/myproblem":
		body := []byte(`<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
	</html>`)
		h := response.GetDefaultHeaders(len(body))
		h.Set("Content-Type", "text/html")
		res.WriteStatusLine(response.StatusOk)
		res.WriteHeaders(h)
		res.WriteBody(body)
		return nil

	default:
		body := []byte(`<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>Success!</h1>
			<p>Your request was an absolute banger.</p>
		</body>
	</html>`)
		h := response.GetDefaultHeaders(len(body))
		h.Set("Content-Type", "text/html")
		res.WriteStatusLine(response.StatusOk)
		res.WriteHeaders(h)
		res.WriteBody(body)
		return nil

	}

}

func youtubeHandler(res *response.Writer, req *request.Request) *server.HandlerError {

	vid, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    "Can't open the file:" + err.Error(),
		}
	}

	h := response.GetDefaultHeaders(len(vid))
	h.Set("Content-Type", "video/mp4")
	res.WriteStatusLine(response.StatusOk)
	res.WriteHeaders(h)
	res.WriteBody(vid)
	return nil

}

func proxyhandler(res *response.Writer, req *request.Request) *server.HandlerError {

	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	url := "https://httpbin.org" + path

	resp, err := http.Get(url)
	if err != nil {
		return &server.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Message:    "cant get to httpbin" + err.Error(),
		}
	}
	defer resp.Body.Close()

	// build response
	h := response.GetDefaultHeaders(0)
	h.Delete("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Content-Type", resp.Header.Get("Content-Type"))
	h.Set("Trailer", "X-Content-SHA256, X-Content-Length")

	res.WriteStatusLine(response.StatusOk)
	res.WriteHeaders(h)

	buf := make([]byte, 1024)
	dataLen := 0
	hasaer := sha256.New()
	for {

		n, err := resp.Body.Read(buf)
		if n > 0 {
			log.Printf("Read %d bytes from httpbin", n)
			hasaer.Write(buf[:n])
			dataLen += n
			_, wErr := res.WriteChunkedBody(buf[:n])
			if wErr != nil {
				log.Printf("Error writing chunk: %v", wErr)
				break
			}
		}

		// EOF
		if err != nil {
			break
		}
	}

	t := headers.NewHeaders()
	t.Set("X-Content-SHA256", fmt.Sprintf("%x", hasaer.Sum(nil)))
	t.Set("X-Content-Length", strconv.Itoa(dataLen))

	res.WriteChunkedBodyDone()
	res.WriteTrailers(t)
	return nil
}

func defaultHandler(w *response.Writer, r *request.Request) *server.HandlerError {
	body := []byte(`<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>Success!</h1>
			<p>Your request was an absolute banger.</p>
		</body>
	</html>`)
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(h)
	w.WriteBody(body)
	return nil
}
