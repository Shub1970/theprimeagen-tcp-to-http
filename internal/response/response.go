package response

import (
	"TCPtoHTTP/internal/headers"
	"fmt"
	"io"
)

type Response struct {
}

type StatusCode int

const (
	StatusOk          StatusCode = 200
	StatusBadRequest  StatusCode = 400
	StatusServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusOut := []byte{}
	switch statusCode {
	case StatusOk:
		statusOut = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
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

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusOut := []byte{}
	switch statusCode {
	case StatusOk:
		statusOut = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusOut = []byte("HTTP/1.1 400 Bad Request\r\n")

	case StatusServerError:
		statusOut = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("reach unexpected status")
	}

	_, err := w.writer.Write(statusOut)
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	headersFinal := []byte{}

	for key, value := range headers.Headers {
		headersFinal = fmt.Appendf(headersFinal, fmt.Sprintf("%s:%s\r\n", key, value))
	}

	headersFinal = fmt.Appendf(headersFinal, "\r\n")

	_, err := w.writer.Write(headersFinal)

	return err
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	return n, err
}
