package main

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Files(t *testing.T) {

	f, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE, 0644)
	assert.NoError(t, err)

	n, err := io.WriteString(f, "test "+time.Now().Format("2006-01-02 15:04:05"))
	assert.NoError(t, err)
	assert.NotZero(t, n)
	assert.True(t, (n > 10))

	// io.WriteString into a file
	io.WriteString(f, "\r\n")

	// or File.WriteString
	n, err = f.WriteString("Writing String")
	assert.NoError(t, err)
	assert.Equal(t, len("Writing String"), n)

	f.Close()

	f, err = os.OpenFile("test.txt", os.O_RDWR|os.O_APPEND, 0644)
	defer f.Close()
	assert.NoError(t, err)

	f.WriteString("\r\n")
	f.WriteString("Appended @: " + time.Now().Format("2006-01-02 15:04:05"))

}
