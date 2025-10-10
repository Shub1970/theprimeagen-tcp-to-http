package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var REGISTER_NEURSE = "\r\n"
var MAILFORMATE_IN_FIELD_LINE = fmt.Errorf("malformate in field line")

func NewHeaders() Headers {
	return map[string]string{}
}

func parsingFieldLine(fieldLine []byte) (string, string, error) {
	fmt.Printf("fieldLine:%s\n", fieldLine)
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

func (h Headers) Parse(data []byte) (int, bool, error) {
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
		h[key] = value

		read += idx + len(REGISTER_NEURSE)
	}

	return read, status, nil
}
