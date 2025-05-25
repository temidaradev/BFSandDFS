package ui

import (
	"image/color"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ContextMenuItem represents a single menu option in a context menu
type ContextMenuItem struct {
	Label        string
	Action       func()
	Hover        bool
	HoverOpacity float64 // For smooth hover animation
	Icon         string  // Optional icon character
}

// ContextMenu represents a right-click context menu
type ContextMenu struct {
	X, Y              int
	Width, ItemHeight int
	Items             []*ContextMenuItem
	Visible           bool
	TargetNode        int // Node index that was right-clicked, -1 if not on a node

	// Animation properties
	AnimationProgress float64
	LastUpdate        time.Time
}

// NewContextMenu creates a new context menu
func NewContextMenu() *ContextMenu {
	return &ContextMenu{
		Width:             180,
		ItemHeight:        28,
		Visible:           false,
		TargetNode:        -1,
		AnimationProgress: 0,
		LastUpdate:        time.Now(),
	}
}

// Show displays the context menu at the specified position
func (m *ContextMenu) Show(x, y int, targetNode int) {
	// Adjust position if menu would go off screen
	w, h := ebiten.WindowSize()

	// Check right edge
	if x+m.Width > w {
		x = w - m.Width - 5
	}

	// Check bottom edge
	totalHeight := len(m.Items) * m.ItemHeight
	if y+totalHeight > h {
		y = h - totalHeight - 5
	}

	m.X = x
	m.Y = y
	m.Visible = true
	m.TargetNode = targetNode
	m.AnimationProgress = 0
	m.LastUpdate = time.Now()
}

// Hide hides the context menu
func (m *ContextMenu) Hide() {
	m.Visible = false
	m.TargetNode = -1
}

// AddItem adds a menu item to the context menu
func (m *ContextMenu) AddItem(label string, action func()) {
	// Automatically assign appropriate icons based on action name
	var icon string

	switch {
	case contains(label, "Start"):
		icon = "▶ "
	case contains(label, "Set"):
		icon = "⚑ "
	case contains(label, "Delete"):
		icon = "✕ "
	case contains(label, "Remove"):
		icon = "✕ "
	case contains(label, "Clear"):
		icon = "🗑 "
	case contains(label, "Add"):
		icon = "✚ "
	case contains(label, "Edge"):
		icon = "↔ "
	case contains(label, "Save"):
		icon = "💾 "
	case contains(label, "Load"):
		icon = "📂 "
	case contains(label, "Create"):
		icon = "🔄 "
	default:
		icon = "• "
	}

	m.Items = append(m.Items, &ContextMenuItem{
		Label:        label,
		Action:       action,
		Hover:        false,
		HoverOpacity: 0,
		Icon:         icon,
	})
}

// ClearItems removes all items from the context menu
func (m *ContextMenu) ClearItems() {
	m.Items = nil
}

// Update updates the menu animations
func (m *ContextMenu) Update() {
	if !m.Visible {
		return
	}

	// Update animation progress
	elapsed := time.Since(m.LastUpdate).Seconds()
	m.LastUpdate = time.Now()

	if m.AnimationProgress < 1.0 {
		m.AnimationProgress += elapsed * 5 // Full animation in 0.2 seconds
		if m.AnimationProgress > 1.0 {
			m.AnimationProgress = 1.0
		}
	}

	// Update hover animations
	for _, item := range m.Items {
		if item.Hover && item.HoverOpacity < 1.0 {
			item.HoverOpacity += elapsed * 8 // Hover animation in 0.125 seconds
			if item.HoverOpacity > 1.0 {
				item.HoverOpacity = 1.0
			}
		} else if !item.Hover && item.HoverOpacity > 0 {
			item.HoverOpacity -= elapsed * 8
			if item.HoverOpacity < 0 {
				item.HoverOpacity = 0
			}
		}
	}
}

// HandleClick processes a click within the context menu
// Returns true if the click was handled by the menu
func (m *ContextMenu) HandleClick(x, y int) bool {
	if !m.Visible {
		return false
	}

	// Check if click is within the menu bounds
	if x < m.X || x > m.X+m.Width {
		m.Hide()
		return false
	}

	for i, item := range m.Items {
		itemY := m.Y + i*m.ItemHeight
		if y >= itemY && y < itemY+m.ItemHeight {
			if item.Action != nil {
				item.Action()
			}
			m.Hide()
			return true
		}
	}

	// Click outside menu items but within X bounds
	m.Hide()
	return true
}

// UpdateHoverState updates which menu item the mouse is hovering over
func (m *ContextMenu) UpdateHoverState(x, y int) {
	if !m.Visible {
		return
	}

	for i, item := range m.Items {
		itemY := m.Y + i*m.ItemHeight
		item.Hover = y >= itemY && y < itemY+m.ItemHeight &&
			x >= m.X && x <= m.X+m.Width
	}
}

// Draw renders the context menu on the screen
func (m *ContextMenu) Draw(screen *ebiten.Image) {
	if !m.Visible {
		return
	}

	// Calculate total height
	totalHeight := len(m.Items) * m.ItemHeight

	// Apply animation scaling effect
	scale := m.AnimationProgress
	scaledWidth := int(float64(m.Width) * scale)
	scaledHeight := int(float64(totalHeight) * scale)

	// Skip drawing if animation just started
	if scaledWidth < 10 || scaledHeight < 10 {
		return
	}

	// Center the scaled menu around original position
	scaledX := m.X + (m.Width-scaledWidth)/2
	scaledY := m.Y + (totalHeight-scaledHeight)/2

	// Draw background with modern semi-transparent effect
	bg := ebiten.NewImage(scaledWidth, scaledHeight)
	bg.Fill(color.RGBA{40, 42, 54, 245})

	// Add modern gradient border
	borderTop := color.RGBA{100, 140, 200, 255}
	borderSides := color.RGBA{70, 90, 120, 255}

	// Top gradient
	for i := 0; i < scaledWidth; i++ {
		intensity := uint8(float64(i) / float64(scaledWidth) * 60)
		bg.Set(i, 0, color.RGBA{borderTop.R, borderTop.G - intensity, borderTop.B, borderTop.A})
		bg.Set(i, 1, color.RGBA{borderTop.R - 10, borderTop.G - 10 - intensity, borderTop.B - 10, borderTop.A})
	}

	// Side and bottom borders
	for i := 2; i < scaledHeight; i++ {
		bg.Set(0, i, borderSides)
		bg.Set(1, i, borderSides)
		bg.Set(scaledWidth-1, i, borderSides)
		bg.Set(scaledWidth-2, i, borderSides)
	}

	for i := 0; i < scaledWidth; i++ {
		bg.Set(i, scaledHeight-1, borderSides)
		bg.Set(i, scaledHeight-2, borderSides)
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(scaledX), float64(scaledY))
	screen.DrawImage(bg, opts)

	// Skip drawing items if animation isn't far enough along
	if m.AnimationProgress < 0.8 {
		return
	}

	// Draw items
	for i, item := range m.Items {
		itemY := m.Y + i*m.ItemHeight

		// Draw separator lines between items (softer, more modern look)
		if i > 0 {
			separator := ebiten.NewImage(m.Width-8, 1)
			separator.Fill(color.RGBA{70, 80, 100, 150})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(m.X+4), float64(itemY))
			screen.DrawImage(separator, opts)
		}

		// Draw hover highlight with smooth animation
		if item.HoverOpacity > 0 {
			hoverBg := ebiten.NewImage(m.Width-4, m.ItemHeight-1)
			alpha := uint8(item.HoverOpacity * 255)
			hoverBg.Fill(color.RGBA{70, 100, 160, alpha})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(m.X+2), float64(itemY+1))
			screen.DrawImage(hoverBg, opts)
		}

		// Draw item text with shadow for better visibility
		textColor := color.RGBA{220, 220, 240, 255}
		shadowColor := color.RGBA{0, 0, 0, 100}

		// Draw icon and label
		fullText := item.Icon + item.Label
		text.Draw(screen, fullText, basicfont.Face7x13, m.X+11, itemY+m.ItemHeight/2+5, shadowColor)
		text.Draw(screen, fullText, basicfont.Face7x13, m.X+10, itemY+m.ItemHeight/2+4, textColor)
	}
}

// Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
