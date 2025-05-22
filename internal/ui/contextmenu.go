package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ContextMenuItem represents a single menu option in a context menu
type ContextMenuItem struct {
	Label  string
	Action func()
	Hover  bool
}

// ContextMenu represents a right-click context menu
type ContextMenu struct {
	X, Y              int
	Width, ItemHeight int
	Items             []*ContextMenuItem
	Visible           bool
	TargetNode        int // Node index that was right-clicked, -1 if not on a node
}

// NewContextMenu creates a new context menu
func NewContextMenu() *ContextMenu {
	return &ContextMenu{
		Width:      150,
		ItemHeight: 25,
		Visible:    false,
		TargetNode: -1,
	}
}

// Show displays the context menu at the specified position
func (m *ContextMenu) Show(x, y int, targetNode int) {
	m.X = x
	m.Y = y
	m.Visible = true
	m.TargetNode = targetNode
}

// Hide hides the context menu
func (m *ContextMenu) Hide() {
	m.Visible = false
	m.TargetNode = -1
}

// AddItem adds a menu item to the context menu
func (m *ContextMenu) AddItem(label string, action func()) {
	m.Items = append(m.Items, &ContextMenuItem{
		Label:  label,
		Action: action,
		Hover:  false,
	})
}

// ClearItems removes all items from the context menu
func (m *ContextMenu) ClearItems() {
	m.Items = nil
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

	// Draw background with semi-transparent effect
	bg := ebiten.NewImage(m.Width, totalHeight)
	bg.Fill(color.RGBA{40, 40, 40, 230})

	// Add a subtle border
	borderColor := color.RGBA{80, 80, 80, 255}
	for i := 0; i < m.Width; i++ {
		bg.Set(i, 0, borderColor)             // Top
		bg.Set(i, totalHeight-1, borderColor) // Bottom
	}
	for i := 0; i < totalHeight; i++ {
		bg.Set(0, i, borderColor)         // Left
		bg.Set(m.Width-1, i, borderColor) // Right
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(m.X), float64(m.Y))
	screen.DrawImage(bg, opts)

	// Draw items
	for i, item := range m.Items {
		itemY := m.Y + i*m.ItemHeight

		// Draw separator lines between items
		if i > 0 {
			separator := ebiten.NewImage(m.Width-2, 1)
			separator.Fill(color.RGBA{60, 60, 60, 255})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(m.X+1), float64(itemY))
			screen.DrawImage(separator, opts)
		}

		// Draw hover highlight
		if item.Hover {
			hoverBg := ebiten.NewImage(m.Width-2, m.ItemHeight-1)
			hoverBg.Fill(color.RGBA{70, 90, 120, 255})
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(m.X+1), float64(itemY+1))
			screen.DrawImage(hoverBg, opts)
		}

		// Draw item text with shadow for better visibility
		textColor := color.RGBA{220, 220, 220, 255}
		shadowColor := color.RGBA{0, 0, 0, 100}
		text.Draw(screen, item.Label, basicfont.Face7x13, m.X+11, itemY+m.ItemHeight/2+6, shadowColor)
		text.Draw(screen, item.Label, basicfont.Face7x13, m.X+10, itemY+m.ItemHeight/2+5, textColor)
	}
}
