package draw

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawRect draws a rectangle on the screen
func DrawRect(screen *ebiten.Image, x, y, width, height float64, color color.RGBA) {
	// Draw top border
	DrawLine(screen, x, y, x+width, y, color)
	// Draw bottom border
	DrawLine(screen, x, y+height, x+width, y+height, color)
	// Draw left border
	DrawLine(screen, x, y, x, y+height, color)
	// Draw right border
	DrawLine(screen, x+width, y, x+width, y+height, color)

	// Fill the inside of the rectangle (simple fill for now)
	for i := 0; i < int(height); i++ {
		DrawLine(screen, x, y+float64(i), x+width, y+float64(i), color)
	}
}

// DrawCircle draws a circle on the screen using the midpoint circle algorithm
// ... existing code ...
