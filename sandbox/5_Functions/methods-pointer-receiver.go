package main

import (
	"fmt"
	"math"
)

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
// func Scale(v *Vertex, f float64) {
// 	v.X = v.X * f
// 	v.Y = v.Y * f
// }
// called like this:
// Scale(&v, 10)

func main() {
	v := Vertex{3, 4}
	v.Scale(10)
	fmt.Println(v.Abs())
}
