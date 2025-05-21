package draw

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawLine draws a line from (x0,y0) to (x1,y1) with the given color
func DrawLine(img *ebiten.Image, x0, y0, x1, y1 float64, clr color.Color) {
	steps := int(math.Max(math.Abs(x1-x0), math.Abs(y1-y0)))
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x0 + (x1-x0)*t
		y := y0 + (y1-y0)*t
		img.Set(int(x), int(y), clr)
	}
}

// DrawCircle draws a filled circle with center (cx,cy) and radius r
func DrawCircle(img *ebiten.Image, cx, cy, r int, clr color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(cx+x, cy+y, clr)
			}
		}
	}
}
