package main

import (
	"github.com/ungerik/go-cairo"
	"math"
)

// RGB is an RGB value
type RGB struct {
	Red   float64
	Green float64
	Blue  float64
}

// Position is an (x,y) coordinate
type Position struct {
	X int
	Y int
}

// BasicTextConfig provides basic style information for Cairo text
type BasicTextConfig struct {
	FontSize  int
	Pos       Position
	Color     RGB
	TopColor  RGB
	MidColor  RGB
	BotColor  RGB
	Angle     int
	TextFmt   string
	UpperCase bool
	Height    int
	Width     int
	MaxWidth  int
}

// SkinParams represents parameters given to Skin objects
type SkinParams struct {
	Background        string
	BackgroundColor   RGB
	Overlay           string
	Font              string
	Width             int
	Height            int
	NumGameTypes      int
	NickConfig        BasicTextConfig
	GameTypeConfig    BasicTextConfig
	NoStatsConfig     BasicTextConfig
	EloConfig         BasicTextConfig
	RankConfig        BasicTextConfig
	WinConfig         BasicTextConfig
	WinPctConfig      BasicTextConfig
	LossConfig        BasicTextConfig
	KDConfig          BasicTextConfig
	KDRatio           BasicTextConfig
	KillsConfig       BasicTextConfig
	DeathsConfig      BasicTextConfig
	PlayingTimeConfig BasicTextConfig
}

// Skin represents the look and feel of a XonStat badge
type Skin struct {
	Name    string
	Params  SkinParams
	surface *cairo.Surface
}

// String representation of a Skin
func (s *Skin) String() string {
	return s.Name
}

// ShowText shows the given text on the surface with alignment and angling.
func (s *Skin) ShowText(text string, pos Position, align int, angle int, offsetX int, offsetY int) {
	te := s.surface.TextExtents(text)

	if align > 0 {
		s.surface.MoveTo(float64(pos.X+offsetX)-te.Xbearing, float64(pos.Y+offsetY)-te.Ybearing)
	} else if align < 0 {
		s.surface.MoveTo(float64(pos.X+offsetX)-te.Xbearing-te.Width, float64(pos.Y+offsetY)-te.Ybearing)
	} else {
		s.surface.MoveTo(float64(pos.X+offsetX)-te.Xbearing-te.Width/2, float64(pos.Y+offsetY)-te.Ybearing)
	}
	s.surface.Save()

	if angle > 0 {
		s.surface.Rotate(float64(angle) * math.Pi / 180.0)
	}

	s.surface.ShowText(text)
	s.surface.Restore()
}
