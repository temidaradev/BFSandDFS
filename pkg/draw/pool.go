package draw

import (
	"image/color"
	"sync"
)

// ObjectPool manages reusable objects to reduce GC pressure
type ObjectPool struct {
	circles    sync.Pool
	lines      sync.Pool
	rectangles sync.Pool
}

// CircleData represents cached circle drawing data
type CircleData struct {
	X, Y, Radius int
	Color        color.Color
	Filled       bool
}

// LineData represents cached line drawing data
type LineData struct {
	X1, Y1, X2, Y2 float64
	Color          color.Color
	Width          float64
}

// RectData represents cached rectangle drawing data
type RectData struct {
	X, Y, Width, Height int
	Color               color.Color
	Filled              bool
}

// NewObjectPool creates a new object pool
func NewObjectPool() *ObjectPool {
	return &ObjectPool{
		circles: sync.Pool{
			New: func() interface{} {
				return &CircleData{}
			},
		},
		lines: sync.Pool{
			New: func() interface{} {
				return &LineData{}
			},
		},
		rectangles: sync.Pool{
			New: func() interface{} {
				return &RectData{}
			},
		},
	}
}

// GetCircle gets a circle object from the pool
func (p *ObjectPool) GetCircle() *CircleData {
	return p.circles.Get().(*CircleData)
}

// PutCircle returns a circle object to the pool
func (p *ObjectPool) PutCircle(c *CircleData) {
	c.X, c.Y, c.Radius = 0, 0, 0
	c.Color = nil
	c.Filled = false
	p.circles.Put(c)
}

// GetLine gets a line object from the pool
func (p *ObjectPool) GetLine() *LineData {
	return p.lines.Get().(*LineData)
}

// PutLine returns a line object to the pool
func (p *ObjectPool) PutLine(l *LineData) {
	l.X1, l.Y1, l.X2, l.Y2 = 0, 0, 0, 0
	l.Color = nil
	l.Width = 0
	p.lines.Put(l)
}

// GetRect gets a rectangle object from the pool
func (p *ObjectPool) GetRect() *RectData {
	return p.rectangles.Get().(*RectData)
}

// PutRect returns a rectangle object to the pool
func (p *ObjectPool) PutRect(r *RectData) {
	r.X, r.Y, r.Width, r.Height = 0, 0, 0, 0
	r.Color = nil
	r.Filled = false
	p.rectangles.Put(r)
}

// Global object pool instance
var globalPool = NewObjectPool()
