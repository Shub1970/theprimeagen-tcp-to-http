// package headers

package main

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var REGISTER_NURSE = "\r\n"

func parseHeader(headerLine []byte) (string, string, error) {

	splitData := bytes.SplitN(headerLine, []byte(":"), 2)

	fmt.Printf("data 1: %s, data 2: %s \n", string(splitData[0]), string(splitData[1]))
	if len(splitData) != 2 {
		return "", "", fmt.Errorf("data is not fill")
	}

	name := splitData[0]
	value := bytes.TrimSpace(splitData[1])

	if bytes.HasSuffix(value, []byte(" ")) {
		return "", "", fmt.Errorf("name contain space")
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {

	read := 0
	done := false

	fmt.Printf("data input: %s\n", string(data))
	for {
		idx := bytes.Index(data[read:], []byte(REGISTER_NURSE))

		if idx == -1 {
			return 0, false, nil
		}

		if idx == 1 {
			done = true
			return 0, false, nil
		}

		name, value, err := parseHeader(data[read : read+idx])

		if err != nil {
			return 0, false, err
		}

		dat := {
			return 0, false, err
		}

		read += idx + len(REGISTER_NURSE)

		h[name] = value
	}

	return read, done, nil
}

func main() {
	h := Headers{}
	data := []byte("Content-Type: text/html\r\nContent-Length: 123\r\n\r\n")
	n, done, err := h.Parse(data)
	fmt.Println("Parsed bytes:", n, "Done:", done, "Headers:", h, "Error:", err)
}
