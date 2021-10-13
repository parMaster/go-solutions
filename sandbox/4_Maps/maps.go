package sandbox

import (
	"fmt"
	"testing"
)

type Vertex struct {
	Lat, Long float64
}

var m map[string]Vertex

func TestMaps(t *testing.T) {
	m = make(map[string]Vertex)
	m["Bell Labs"] = Vertex{
		40.68433, -74.39967,
	}
	m["Google"] = Vertex{
		41.68433, -72.39967,
	}
	fmt.Println(m["Google"])

	var ms = map[string]Vertex{
		"Key 1": {40.68433, -74.39967},
		"Key 2": {41.68433, -72.39967},
	}

	ms["Key 3"] = Vertex{3.33333, 4.44444}
	ms["Key 3"] = Vertex{5.55555, 6.66666}

	fmt.Println(ms)

	// Test that a key is present with a two-value assignment:
	// elem, ok = m[key]

	elem, ok := ms["Key 3"]
	fmt.Println("The value:", elem, "Present?", ok)

	delete(ms, "Key 3")

	elem, ok = ms["Key 3"]
	fmt.Println("The value:", elem, "Present?", ok)

	fmt.Println(ms)

}
