package main

import (
	"image"
	"image/color"
	"log"
	"testing"

	"golang.org/x/tour/pic"
)

type Image struct{}

// ColorModel returns the Image's color model.
func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, 256, 256)
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (i Image) At(x, y int) color.Color {
	r := uint8(float32(x) / 640 * 250)
	g := uint8(float32(y) / 640 * 250)
	b := uint8(float32(x+y) / 640 * 250)
	return color.RGBA{r, g, b / 3, 100}
}

func Test_At(t *testing.T) {
	m := Image{}
	log.Println(m.At(0, 0))
}

func main() {
	m := Image{}
	pic.ShowImage(m)
}
