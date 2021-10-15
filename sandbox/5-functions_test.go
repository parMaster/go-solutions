package sandbox

import (
	"fmt"
	"math"
	"testing"
)

func stub() bool {
	return false
}

// func flushICache(begin, end int) // implemented externally, fuction body ommited

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func TestFunctions(t *testing.T) {
	// flushICache(0, 0)
	fmt.Println(stub())
	fmt.Println(min(1, 2))
}

// func IndexRune(s string, r rune) int {
// 	for i, c := range s {
// 		if c == r {
// 			return i
// 		}
// 	}
// 	// invalid: missing return statement
// 	return false
// }

// Pointer-Receiver
type Vertex struct {
	X, Y float64
}

func (rec Vertex) Abs() float64 {
	return math.Sqrt(rec.X*rec.X + rec.Y*rec.Y)
}

func (rec *Vertex) Scale(f float64) {
	rec.X = rec.X * f
	rec.Y = rec.Y * f
}

// Same as this function
func Scale2(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func TestPointerReceiver(t *testing.T) {
	v := Vertex{3, 4}
	v.Scale(10)
	fmt.Println(v.Abs())

	Scale2(&v, 10)
	fmt.Println(v.Abs())
}
