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

	pp, err := NewPlayerDataFetcher(config.Config.ConnStr)
	if err != nil {
		log.Fatal(err)
	}

	// negate the delta value if we want all players generated
	if *all {
		*delta = -1
	}

	var pids []int
	if *pid == -1 {
		// Find players by delta, with an optional limit
		pids, err = pp.FindPlayers(*delta, *limit)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Use just the one player ID for fetching
		pids = []int{*pid}
	}

	for _, pid := range pids {
		pd, err := pp.GetPlayerData(pid)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Rendering image for player #%d\n", pid)

		filename := fmt.Sprintf("%d.png", pid)
		ArcherSkin.Render(pd, filename)
	}
}
