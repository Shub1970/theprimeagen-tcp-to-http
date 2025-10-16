package server

import (
	"TCPtoHTTP/internal/response"
	"fmt"
	"io"
	"net"
)

type Server struct {
	closed bool
}

func runConnection(_s *Server, conn io.ReadWriteCloser) {
	headersData := response.GetDefaultHeaders(0)
	response.WriteStatusLine(conn, 200)
	response.WriteHeaders(conn, headersData)
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
