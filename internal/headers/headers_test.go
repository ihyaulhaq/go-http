package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders_Parse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Capital
	headers = NewHeaders()
	data = []byte("HOST: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid 2 Headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nContent-Type: application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 23, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers.Get("host"))
	assert.Equal(t, "application/json", headers.Get("content-type"))
	assert.Equal(t, 32, n)
	assert.False(t, done)

	// Test: Valid 2 samw Headers
	headers = NewHeaders()
	data = []byte("Set-Person: localhost:42069\r\nSet-Person: application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 29, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069,application/json", headers.Get("set-person"))
	assert.Equal(t, 30, n)
	assert.False(t, done)

	// Test: vaid Done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid Char
	headers = NewHeaders()
	data = []byte("H@ST: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

}
