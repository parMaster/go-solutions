package kindergarten

import (
	"errors"
	"slices"
	"strings"
)

// Define the Garden type here.

// The diagram argument starts each row with a '\n'.  This allows Go's
// raw string literals to present diagrams in source code nicely as two
// rows flush left, for example,
//
//     diagram := `
//     VVCCGG
//     VVCCGG`

type Garden struct {
	diagram  []string
	children map[string]int
}

var plants = map[byte]string{
	'G': "grass",
	'C': "clover",
	'R': "radishes",
	'V': "violets",
}

var validCupCodes = []rune{
	'G', 'C', 'R', 'V',
}

func NewGarden(diagram string, childrenInput []string) (*Garden, error) {
	if diagram[0] != '\n' {
		return nil, errors.New("error diagram format")
	}
	lines := strings.Split(strings.Trim(strings.TrimSpace(diagram), "\n"), "\n")
	for _, line := range lines {
		if len(line) != len(lines[0]) {
			return nil, errors.New("mismatched lines")
		}
		if len(line)%2 != 0 {
			return nil, errors.New("odd number of cups")
		}
		for _, cupCode := range line {
			if !slices.Contains(validCupCodes, cupCode) {
				return nil, errors.New("invalid cup code")
			}
		}
	}
	children := slices.Clone(childrenInput)
	slices.Sort(children)
	cMap := map[string]int{}

	for i, child := range children {
		if _, ok := cMap[child]; ok {
			return nil, errors.New("duplicate name")
		}
		cMap[child] = i
	}

	g := &Garden{
		diagram:  lines,
		children: cMap,
	}
	return g, nil
}

func (g *Garden) Plants(child string) ([]string, bool) {
	ci, ok := g.children[child]
	if !ok {
		return nil, false
	}

	res := []string{}
	for _, line := range g.diagram {
		res = append(res, plants[line[ci*2]])
		res = append(res, plants[line[ci*2+1]])
	}

	return res, true
}
