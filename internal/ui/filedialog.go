package ui

import (
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// FileDialog represents a file dialog for saving/loading graphs
type FileDialog struct {
	X, Y            int
	Width, Height   int
	Visible         bool
	IsSaveDialog    bool
	CurrentDir      string
	Files           []string
	SelectedFile    int
	HoveredFile     int
	FileName        string
	CursorPos       int
	SaveLabel       string
	ScrollOffset    int
	MaxVisibleFiles int

	// Animation and hover effects
	AnimationProgress float64
	OkButtonHover     bool
	CancelButtonHover bool
	LastCursorBlink   time.Time
	ShowCursor        bool
}

// NewFileDialog creates a new file dialog
func NewFileDialog(isSaveDialog bool) *FileDialog {
	// Create default save directory if it doesn't exist
	saveDir := filepath.Join("saves")
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		os.MkdirAll(saveDir, 0755)
	}

	// Get screen dimensions for centering
	w, h := ebiten.WindowSize()
	dialogWidth := 500
	dialogHeight := 450

	dialog := &FileDialog{
		X:                 (w - dialogWidth) / 2,  // Center horizontally
		Y:                 (h - dialogHeight) / 2, // Center vertically
		Width:             dialogWidth,
		Height:            dialogHeight,
		IsSaveDialog:      isSaveDialog,
		CurrentDir:        saveDir,
		Files:             []string{},
		SelectedFile:      -1,
		HoveredFile:       -1,
		MaxVisibleFiles:   12,
		AnimationProgress: 0,
		LastCursorBlink:   time.Now(),
		ShowCursor:        true,
	}

	if isSaveDialog {
		dialog.SaveLabel = "Save Graph"
		dialog.FileName = "graph.json"
	} else {
		dialog.SaveLabel = "Load Graph"
	}

	dialog.RefreshFiles()
	return dialog
}

// Show displays the file dialog
func (fd *FileDialog) Show() {
	fd.Visible = true
	fd.AnimationProgress = 0
	fd.RefreshFiles()

	// Recenter dialog based on current window size
	w, h := ebiten.WindowSize()
	fd.X = (w - fd.Width) / 2
	fd.Y = (h - fd.Height) / 2
}

// Hide hides the file dialog
func (fd *FileDialog) Hide() {
	fd.Visible = false
}

// RefreshFiles updates the list of files in the current directory
func (fd *FileDialog) RefreshFiles() {
	fd.Files = []string{}

	// Add parent directory option if not in root
	parentDir := filepath.Dir(fd.CurrentDir)
	if parentDir != fd.CurrentDir {
		fd.Files = append(fd.Files, "../")
	}

	// Read directory contents
	files, err := os.ReadDir(fd.CurrentDir)
	if err == nil {
		// First add directories
		for _, file := range files {
			if file.IsDir() {
				fd.Files = append(fd.Files, file.Name()+"/")
			}
		}

		// Then add JSON files
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
				fd.Files = append(fd.Files, file.Name())
			}
		}
	}

	// Reset scroll position when changing directories
	fd.ScrollOffset = 0
}

// Update updates the file dialog state (animations, cursor blinking)
func (fd *FileDialog) Update() {
	if !fd.Visible {
		return
	}

	// Update animation
	if fd.AnimationProgress < 1.0 {
		fd.AnimationProgress += 0.1
		if fd.AnimationProgress > 1.0 {
			fd.AnimationProgress = 1.0
		}
	}

	// Blink cursor every 500ms
	if time.Since(fd.LastCursorBlink).Milliseconds() > 500 {
		fd.ShowCursor = !fd.ShowCursor
		fd.LastCursorBlink = time.Now()
	}

	// Update button hover states
	mx, my := ebiten.CursorPosition()

	// OK button hover
	fd.OkButtonHover = mx >= fd.X+fd.Width-180 && mx <= fd.X+fd.Width-100 &&
		my >= fd.Y+fd.Height-30 && my <= fd.Y+fd.Height

	// Cancel button hover
	fd.CancelButtonHover = mx >= fd.X+fd.Width-90 && mx <= fd.X+fd.Width-10 &&
		my >= fd.Y+fd.Height-30 && my <= fd.Y+fd.Height

	// Update hovered file
	fileListY := fd.Y + 60
	fileHeight := 24

	if mx >= fd.X+10 && mx <= fd.X+fd.Width-10 &&
		my >= fileListY && my <= fileListY+fd.MaxVisibleFiles*fileHeight {
		hoveredIndex := fd.ScrollOffset + (my-fileListY)/fileHeight
		if hoveredIndex >= 0 && hoveredIndex < len(fd.Files) {
			fd.HoveredFile = hoveredIndex
		} else {
			fd.HoveredFile = -1
		}
	} else {
		fd.HoveredFile = -1
	}
}

// Draw renders the file dialog
func (fd *FileDialog) Draw(screen *ebiten.Image) {
	if !fd.Visible {
		return
	}

	// Apply animation scaling effect
	scale := fd.AnimationProgress
	scaledWidth := int(float64(fd.Width) * scale)
	scaledHeight := int(float64(fd.Height) * scale)

	// Center the scaled dialog
	scaledX := fd.X + (fd.Width-scaledWidth)/2
	scaledY := fd.Y + (fd.Height-scaledHeight)/2

	// Draw semi-transparent overlay
	overlayColor := color.RGBA{0, 0, 0, 180}
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(overlayColor)

	// Only draw overlay if animation is far enough along
	if fd.AnimationProgress > 0.5 {
		opacity := (fd.AnimationProgress - 0.5) * 2.0
		if opacity > 1.0 {
			opacity = 1.0
		}

		opts := &ebiten.DrawImageOptions{}
		opts.ColorM.Scale(1, 1, 1, opacity)
		screen.DrawImage(overlay, opts)
	}

	// Skip drawing the rest if animation is just starting
	if fd.AnimationProgress < 0.2 {
		return
	}

	// Draw dialog background with modern effect
	bg := ebiten.NewImage(scaledWidth, scaledHeight)
	bg.Fill(color.RGBA{40, 42, 54, 245})

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(scaledX), float64(scaledY))
	screen.DrawImage(bg, opts)

	// Only proceed with detailed drawing if animation is far enough along
	if fd.AnimationProgress < 0.8 {
		return
	}

	// Draw modern border with gradient effect
	borderColor := color.RGBA{80, 100, 140, 255}
	accentColor := color.RGBA{130, 160, 190, 255}

	// Top border with accent
	for i := 0; i < fd.Width; i++ {
		intensity := uint8(float64(i) / float64(fd.Width) * 70)
		screen.Set(fd.X+i, fd.Y, color.RGBA{accentColor.R, accentColor.G - intensity, accentColor.B, accentColor.A})
		screen.Set(fd.X+i, fd.Y+1, color.RGBA{accentColor.R, accentColor.G - intensity, accentColor.B, accentColor.A})
	}

	// Side and bottom borders
	for i := 2; i < fd.Height; i++ {
		screen.Set(fd.X, fd.Y+i, borderColor)
		screen.Set(fd.X+1, fd.Y+i, borderColor)
		screen.Set(fd.X+fd.Width-1, fd.Y+i, borderColor)
		screen.Set(fd.X+fd.Width-2, fd.Y+i, borderColor)
	}

	for i := 0; i < fd.Width; i++ {
		screen.Set(fd.X+i, fd.Y+fd.Height-1, borderColor)
		screen.Set(fd.X+i, fd.Y+fd.Height-2, borderColor)
	}

	// Draw title with shadow and gradient
	titleColor := color.RGBA{220, 230, 255, 255}
	shadowColor := color.RGBA{0, 0, 0, 100}
	text.Draw(screen, fd.SaveLabel, basicfont.Face7x13, fd.X+fd.Width/2-text.BoundString(basicfont.Face7x13, fd.SaveLabel).Dx()/2+1, fd.Y+21, shadowColor)
	text.Draw(screen, fd.SaveLabel, basicfont.Face7x13, fd.X+fd.Width/2-text.BoundString(basicfont.Face7x13, fd.SaveLabel).Dx()/2, fd.Y+20, titleColor)

	// Draw current directory with better styling
	dirText := fd.CurrentDir
	// Truncate with ellipsis if too long
	if len(dirText) > (fd.Width-40)/7 { // Rough estimate for character width
		dirText = "..." + dirText[len(dirText)-(fd.Width-40)/7+3:]
	}

	dirColor := color.RGBA{180, 200, 255, 255}
	dirBg := ebiten.NewImage(fd.Width-20, 20)
	dirBg.Fill(color.RGBA{30, 35, 45, 255})

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+10), float64(fd.Y+35))
	screen.DrawImage(dirBg, opts)

	text.Draw(screen, dirText, basicfont.Face7x13, fd.X+15, fd.Y+48, dirColor)

	// Draw separator
	separator := ebiten.NewImage(fd.Width-20, 2)
	separator.Fill(color.RGBA{70, 80, 100, 255})

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+10), float64(fd.Y+55))
	screen.DrawImage(separator, opts)

	// Draw file list with improved styling
	fileListY := fd.Y + 60
	fileHeight := 24
	fileListBg := ebiten.NewImage(fd.Width-20, fd.MaxVisibleFiles*fileHeight)
	fileListBg.Fill(color.RGBA{25, 28, 38, 255})

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+10), float64(fileListY))
	screen.DrawImage(fileListBg, opts)

	// Draw scrollbar background
	scrollbarBg := ebiten.NewImage(8, fd.MaxVisibleFiles*fileHeight)
	scrollbarBg.Fill(color.RGBA{40, 45, 60, 255})

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+fd.Width-18), float64(fileListY))
	screen.DrawImage(scrollbarBg, opts)

	// Draw scrollbar handle if needed
	if len(fd.Files) > fd.MaxVisibleFiles {
		scrollbarHeight := float64(fd.MaxVisibleFiles*fileHeight) * float64(fd.MaxVisibleFiles) / float64(len(fd.Files))
		scrollbarY := float64(fileListY) + float64(fd.MaxVisibleFiles*fileHeight)*float64(fd.ScrollOffset)/float64(len(fd.Files))

		scrollbar := ebiten.NewImage(8, int(scrollbarHeight))
		scrollbar.Fill(color.RGBA{100, 120, 160, 255})

		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(fd.X+fd.Width-18), scrollbarY)
		screen.DrawImage(scrollbar, opts)
	}

	// Draw files with improved styling
	endIdx := fd.ScrollOffset + fd.MaxVisibleFiles
	if endIdx > len(fd.Files) {
		endIdx = len(fd.Files)
	}

	for i := fd.ScrollOffset; i < endIdx; i++ {
		y := fileListY + (i-fd.ScrollOffset)*fileHeight

		// Draw selection or hover highlight with animation effects
		if i == fd.SelectedFile {
			selectionBg := ebiten.NewImage(fd.Width-28, fileHeight)

			// Gradient for selected item
			for j := 0; j < fileHeight; j++ {
				progress := float64(j) / float64(fileHeight)
				r := uint8(80 - progress*20)
				g := uint8(110 - progress*20)
				b := uint8(180 - progress*20)

				for k := 0; k < fd.Width-28; k++ {
					selectionBg.Set(k, j, color.RGBA{r, g, b, 255})
				}
			}

			// Add highlight to left edge for emphasis
			for j := 0; j < fileHeight; j++ {
				for k := 0; k < 3; k++ {
					selectionBg.Set(k, j, color.RGBA{130, 170, 255, 255})
				}
			}

			opts = &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(fd.X+10), float64(y))
			screen.DrawImage(selectionBg, opts)
		} else if i == fd.HoveredFile {
			hoverBg := ebiten.NewImage(fd.Width-28, fileHeight)
			hoverBg.Fill(color.RGBA{60, 70, 110, 255})

			// Animate hover effect
			now := time.Now()
			pulseIntensity := float64(now.UnixNano()%1000000000) / 1000000000.0
			pulseIntensity = (math.Sin(pulseIntensity*math.Pi*2) + 1) / 2 * 20

			// Add subtle highlight to left edge
			for j := 0; j < fileHeight; j++ {
				for k := 0; k < 2; k++ {
					hoverBg.Set(k, j, color.RGBA{80 + uint8(pulseIntensity), 100 + uint8(pulseIntensity), 150 + uint8(pulseIntensity), 255})
				}
			}

			opts = &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(fd.X+10), float64(y))
			screen.DrawImage(hoverBg, opts)
		}

		// Draw file icon based on type
		fileName := fd.Files[i]
		var fileIcon string
		var fileColor color.RGBA

		if strings.HasSuffix(fileName, "/") {
			if fileName == "../" {
				fileIcon = "↑ "
				fileColor = color.RGBA{180, 210, 255, 255}
			} else {
				fileIcon = "📁 "
				fileColor = color.RGBA{180, 210, 255, 255}
			}
		} else if strings.HasSuffix(fileName, ".json") {
			fileIcon = "📄 "
			fileColor = color.RGBA{230, 230, 240, 255}
		} else {
			fileIcon = "📄 "
			fileColor = color.RGBA{200, 200, 210, 255}
		}

		// Draw file name with icon and shadow for better readability
		displayName := fileIcon + fileName
		text.Draw(screen, displayName, basicfont.Face7x13, fd.X+16, y+fileHeight/2+5, color.RGBA{20, 20, 30, 100})
		text.Draw(screen, displayName, basicfont.Face7x13, fd.X+15, y+fileHeight/2+4, fileColor)
	}

	// Draw bottom separator
	separator = ebiten.NewImage(fd.Width-20, 2)
	separator.Fill(color.RGBA{70, 80, 100, 255})

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+10), float64(fd.Y+fd.Height-80))
	screen.DrawImage(separator, opts)

	// Draw filename input for save dialog with improved styling
	if fd.IsSaveDialog {
		text.Draw(screen, "Filename:", basicfont.Face7x13, fd.X+15, fd.Y+fd.Height-59, color.RGBA{200, 210, 255, 255})

		// Input field background with better styling
		inputBg := ebiten.NewImage(fd.Width-30, 28)
		inputBg.Fill(color.RGBA{25, 28, 38, 255})

		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(fd.X+15), float64(fd.Y+fd.Height-55))
		screen.DrawImage(inputBg, opts)

		// Draw input field border with focus effect
		inputBorderColor := color.RGBA{80, 100, 140, 255}
		for i := 0; i < fd.Width-30; i++ {
			screen.Set(fd.X+15+i, fd.Y+fd.Height-55, inputBorderColor)
			screen.Set(fd.X+15+i, fd.Y+fd.Height-28, inputBorderColor)
		}
		for i := 0; i < 28; i++ {
			screen.Set(fd.X+15, fd.Y+fd.Height-55+i, inputBorderColor)
			screen.Set(fd.X+15, fd.Y+fd.Height-55+i, inputBorderColor)
		}

		// Draw filename text
		text.Draw(screen, fd.FileName, basicfont.Face7x13, fd.X+20, fd.Y+fd.Height-38, color.RGBA{220, 230, 250, 255})

		// Draw blinking cursor with improved visibility
		if fd.ShowCursor {
			cursorPos := text.BoundString(basicfont.Face7x13, fd.FileName[:fd.CursorPos]).Dx()
			cursorHeight := 18
			cursorColor := color.RGBA{220, 230, 250, 255}
			for i := 0; i < cursorHeight; i++ {
				screen.Set(fd.X+20+cursorPos, fd.Y+fd.Height-48+i, cursorColor)
			}
		}
	}

	// Draw OK and Cancel buttons with improved styling and hover effects
	okBtnColor := color.RGBA{60, 130, 180, 255}
	if fd.OkButtonHover {
		okBtnColor = color.RGBA{80, 150, 200, 255}
	}

	cancelBtnColor := color.RGBA{180, 70, 70, 255}
	if fd.CancelButtonHover {
		cancelBtnColor = color.RGBA{200, 90, 90, 255}
	}

	// Draw button backgrounds
	okBtnBg := ebiten.NewImage(80, 30)
	okBtnBg.Fill(okBtnColor)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+fd.Width-180), float64(fd.Y+fd.Height-30))
	screen.DrawImage(okBtnBg, opts)

	cancelBtnBg := ebiten.NewImage(80, 30)
	cancelBtnBg.Fill(cancelBtnColor)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+fd.Width-90), float64(fd.Y+fd.Height-30))
	screen.DrawImage(cancelBtnBg, opts)

	// Draw button text centered
	okText := "OK"
	okWidth := text.BoundString(basicfont.Face7x13, okText).Dx()
	okX := fd.X + fd.Width - 180 + (80-okWidth)/2

	cancelText := "Cancel"
	cancelWidth := text.BoundString(basicfont.Face7x13, cancelText).Dx()
	cancelX := fd.X + fd.Width - 90 + (80-cancelWidth)/2

	// Draw button text with shadows
	text.Draw(screen, okText, basicfont.Face7x13, okX+1, fd.Y+fd.Height-11, shadowColor)
	text.Draw(screen, okText, basicfont.Face7x13, okX, fd.Y+fd.Height-12, color.White)

	text.Draw(screen, cancelText, basicfont.Face7x13, cancelX+1, fd.Y+fd.Height-11, shadowColor)
	text.Draw(screen, cancelText, basicfont.Face7x13, cancelX, fd.Y+fd.Height-12, color.White)
}

// HandleClick processes clicks within the file dialog
// Returns true if the dialog handled the click, false otherwise
func (fd *FileDialog) HandleClick(x, y int) bool {
	if !fd.Visible {
		return false
	}

	// Check if click is outside the dialog
	if x < fd.X || x > fd.X+fd.Width || y < fd.Y || y > fd.Y+fd.Height {
		return false
	}

	// File list area
	fileListY := fd.Y + 60
	fileHeight := 24
	fileListHeight := fd.MaxVisibleFiles * fileHeight

	if x >= fd.X+10 && x <= fd.X+fd.Width-20 &&
		y >= fileListY && y <= fileListY+fileListHeight {
		// Clicked on file list
		clickedIndex := fd.ScrollOffset + (y-fileListY)/fileHeight
		if clickedIndex >= 0 && clickedIndex < len(fd.Files) {
			// If clicking on a directory
			if strings.HasSuffix(fd.Files[clickedIndex], "/") {
				if fd.Files[clickedIndex] == "../" {
					// Go up one directory
					fd.CurrentDir = filepath.Dir(fd.CurrentDir)
				} else {
					// Enter subdirectory
					fd.CurrentDir = filepath.Join(fd.CurrentDir, fd.Files[clickedIndex][:len(fd.Files[clickedIndex])-1])
				}
				fd.RefreshFiles()
				fd.SelectedFile = -1
				return true
			}

			fd.SelectedFile = clickedIndex
			if !fd.IsSaveDialog && fd.SelectedFile >= 0 && fd.SelectedFile < len(fd.Files) {
				fd.FileName = fd.Files[fd.SelectedFile]
			}
			return true
		}
	}

	// Check for scrollbar clicks
	scrollbarX := fd.X + fd.Width - 18
	if x >= scrollbarX && x <= scrollbarX+8 &&
		y >= fileListY && y <= fileListY+fd.MaxVisibleFiles*fileHeight {

		// Calculate relative position in scrollbar
		relativeY := y - fileListY
		relativePos := float64(relativeY) / float64(fd.MaxVisibleFiles*fileHeight)

		// Set scroll position
		if len(fd.Files) > fd.MaxVisibleFiles {
			fd.ScrollOffset = int(relativePos * float64(len(fd.Files)-fd.MaxVisibleFiles))
			if fd.ScrollOffset < 0 {
				fd.ScrollOffset = 0
			}
			if fd.ScrollOffset > len(fd.Files)-fd.MaxVisibleFiles {
				fd.ScrollOffset = len(fd.Files) - fd.MaxVisibleFiles
			}
		}
		return true
	}

	// OK button
	if fd.OkButtonHover {
		// OK button was clicked
		return true
	}

	// Cancel button
	if fd.CancelButtonHover {
		fd.Hide()
		return true
	}

	return true
}

// GetSelectedFilePath returns the full path to the selected file
func (fd *FileDialog) GetSelectedFilePath() string {
	// For load dialog, get the selected file
	if !fd.IsSaveDialog && fd.SelectedFile >= 0 && fd.SelectedFile < len(fd.Files) {
		return filepath.Join(fd.CurrentDir, fd.Files[fd.SelectedFile])
	}

	// For save dialog, use the entered filename
	if fd.IsSaveDialog {
		filename := fd.FileName
		if !strings.HasSuffix(filename, ".json") {
			filename += ".json"
		}
		return filepath.Join(fd.CurrentDir, filename)
	}

	return ""
}

// TypeCharacter adds a character to the filename
func (fd *FileDialog) TypeCharacter(ch rune) {
	if !fd.IsSaveDialog {
		return
	}

	if fd.CursorPos < len(fd.FileName) {
		fd.FileName = fd.FileName[:fd.CursorPos] + string(ch) + fd.FileName[fd.CursorPos:]
	} else {
		fd.FileName += string(ch)
	}
	fd.CursorPos++
}

// DeleteCharacter deletes a character from the filename
func (fd *FileDialog) DeleteCharacter() {
	if !fd.IsSaveDialog || len(fd.FileName) == 0 || fd.CursorPos == 0 {
		return
	}

	fd.FileName = fd.FileName[:fd.CursorPos-1] + fd.FileName[fd.CursorPos:]
	fd.CursorPos--
}

// MoveCursor moves the cursor position
func (fd *FileDialog) MoveCursor(offset int) {
	fd.CursorPos += offset
	if fd.CursorPos < 0 {
		fd.CursorPos = 0
	}
	if fd.CursorPos > len(fd.FileName) {
		fd.CursorPos = len(fd.FileName)
	}
}

// ScrollFiles scrolls the file list up or down
func (fd *FileDialog) ScrollFiles(amount int) {
	fd.ScrollOffset += amount
	if fd.ScrollOffset < 0 {
		fd.ScrollOffset = 0
	}
	if fd.ScrollOffset > len(fd.Files)-fd.MaxVisibleFiles {
		if len(fd.Files) > fd.MaxVisibleFiles {
			fd.ScrollOffset = len(fd.Files) - fd.MaxVisibleFiles
		} else {
			fd.ScrollOffset = 0
		}
	}
}
