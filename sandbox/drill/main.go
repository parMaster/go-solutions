package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Parser struct {
	delim byte
}

type Option func(*Parser) error

func NewParser(opts ...Option) *Parser {

	p := Parser{}

	for _, opt := range opts {
		opt(&p)
	}

	return &p
}

func UseDelim(d byte) Option {
	return func(p *Parser) error {
		if string(d) == "" {
			return fmt.Errorf("empty delimiter: %b", d)
		}

		p.delim = d
		return nil
	}
}

func (p *Parser) Parse(file string) error {

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("can't open '%s': %e", file, err)
	}

	b := bufio.NewReader(f)
	for {
		line, err := b.ReadString(p.delim)
		if err != nil {
			break
		}

		// fmt.Printf(line)

		fields := strings.Fields(line)
		fmt.Printf("%s - %s \n", fields[0][1:5], strings.Join(fields[7:], " "))
	}

	return nil
}

// I call it "horisontal" approach, because I don't know the right word
// Pretty much the same, but different
func New() *Parser {
	p := Parser{}
	return &p
}

func (p *Parser) WithDelimiter(delim byte) *Parser {
	p.delim = delim
	return p
}

func main() {

	// kind of "vertical" approach
	p := NewParser(UseDelim('\n'))
	p.Parse("smc.txt")

	// "horisontal" approach
	New().WithDelimiter('\n').Parse("smc.txt")
}
