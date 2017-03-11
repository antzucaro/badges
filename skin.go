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
	Align     int
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

// The "classic" skin theme
var DefaultSkin = Skin{
	Name: "classic",
	Params: SkinParams{
		Background:      "",
		BackgroundColor: RGB{Red: 0.00, Green: 0.00, Blue: 0.00},
		Overlay:         "",
		Font:            "Xolonium",
		Width:           560,
		Height:          720,
		NumGameTypes:    3,
		NickConfig: BasicTextConfig{
			FontSize: 22,
			Pos:      Position{X: 53, Y: 20},
			MaxWidth: 270,
		},
		GameTypeConfig: BasicTextConfig{
			FontSize: 10,
			Pos:      Position{X: 101, Y: 33},
			Color:    RGB{Red: 0.09, Green: 0.09, Blue: 0.09},
			TextFmt:  "%s",
			Width:    94,
		},
		NoStatsConfig: BasicTextConfig{
			FontSize: 12,
			Pos:      Position{X: 101, Y: 59},
			Color:    RGB{Red: 0.8, Green: 0.2, Blue: 0.1},
			Angle:    -10,
			TextFmt:  "no stats yet!",
		},
		EloConfig: BasicTextConfig{
			FontSize: 10,
			Pos:      Position{X: 101, Y: 47},
			Color:    RGB{Red: 1.0, Green: 1.0, Blue: 0.5},
			TextFmt:  "Elo %.0f",
		},
		RankConfig: BasicTextConfig{
			FontSize: 8,
			Pos:      Position{X: 101, Y: 58},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 1.0},
			TextFmt:  "Rank %d of %d",
		},
		WinConfig: BasicTextConfig{
			FontSize:  10,
			Pos:       Position{X: 508, Y: 3},
			Color:     RGB{Red: 0.8, Green: 0.8, Blue: 0.8},
			TextFmt:   "Win Percentage",
			UpperCase: true,
		},
		WinPctConfig: BasicTextConfig{
			FontSize: 15,
			Pos:      Position{X: 509, Y: 18},
			Color:    RGB{Red: 0.00, Green: 0.00, Blue: 0.00},
			TopColor: RGB{Red: 0.2, Green: 1.0, Blue: 1.0},
			MidColor: RGB{Red: 0.4, Green: 0.8, Blue: 0.4},
			BotColor: RGB{Red: 1.0, Green: 1.0, Blue: 0.2},
		},
		LossConfig: BasicTextConfig{
			FontSize: 9,
			Pos:      Position{X: 508, Y: 44},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 0.6},
		},
		KDConfig: BasicTextConfig{
			FontSize: 10,
			Pos:      Position{X: 390, Y: 3},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 0.8},
			TextFmt:  "Kill Ratio",
			Width:    102,
		},
		KDRatio: BasicTextConfig{
			FontSize: 15,
			Pos:      Position{X: 392, Y: 18},
		},
		KillsConfig: BasicTextConfig{
			FontSize: 9,
			Pos:      Position{X: 392, Y: 33},
			Color:    RGB{Red: 0.6, Green: 0.8, Blue: 0.6},
		},
		DeathsConfig: BasicTextConfig{
			FontSize: 9,
			Pos:      Position{X: 392, Y: 44},
			Color:    RGB{Red: 0.8, Green: 0.6, Blue: 0.6},
		},
		PlayingTimeConfig: BasicTextConfig{
			FontSize: 10,
			Pos:      Position{X: 451, Y: 59},
			Color:    RGB{Red: 0.1, Green: 0.1, Blue: 0.1},
			TextFmt:  "Playing Time %s",
		},
	},
	surface: nil,
}
