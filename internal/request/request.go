package request

import (
	"bytes"
	"fmt"
	"io"
	"unicode"
)

type parseState string

const (
	StateInit  parseState = "initalized"
	StateDone  parseState = "done"
	StateError parseState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parseState
}

var BUFFERSIZE = 1024
var ERROR_BAD_START_LINE = fmt.Errorf("bad start line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request error")
var SPLITING_BYTES = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func capitalCheck(str string) bool {
	isUpper := true
	for _, ch := range str {
		if !unicode.IsUpper(ch) {
			isUpper = false
		}
	}

	return isUpper
}

func parseRequestLine(requestBytes []byte) (*RequestLine, int, error) {
	splitingIndex := bytes.Index(requestBytes, SPLITING_BYTES)

	if splitingIndex == -1 {
		return nil, 0, nil
	}
	remaning_message := requestBytes[splitingIndex+len(SPLITING_BYTES):]

	requestLine := requestBytes[:splitingIndex]

	requestLine_parts := bytes.Split(requestLine, []byte(" "))
	if len(requestLine_parts) != 3 {
		return nil, 0, ERROR_BAD_START_LINE
	}
	method_part := string(requestLine_parts[0])

	http_version_part := bytes.Split(requestLine_parts[2], []byte("/"))

	if !capitalCheck(method_part) || string(http_version_part[0]) != "HTTP" || string(http_version_part[1]) != "1.1" {
		return nil, 0, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	return &RequestLine{
		Method:        method_part,
		RequestTarget: string(requestLine_parts[1]),
		HttpVersion:   string(http_version_part[1]),
	}, len(remaning_message), nil

}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(data)
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			r.state = StateDone
			read += n

		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	newReq := newRequest()
	buf := make([]byte, BUFFERSIZE)
	bufIndex := 0

	for !newReq.done() {
		readLeng, err := reader.Read(buf[bufIndex:])
		if err != nil {
			return nil, err
		}
		bufIndex += readLeng
		consumeLeng, err := newReq.parse(buf[:bufIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[consumeLeng:bufIndex])
		bufIndex -= consumeLeng
	}

	return newReq, nil
}
