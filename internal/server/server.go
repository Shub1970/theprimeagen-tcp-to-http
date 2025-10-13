package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	closed bool
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	body := []byte("Hello World!\n")
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+"Content-Type: text/plain\r\n"+"\r\n%s",
		len(body), body,
	)

	conn.Write([]byte(response))
	conn.Close()
}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}

		go runConnection(s, conn)
	}
}

func Serve(port uint16) (*Server, error) {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{closed: false}

	go runServer(server, listner)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
