package main

import (
	"boot.shubeagen.tv/internal/request"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string, 1)

	go func() {

		defer close(ch)
		defer f.Close()
		str := ""

		for {
			buff := make([]byte, 8)

			_, err := f.Read(buff)
			if err != nil {
				break
			}

			if i := bytes.IndexByte(buff, '\n'); i != -1 {
				str += string(buff[:i])
				buff = buff[i+1:]
				ch <- str
				str = ""
			}
			str += string(buff)
		}

		if len(str) != 0 {
			ch <- str
		}

	}()

	return ch
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
		}

		rl, err := request.RequestFromReader(conn)

		if err != nil {
			fmt.Printf("erro while running RequestFromReader")
		}
		fmt.Printf("Request line: \n")
		fmt.Printf("- Method: %s \n", rl.RequestLine.Method)
		fmt.Printf("- Target: %s \n", rl.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s \n", rl.RequestLine.HttpVersion)
	}

}
