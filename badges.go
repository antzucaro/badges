package main

import (
	"flag"
	"fmt"
	"github.com/antzucaro/badges/config"
	"log"
)

func main() {
	all := flag.Bool("all", false, "Generate badges for all ranked players")
	delta := flag.Int("delta", 6, "Generate for players having activity in this number of hours")
	pid := flag.Int("pid", -1, "Generate a badge for this player ID")
	limit := flag.Int("limit", -1, "Only generate badges for this many players")
	// verbose := flag.Bool("verbose", false, "Turn on verbose output and timings")
	flag.Parse()

	pp, err := NewPlayerProcessor(config.Config.ConnStr)
	if err != nil {
		log.Fatal(err)
	}

	// negate the delta value if we want all players generated
	if *all {
		*delta = -1
	}

	if *pid == -1 {
        // Find players by delta, with an optional limit
		pids, err := pp.FindPlayers(*delta, *limit)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range pids {
			pd, err := pp.GetPlayerData(v)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%+v\n", *pd)
		}
	} else {
        // Find a specific player by ID
        pd, err := pp.GetPlayerData(*pid)
        if err != nil {
			log.Fatal(err)
        }
        fmt.Printf("%+v\n", *pd)
    }
}
