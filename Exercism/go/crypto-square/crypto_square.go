package cryptosquare

import (
	"math"
	"strings"
)

func Encode(pt string) string {
	pt = strings.ToLower(pt)
	var s string
	// screw regexp
	for _, p := range pt {
		if (p >= '0' && p <= '9') || (p >= 'a' && p <= 'z') {
			s += string(p)
		}
	}

	// finding rectangle size
	cols := int(math.Ceil(math.Sqrt(float64(len(s)))))
	rows := cols - 1

	if len(s) > rows*cols {
		rows++
	}

	if len(s) < cols*rows {
		for i := len(s); i <= cols*rows; i++ {
			s += " "
		}
	}

	// encrypting
	vlines := []string{}
	for i := 0; i < cols; i++ {
		vline := ""
		for j := 0; j < rows; j++ {
			if j*cols+i < len(s) {
				vline = vline + string(s[j*cols+i])
			}
		}
		vlines = append(vlines, vline)
	}

	return strings.Join(vlines, " ")
}
