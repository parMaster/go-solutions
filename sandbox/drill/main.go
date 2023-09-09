// Drilling Read file line by line and functiona options
// smc.txt is from https://github.com/acidanthera/VirtualSMC/blob/master/Docs/SMCSensorKeys.txt
package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Record string

type fr struct {
	data   map[string]Record
	delim  byte
	sorted bool
}

type Option func(*fr)

func NewFileReader(opts ...Option) *fr {
	r := fr{}

	r.data = make(map[string]Record)

	for _, opt := range opts {
		opt(&r)
	}
	return &r
}

func WithDelim(d byte) Option {
	return func(r *fr) {
		r.delim = d
	}
}

func Sorted() Option {
	return func(r *fr) {
		r.sorted = true
	}
}

func (r *fr) Read(path string) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fb := bufio.NewReader(f)

	for {
		line, err := fb.ReadString(r.delim)

		if err != nil {
			break
		}

		fields := strings.Fields(string(line))

		if len(fields) == 0 {
			continue
		}

		r.data[fields[0][1:5]] = Record(strings.Join(fields[7:], " "))
	}

	return nil
}

func (r *fr) Display() {

	keys := Keys(r.data)
	if r.sorted {
		slices.Sort(keys)
	}

	for _, k := range keys {
		fmt.Printf("%s | %s \n", k, r.data[k])
	}
}

func Keys[V any](m map[string]V) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {

	r := NewFileReader(WithDelim('\n'), Sorted())

	r.Read("smc.txt")

	r.Display()
}
