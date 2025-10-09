package main

import (
	"TCPtoHTTP/internal/request"
	"fmt"
	"io"
	"log"
	"net"
)

func handleConnection(conn io.ReadCloser) {
	defer conn.Close()
	request, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Errorf("error while handleConnection")
	}
	fmt.Print("Request line:\n")
	fmt.Printf("Method: %s\n", request.RequestLine.Method)
	fmt.Printf("Target: %s\n", request.RequestLine.RequestTarget)
	fmt.Printf("Version: %s\n", request.RequestLine.HttpVersion)

}

func main() {
	tcpConnect, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for {
		conn, err := tcpConnect.Accept()
		if err != nil {
			fmt.Printf("erro while reading")
			continue
		}

		fmt.Println("connection is setup")
		go handleConnection(conn)

	}

}
