// Package draw provides optimized drawing utilities for the application
package draw

import (
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Cache for frequently used shapes to avoid recreating them
var (
	lineImageCache   = make(map[color.RGBA]*ebiten.Image)
	circleImageCache = make(map[circleKey]*ebiten.Image)
	cacheMutex       sync.RWMutex
)

type circleKey struct {
	radius int
	r, g, b, a uint8
}

// DrawCachedLine draws a line from (x0,y0) to (x1,y1) with the given color
// Uses an optimized vector drawing approach with caching for better performance
func DrawCachedLine(img *ebiten.Image, x0, y0, x1, y1 float64, clr color.Color) {
	// Convert color to RGBA for cache key
	r, g, b, a := clr.RGBA()
	rgba := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	
	// Get or create the cached 1x1 pixel line image
	lineImg := getOrCreateLineImage(rgba)
	
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

// DrawCachedCircle draws a filled circle with center (cx,cy) and radius r
// Uses cached circle images for better performance
func DrawCachedCircle(img *ebiten.Image, cx, cy, r int, clr color.Color) {
	// Convert color to RGBA for cache key
	r8, g, b, a := clr.RGBA()
	rgba := color.RGBA{uint8(r8 >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	
	// Get or create the cached circle image
	circleImg := getOrCreateCircleImage(r, rgba)
	
	// Create transform options
	opts := &ebiten.DrawImageOptions{}
	
	// Position the circle
	opts.GeoM.Translate(float64(cx-r), float64(cy-r))
	
	// Draw the circle
	img.DrawImage(circleImg, opts)
}

// getOrCreateLineImage retrieves a line image from cache or creates a new one
func getOrCreateLineImage(clr color.RGBA) *ebiten.Image {
	cacheMutex.RLock()
	lineImg, exists := lineImageCache[clr]
	cacheMutex.RUnlock()
	
	if exists {
		return lineImg
	}
	
	// Create a new 1x1 pixel line image
	lineImg = ebiten.NewImage(1, 1)
	lineImg.Fill(clr)
	
	// Cache the new image
	cacheMutex.Lock()
	lineImageCache[clr] = lineImg
	cacheMutex.Unlock()
	
	return lineImg
}

// getOrCreateCircleImage retrieves a circle image from cache or creates a new one
func getOrCreateCircleImage(radius int, clr color.RGBA) *ebiten.Image {
	key := circleKey{radius, clr.R, clr.G, clr.B, clr.A}
	
	cacheMutex.RLock()
	circleImg, exists := circleImageCache[key]
	cacheMutex.RUnlock()
	
	if exists {
		return circleImg
	}
	
	// Create a new circle image
	diameter := radius * 2
	circleImg = ebiten.NewImage(diameter, diameter)
	
	// Draw the circle into the image
	for y := 0; y < diameter; y++ {
		for x := 0; x < diameter; x++ {
			dx := float64(x - radius)
			dy := float64(y - radius)
			if dx*dx+dy*dy <= float64(radius*radius) {
				circleImg.Set(x, y, clr)
			}
		}
	}
	
	// Cache the new image
	cacheMutex.Lock()
	circleImageCache[key] = circleImg
	cacheMutex.Unlock()
	
	return circleImg
}

// ClearCaches clears the image caches to free memory
// Call this when changing themes or when exiting the application
func ClearCaches() {
	cacheMutex.Lock()
	lineImageCache = make(map[color.RGBA]*ebiten.Image)
	circleImageCache = make(map[circleKey]*ebiten.Image)
	cacheMutex.Unlock()
}
