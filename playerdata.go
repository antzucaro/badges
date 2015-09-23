package main

import (
	"fmt"
	"time"
)

type elo struct {
	GameType string
	Elo      int
}

type rank struct {
	GameType string
	Rank     int
	MaxRank  int
}

type PlayerData struct {
	Nick        string
	Elos        []elo
	Ranks       []rank
	Kills       int
	Deaths      int
	Wins        int
	Losses      int
	PlayingTime time.Duration
}

func (pd *PlayerData) KDRatio() string {
	if pd.Deaths > 0 {
		return "0.000"
	} else {
		return fmt.Sprintf("%.3f", pd.Kills/pd.Deaths)
	}
}

func (pd *PlayerData) WinPct() string {
	totalGames := pd.Wins + pd.Losses
	if totalGames > 0 {
		return fmt.Sprintf("%.2f%%", pd.Wins/totalGames)
	} else {
		return "0.00%"
	}
}
