package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

type PlayerProcessor struct {
	db *sql.DB
}

func NewPlayerProcessor(connStr string) (*PlayerProcessor, error) {
	// establish a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// connection pooling
	db.SetMaxIdleConns(5)

	pp := new(PlayerProcessor)
	pp.db = db
	return pp, nil
}

func (pp *PlayerProcessor) FindPlayers(delta int, limit int) ([]int, error) {
	playersSQL := `SELECT distinct p.player_id 
	FROM players p JOIN player_game_stats pgs on p.player_id = pgs.player_id
    JOIN player_elos pe on p.player_id = pe.player_id
	WHERE p.active_ind = true
	AND p.player_id > 2
	AND p.nick IS NOT NULL`

	// constrain the time window if needed
	if delta > 0 {
		playersSQL += " AND pgs.create_dt > now() - interval '" + fmt.Sprintf("%d", delta) + " hours'"
	}

	// limit the number of players if needed
	if limit > 0 {
		playersSQL += " LIMIT " + fmt.Sprintf("%d", limit)
	}

	// DEBUG
	// fmt.Println(playersSQL)

	rows, err := pp.db.Query(playersSQL)
	if err != nil {
		return nil, err
	}

	pids := make([]int, 0, 100)
	var pid int
	for rows.Next() {
		rows.Scan(&pid)
		pids = append(pids, pid)
	}

	return pids, nil
}
