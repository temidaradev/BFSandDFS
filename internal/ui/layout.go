package ui

import "github.com/hajimehoshi/ebiten/v2"

// HandleResize is called during Layout to handle window resizing
func (g *Game) HandleResize(outsideWidth, outsideHeight int) {
	// Force redraw when window size changes
	screenWidth, screenHeight := ebiten.WindowSize()
	if g.graphCanvas != nil && (g.graphCanvas.Bounds().Dx() != screenWidth || g.graphCanvas.Bounds().Dy() != screenHeight) {
		g.canvasNeedsRedraw = true
	}
}
