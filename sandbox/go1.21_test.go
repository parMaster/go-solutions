package sandbox

import (
	"fmt"
	"log"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MinMax(t *testing.T) {
	// t.Fatal("not implemented")

	var x, y int

	y = 5

	m := min(x) // m == x
	assert.Equal(t, m, x)

	m = min(x, y) // m is the smaller of x and y
	assert.Equal(t, m, x)

	m = max(x, y) // m is the larger of x and y
	assert.Equal(t, m, y)

	m = max(x, y, 10) // m is the larger of x and y but at least 10
	assert.Equal(t, m, 10)

	// For numeric arguments, assuming all NaNs are equal, min and max are commutative and associative:
	assert.Equal(t, min(x, y), min(y, x))
	assert.Equal(t, max(x, y), max(y, x))

	//min(x, y, z) == min(min(x, y), z) == min(x, min(y, z))
	z := 10
	assert.Equal(t, min(x, y, z), min(min(x, y), z))
	assert.Equal(t, min(x, y, z), min(x, min(y, z)))

	c := max(1, 2.0, 10) // c == 10.0 (floating-point kind)
	assert.Equal(t, c, 10.0)
	assert.IsType(t, c, float64(10.0))

	f := max(0, float32(x)) // type of f is float32
	assert.Equal(t, f, float32(0))
	assert.IsType(t, f, float32(0))

	// var s []string
	// _ = min(s...)               // invalid: slice arguments are not permitted
	mx := max("", "foo", "bar") // mx == "foo" (string kind)
	assert.Equal(t, mx, "foo")

	// strings
	assert.Equal(t, min("a", "b", "c"), "a")
	assert.Equal(t, max("a", "b", "c"), "c")

	assert.Equal(t, min("foo", "bar"), min("bar", "foo"))
	assert.Equal(t, max("foo", "bar"), max("bar", "foo"))

	slices.Contains([]string{"a", "b", "c"}, "a")
}

func Test_clear(t *testing.T) {
	//

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	assert.Equal(t, m["a"], 1)
	assert.Equal(t, m["b"], 2)
	assert.Equal(t, m["c"], 3)
	assert.Equal(t, len(m), 3)

	clear(m)

	log.Printf("m: %v", m) // "m: map[]" (the map has been cleared)
	v, ok := m["a"]
	log.Printf("m['a']: %v", v) // "m['a']: 0" (note: ok is false)
	assert.Equal(t, ok, false)
	assert.Equal(t, v, 0)

	log.Printf("m['f']: %v", m["f"]) // "m['f']: 0"
	assert.Equal(t, m["a"], 0)
	assert.Equal(t, m["b"], 0)
	assert.Equal(t, m["c"], 0)
	assert.Equal(t, len(m), 0)

	s := []int{1, 2, 3}

	assert.Equal(t, s[0], 1)
	assert.Equal(t, s[1], 2)
	assert.Equal(t, s[2], 3)
	log.Printf("&s[0]: %v", &s[0])

	clear(s)

	log.Printf("s: %v", s) // "s: [0 0 0]"
	log.Printf("&s[0]: %v", &s[0])
	assert.Equal(t, s[0], 0)
	assert.Equal(t, s[1], 0)
	assert.Equal(t, s[2], 0)
	assert.Equal(t, len(s), 3)
	assert.Equal(t, cap(s), 3)
}

// Test for values
func Test_forValues(t *testing.T) {
	done := make(chan bool)

	values := []string{"a", "b", "c"}
	for _, v := range values {
		// v := v // create a new 'v'. - not necessary since 1.21 with GOEXPERIMENT=loopvar environment variable
		go func() {
			fmt.Println(v) // nolint: // loop variable v captured by func literal
			done <- true
		}()
	}

	// wait for all goroutines to complete before exiting
	for range values {
		<-done
	}
}

// New testing/slogtest package
// The new testing/slogtest package can help to validate slog.Handler implementations.

// New slices package
// The new slices package provides many common operations on slices, using generic functions that work with slices of any element type.

// New maps package
// The new maps package provides several common operations on maps, using generic functions that work with maps of any key or element type.

// New cmp package
// The new cmp package defines the type constraint Ordered and two new generic functions Less and Compare that are useful with ordered types.
