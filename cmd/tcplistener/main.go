package main

import (
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
			fmt.Errorf("not accepted connection")
		}
		channel := getLinesChannel(conn)
		for line := range channel {
			fmt.Printf("read: %s\n", string(line))
		}
	}

}
