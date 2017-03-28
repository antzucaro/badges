package main

import (
	"fmt"
	"github.com/antzucaro/qstr"
	"github.com/fogleman/gg"
)

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
	Color []qstr.RGBColor

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
	Background         string
	BackgroundColor    qstr.RGBColor
	Overlay            string
	Font               string
	Width              int
	Height             int
	NumGameTypes       int
	NickConfig         TextConfig
	GameTypeConfig     TextConfig
	NoStatsConfig      TextConfig
	EloConfig          TextConfig
	RankConfig         TextConfig
	WinPctLabelConfig  TextConfig
	WinPctConfig       TextConfig
	WinConfig          TextConfig
	LossConfig         TextConfig
	KDRatioLabelConfig TextConfig
	KDRatio            TextConfig
	KillsConfig        TextConfig
	DeathsConfig       TextConfig
	PlayingTimeConfig  TextConfig
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
	s.context.SetRGB(config.Color[0].R, config.Color[0].G, config.Color[0].B)
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
func (s *Skin) placeQStr(text qstr.QStr, config TextConfig, lightnessFloor float64, lightnessCeiling float64) {
	s.context.LoadFontFace(config.Font, config.FontSize)

	// shrink the nick until it fits within the allotted space
	stripped := text.Stripped()
	for w, _ := s.context.MeasureString(stripped); int(w) > config.MaxWidth; {
		// decrease the fontsize by two points and try again
		config.FontSize -= 2

		s.context.LoadFontFace(config.Font, config.FontSize)
		w, _ = s.context.MeasureString(stripped)
	}

	x := config.Pos.X
	var cappedColor qstr.RGBColor
	for _, colorPart := range text.ColorParts() {
		// allow capping the lightness at the high or low end depending on the background
		if lightnessFloor > 0 || lightnessCeiling < 1 {
			cappedColor = colorPart.Color.CapLightness(lightnessFloor, lightnessCeiling)
		} else {
			cappedColor = colorPart.Color
		}

		s.context.SetRGB(cappedColor.R, cappedColor.G, cappedColor.B)
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
	s.placeQStr(pd.Nick, s.Params.NickConfig, 0.4, 1)

	// Gametype labels
	gameTypePositions := []Position{{100.0, 35.0}, {195.0, 35.0}, {290.0, 35.0}}
	for i, elo := range pd.Elos {
		s.Params.GameTypeConfig.Pos = gameTypePositions[i]
		s.placeText(elo.GameType, s.Params.GameTypeConfig)
	}

	// Elos for those game types
	eloPositions := []Position{{100.0, 50.0}, {195.0, 50.0}, {290.0, 50.0}}
	for i, elo := range pd.Elos {
		s.Params.EloConfig.Pos = eloPositions[i]
		s.placeText(fmt.Sprintf("Elo %d", elo.Elo), s.Params.EloConfig)
	}

	// Ranks for those game types
	rankPositions := []Position{{100.0, 60.0}, {195.0, 60.0}, {290.0, 60.0}}
	for i, pos := range rankPositions {
		s.Params.RankConfig.Pos = pos
		if i < len(pd.Ranks) {
			s.placeText(fmt.Sprintf("Rank %d of %d", pd.Ranks[i].Rank, pd.Ranks[i].MaxRank), s.Params.RankConfig)
		} else {
			s.placeText("(preliminary)", s.Params.RankConfig)
		}
	}

	// Kill Ratio and its details
	s.placeText("Kill Ratio", s.Params.KDRatioLabelConfig)
	s.placeText(fmt.Sprintf("%.3f", pd.KDRatio()), s.Params.KDRatio)
	s.placeText(fmt.Sprintf("%d kills", pd.Kills), s.Params.KillsConfig)
	s.placeText(fmt.Sprintf("%d deaths", pd.Deaths), s.Params.DeathsConfig)

	// Win Percentage and its details
	s.placeText("Win Percentage", s.Params.WinPctLabelConfig)
	s.placeText(fmt.Sprintf("%.2f%%", pd.WinPct()), s.Params.WinPctConfig)
	s.placeText(fmt.Sprintf("%d wins", pd.Wins), s.Params.WinConfig)
	s.placeText(fmt.Sprintf("%d losses", pd.Losses), s.Params.LossConfig)

	// Playing time
	s.placeText(fmt.Sprintf("Playing Time: %s", pd.PlayingTime), s.Params.PlayingTimeConfig)

	s.context.SavePNG(filename)
}

// The "classic" skin theme
var ArcherSkin = Skin{
	Name: "archer",
	Params: SkinParams{
		Background:      "images/background_archer-v3.png",
		BackgroundColor: qstr.RGBColor{0.00, 0.00, 0.00},
		Overlay:         "",
		Font:            "Xolonium",
		Width:           560,
		Height:          720,
		NumGameTypes:    3,
		NickConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 22,
			Pos:      Position{X: 53.0, Y: 20.0},
			Color:    []qstr.RGBColor{{0.5, 0.5, 0.5}},
			MaxWidth: 270,
		},
		GameTypeConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 101.0, Y: 33.0},
			Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
			Width:    94,
			Align:    "center",
		},
		NoStatsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 12,
			Pos:      Position{X: 101.0, Y: 59.0},
			Color:    []qstr.RGBColor{{0.8, 0.2, 0.1}},
			Angle:    -10,
		},
		EloConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 101.0, Y: 47.0},
			Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
			Align:    "center",
		},
		RankConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 8,
			Pos:      Position{X: 101.0, Y: 58.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 1.0}},
			Align:    "center",
		},
		WinPctLabelConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 508.0, Y: 6.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.8}},
			Align:    "center",
		},
		WinPctConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 509.0, Y: 24.0},
			Color:    []qstr.RGBColor{{1.00, 1.00, 1.00}},
			Align:    "center",
		},
		WinConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 508.0, Y: 37.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.6}},
			Align:    "center",
		},
		LossConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 508.0, Y: 47.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.6}},
			Align:    "center",
		},
		KDRatioLabelConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 390.0, Y: 6.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.8}},
			Width:    102,
			Align:    "center",
		},
		KDRatio: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 392.0, Y: 24.0},
			Color:    []qstr.RGBColor{{1.00, 1.00, 1.00}},
			Align:    "center",
		},
		KillsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392.0, Y: 37.0},
			Color:    []qstr.RGBColor{{0.6, 0.8, 0.6}},
			Align:    "center",
		},
		DeathsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392.0, Y: 47.0},
			Color:    []qstr.RGBColor{{0.8, 0.6, 0.6}},
			Align:    "center",
		},
		PlayingTimeConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 451.0, Y: 63.0},
			Color:    []qstr.RGBColor{{0.1, 0.1, 0.1}},
			Align:    "center",
		},
	},
}
