package main

import (
	"testing"

	"golang.org/x/tour/reader"
)

type MyReader struct{}

func (mr MyReader) Read(buf []byte) (int, error) {
	size := len(buf)
	for i := range size {
		buf[i] = byte('A')
	}
	return size, nil
}

func Test_ExerciseReader(t *testing.T) {
	reader.Validate(MyReader{})
}
