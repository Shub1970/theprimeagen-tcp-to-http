package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var REGISTER_NEURSE = "\r\n"
var MAILFORMATE_IN_FIELD_LINE = fmt.Errorf("malformate in field line")

type Headers struct {
	Headers map[string]string
}

func (h *Headers) Get(key string) (string, bool) {
	toLow := strings.ToLower(key)
	val, ok := h.Headers[toLow]
	fmt.Printf("header get access:%s,output val :%s\n", toLow, val)
	return val, ok
}

func (h *Headers) Set(key string, val string) {
	lowKey := strings.ToLower(key)
	if mapValue, ok := h.Headers[lowKey]; ok {
		h.Headers[lowKey] = fmt.Sprintf("%s,%s", mapValue, val)
	} else {
		h.Headers[lowKey] = val
	}
}

func NewHeaders() *Headers {
	return &Headers{
		Headers: map[string]string{},
	}
}

func isValidFieldName(name string) bool {
	if name == "" {
		return false
	}
	for _, ch := range name {
		if ch >= 'A' && ch <= 'Z' {
			continue
		}
		if ch >= 'a' && ch <= 'z' {
			continue
		}
		if ch >= '0' && ch <= '9' {
			continue
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			return false
		}
	}
	return true
}

func parsingFieldLine(fieldLine []byte) (string, string, error) {
	fieldLinePart := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(fieldLinePart) != 2 {
		return "", "", MAILFORMATE_IN_FIELD_LINE
	}

	field_name, field_value := fieldLinePart[0], bytes.TrimSpace(fieldLinePart[1])

	if bytes.HasSuffix(field_name, []byte(" ")) {
		return "", "", MAILFORMATE_IN_FIELD_LINE
	}

	return string(field_name), string(field_value), nil
}

func (head *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	status := false
	for {
		idx := bytes.Index(data[read:], []byte(REGISTER_NEURSE))
		if idx == -1 {
			break
		}

		if idx == 0 {
			status = true
			read += len(REGISTER_NEURSE)
			break
		}

		key, value, err := parsingFieldLine(data[read : read+idx])
		if err != nil {
			return read, false, err
		}

		if !isValidFieldName(key) {
			return read, false, fmt.Errorf("field name is not following the RFC")
		}
		head.Set(key, value)
		read += idx + len(REGISTER_NEURSE)
	}

	return read, status, nil
}
