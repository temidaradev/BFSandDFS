package config

// PerformanceTier defines the level of performance optimizations
// PerformanceTier type
// (0 = Low, 1 = Medium, 2 = High)
//
//go:generate stringer -type=PerformanceTier
type PerformanceTier int

const (
	PerformanceLow PerformanceTier = iota
	PerformanceMedium
	PerformanceHigh
)

// Config holds application configuration settings
type Config struct {
	// WindowWidth is the initial window width
	WindowWidth int
	// WindowHeight is the initial window height
	WindowHeight int
	// Title is the window title
	Title string
	// PerformanceTier sets the performance level (Low, Medium, High)
	PerformanceTier PerformanceTier
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		WindowWidth:     1200,
		WindowHeight:    800,
		Title:           "BFS, DFS, and AVL Tree Simulator",
		PerformanceTier: PerformanceMedium, // Default to Medium
	}
}
