package main

import (
	"fmt"
	"github.com/antzucaro/qstr"
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
	X float64
	Y float64
}

// TextConfig provides information used to place text on a canvas
type TextConfig struct {
	// The full path to the TTF font file to be loaded
	Font string

	// The size in points the font will be
	FontSize float64

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

	// Text alignment
	Align string
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
	context *gg.Context
}

// placeText "writes" text on the drawing canvas
func (s *Skin) placeText(text string, config TextConfig) {
	s.context.LoadFontFace(config.Font, config.FontSize)
	s.context.SetRGB(config.Color.Red, config.Color.Green, config.Color.Blue)
	if config.Align == "" {
		s.context.DrawString(text, config.Pos.X, config.Pos.Y)
	} else if config.Align == "center" {
		s.context.DrawStringAnchored(text, config.Pos.X, config.Pos.Y, 0.5, 0.5)
	} else if config.Align == "right" {
		s.context.DrawStringAnchored(text, config.Pos.X, config.Pos.Y, 1, 0.5)
	}

}

// placeQStr does the same thing as placeText does, but with potentially
// colorized QStrs
func (s *Skin) placeQStr(text qstr.QStr, config TextConfig) {
	s.context.LoadFontFace(config.Font, config.FontSize)

	x := config.Pos.X
	for _, colorPart := range text.ColorParts() {
		s.context.SetRGB(colorPart.Color.R, colorPart.Color.G, colorPart.Color.B)
		s.context.DrawString(colorPart.Part, x, config.Pos.Y)

		// the starting point for the next part is the end of the last one
		w, _ := s.context.MeasureString(colorPart.Part)
		x += w
	}
}

// String representation of a Skin
func (s *Skin) String() string {
	return s.Name
}

// Render the provided PlayerData using this Skin
func (s *Skin) Render(pd *PlayerData, filename string) {
	// load the background
	im, err := gg.LoadPNG(s.Params.Background)
	if err != nil {
		panic(err)
	}
	s.context = gg.NewContextForImage(im)

	// Nick
	s.placeQStr(pd.Nick, s.Params.NickConfig)

	// Gametype labels
	gameTypePositions := []Position{
		Position{X: 100.0, Y: 35.0},
		Position{X: 195.0, Y: 35.0},
		Position{X: 290.0, Y: 35.0},
	}

	for i, elo := range pd.Elos {
		s.Params.GameTypeConfig.Pos = gameTypePositions[i]
		s.placeText(elo.GameType, s.Params.GameTypeConfig)
	}

	// Elos for those game types
	eloPositions := []Position{
		Position{X: 100.0, Y: 50.0},
		Position{X: 195.0, Y: 50.0},
		Position{X: 290.0, Y: 50.0},
	}

	for i, elo := range pd.Elos {
		s.Params.EloConfig.Pos = eloPositions[i]
		s.placeText(fmt.Sprintf("Elo %d", elo.Elo), s.Params.EloConfig)
	}

	// Ranks for those game types
	rankPositions := []Position{
		Position{X: 100.0, Y: 60.0},
		Position{X: 195.0, Y: 60.0},
		Position{X: 290.0, Y: 60.0},
	}

	for i, pos := range rankPositions {
		s.Params.RankConfig.Pos = pos
		if i < len(pd.Ranks) {
			s.placeText(fmt.Sprintf("Rank %d of %d", pd.Ranks[i].Rank, pd.Ranks[i].MaxRank), s.Params.RankConfig)
		} else {
			s.placeText("(preliminary)", s.Params.RankConfig)
		}
	}

	s.context.SavePNG(filename)
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
			Pos:      Position{X: 53.0, Y: 20.0},
			Color:    RGB{Red: 0.5, Green: 0.5, Blue: 0.5},
			MaxWidth: 270,
		},
		GameTypeConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 101.0, Y: 33.0},
			Color:    RGB{Red: 0.9, Green: 0.9, Blue: 0.9},
			Width:    94,
			Align:    "center",
		},
		NoStatsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 12,
			Pos:      Position{X: 101.0, Y: 59.0},
			Color:    RGB{Red: 0.8, Green: 0.2, Blue: 0.1},
			Angle:    -10,
		},
		EloConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 101.0, Y: 47.0},
			Color:    RGB{Red: 1.0, Green: 1.0, Blue: 0.5},
			Align:    "center",
		},
		RankConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 8,
			Pos:      Position{X: 101.0, Y: 58.0},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 1.0},
			Align:    "center",
		},
		WinConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 508.0, Y: 3.0},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 0.8},
		},
		WinPctConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 509.0, Y: 18.0},
			Color:    RGB{Red: 0.00, Green: 0.00, Blue: 0.00},
		},
		LossConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 508.0, Y: 44.0},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 0.6},
		},
		KDConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 390.0, Y: 3.0},
			Color:    RGB{Red: 0.8, Green: 0.8, Blue: 0.8},
			Width:    102,
		},
		KDRatio: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 392.0, Y: 18.0},
		},
		KillsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392.0, Y: 33.0},
			Color:    RGB{Red: 0.6, Green: 0.8, Blue: 0.6},
		},
		DeathsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392.0, Y: 44.0},
			Color:    RGB{Red: 0.8, Green: 0.6, Blue: 0.6},
		},
		PlayingTimeConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 451.0, Y: 59.0},
			Color:    RGB{Red: 0.1, Green: 0.1, Blue: 0.1},
		},
	},
}
