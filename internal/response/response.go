package response

import (
	"TCPtoHTTP/internal/headers"
	"fmt"
	"io"
)

type StatusCode int

const (
	StatusOk          StatusCode = 200
	StatusBad         StatusCode = 400
	StatusServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusOut := []byte{}
	switch statusCode {
	case StatusOk:
		statusOut = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBad:
		statusOut = []byte("HTTP/1.1 400 Bad Request\r\n")

	case StatusServerError:
		statusOut = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("reach unexpected status")
	}

	_, err := w.Write(statusOut)
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	defaultsHeader := headers.NewHeaders()
	defaultsHeader.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	defaultsHeader.Set("Connection", "close")
	defaultsHeader.Set("Content-Type", "text/plain")
	return defaultsHeader
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	headersFinal := []byte{}

	for key, value := range headers.Headers {
		headersFinal = fmt.Appendf(headersFinal, fmt.Sprintf("%s:%s\r\n", key, value))
	}

	headersFinal = fmt.Appendf(headersFinal, "\r\n")

	_, err := w.Write(headersFinal)

	return err
}
