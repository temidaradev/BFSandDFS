package assets

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed DMSans.ttf
var MyFont []byte
var FontFaceS text.Face
