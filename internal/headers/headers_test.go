package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRequestLineParse(t *testing.T) {
	// ✅ Test: Valid multiple headers
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\ntooo: tehiso\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "tehiso", headers.Get("tooo"))
	assert.Equal(t, 39, n)
	assert.True(t, done)

	// ❌ Test: Invalid spacing in header
	headers = NewHeaders()
	data = []byte("       Host : localhst:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// ❌ Test: Invalid name validation (contains ©)
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// ✅ Test: Multiple valid headers with same key (should overwrite or last-one-wins)
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go\r\nSet-Person: prime-loves-zig\r\nSet-Person: tj-loves-ocaml\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.True(t, done)
	assert.Greater(t, n, 0)
	assert.Equal(t, "lane-loves-go,prime-loves-zig,tj-loves-ocaml", headers.Get("Set-Person")) // last value wins
}
