package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/antzucaro/badges/config"
	_ "github.com/lib/pq"
	"log"
)

// connect to the database, return the connection
func connect() *sql.DB {
	// establish a database connection
	db, err := sql.Open("postgres", config.Config.ConnStr)
	if err != nil {
		log.Fatal(err)
	}

	// connection pooling
	db.SetMaxIdleConns(5)

	return db
}

func findPlayers(db *sql.DB, delta int, limit int) []int {
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

	rows, err := db.Query(playersSQL)
	if err != nil {
		log.Fatal(err)
	}

	pids := make([]int, 0, 100)
	var pid int
	for rows.Next() {
		rows.Scan(&pid)
		pids = append(pids, pid)
	}

	return pids
}

func main() {
	all := flag.Bool("all", false, "Generate badges for all ranked players")
	delta := flag.Int("delta", 6, "Generate for players having activity in this number of hours")
	pid := flag.Int("pid", -1, "Generate a badge for this player ID")
	limit := flag.Int("limit", -1, "Only generate badges for this many players")
	// verbose := flag.Bool("verbose", false, "Turn on verbose output and timings")
	flag.Parse()

	db := connect()
	defer db.Close()

	// negate the delta value if we want all players generated
	if *all {
		*delta = -1
	}

	if *pid == -1 {
		pids := findPlayers(db, *delta, *limit)

		for _, v := range pids {
			fmt.Println(v)
		}
	}
}
