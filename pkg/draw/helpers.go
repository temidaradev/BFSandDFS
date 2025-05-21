package draw

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawLine draws a line from (x0,y0) to (x1,y1) with the given color
func DrawLine(img *ebiten.Image, x0, y0, x1, y1 float64, clr color.Color) {
	// Use vector graphics instead of pixel-by-pixel operations
	lineImg := ebiten.NewImage(1, 1)
	lineImg.Fill(clr)

	// Calculate line length and angle
	length := math.Sqrt((x1-x0)*(x1-x0) + (y1-y0)*(y1-y0))
	angle := math.Atan2(y1-y0, x1-x0)

	// Create transform options
	opts := &ebiten.DrawImageOptions{}

	// Scale to match line length (horizontal scaling to the length of the line)
	opts.GeoM.Scale(length, 1)

	// Rotate to match line angle
	opts.GeoM.Rotate(angle)

	// Position the line
	opts.GeoM.Translate(x0, y0)

	// Draw the line
	img.DrawImage(lineImg, opts)
}

// DrawCircle draws a filled circle with center (cx,cy) and radius r
func DrawCircle(img *ebiten.Image, cx, cy, r int, clr color.Color) {
	// Create a circle image instead of drawing pixel by pixel
	diameter := r * 2
	circleImg := ebiten.NewImage(diameter, diameter)

	// Draw the circle into the image
	for y := 0; y < diameter; y++ {
		for x := 0; x < diameter; x++ {
			dx := float64(x - r)
			dy := float64(y - r)
			if dx*dx+dy*dy <= float64(r*r) {
				circleImg.Set(x, y, clr)
			}
		}
	}

	// Create transform options
	opts := &ebiten.DrawImageOptions{}

	// Position the circle
	opts.GeoM.Translate(float64(cx-r), float64(cy-r))

	// Draw the circle
	img.DrawImage(circleImg, opts)
}
