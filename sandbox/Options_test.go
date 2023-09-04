package sandbox

// Remembering Options options
// no copilot

import (
	"fmt"
	"testing"

	log "github.com/go-pkgz/lgr"
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

// ===
// === ChatGPTs best take on options

// Config is the struct we want to initialize with various options
type Config struct {
	Option1 string
	Option2 int
	Option3 bool
}

// ConfigOption is a functional option type that modifies the Config struct
type ConfigOption func(*Config)

// WithOption1 sets the Option1 field of the Config struct
func WithOption1(value string) ConfigOption {
	return func(c *Config) {
		c.Option1 = value
	}
}

// WithOption2 sets the Option2 field of the Config struct
func WithOption2(value int) ConfigOption {
	return func(c *Config) {
		c.Option2 = value
	}
}

// WithOption3 sets the Option3 field of the Config struct
func WithOption3(value bool) ConfigOption {
	return func(c *Config) {
		c.Option3 = value
	}
}

// NewConfig initializes a Config instance with the provided options
func NewConfig(options ...ConfigOption) *Config {
	config := &Config{}

	// Apply each option to the config
	for _, option := range options {
		option(config)
	}

	return config
}

func main() {
	// Create a new Config instance with various options
	config := NewConfig(
		WithOption1("Value 1"),
		WithOption2(42),
		WithOption3(true),
	)

	// Print the configured values
	fmt.Println("Option1:", config.Option1)
	fmt.Println("Option2:", config.Option2)
	fmt.Println("Option3:", config.Option3)
}

// In this example, we define a Config struct with three configurable options: Option1, Option2, and Option3.
// We also define functional option types like WithOption1, WithOption2, and WithOption3 that modify the respective
// fields in the Config struct. The NewConfig function takes a variadic number of these options and applies them
// to create a new Config instance.

// This approach allows you to create instances of the Config struct with different combinations of options,
// making the code more readable and maintaining a clear separation of concerns.

//

//

/// Once again

type M struct {
	config string
}

func (m *M) Report() {
	log.Printf("m.config=%s", m.config)
}

type FOpt func(*M)

func NewM(opts ...FOpt) *M {
	m := M{}
	for _, o := range opts {
		o(&m)
	}
	return &m
}

func WithOpt1(opt string) FOpt {
	return func(m *M) {
		m.config = opt
	}
}
