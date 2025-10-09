package request

import (
	"bytes"
	"fmt"
	"io"
)

type parsingStatus string

const (
	initialState parsingStatus = "initial"
	doneState    parsingStatus = "done"
	errorState   parsingStatus = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Status      parsingStatus
}

var NOT_FOUND_NEW_LINE_INDECTAR = fmt.Errorf("request line is not complete")

var NEW_LINE_INDECTAR = "\r\n"

func isCapital(str string) bool {
	for _, char := range str {
		if char < 'A' || char > 'Z' {
			return false
		}
	}
	return true
}

func parseRequestLine(requestLine []byte) (*RequestLine, int, error) {
	fmt.Printf("running parseRequestLine: %s\n", string(requestLine))
	endIndex := bytes.Index(requestLine, []byte(NEW_LINE_INDECTAR))
	if endIndex == -1 {
		return nil, 0, NOT_FOUND_NEW_LINE_INDECTAR
	}

	reqLineBytes := requestLine[:endIndex]
	readDataLength := endIndex + len(NEW_LINE_INDECTAR)
	split_reqLineBytes := bytes.SplitN(reqLineBytes, []byte(" "), 3)
	if len(split_reqLineBytes) != 3 {
		return nil, 0, NOT_FOUND_NEW_LINE_INDECTAR
	}
	method, request_target, HTTP_txt := split_reqLineBytes[0], split_reqLineBytes[1], split_reqLineBytes[2]
	if !isCapital(string(method)) {
		return nil, 0, fmt.Errorf("wrong method")
	}
	HTTP_part := bytes.SplitN(HTTP_txt, []byte("/"), 2)
	if len(HTTP_part) != 2 || string(HTTP_part[0]) != "HTTP" || string(HTTP_part[1]) != "1.1" {
		return nil, 0, fmt.Errorf("http version is wrong")
	}

	reqLine := &RequestLine{
		HttpVersion:   string(HTTP_part[1]),
		RequestTarget: string(request_target),
		Method:        string(method),
	}

	return reqLine, readDataLength, nil
}

func (r *Request) doneCheck() bool {
	return r.Status == doneState || r.Status == errorState
}

// state machine
func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.Status {
		case initialState:
			requestLine, parseLen, err := parseRequestLine(data)
			if err != nil {
				if err == NOT_FOUND_NEW_LINE_INDECTAR {
					break outer
				}
				return 0, err
			}
			if parseLen == 0 {
				break outer
			}
			read += parseLen
			r.RequestLine = *requestLine
			r.Status = doneState

		case doneState:
			break outer
		case errorState:
			return 0, fmt.Errorf("reach error state")
		default:
			return read, fmt.Errorf("unknown parser state")
		}
	}

	return read, nil
}

func newRequest() *Request {
	return &Request{
		Status: initialState,
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// craete new request
	req := newRequest()
	buff := make([]byte, 1024)

	readToIndex := 0

	for !req.doneCheck() {
		read, err := reader.Read(buff[readToIndex:])
		if err != nil {
			return nil, fmt.Errorf("error while reading")
		}

		readToIndex += read

		parseRead, err := req.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[parseRead:readToIndex])
		readToIndex -= parseRead
	}

	return req, nil
}
