// Package draw provides optimized drawing utilities for the application
package draw

import (
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Performance optimizations
var (
	// Enhanced caching with LRU eviction
	lineImageCache   = make(map[color.RGBA]*ebiten.Image)
	circleImageCache = make(map[circleKey]*ebiten.Image)
	rectImageCache   = make(map[rectKey]*ebiten.Image)
	cacheMutex       sync.RWMutex

	// Object pools to reduce GC pressure
	drawOptionsPool = sync.Pool{
		New: func() interface{} {
			return &ebiten.DrawImageOptions{}
		},
	}

	// Cache hit counters for performance monitoring
	cacheHits   uint64
	cacheMisses uint64

	// Cache size limits
	maxCacheSize = 1000
)

type circleKey struct {
	radius     int
	r, g, b, a uint8
}

type rectKey struct {
	width, height int
	r, g, b, a    uint8
}

// GetDrawOptions gets a reusable DrawImageOptions from the pool
func GetDrawOptions() *ebiten.DrawImageOptions {
	opts := drawOptionsPool.Get().(*ebiten.DrawImageOptions)
	opts.GeoM.Reset()
	opts.ColorM.Reset()
	return opts
}

// PutDrawOptions returns a DrawImageOptions to the pool for reuse
func PutDrawOptions(opts *ebiten.DrawImageOptions) {
	drawOptionsPool.Put(opts)
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

	// Get pooled options for better performance
	opts := GetDrawOptions()
	defer PutDrawOptions(opts)

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

	// Get pooled options for better performance
	opts := GetDrawOptions()
	defer PutDrawOptions(opts)

	// Position the circle
	opts.GeoM.Translate(float64(cx-r), float64(cy-r))

	// Draw the circle
	img.DrawImage(circleImg, opts)
}

// DrawCachedRect draws a filled rectangle with optimized caching
func DrawCachedRect(img *ebiten.Image, x, y, width, height int, clr color.Color) {
	// Convert color to RGBA for cache key
	r, g, b, a := clr.RGBA()
	rgba := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}

	// Get or create the cached rectangle image
	rectImg := getOrCreateRectImage(width, height, rgba)

	// Get pooled options for better performance
	opts := GetDrawOptions()
	defer PutDrawOptions(opts)

	// Position the rectangle
	opts.GeoM.Translate(float64(x), float64(y))

	// Draw the rectangle
	img.DrawImage(rectImg, opts)
}

// BatchDrawCircles draws multiple circles in one operation for better performance
func BatchDrawCircles(img *ebiten.Image, circles []struct {
	X, Y, R int
	Color   color.Color
}) {
	// Group circles by radius and color for optimal caching
	circleGroups := make(map[circleKey][]struct{ X, Y int })

	for _, circle := range circles {
		r8, g, b, a := circle.Color.RGBA()
		rgba := color.RGBA{uint8(r8 >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
		key := circleKey{circle.R, rgba.R, rgba.G, rgba.B, rgba.A}

		circleGroups[key] = append(circleGroups[key], struct{ X, Y int }{circle.X, circle.Y})
	}

	// Draw each group using the same cached image
	for key, positions := range circleGroups {
		circleImg := getOrCreateCircleImageFromKey(key)

		for _, pos := range positions {
			opts := GetDrawOptions()
			opts.GeoM.Translate(float64(pos.X-key.radius), float64(pos.Y-key.radius))
			img.DrawImage(circleImg, opts)
			PutDrawOptions(opts)
		}
	}
}

// getOrCreateLineImage retrieves a line image from cache or creates a new one
func getOrCreateLineImage(clr color.RGBA) *ebiten.Image {
	cacheMutex.RLock()
	lineImg, exists := lineImageCache[clr]
	cacheMutex.RUnlock()

	if exists {
		cacheHits++
		return lineImg
	}

	cacheMisses++

	// Create a new 1x1 pixel line image
	lineImg = ebiten.NewImage(1, 1)
	lineImg.Fill(clr)

	// Cache the new image with size limit check
	cacheMutex.Lock()
	if len(lineImageCache) >= maxCacheSize {
		// Simple eviction: remove a random element
		for k := range lineImageCache {
			delete(lineImageCache, k)
			break
		}
	}
	lineImageCache[clr] = lineImg
	cacheMutex.Unlock()

	return lineImg
}

// getOrCreateCircleImage retrieves a circle image from cache or creates a new one
func getOrCreateCircleImage(radius int, clr color.RGBA) *ebiten.Image {
	key := circleKey{radius, clr.R, clr.G, clr.B, clr.A}
	return getOrCreateCircleImageFromKey(key)
}

// getOrCreateCircleImageFromKey retrieves a circle image from cache using a key
func getOrCreateCircleImageFromKey(key circleKey) *ebiten.Image {
	cacheMutex.RLock()
	circleImg, exists := circleImageCache[key]
	cacheMutex.RUnlock()

	if exists {
		cacheHits++
		return circleImg
	}

	cacheMisses++

	// Create a new circle image with optimized algorithm
	diameter := key.radius * 2
	circleImg = ebiten.NewImage(diameter, diameter)

	// Use more efficient circle drawing algorithm
	clr := color.RGBA{key.r, key.g, key.b, key.a}
	radiusSquared := float64(key.radius * key.radius)

	for y := 0; y < diameter; y++ {
		for x := 0; x < diameter; x++ {
			dx := float64(x - key.radius)
			dy := float64(y - key.radius)
			if dx*dx+dy*dy <= radiusSquared {
				circleImg.Set(x, y, clr)
			}
		}
	}

	// Cache the new image with size limit check
	cacheMutex.Lock()
	if len(circleImageCache) >= maxCacheSize {
		// Simple eviction: remove a random element
		for k := range circleImageCache {
			delete(circleImageCache, k)
			break
		}
	}
	circleImageCache[key] = circleImg
	cacheMutex.Unlock()

	return circleImg
}

// getOrCreateRectImage retrieves a rectangle image from cache or creates a new one
func getOrCreateRectImage(width, height int, clr color.RGBA) *ebiten.Image {
	key := rectKey{width, height, clr.R, clr.G, clr.B, clr.A}

	cacheMutex.RLock()
	rectImg, exists := rectImageCache[key]
	cacheMutex.RUnlock()

	if exists {
		cacheHits++
		return rectImg
	}

	cacheMisses++

	// Create a new rectangle image
	rectImg = ebiten.NewImage(width, height)
	rectImg.Fill(clr)

	// Cache the new image with size limit check
	cacheMutex.Lock()
	if len(rectImageCache) >= maxCacheSize {
		// Simple eviction: remove a random element
		for k := range rectImageCache {
			delete(rectImageCache, k)
			break
		}
	}
	rectImageCache[key] = rectImg
	cacheMutex.Unlock()

	return rectImg
}

// GetCacheStats returns cache performance statistics
func GetCacheStats() (hits, misses uint64, hitRatio float64) {
	total := cacheHits + cacheMisses
	if total == 0 {
		return 0, 0, 0
	}
	return cacheHits, cacheMisses, float64(cacheHits) / float64(total)
}

// ClearCaches clears the image caches to free memory
// Call this when changing themes or when exiting the application
func ClearCaches() {
	cacheMutex.Lock()
	lineImageCache = make(map[color.RGBA]*ebiten.Image)
	circleImageCache = make(map[circleKey]*ebiten.Image)
	rectImageCache = make(map[rectKey]*ebiten.Image)
	cacheHits = 0
	cacheMisses = 0
	cacheMutex.Unlock()
}
