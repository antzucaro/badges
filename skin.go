package main

import (
	"github.com/fogleman/gg"
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

// TextConfig provides information used to place text on a canvas
type TextConfig struct {
	// The full path to the TTF font file to be loaded
	Font string

	// The size in points the font will be
	FontSize int

	// The X and Y location where the bottom-left of the text will begin
	Pos Position

	// The color the text will be
	Color RGB

	// Rotation
	Angle int

	// Dimensions
	Height   int
	Width    int
	MaxWidth int

	// Text alignment: left: -1, center: 0, right: 1
	Align int
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
	NickConfig        TextConfig
	GameTypeConfig    TextConfig
	NoStatsConfig     TextConfig
	EloConfig         TextConfig
	RankConfig        TextConfig
	WinConfig         TextConfig
	WinPctConfig      TextConfig
	LossConfig        TextConfig
	KDConfig          TextConfig
	KDRatio           TextConfig
	KillsConfig       TextConfig
	DeathsConfig      TextConfig
	PlayingTimeConfig TextConfig
}

// Skin represents the look and feel of a XonStat badge
type Skin struct {
	Name    string
	Params  SkinParams
	Context *gg.Context
}

// String representation of a Skin
func (s *Skin) String() string {
	return s.Name
}

// Render the provided PlayerData using this Skin
func (s *Skin) Render(pd *PlayerData, output string) {
}

// The "classic" skin theme
var ArcherSkin = Skin{
	Name: "archer",
	Params: SkinParams{
		Background:      "images/background_archer-v3.png",
		BackgroundColor: RGB{Red: 0.00, Green: 0.00, Blue: 0.00},
		Overlay:         "",
		Font:            "Xolonium",
		Width:           560,
		Height:          720,
		NumGameTypes:    3,
		NickConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 22,
			Pos:      Position{X: 53, Y: 20},
			MaxWidth: 270,
		},
		GameTypeConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 101, Y: 33},
			Color:    RGB{Red: 0.09, Green: 0.09, Blue: 0.09},
			Width:    94,
		},
		NoStatsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 12,
			Pos:      Position{X: 101, Y: 59},
			Color:    RGB{Red: 0.8, Green: 0.2, Blue: 0.1},
			Angle:    -10,
		},
		EloConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 101, Y: 47},
			Color:    RGB{Red: 1.0, Green: 1.0, Blue: 0.5},
		},
		RankConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 8,
			Pos:      Position{X: 101, Y: 58},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 1.0},
		},
		WinConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 508, Y: 3},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 0.8},
		},
		WinPctConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 509, Y: 18},
			Color:    RGB{Red: 0.00, Green: 0.00, Blue: 0.00},
		},
		LossConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 508, Y: 44},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 0.6},
		},
		KDConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 390, Y: 3},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 0.8},
			Width:    102,
		},
		KDRatio: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 392, Y: 18},
		},
		KillsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392, Y: 33},
			Color:    RGB{Red: 0.6, Green: 0.8, Blue: 0.6},
		},
		DeathsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392, Y: 44},
			Color:    RGB{Red: 0.8, Green: 0.6, Blue: 0.6},
		},
		PlayingTimeConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 451, Y: 59},
			Color:    RGB{Red: 0.1, Green: 0.1, Blue: 0.1},
		},
	},
}
