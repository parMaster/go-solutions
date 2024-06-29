package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Opts struct {
	cmd   string
	flags []string
	twist bool
	cnt   int
}

type OptsFunc func(*Opts)

func DefaultOpts() Opts {
	return Opts{
		cmd:   "echo",
		flags: []string{"-v"},
		twist: false,
		cnt:   1,
	}
}

func WithOpts(newOpts Opts) OptsFunc {
	return func(opts *Opts) {
		opts.cmd = newOpts.cmd
		opts.twist = newOpts.twist
		opts.flags = newOpts.flags
		opts.cnt = newOpts.cnt
	}
}

func WithATwist(opts *Opts) {
	opts.twist = true
}

func WithFlags(flags []string) OptsFunc {
	return func(opts *Opts) {
		opts.flags = flags
	}
}

func WithCmd(cmd string) OptsFunc {
	return func(opts *Opts) {
		opts.cmd = cmd
	}
}

type Worker struct {
	Opts
}

func NewWorker(opts ...OptsFunc) *Worker {
	o := DefaultOpts()

	for _, optFunc := range opts {
		optFunc(&o)
	}

	return &Worker{
		Opts: o,
	}
}

func Test_Opts(t *testing.T) {

	defaultOpts := DefaultOpts()

	defaultW := NewWorker()
	assert.Equal(t, defaultW.Opts, defaultOpts)

	defaultOpts.twist = true
	w := NewWorker(WithATwist)
	assert.Equal(t, w.Opts, defaultOpts)

	defaultOpts.twist = true
	defaultOpts.cmd = "ls"
	defaultOpts.flags = []string{"one", "two", "three"}

	w = NewWorker(WithATwist,
		WithFlags([]string{"one", "two", "three"}),
		WithCmd("ls"))
	assert.Equal(t, w.Opts, defaultOpts)
}

func TestWithOpts(t *testing.T) {

	newOpts := Opts{
		cmd:   "ps",
		flags: []string{"-a", "-u", "-x"},
		twist: true,
		cnt:   16,
	}

	w := NewWorker(WithOpts(newOpts))
	assert.Equal(t, newOpts, w.Opts)

}
