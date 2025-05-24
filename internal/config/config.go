package config

// Config holds application configuration settings
type Config struct {
	// WindowWidth is the initial window width
	WindowWidth int
	// WindowHeight is the initial window height
	WindowHeight int
	// Title is the window title
	Title string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		WindowWidth:  1200,
		WindowHeight: 800,
		Title:        "BFS, DFS, and AVL Tree Simulator",
	}
}
