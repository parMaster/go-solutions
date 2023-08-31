package sandbox

// Remembering Options options
// no copilot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Worker struct {
	cmd   string
	flags []string
	twist bool
	cnt   int
}

func NewWorker() *Worker {
	return &Worker{}
}

// Getters-setters kind of options:

func (w *Worker) WithCmd(cmd string) error {
	if cmd == "" {
		return fmt.Errorf("WithCmd received empty cmd: %s", cmd)
	}
	w.cmd = cmd
	return nil
}

// Trigger option

func (w *Worker) Twist() {
	w.twist = true
}

func Test_Options(t *testing.T) {

	w := NewWorker()
	w.WithCmd("/usr/sbin/cmd")

}

// Func option

type Option func(w *Worker)

func WithFlags(flags []string) Option {
	return func(w *Worker) {
		w.flags = append(w.flags, flags...)
	}
}

func NewOptsWorker(opts ...Option) *Worker {
	w := Worker{}
	for _, opt := range opts {
		opt(&w)
	}
	return &w
}

func (w *Worker) InspectFlags() []string {
	return w.flags
}

func Test_FuncOptions(t *testing.T) {

	w := NewOptsWorker(
		WithFlags([]string{"test"}),
		WithFlags([]string{"test", "test"}),
	)

	assert.Equal(t, []string{"test", "test", "test"}, w.InspectFlags())
}

// Once again

type Opt func(w *Worker)

func Inc(n int) Opt {
	return func(w *Worker) {
		if n > 1 {
			w.cnt += n
		} else {
			w.cnt++
		}
	}
}

func NewWOpt(opts ...Opt) *Worker {
	w := Worker{}

	for _, opt := range opts {
		opt(&w)
	}

	return &w
}

func Test_FuncOpts(t *testing.T) {

	w := NewWOpt(Inc(1), Inc(1), Inc(10))

	assert.Equal(t, 12, w.cnt)
}
