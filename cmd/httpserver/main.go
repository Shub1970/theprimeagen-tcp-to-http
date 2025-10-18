package main

import (
	"TCPtoHTTP/internal/request"
	"TCPtoHTTP/internal/response"
	"TCPtoHTTP/internal/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func requestStatusResponse(status response.StatusCode) ([]byte, int) {
	var res []byte
	switch status {
	case response.StatusOk:
		res = []byte(`<html>
		<head>
		<title>200 OK</title>
		</head>
		<body>
		<h1>Success!</h1>
		<p>Your request was an absolute banger.</p>
		</body>
		</html>`)
	case response.StatusBadRequest:
		res = []byte(`<html>
		<head>
		<title>400 Bad Request</title>
		</head>
		<body>
		<h1>Bad Request</h1>
		<p>Your request honestly kinda sucked.</p>
		</body>
		</html>`)
	case response.StatusServerError:
		res = []byte(`<html>
		<head>
		<title>500 Internal Server Error</title>
		</head>
		<body>
		<h1>Internal Server Error</h1>
		<p>Okay, you know what? This one is on me.</p>
		</body>
		</html>`)
	}

	return res, len(res)
}

func checker(w *response.Writer, req *request.Request) {
	h := response.GetDefaultHeaders(0)
	status := response.StatusOk
	body, length := requestStatusResponse(status)
	if req.RequestLine.RequestTarget == "/yourproblem" {
		status = response.StatusBadRequest
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		status = response.StatusServerError
	}
	body, length = requestStatusResponse(status)
	h.Replace("Content-Type", "text/html; charset=utf-8")
	h.Replace("Content-Length", fmt.Sprintf("%d", length))
	w.WriteStatusLine(status)
	w.WriteHeaders(*h)
	w.WriteBody(body)
}

func main() {
	server, err := server.Serve(port, checker)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
