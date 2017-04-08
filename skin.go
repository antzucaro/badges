package main

import (
	"fmt"
	"github.com/antzucaro/qstr"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"math"
	"path/filepath"
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
	NoStatsConfig      TextConfig
	GameTypeConfig     []TextConfig
	EloConfig          []TextConfig
	RankConfig         []TextConfig
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
	Name      string
	Params    SkinParams
	context   *gg.Context
	fontCache map[string]font.Face
}

// setFontFace loads a font either from the cache or from the filesystem
func (s *Skin) setFontFace(path string, points float64) error {
	// do we even have a cache yet?
	if len(s.fontCache) == 0 {
		s.fontCache = make(map[string]font.Face)
	}

	_, file := filepath.Split(path)
	key := fmt.Sprintf("%s %f", file, points)
	if ff, ok := s.fontCache[key]; ok {
		s.context.SetFontFace(ff)
		return nil
	} else {
		ff, err := gg.LoadFontFace(path, points)
		if err != nil {
			return err
		} else {
			s.context.SetFontFace(ff)
			s.fontCache[key] = ff
			return nil
		}
	}
}

// placeText "writes" text on the drawing canvas
func (s *Skin) placeText(text string, config TextConfig) {
	if config.FontSize == 0 {
		return
	}

	s.setFontFace(config.Font, config.FontSize)

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
	s.setFontFace(config.Font, config.FontSize)

	// shrink the nick until it fits within the allotted space
	stripped := text.Stripped()
	for w, _ := s.context.MeasureString(stripped); int(w) > config.MaxWidth; {
		// decrease the fontsize by two points and try again
		config.FontSize -= 2

		s.setFontFace(config.Font, config.FontSize)
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

func (s *Skin) ShadeKDRatio(kdRatio float64, hiColor, midColor, loColor *qstr.RGBColor) qstr.RGBColor {
	var nr float64
	var c1, c2 *qstr.RGBColor
	if kdRatio >= 1.0 {
		nr = kdRatio - 1.0
		if nr > 1 {
			nr = 1.0
		}
		c1 = hiColor
		c2 = midColor
	} else {
		nr = kdRatio
		c1 = midColor
		c2 = loColor
	}

	// shade the KDRatio according to how good it is
	r := nr*c1.R + (1-nr)*c2.R
	g := nr*c1.G + (1-nr)*c2.G
	b := nr*c1.B + (1-nr)*c2.B

	return qstr.RGBColor{r, g, b}
}

func (s *Skin) ShadeWinPct(winPct float64, hiColor, midColor, loColor *qstr.RGBColor) qstr.RGBColor {
	var nr float64
	var c1, c2 *qstr.RGBColor

	if winPct > 50.0 {
		nr = 2 * (winPct/100 - 0.5)
		c1 = &s.Params.WinPctConfig.Color[0]
		c2 = &s.Params.WinPctConfig.Color[1]
	} else {
		nr = 2 * (winPct / 100)
		c1 = &s.Params.WinPctConfig.Color[1]
		c2 = &s.Params.WinPctConfig.Color[2]
	}

	// shade the WinPct according to how good it is
	r := nr*c1.R + (1-nr)*c2.R
	g := nr*c1.G + (1-nr)*c2.G
	b := nr*c1.B + (1-nr)*c2.B

	return qstr.RGBColor{r, g, b}
}

// Render the provided PlayerData using this Skin
func (s *Skin) Render(pd *PlayerData, filename string) {
	s.context = gg.NewContext(s.Params.Width, s.Params.Height)

	// load the background
	if s.Params.Background != "" {
		bg, err := gg.LoadPNG(s.Params.Background)

		// the background can be a small image that can be repeated (tiled)
		bgW := bg.Bounds().Size().X
		bgH := bg.Bounds().Size().Y
		if err != nil {
			panic(err)
		}

		repeatX := int(math.Ceil(float64(s.Params.Width) / float64(bgW)))
		repeatY := int(math.Ceil(float64(s.Params.Height) / float64(bgH)))
		for i := 0; i < repeatX; i++ {
			for j := 0; j < repeatY; j++ {
				s.context.DrawImage(bg, bgW*i, bgH*j)
			}
		}
	}

	// load the overlay
	if s.Params.Overlay != "" {
		overlay, err := gg.LoadPNG(s.Params.Overlay)
		if err != nil {
			panic(err)
		}
		s.context.DrawImage(overlay, 0, 0)
	}

	// Nick
	s.placeQStr(pd.Nick, s.Params.NickConfig, 0.4, 1)

	// Game type labels along with Elos for those game types
	for i, elo := range pd.Elos {
		s.placeText(elo.GameType, s.Params.GameTypeConfig[i])
		s.placeText(fmt.Sprintf("Elo %d", elo.Elo), s.Params.EloConfig[i])
	}

	// Ranks for those game types
	for i, textConfig := range s.Params.RankConfig {
		if i < len(pd.Ranks) {
			s.placeText(fmt.Sprintf("Rank %d of %d", pd.Ranks[i].Rank, pd.Ranks[i].MaxRank), textConfig)
		} else {
			s.placeText("(preliminary)", textConfig)
		}
	}

	// Kill Ratio and its details
	s.placeText("Kill Ratio", s.Params.KDRatioLabelConfig)

	kdRatio := pd.KDRatio()
	s.Params.KDRatio.Color[0] = s.ShadeKDRatio(kdRatio, &s.Params.KDRatio.Color[0], &s.Params.KDRatio.Color[1],
		&s.Params.KDRatio.Color[2])
	s.placeText(fmt.Sprintf("%.3f", kdRatio), s.Params.KDRatio)

	s.placeText(fmt.Sprintf("%d kills", pd.Kills), s.Params.KillsConfig)
	s.placeText(fmt.Sprintf("%d deaths", pd.Deaths), s.Params.DeathsConfig)

	// Win Percentage and its details
	s.placeText("Win Percentage", s.Params.WinPctLabelConfig)

	winPct := pd.WinPct()
	s.Params.WinPctConfig.Color[0] = s.ShadeWinPct(winPct, &s.Params.WinPctConfig.Color[0],
		&s.Params.WinPctConfig.Color[1], &s.Params.WinPctConfig.Color[2])
	s.placeText(fmt.Sprintf("%.2f%%", winPct), s.Params.WinPctConfig)
	s.placeText(fmt.Sprintf("%d wins", pd.Wins), s.Params.WinConfig)
	s.placeText(fmt.Sprintf("%d losses", pd.Losses), s.Params.LossConfig)

	// Playing time
	s.placeText(fmt.Sprintf("Playing Time: %s", pd.PlayingTime), s.Params.PlayingTimeConfig)

	s.context.SavePNG(filename)
}

// The "Archer" skin theme
var ArcherSkin = Skin{
	Name: "archer",
	Params: SkinParams{
		Background:      "images/background_archer-v3.png",
		BackgroundColor: qstr.RGBColor{0.00, 0.00, 0.00},
		Overlay:         "",
		Font:            "Xolonium",
		Width:           560,
		Height:          70,
		NumGameTypes:    3,
		NickConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 22,
			Pos:      Position{X: 53.0, Y: 20.0},
			Color:    []qstr.RGBColor{{0.5, 0.5, 0.5}},
			MaxWidth: 270,
		},
		GameTypeConfig: []TextConfig{
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 100.0, Y: 35.0},
				Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
				Width:    94,
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 195.0, Y: 35.0},
				Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
				Width:    94,
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 290.0, Y: 35.0},
				Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
				Width:    94,
				Align:    "center",
			},
		},
		NoStatsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 12,
			Pos:      Position{X: 101.0, Y: 59.0},
			Color:    []qstr.RGBColor{{0.8, 0.2, 0.1}},
			Angle:    -10,
		},
		EloConfig: []TextConfig{
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 100.0, Y: 50.0},
				Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 195.0, Y: 50.0},
				Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
				Align:    "center",
			},
						{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 290.0, Y: 50.0},
				Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
				Align:    "center",
			},
		},
		RankConfig: []TextConfig{
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 8,
				Pos:      Position{X: 100.0, Y: 60.0},
				Color:    []qstr.RGBColor{{0.8, 0.8, 1.0}},
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 8,
				Pos:      Position{X: 195.0, Y: 60.0},
				Color:    []qstr.RGBColor{{0.8, 0.8, 1.0}},
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 8,
				Pos:      Position{X: 290.0, Y: 60.0},
				Color:    []qstr.RGBColor{{0.8, 0.8, 1.0}},
				Align:    "center",
			},
		},
		WinPctLabelConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 508.0, Y: 5.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.8}},
			Align:    "center",
		},
		WinPctConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 509.0, Y: 21.0},
			Color:    []qstr.RGBColor{{0.2, 1.0, 1.0}, {0.4, 0.8, 0.4}, {1.0, 1.0, 0.2}},
			Align:    "center",
		},
		WinConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 508.0, Y: 34.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.6}},
			Align:    "center",
		},
		LossConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 508.0, Y: 44.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.6}},
			Align:    "center",
		},
		KDRatioLabelConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 390.0, Y: 5.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.8}},
			Width:    102,
			Align:    "center",
		},
		KDRatio: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 392.0, Y: 21.0},
			Color:    []qstr.RGBColor{{0.2, 1.0, 0.2}, {0.8, 0.8, 0.4}, {1.0, 0.2, 0.2}},
			Align:    "center",
		},
		KillsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392.0, Y: 34.0},
			Color:    []qstr.RGBColor{{0.6, 0.8, 0.6}},
			Align:    "center",
		},
		DeathsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392.0, Y: 44.0},
			Color:    []qstr.RGBColor{{0.8, 0.6, 0.6}},
			Align:    "center",
		},
		PlayingTimeConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 451.0, Y: 61.0},
			Color:    []qstr.RGBColor{{0.1, 0.1, 0.1}},
			Align:    "center",
		},
	},
}

// The "Default" skin theme
var DefaultSkin = Skin{
	Name: "default",
	Params: SkinParams{
		Background:      "images/broken_noise.png",
		BackgroundColor: qstr.RGBColor{0.00, 0.00, 0.00},
		Overlay:         "images/overlay_classic.png",
		Font:            "Xolonium",
		Width:           560,
		Height:          70,
		NumGameTypes:    3,
		NickConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 22,
			Pos:      Position{X: 53.0, Y: 20.0},
			Color:    []qstr.RGBColor{{0.5, 0.5, 0.5}},
			MaxWidth: 270,
		},
		NoStatsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 12,
			Pos:      Position{X: 101.0, Y: 59.0},
			Color:    []qstr.RGBColor{{0.8, 0.2, 0.1}},
			Angle:    -10,
		},
		GameTypeConfig: []TextConfig{
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 100.0, Y: 35.0},
				Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
				Width:    94,
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 195.0, Y: 35.0},
				Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
				Width:    94,
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 290.0, Y: 35.0},
				Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
				Width:    94,
				Align:    "center",
			},
		},
		EloConfig: []TextConfig{
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 100.0, Y: 50.0},
				Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 195.0, Y: 50.0},
				Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
				Align:    "center",
			},
						{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 290.0, Y: 50.0},
				Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
				Align:    "center",
			},
		},
		RankConfig: []TextConfig{
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 8,
				Pos:      Position{X: 100.0, Y: 60.0},
				Color:    []qstr.RGBColor{{0.8, 0.8, 1.0}},
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 8,
				Pos:      Position{X: 195.0, Y: 60.0},
				Color:    []qstr.RGBColor{{0.8, 0.8, 1.0}},
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 8,
				Pos:      Position{X: 290.0, Y: 60.0},
				Color:    []qstr.RGBColor{{0.8, 0.8, 1.0}},
				Align:    "center",
			},
		},
		WinPctLabelConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 508.0, Y: 5.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.8}},
			Align:    "center",
		},
		WinPctConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 509.0, Y: 21.0},
			Color:    []qstr.RGBColor{{0.2, 1.0, 1.0}, {0.4, 0.8, 0.4}, {1.0, 1.0, 0.2}},
			Align:    "center",
		},
		WinConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 508.0, Y: 34.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.6}},
			Align:    "center",
		},
		LossConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 508.0, Y: 44.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.6}},
			Align:    "center",
		},
		KDRatioLabelConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 390.0, Y: 5.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.8}},
			Width:    102,
			Align:    "center",
		},
		KDRatio: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 392.0, Y: 21.0},
			Color:    []qstr.RGBColor{{0.2, 1.0, 0.2}, {0.8, 0.8, 0.4}, {1.0, 0.2, 0.2}},
			Align:    "center",
		},
		KillsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392.0, Y: 34.0},
			Color:    []qstr.RGBColor{{0.6, 0.8, 0.6}},
			Align:    "center",
		},
		DeathsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 9,
			Pos:      Position{X: 392.0, Y: 44.0},
			Color:    []qstr.RGBColor{{0.8, 0.6, 0.6}},
			Align:    "center",
		},
		PlayingTimeConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 451.0, Y: 61.0},
			Color:    []qstr.RGBColor{{0.1, 0.1, 0.1}},
			Align:    "center",
		},
	},
}

// The "Minimal" skin theme
var MinimalSkin = Skin{
	Name: "minimal",
	Params: SkinParams{
		BackgroundColor: qstr.RGBColor{0.04, 0.04, 0.04},
		Overlay:         "images/overlay_minimal.png",
		Font:            "Xolonium",
		Width:           560,
		Height:          40,
		NumGameTypes:    3,
		NickConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 22,
			Pos:      Position{X: 36.0, Y: 17.0},
			MaxWidth: 280,
		},
		NoStatsConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 12,
			Pos:      Position{X: 101.0, Y: 59.0},
			Color:    []qstr.RGBColor{{0.8, 0.2, 0.1}},
			Angle:    -10,
		},
		GameTypeConfig: []TextConfig{
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 52.0, Y: 37.0},
				Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 142.0, Y: 37.0},
				Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 251.0, Y: 37.0},
				Color:    []qstr.RGBColor{{0.9, 0.9, 0.9}},
				Width:    94,
				Align:    "center",
			},
		},
		EloConfig: []TextConfig{
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 100.0, Y: 50.0},
				Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
				Align:    "center",
			},
			{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 195.0, Y: 50.0},
				Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
				Align:    "center",
			},
						{
				Font:     "fonts/xolonium.ttf",
				FontSize: 10,
				Pos:      Position{X: 290.0, Y: 50.0},
				Color:    []qstr.RGBColor{{1.0, 1.0, 0.5}},
				Align:    "center",
			},
		},
		WinPctLabelConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 508.0, Y: 5.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.8}},
			Align:    "center",
		},
		WinPctConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 509.0, Y: 15.0},
			Color:    []qstr.RGBColor{{0.6, 0.8, 0.8}, {0.6, 0.6, 0.6}, {0.8, 0.8, 0.6}},
			Align:    "center",
		},
		KDRatioLabelConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 11,
			Pos:      Position{X: 392.0, Y: 15.0},
			Color:    []qstr.RGBColor{{0.8, 0.8, 0.8}},
			Width:    102,
			Align:    "center",
		},
		KDRatio: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 15,
			Pos:      Position{X: 392.0, Y: 21.0},
			Color:    []qstr.RGBColor{{0.6, 0.8, 0.6}, {0.6, 0.6, 0.6}, {0.8, 0.6, 0.6}},
			Align:    "center",
		},
		PlayingTimeConfig: TextConfig{
			Font:     "fonts/xolonium.ttf",
			FontSize: 10,
			Pos:      Position{X: 451.0, Y: 30.0},
			Color:    []qstr.RGBColor{{0.7, 0.7, 0.7}},
			Align:    "center",
		},
	},
}
