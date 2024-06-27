// 100 mistakes in Go and how to avoid them. Tests and memos.

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Runes(t *testing.T) {

	s := "hêllo"

	// ranging over runes
	for i, r := range s {
		fmt.Printf("position %d: %c\n", i, r)
	}

	// ranging over runes reindexed
	runes := []rune(s)
	for i, r := range runes {
		fmt.Printf("position %d: %c\n", i, r)
	}

	// ranging over bytes
	for i := 0; i < len(s); i++ {
		fmt.Printf("position %d: %c\n", i, s[i])
	}

	fmt.Println(s[1:3]) // "ê" [195 170] 2-byte UTF-8 encoded rune
	fmt.Println(s[1])   // 195 (Ã)
	fmt.Println(s[2])   // 170 (ª)

	fmt.Println(string([]rune(s)[1])) // "ê" UTF-8 encoded rune
}

// strings are immutable

func Test_Strings(t *testing.T) {
	list := []string{"h", "ê", "ll", "o"}

	// 1. inefficient strings concatenation
	result := ""
	for _, value := range list {
		// each concatenation creates a new string, allocating memory
		result += value
	}
	assert.Equal(t, "hêllo", result)

	// 2. much more efficient strings concatenation
	sb := strings.Builder{}
	for _, value := range list {
		_, _ = sb.WriteString(value)
	}
	sbResult := sb.String()
	assert.Equal(t, "hêllo", sbResult)

	// 3. the most efficient strings concatenation
	concatResult := concat(list)
	assert.Equal(t, "hêllo", concatResult)
}

// the most efficient strings concatenation
func concat(values []string) string {

	// calculate the total length of the resulting string
	total := 0
	for i := 0; i < len(values); i++ {
		total += len(values[i])
	}

	sb := strings.Builder{}
	sb.Grow(total) // preallocate memory
	for _, value := range values {
		_, _ = sb.WriteString(value)
	}
	return sb.String()
}

// naked returns

// what is returned?
func returnInitialized() (result int, resMap map[int]string) {
	return
}

func Test_NakedReturns(t *testing.T) {
	result, resMap := returnInitialized()
	assert.Equal(t, 0, result)
	assert.Nil(t, resMap)
}

// what is returned?
func returnInitializedModified() (result int, resMap map[int]string) {
	list := []string{"a", "b", "c"}
	resMap = map[int]string{} // so accesing resMap is safe
	for _, s := range list {
		if s == "x" && resMap[1] == "x" {
			result = 1
			resMap[1] = "zzz"
			return
		}
	}
	return
}

func Test_NakedReturnsModified(t *testing.T) {
	result, resMap := returnInitializedModified()
	assert.Equal(t, 0, result)
	assert.NotNil(t, resMap) // resMap is initialized, never nil
}

// nil receiver

type Foo struct {
	list []string
}

// methods are just syntactic sugar for functions with a receiver argument
func (foo *Foo) Bar() string {
	return "bar"
}

// same as
func Bar(foo *Foo) string {
	return "bar"
}

func (foo *Foo) BarList() string {
	return strings.Join(foo.list, ",") // panic if foo is nil
}

func Test_NilReceiver(t *testing.T) {
	var foo *Foo
	assert.Equal(t, "bar", foo.Bar()) // no panic, nil receiver is safe
	assert.Equal(t, "bar", Bar(foo))  // no panic, nil receiver is safe
	assert.Panics(t, func() {
		foo.BarList()
	}, "panic expected") // panic, nil receiver is not safe
}

// but wrapped in an interface, nil receiver is safe

// satisfies error interface because of Error()string method
type MultiError struct {
	errs []string
}

func (m *MultiError) Error() string {
	return strings.Join(m.errs, ";")
}

func (m *MultiError) Validate() error {
	var er *MultiError
	return er
}

func Test_NilReceiverInterface(t *testing.T) {
	var foo *MultiError
	var err error = foo // wrapped in an interface

	// not nil interface
	assert.NotEqual(t, nil, err)
	// nil value wrapped in an interface
	assert.Equal(t, "<nil>", fmt.Sprintf("%v", err))
	assert.NotEqual(t, nil, foo.Validate()) // nil wrapped in an interface is not nil

	// in short:
	var ef = func() error {
		var er *MultiError
		return er
	}
	// nil value returned from a function is wrapped in an interface
	assert.NotEqual(t, nil, ef())

	// but calling a method on nil receiver is not safe
	assert.Panics(t, func() {
		_ = err.Error()
	}, "still panic expected")
}

// defer evaluation, nifty trick
func Test_DeferEvaluation(t *testing.T) {
	var i = 0
	var j = 0
	defer func(out int) {
		assert.Equal(t, out, 0) // evaluated at the time of defer
		// reference to j outside of the closure
		assert.Equal(t, j, 1) // evaluated at the func end
	}(i)
	j = 1
}

// error handling

type transientError struct {
	err error
}

func (t transientError) Error() string {
	return fmt.Sprintf("transient error: %v", t.err)
}

func Validate(i int) error {
	if i == 0 {
		return fmt.Errorf("id validation fail: %d", i)
	}

	return nil
}

// switch on error type:
func Test_TransientErrorSwitchOnType(t *testing.T) {

	var trans = func(val int) error {
		if val == 5 {
			return fmt.Errorf("val == 5")
		}
		err := Validate(0)
		if err != nil {
			return transientError{err}
		}
		return nil
	}

	err := trans(0)
	switch err.(type) {
	case transientError:
		log.Printf("got error: %v", err)
	default:
		log.Printf("other error: %v", err)
	}

	err = trans(5)
	switch err.(type) {
	case transientError:
		log.Printf("got error: %v", err)
	default:
		log.Printf("other error: %v", err)
	}
}

// (!) Expected errors should be designed as error values (sentinel errors):
//		var ErrFoo = errors.New("foo").
// (!) Unexpected errors should be designed as error types:
// 		type BarError struct { ... }, with BarError implementing the error interface.

var ErrFoo = errors.New("foo")

func Test_ErrorIs(t *testing.T) {
	err := fmt.Errorf("no rows in the result: %w", sql.ErrNoRows)
	assert.True(t, errors.Is(err, sql.ErrNoRows))

	fooErr := fmt.Errorf("foo error: %w", ErrFoo)
	assert.False(t, errors.Is(fooErr, errors.New("foo")))
	// but
	assert.True(t, errors.Is(fooErr, ErrFoo))
}

func Test_DeferError(t *testing.T) {

	errf := func(i int) (err error) {
		defer func() {
			// err = connection.Close or something that can return error
			if i == 1 {
				err = errors.New("error from defer")
			}
		}()

		if i == 0 {
			return errors.New("error i==0")
		}

		return
	}

	assert.NoError(t, errf(2))
	fmt.Println(errf(2)) // nil

	assert.Error(t, errf(1))
	fmt.Println(errf(1)) // error from defer

	assert.Error(t, errf(0))
	fmt.Println(errf(0)) // error i==0
}

// quick reminder of channel owner worker with a context
func worker(ctx context.Context, input <-chan int) chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for {
			select {
			case in, ok := <-input:
				if !ok {
					return
				}
				log.Println("input received", in)
				out <- in * 2
			case <-ctx.Done():
				log.Println("context done", ctx.Err())
				return
			}
		}
	}()

	return out
}

func Test_Worker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int)
	results := worker(ctx, ch)
	go func() {
		for i := range 10 {
			ch <- i
		}
		close(ch)
	}()

	received := 0
	for res := range results {
		log.Println("result received", res)
		received++
	}
	assert.Equal(t, 10, received)
}

// todo: exercise: read bytes from file, spin up a pool of goroutines that count a number of
// some specific byte
