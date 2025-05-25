package draw

import (
	"image/color"
	"math"
	"sync"

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

// Performance-optimized grid drawing with caching
var (
	gridCache      = make(map[gridCacheKey]*ebiten.Image)
	gridCacheMutex sync.RWMutex
)

type gridCacheKey struct {
	width, height  int
	cellSize       int
	majorLineEvery int
	minorColor     color.RGBA
	majorColor     color.RGBA
}

// DrawGrid renders a grid on the screen with full caching
func DrawGrid(screen *ebiten.Image, width, height int, config GridConfig) {
	// Create cache key
	key := gridCacheKey{
		width:          width,
		height:         height,
		cellSize:       config.CellSize,
		majorLineEvery: config.MajorLineEvery,
		minorColor:     config.MinorColor,
		majorColor:     config.MajorColor,
	}

	// Check cache first
	gridCacheMutex.RLock()
	cachedGrid, exists := gridCache[key]
	gridCacheMutex.RUnlock()

	if exists {
		opts := GetDrawOptions()
		defer PutDrawOptions(opts)
		screen.DrawImage(cachedGrid, opts)
		return
	}

	// Create new grid image
	gridImg := ebiten.NewImage(width, height)
	gridImg.Fill(color.RGBA{240, 240, 240, 255}) // Background

	// Get cached line images
	minorLineImg := getOrCreateLineImage(config.MinorColor)
	majorLineImg := getOrCreateLineImage(config.MajorColor)

	// Draw horizontal grid lines
	for y := 0; y < height; y += config.CellSize {
		lineImg := minorLineImg
		if y%(config.CellSize*config.MajorLineEvery) == 0 {
			lineImg = majorLineImg
		}

		opts := GetDrawOptions()
		opts.GeoM.Scale(float64(width), 1)
		opts.GeoM.Translate(0, float64(y))
		gridImg.DrawImage(lineImg, opts)
		PutDrawOptions(opts)
	}

	// Draw vertical grid lines
	for x := 0; x < width; x += config.CellSize {
		lineImg := minorLineImg
		if x%(config.CellSize*config.MajorLineEvery) == 0 {
			lineImg = majorLineImg
		}

		opts := GetDrawOptions()
		opts.GeoM.Scale(1, float64(height))
		opts.GeoM.Translate(float64(x), 0)
		gridImg.DrawImage(lineImg, opts)
		PutDrawOptions(opts)
	}

	// Cache the grid with size limit
	gridCacheMutex.Lock()
	if len(gridCache) >= 50 { // Limit cache size
		// Simple eviction: remove oldest entry
		for k := range gridCache {
			delete(gridCache, k)
			break
		}
	}
	gridCache[key] = gridImg
	gridCacheMutex.Unlock()

	// Draw the grid
	opts := GetDrawOptions()
	defer PutDrawOptions(opts)
	screen.DrawImage(gridImg, opts)
}

// DrawOptimizedGrid renders a grid with viewport culling for massive performance gains
func DrawOptimizedGrid(screen *ebiten.Image, width, height int, config GridConfig) {
	// Only draw visible grid lines based on viewport
	screenBounds := screen.Bounds()

	minorLineImg := getOrCreateLineImage(config.MinorColor)
	majorLineImg := getOrCreateLineImage(config.MajorColor)

	// Calculate visible line ranges with proper bounds checking
	startX := max(0, (screenBounds.Min.X/config.CellSize)*config.CellSize)
	endX := min(width, ((screenBounds.Max.X/config.CellSize + 1) * config.CellSize))
	startY := max(0, (screenBounds.Min.Y/config.CellSize)*config.CellSize)
	endY := min(height, ((screenBounds.Max.Y/config.CellSize + 1) * config.CellSize))

	// Draw only visible horizontal lines
	for y := startY; y < endY; y += config.CellSize {
		lineImg := minorLineImg
		if y%(config.CellSize*config.MajorLineEvery) == 0 {
			lineImg = majorLineImg
		}

		opts := GetDrawOptions()
		opts.GeoM.Scale(float64(width), 1)
		opts.GeoM.Translate(0, float64(y))
		screen.DrawImage(lineImg, opts)
		PutDrawOptions(opts)
	}

	// Draw only visible vertical lines
	for x := startX; x < endX; x += config.CellSize {
		lineImg := minorLineImg
		if x%(config.CellSize*config.MajorLineEvery) == 0 {
			lineImg = majorLineImg
		}

		opts := GetDrawOptions()
		opts.GeoM.Scale(1, float64(height))
		opts.GeoM.Translate(float64(x), 0)
		screen.DrawImage(lineImg, opts)
		PutDrawOptions(opts)
	}
}

// SnapToGrid aligns coordinates to the nearest grid intersection
func SnapToGrid(x, y int, cellSize int) (int, int) {
	return int(math.Round(float64(x)/float64(cellSize))) * cellSize,
		int(math.Round(float64(y)/float64(cellSize))) * cellSize
}

// ClearGridCache clears the grid cache
func ClearGridCache() {
	gridCacheMutex.Lock()
	gridCache = make(map[gridCacheKey]*ebiten.Image)
	gridCacheMutex.Unlock()
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
