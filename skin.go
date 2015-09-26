package main

import (
	"github.com/ungerik/go-cairo"
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
