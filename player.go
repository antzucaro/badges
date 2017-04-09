package main

import (
	"database/sql"
	"fmt"
	"github.com/antzucaro/qstr"
	_ "github.com/lib/pq"
	"time"
)

// playerElo holds records coming from the player_elos table in stats
type playerElo struct {
	GameType string
	Elo      int64
}

// playerRank holds records coming from the player_ranks table in stats
type playerRank struct {
	GameType string
	Rank     int64
	MaxRank  int64
}

// PlayerData holds aggregate statistics for players
type PlayerData struct {
	Nick         qstr.QStr
	StrippedNick string
	Elos         []playerElo
	Ranks        []playerRank
	Kills        int
	Deaths       int
	Wins         int
	Losses       int
	PlayingTime  time.Duration
}

// KDRatio returns the player's Kill:Death ratio as a string
func (pd *PlayerData) KDRatio() float64 {
	if pd.Deaths > 0 {
		return float64(pd.Kills) / float64(pd.Deaths)
	} else {
		return 0.000
	}
}

// WinPct returns the player's win percentage as a string
func (pd *PlayerData) WinPct() float64 {
	totalGames := pd.Wins + pd.Losses
	if totalGames > 0 {
		return float64(pd.Wins) / float64(totalGames) * 100
	} else {
		return 0.00
	}
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros.  it omits the decimal
// point too when the fraction is 0.  It returns the index where the
// output bytes begin and the value v/10**prec.
// Modified from https://golang.org/src/time/time.go?s=15202:15235#L462
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)

	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}

	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
// Modified from https://golang.org/src/time/time.go?s=15202:15235#L462
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

// DurationString creates a human-readable duration string with a days component.
// Modified from https://golang.org/src/time/time.go?s=15202:15235#L462
func DurationString(d time.Duration) string {
	// Largest time is 2540400h10m10.000000000s
	var buf [32]byte

	w := len(buf)
	u := uint64(d)

	neg := d < 0
	if neg {
		u = -u
	}

	w--
	buf[w] = 's'
	w, u = fmtFrac(buf[:w], u, 9)

	// u is now integer seconds
	w = fmtInt(buf[:w], u%60)
	u /= 60

	// u is now integer minutes
	if u > 0 {
		w--
		buf[w] = 'm'
		w = fmtInt(buf[:w], u%60)
		u /= 60
		// u is now integer hours

		// Stop at hours because days can be different lengths.
		if u > 0 {
			w--
			buf[w] = 'h'
			w = fmtInt(buf[:w], u%24)
			u /= 24

			// u is now integer days
			if u > 0 {
				w--
				buf[w] = 'd'
				w = fmtInt(buf[:w], u)
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

// PlayingTime constructs a human-readable duration string with a day component.
func (pd *PlayerData) PlayingTimeString() string {
	return DurationString(pd.PlayingTime)
}

// PlayerDataFetcher fetches player information from the database
type PlayerDataFetcher struct {
	db *sql.DB
}

// NewPlayerDataFetcher creates a new PlayerDataFetcher for obtaining
// player information from the database
func NewPlayerDataFetcher(connStr string) (*PlayerDataFetcher, error) {
	// establish a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// connection pooling
	db.SetMaxIdleConns(5)

	pp := new(PlayerDataFetcher)
	pp.db = db
	return pp, nil
}

// FindPlayers finds a list of player_id values according to certain criteria.
// If delta is set, it will look for players who have had activity in the last
// $delta hours. If limit is set, the total number of player_ids returned is
// limited to that amount.
func (pp *PlayerDataFetcher) FindPlayers(delta int, limit int) ([]int, error) {
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

// initPlayerDataStmt generates the SQL statement string used to fetch
// the information used to populate PlayerData objects
func (pp *PlayerDataFetcher) genPlayerDataStmt(playerID int) string {
	query := `SELECT
   p.nick,
   p.stripped_nick,
   UPPER(agg_stats.game_type_cd) game_type_cd,
   ROUND(pe.elo) elo,
   pr.rank,
   pr.max_rank,
   SUM(win) wins,
   SUM(loss) losses,
   SUM(kills) kills,
   SUM(deaths) deaths,
   round(Sum(alivetime)/60) alivetime    
FROM
   (SELECT
      pgs.player_id,
      g.game_id,
      g.game_type_cd,
      CASE                      
         WHEN g.winner = pgs.team THEN 1                      
         WHEN pgs.scoreboardpos = 1 THEN 1                      
         ELSE 0                    
      END win,
      CASE                      
         WHEN g.winner = pgs.team THEN 0                      
         WHEN pgs.scoreboardpos = 1 THEN 0                      
         ELSE 1                    
      END loss,
      pgs.kills,
      pgs.deaths,
      extract(epoch 
   from
      pgs.alivetime) alivetime            
   FROM
      games g,
      player_game_stats pgs             
   WHERE
      g.game_id = pgs.game_id             
      AND pgs.player_id = %d
      AND g.players @> ARRAY[%d]) agg_stats
JOIN
   players p 
      on p.player_id = agg_stats.player_id            
JOIN
   player_elos pe 
      on agg_stats.game_type_cd = pe.game_type_cd 
      and pe.player_id = agg_stats.player_id            
LEFT OUTER JOIN
   (
      select
         pr.game_type_cd,
         pr.rank,
         overall.max_rank                   
      from
         player_ranks pr,
         (select
            game_type_cd,
            max(rank) max_rank                       
         from
            player_ranks                        
         group by
            game_type_cd) overall                   
      where
         pr.game_type_cd = overall.game_type_cd                    
         and max_rank > 1                   
         and player_id = %d
      ) pr 
         on pr.game_type_cd = pe.game_type_cd            
   GROUP BY
      p.nick,
      p.stripped_nick,
      agg_stats.game_type_cd,
      pe.elo,
      pr.rank,
      pr.max_rank            
   ORDER BY
      pe.elo desc NULLS LAST
   LIMIT 3`

	return fmt.Sprintf(query, playerID, playerID, playerID)
}

// GetPlayerData retrieves player information for the given player_id
func (pp *PlayerDataFetcher) GetPlayerData(playerID int) (*PlayerData, error) {
	sqlQuery := pp.genPlayerDataStmt(playerID)

	rows, err := pp.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}

	pd := new(PlayerData)

	filled := false
	var nick, strippedNick, gameType string
	var wins, losses, kills, deaths, alivetime int
	var elo, rank, maxRank sql.NullInt64
	var totalWins, totalLosses, totalKills, totalDeaths, totalAlivetime int
	elos := make([]playerElo, 0, 5)
	ranks := make([]playerRank, 0, 5)

	for rows.Next() {
		err := rows.Scan(&nick, &strippedNick, &gameType, &elo, &rank, &maxRank, &wins, &losses, &kills, &deaths, &alivetime)
		if err != nil {
			panic(err)
		}

		// did we fill in the player information yet?
		if !filled {
			pd.Nick = qstr.QStr(nick)
			pd.StrippedNick = strippedNick
			filled = true
		}

		// elo and rank are outer joins, thus may be NULL
		if elo.Valid {
			elos = append(elos, playerElo{GameType: gameType, Elo: elo.Int64})
		}
		if rank.Valid && maxRank.Valid {
			ranks = append(ranks, playerRank{GameType: gameType, Rank: rank.Int64, MaxRank: maxRank.Int64})
		}

		totalWins += wins
		totalLosses += losses
		totalKills += kills
		totalDeaths += deaths
		totalAlivetime += alivetime
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	pd.Elos = elos
	pd.Ranks = ranks
	pd.Kills = totalKills
	pd.Deaths = totalDeaths
	pd.Wins = totalWins
	pd.Losses = totalLosses
	pd.PlayingTime = time.Duration(totalAlivetime) * time.Minute

	return pd, nil
}
