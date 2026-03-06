# go-http

`go-http` is a custom, from-scratch HTTP server implementation written in Go. Instead of relying on Go's standard `net/http` package for its core server logic, this project builds an HTTP server directly on top of raw TCP connections using the `net` package. It serves as an educational and foundational project to understand the inner workings of the HTTP protocol.

## Features

- **Custom TCP Listener:** Manages raw TCP connections directly via the `net` package (`internal/server`).
- **Protocol Parsing:** Custom parsing for HTTP requests (`internal/request`) and constructing HTTP responses (`internal/response`).
- **Header Management:** A custom implementation for managing HTTP headers (`internal/headers`).
- **Routing System:** A minimalistic router supporting basic HTTP methods like `GET`, `POST`, and `DELETE`.
- **Chunked Transfer Encoding:** Includes support for reading and writing chunked payloads.
- **Graceful Shutdown:** Handles system signals (`SIGINT`, `SIGTERM`) to cleanly close connections and shut down the server.
- **Proxy/Streaming Capabilities:** Demonstrates streaming data (e.g., video files) and proxying requests to external services like `httpbin`.

## Project Structure

```
.
├── cmd/
│   └── httpserver/        # The entry point of the application
│       └── main.go        # Contains the setup and custom handlers
├── internal/
│   ├── headers/           # Logic for HTTP header manipulation
│   ├── request/           # Parsing of raw incoming HTTP requests
│   ├── response/          # Construction of HTTP responses
│   └── server/            # Core TCP server, router, and handler definitions
└── go.mod                 # Go module file
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/) (Version 1.20+ recommended)

### Running the Server

1. Clone or navigate to the project directory.
2. Run the main server application:

```bash
go run cmd/httpserver/main.go
```

The server will start on port `42069` by default.

## Handlers & Endpoints

In the provided example (`cmd/httpserver/main.go`), the following endpoints are available:

- `GET /`: Returns a simple `200 OK` HTML success page.
- `GET /video`: Streams a local `.mp4` video file (`assets/vim.mp4`) to the client.
- `GET /yourproblem`: Simulates a `400 Bad Request` with custom HTML.
- `GET /myproblem`: Simulates a `500 Internal Server Error` with custom HTML.
- `GET /httpbin/*`: Acts as a proxy to `httpbin.org`. It fetches the corresponding path from HTTPBin and streams the response back to the client using **Chunked Transfer Encoding**, alongside appending `Trailing-Headers` (`X-Content-SHA256` and `X-Content-Length`).

## How it Works

1. **Accepting Connections:** The `Server` listens on a specified port, accepting incoming TCP connections.
2. **Handling Requests:** Each connection is handled in a separate goroutine. 
3. **Parsing Data:** Incoming bytes are read and parsed by `request.RequestFromReader()` according to the HTTP/1.1 specification (Request-Line, Headers, Body).
4. **Routing:** The `Router` maps the request method and path to a specific `Handler`.
5. **Formulating Responses:** The handler constructs a response mapping using the `response.Writer`, sending the Status-Line, Headers, and Response Body back over the TCP connection.
6. **Error Handling:** If any step fails, or the handler encounters an issue, a custom `HandlerError` is returned and formulated as a clean error response to the client.


