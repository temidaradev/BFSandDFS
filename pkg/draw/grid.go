package draw

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// GridConfig defines the appearance and behavior of a grid
type GridConfig struct {
	CellSize        int
	MajorLineEvery  int
	MinorColor      color.RGBA
	MajorColor      color.RGBA
	ShowCoordinates bool
}

// DefaultGridConfig returns a default grid configuration
func DefaultGridConfig() GridConfig {
	return GridConfig{
		CellSize:        20,
		MajorLineEvery:  5,
		MinorColor:      color.RGBA{220, 220, 220, 255},
		MajorColor:      color.RGBA{180, 180, 180, 255},
		ShowCoordinates: false,
	}
}

// DrawGrid renders a grid on the screen
func DrawGrid(screen *ebiten.Image, width, height int, config GridConfig) {
	// Create a single pixel image for lines
	lineImg := ebiten.NewImage(1, 1)

	// Draw horizontal grid lines
	for y := 0; y < height; y += config.CellSize {
		// Choose line color
		lineColor := config.MinorColor
		if y%(config.CellSize*config.MajorLineEvery) == 0 {
			lineColor = config.MajorColor
		}

		// Fill the pixel with the line color
		lineImg.Fill(lineColor)

		// Create transform options
		opts := &ebiten.DrawImageOptions{}

		// Scale to match screen width
		opts.GeoM.Scale(float64(width), 1)

		// Position the line
		opts.GeoM.Translate(0, float64(y))

		// Draw the line
		screen.DrawImage(lineImg, opts)
	}

	// Draw vertical grid lines
	for x := 0; x < width; x += config.CellSize {
		// Choose line color
		lineColor := config.MinorColor
		if x%(config.CellSize*config.MajorLineEvery) == 0 {
			lineColor = config.MajorColor
		}

		// Fill the pixel with the line color
		lineImg.Fill(lineColor)

		// Create transform options
		opts := &ebiten.DrawImageOptions{}

		// Scale to match screen height
		opts.GeoM.Scale(1, float64(height))

		// Position the line
		opts.GeoM.Translate(float64(x), 0)

		// Draw the line
		screen.DrawImage(lineImg, opts)
	}
}

// DrawOptimizedGrid renders a grid on the screen with optimizations
func DrawOptimizedGrid(screen *ebiten.Image, width, height int, config GridConfig) {
	// Create cached line images for minor and major lines
	minorLineImg := ebiten.NewImage(1, 1)
	minorLineImg.Fill(config.MinorColor)
	majorLineImg := ebiten.NewImage(1, 1)
	majorLineImg.Fill(config.MajorColor)

	// Draw horizontal grid lines
	for y := 0; y < height; y += config.CellSize {
		// Choose line image based on whether it's a major line
		lineImg := minorLineImg
		if y%(config.CellSize*config.MajorLineEvery) == 0 {
			lineImg = majorLineImg
		}

		// Create transform options
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(width), 1)
		opts.GeoM.Translate(0, float64(y))

		// Draw the line
		screen.DrawImage(lineImg, opts)
	}

	// Draw vertical grid lines
	for x := 0; x < width; x += config.CellSize {
		// Choose line image based on whether it's a major line
		lineImg := minorLineImg
		if x%(config.CellSize*config.MajorLineEvery) == 0 {
			lineImg = majorLineImg
		}

		// Create transform options
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(1, float64(height))
		opts.GeoM.Translate(float64(x), 0)

		// Draw the line
		screen.DrawImage(lineImg, opts)
	}
}

// SnapToGrid aligns coordinates to the nearest grid intersection
func SnapToGrid(x, y int, cellSize int) (int, int) {
	return int(math.Round(float64(x)/float64(cellSize))) * cellSize,
		int(math.Round(float64(y)/float64(cellSize))) * cellSize
}
