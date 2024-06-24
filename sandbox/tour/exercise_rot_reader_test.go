package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type rot13Reader struct {
	r io.Reader
}

func (rot rot13Reader) Read(out []byte) (int, error) {
	n, err := rot.r.Read(out)

	for i := range n {
		if out[i] >= 'A' && out[i] <= 'Z' {
			out[i] = 'A' + (out[i]-'A'+13)%26
		} else if out[i] >= 'a' && out[i] <= 'z' {
			out[i] = 'a' + (out[i]-'a'+13)%26
		}
	}

	return n, err
}

func Test_Decode(t *testing.T) {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)

	s = strings.NewReader("Lbh penpxrq gur pbqr!")
	r = rot13Reader{s}
	buf := make([]byte, len("Lbh penpxrq gur pbqr!"))
	b := bufio.NewReader(&r)
	b.Read(buf)
	assert.Equal(t, "You cracked the code!", string(buf))

}
