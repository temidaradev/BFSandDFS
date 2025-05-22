package ui

import (
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// FileDialog represents a simple file dialog for saving/loading graphs
type FileDialog struct {
	X, Y            int
	Width, Height   int
	Visible         bool
	IsSaveDialog    bool
	CurrentDir      string
	Files           []string
	SelectedFile    int
	FileName        string
	CursorPos       int
	SaveLabel       string
	ScrollOffset    int
	MaxVisibleFiles int
}

// NewFileDialog creates a new file dialog
func NewFileDialog(isSaveDialog bool) *FileDialog {
	// Create default save directory if it doesn't exist
	saveDir := filepath.Join("saves")
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		os.MkdirAll(saveDir, 0755)
	}

	dialog := &FileDialog{
		X:               150,
		Y:               100,
		Width:           400,
		Height:          400,
		IsSaveDialog:    isSaveDialog,
		CurrentDir:      saveDir,
		Files:           []string{},
		SelectedFile:    -1,
		MaxVisibleFiles: 10,
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
	fd.RefreshFiles()
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
		fd.Files = append(fd.Files, "..")
	}

	// Read directory contents
	files, err := os.ReadDir(fd.CurrentDir)
	if err == nil {
		for _, file := range files {
			if file.IsDir() || (strings.HasSuffix(file.Name(), ".json") && !file.IsDir()) {
				name := file.Name()
				if file.IsDir() {
					name += "/"
				}
				fd.Files = append(fd.Files, name)
			}
		}
	}
}

// Draw renders the file dialog
func (fd *FileDialog) Draw(screen *ebiten.Image) {
	if !fd.Visible {
		return
	}

	// Draw dialog background with semi-transparent effect
	bg := ebiten.NewImage(fd.Width, fd.Height)
	bg.Fill(color.RGBA{40, 40, 40, 230})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X), float64(fd.Y))
	screen.DrawImage(bg, opts)

	// Draw border with subtle effect
	borderColor := color.RGBA{80, 80, 80, 255}
	for i := 0; i < fd.Width; i++ {
		screen.Set(fd.X+i, fd.Y, borderColor)
		screen.Set(fd.X+i, fd.Y+fd.Height-1, borderColor)
	}
	for i := 0; i < fd.Height; i++ {
		screen.Set(fd.X, fd.Y+i, borderColor)
		screen.Set(fd.X+fd.Width-1, fd.Y+i, borderColor)
	}

	// Draw title with shadow
	titleColor := color.RGBA{220, 220, 220, 255}
	shadowColor := color.RGBA{0, 0, 0, 100}
	text.Draw(screen, fd.SaveLabel, basicfont.Face7x13, fd.X+11, fd.Y+21, shadowColor)
	text.Draw(screen, fd.SaveLabel, basicfont.Face7x13, fd.X+10, fd.Y+20, titleColor)

	// Draw current directory with improved styling
	dirText := "Directory: " + fd.CurrentDir
	if len(dirText) > 50 {
		dirText = "..." + dirText[len(dirText)-47:]
	}
	dirColor := color.RGBA{180, 180, 255, 255}
	text.Draw(screen, dirText, basicfont.Face7x13, fd.X+11, fd.Y+41, shadowColor)
	text.Draw(screen, dirText, basicfont.Face7x13, fd.X+10, fd.Y+40, dirColor)

	// Draw separator
	separator := ebiten.NewImage(fd.Width-20, 1)
	separator.Fill(color.RGBA{60, 60, 60, 255})
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+10), float64(fd.Y+45))
	screen.DrawImage(separator, opts)

	// Draw file list
	fileListY := fd.Y + 60
	fileHeight := 20
	endIdx := fd.ScrollOffset + fd.MaxVisibleFiles
	if endIdx > len(fd.Files) {
		endIdx = len(fd.Files)
	}

	for i := fd.ScrollOffset; i < endIdx; i++ {
		y := fileListY + (i-fd.ScrollOffset)*fileHeight

		// Draw selection highlight
		if i == fd.SelectedFile {
			selectionBg := ebiten.NewImage(fd.Width-20, fileHeight)
			selectionBg.Fill(color.RGBA{70, 90, 120, 255})
			opts = &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(fd.X+10), float64(y))
			screen.DrawImage(selectionBg, opts)
		}

		// Draw file name with shadow
		fileName := fd.Files[i]
		fileColor := color.RGBA{220, 220, 220, 255}
		if strings.HasSuffix(fileName, "/") {
			fileColor = color.RGBA{180, 180, 255, 255}
		}
		text.Draw(screen, fileName, basicfont.Face7x13, fd.X+16, y+16, shadowColor)
		text.Draw(screen, fileName, basicfont.Face7x13, fd.X+15, y+15, fileColor)
	}

	// Draw separator
	separator = ebiten.NewImage(fd.Width-20, 1)
	separator.Fill(color.RGBA{60, 60, 60, 255})
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+10), float64(fd.Y+fd.Height-80))
	screen.DrawImage(separator, opts)

	// Draw filename input for save dialog
	if fd.IsSaveDialog {
		text.Draw(screen, "Filename:", basicfont.Face7x13, fd.X+11, fd.Y+fd.Height-59, shadowColor)
		text.Draw(screen, "Filename:", basicfont.Face7x13, fd.X+10, fd.Y+fd.Height-60, titleColor)

		// Input field background
		inputBg := ebiten.NewImage(fd.Width-20, 25)
		inputBg.Fill(color.RGBA{30, 30, 30, 255})
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(fd.X+10), float64(fd.Y+fd.Height-55))
		screen.DrawImage(inputBg, opts)

		// Draw input field border
		inputBorderColor := color.RGBA{80, 80, 80, 255}
		for i := 0; i < fd.Width-20; i++ {
			screen.Set(fd.X+10+i, fd.Y+fd.Height-55, inputBorderColor)
			screen.Set(fd.X+10+i, fd.Y+fd.Height-30, inputBorderColor)
		}
		for i := 0; i < 25; i++ {
			screen.Set(fd.X+10, fd.Y+fd.Height-55+i, inputBorderColor)
			screen.Set(fd.X+fd.Width-10, fd.Y+fd.Height-55+i, inputBorderColor)
		}

		// Draw filename text with shadow
		text.Draw(screen, fd.FileName, basicfont.Face7x13, fd.X+16, fd.Y+fd.Height-37, shadowColor)
		text.Draw(screen, fd.FileName, basicfont.Face7x13, fd.X+15, fd.Y+fd.Height-38, color.RGBA{220, 220, 220, 255})

		// Draw cursor with improved visibility
		cursorPos := text.BoundString(basicfont.Face7x13, fd.FileName[:fd.CursorPos]).Dx()
		cursorHeight := 18
		cursorColor := color.RGBA{220, 220, 220, 255}
		for i := 0; i < cursorHeight; i++ {
			screen.Set(fd.X+15+cursorPos, fd.Y+fd.Height-48+i, cursorColor)
		}
	}

	// Draw OK and Cancel buttons with improved styling
	okBtnBg := ebiten.NewImage(80, 30)
	okBtnBg.Fill(color.RGBA{70, 130, 180, 255})
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+fd.Width-180), float64(fd.Y+fd.Height-30))
	screen.DrawImage(okBtnBg, opts)

	cancelBtnBg := ebiten.NewImage(80, 30)
	cancelBtnBg.Fill(color.RGBA{180, 70, 70, 255})
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(fd.X+fd.Width-90), float64(fd.Y+fd.Height-30))
	screen.DrawImage(cancelBtnBg, opts)

	// Draw button text with shadows
	text.Draw(screen, "OK", basicfont.Face7x13, fd.X+fd.Width-149, fd.Y+fd.Height-9, shadowColor)
	text.Draw(screen, "OK", basicfont.Face7x13, fd.X+fd.Width-150, fd.Y+fd.Height-10, color.White)

	text.Draw(screen, "Cancel", basicfont.Face7x13, fd.X+fd.Width-69, fd.Y+fd.Height-9, shadowColor)
	text.Draw(screen, "Cancel", basicfont.Face7x13, fd.X+fd.Width-70, fd.Y+fd.Height-10, color.White)
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
	fileHeight := 20
	fileListHeight := fd.MaxVisibleFiles * fileHeight

	if x >= fd.X+10 && x <= fd.X+fd.Width-10 &&
		y >= fileListY && y <= fileListY+fileListHeight {
		// Clicked on file list
		clickedIndex := fd.ScrollOffset + (y-fileListY)/fileHeight
		if clickedIndex >= 0 && clickedIndex < len(fd.Files) {
			// If clicking on a directory
			if strings.HasSuffix(fd.Files[clickedIndex], "/") {
				if fd.Files[clickedIndex] == ".." {
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

	// OK button
	if x >= fd.X+fd.Width-180 && x <= fd.X+fd.Width-100 &&
		y >= fd.Y+fd.Height-30 && y <= fd.Y+fd.Height {
		// Handle OK button
		return true
	}

	// Cancel button
	if x >= fd.X+fd.Width-90 && x <= fd.X+fd.Width-10 &&
		y >= fd.Y+fd.Height-30 && y <= fd.Y+fd.Height {
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
