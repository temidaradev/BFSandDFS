package draw

import (
	"image/color"
	"math"
	
	"github.com/hajimehoshi/ebiten/v2"
)

// GridConfig defines the appearance and behavior of a grid
type GridConfig struct {
	CellSize       int
	MajorLineEvery int
	MinorColor     color.RGBA
	MajorColor     color.RGBA
	ShowCoordinates bool
}

// DefaultGridConfig returns a default grid configuration
func DefaultGridConfig() GridConfig {
	return GridConfig{
		CellSize:       20,
		MajorLineEvery: 5,
		MinorColor:     color.RGBA{220, 220, 220, 255},
		MajorColor:     color.RGBA{180, 180, 180, 255},
		ShowCoordinates: false,
	}
}

// DrawGrid renders a grid on the screen
func DrawGrid(screen *ebiten.Image, width, height int, config GridConfig) {
	// Draw horizontal grid lines
	for y := 0; y < height; y += config.CellSize {
		lineColor := config.MinorColor
		if y%(config.CellSize*config.MajorLineEvery) == 0 {
			lineColor = config.MajorColor
		}
		
		for x := 0; x < width; x++ {
			screen.Set(x, y, lineColor)
		}
	}
	
	// Draw vertical grid lines
	for x := 0; x < width; x += config.CellSize {
		lineColor := config.MinorColor
		if x%(config.CellSize*config.MajorLineEvery) == 0 {
			lineColor = config.MajorColor
		}
		
		for y := 0; y < height; y++ {
			screen.Set(x, y, lineColor)
		}
	}
}

// SnapToGrid aligns coordinates to the nearest grid intersection
func SnapToGrid(x, y int, cellSize int) (int, int) {
	return int(math.Round(float64(x)/float64(cellSize)))*cellSize,
		int(math.Round(float64(y)/float64(cellSize)))*cellSize
}
