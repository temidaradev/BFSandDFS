package draw

import (
	"math"
)

// SpatialGrid provides efficient spatial partitioning for performance optimization
type SpatialGrid struct {
	cellSize float64
	grid     map[int]map[int][]int // grid[x][y] = list of node indices
	bounds   Bounds
}

type Bounds struct {
	MinX, MinY, MaxX, MaxY float64
}

type Point struct {
	X, Y float64
}

// NewSpatialGrid creates a new spatial grid for efficient spatial queries
func NewSpatialGrid(cellSize float64, bounds Bounds) *SpatialGrid {
	return &SpatialGrid{
		cellSize: cellSize,
		grid:     make(map[int]map[int][]int),
		bounds:   bounds,
	}
}

// Clear removes all objects from the spatial grid
func (sg *SpatialGrid) Clear() {
	sg.grid = make(map[int]map[int][]int)
}

// Insert adds a point object to the spatial grid
func (sg *SpatialGrid) Insert(x, y float64, id int) {
	cellX := int(math.Floor(x / sg.cellSize))
	cellY := int(math.Floor(y / sg.cellSize))

	if sg.grid[cellX] == nil {
		sg.grid[cellX] = make(map[int][]int)
	}

	sg.grid[cellX][cellY] = append(sg.grid[cellX][cellY], id)
}

// QueryRadius returns all objects within a radius of the given point
func (sg *SpatialGrid) QueryRadius(x, y, radius float64) []int {
	var result []int

	// Calculate the range of cells to check
	cellRadius := int(math.Ceil(radius / sg.cellSize))
	centerCellX := int(math.Floor(x / sg.cellSize))
	centerCellY := int(math.Floor(y / sg.cellSize))

	for dx := -cellRadius; dx <= cellRadius; dx++ {
		for dy := -cellRadius; dy <= cellRadius; dy++ {
			cellX := centerCellX + dx
			cellY := centerCellY + dy

			if sg.grid[cellX] != nil && sg.grid[cellX][cellY] != nil {
				for _, id := range sg.grid[cellX][cellY] {
					// Note: In practice, you'd need to store actual coordinates
					// This is a simplified version for demonstration
					result = append(result, id)
				}
			}
		}
	}

	return result
}

// QueryRect returns all objects within a rectangular region
func (sg *SpatialGrid) QueryRect(minX, minY, maxX, maxY float64) []int {
	var result []int

	startCellX := int(math.Floor(minX / sg.cellSize))
	startCellY := int(math.Floor(minY / sg.cellSize))
	endCellX := int(math.Floor(maxX / sg.cellSize))
	endCellY := int(math.Floor(maxY / sg.cellSize))

	for cellX := startCellX; cellX <= endCellX; cellX++ {
		for cellY := startCellY; cellY <= endCellY; cellY++ {
			if sg.grid[cellX] != nil && sg.grid[cellX][cellY] != nil {
				result = append(result, sg.grid[cellX][cellY]...)
			}
		}
	}

	return result
}

// Frustum culling for viewport optimization
type ViewportCuller struct {
	viewportBounds Bounds
	spatialGrid    *SpatialGrid
}

// NewViewportCuller creates a new viewport culler for efficient rendering
func NewViewportCuller(bounds Bounds) *ViewportCuller {
	return &ViewportCuller{
		viewportBounds: bounds,
		spatialGrid:    NewSpatialGrid(50.0, bounds), // 50 pixel cells
	}
}

// UpdateViewport updates the viewport bounds for culling
func (vc *ViewportCuller) UpdateViewport(minX, minY, maxX, maxY float64) {
	vc.viewportBounds = Bounds{minX, minY, maxX, maxY}
}

// GetVisibleObjects returns only objects visible in the current viewport
func (vc *ViewportCuller) GetVisibleObjects() []int {
	return vc.spatialGrid.QueryRect(
		vc.viewportBounds.MinX, vc.viewportBounds.MinY,
		vc.viewportBounds.MaxX, vc.viewportBounds.MaxY,
	)
}

// IsPointVisible checks if a point is visible in the current viewport
func (vc *ViewportCuller) IsPointVisible(x, y float64) bool {
	return x >= vc.viewportBounds.MinX && x <= vc.viewportBounds.MaxX &&
		y >= vc.viewportBounds.MinY && y <= vc.viewportBounds.MaxY
}

// IsRectVisible checks if a rectangle intersects with the viewport
func (vc *ViewportCuller) IsRectVisible(x, y, width, height float64) bool {
	return !(x+width < vc.viewportBounds.MinX || x > vc.viewportBounds.MaxX ||
		y+height < vc.viewportBounds.MinY || y > vc.viewportBounds.MaxY)
}

// IsCircleVisible checks if a circle intersects with the viewport
func (vc *ViewportCuller) IsCircleVisible(x, y float64, radius int) bool {
	// Simple bounding box check for circles
	return x+float64(radius) >= vc.viewportBounds.MinX && x-float64(radius) <= vc.viewportBounds.MaxX &&
		y+float64(radius) >= vc.viewportBounds.MinY && y-float64(radius) <= vc.viewportBounds.MaxY
}
