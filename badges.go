package main

import (
	"flag"
	"fmt"
	"github.com/antzucaro/badges/config"
	"github.com/ungerik/go-cairo"
	"log"
	"os"
	"sync"
)

func renderWorker(pids <-chan int, wg *sync.WaitGroup, pp *PlayerDataFetcher, skins map[string]Skin,
	surfaceCache map[string]*cairo.Surface) {

	for pid := range pids {
		pd, err := pp.GetPlayerData(pid)
		if err != nil {
			fmt.Println(err)
		}

		if len(pd.Nick) == 0 {
			fmt.Printf("No data for player #%d!\n", pid)
		} else {
			for name, skin := range skins {
				skin.Render(pd, fmt.Sprintf("output/%s/%d.png", name, pid), surfaceCache)
			}
		}
	}

	wg.Done()
}

func main() {
	all := flag.Bool("all", false, "Generate badges for all ranked players")
	delta := flag.Int("delta", 6, "Generate for players having activity in this number of hours")
	pid := flag.Int("pid", -1, "Generate a badge for this player ID")
	limit := flag.Int("limit", -1, "Only generate badges for this many players")
	workers := flag.Int("workers", 5, "workers")
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

	skins := LoadSkins("skins")
	for name := range skins {
		err := os.MkdirAll(fmt.Sprintf("output/%s", name), os.FileMode(0755))
		if err != nil {
			fmt.Println(err)
		}
	}

	surfaceCache := LoadSurfaces(skins)

	pidsChan := make(chan int)

	var wg sync.WaitGroup

	// start workers
	for w := 1; w <= *workers; w++ {
		wg.Add(1)
		go renderWorker(pidsChan, &wg, pp, skins, surfaceCache)
	}

	// send them work
	for _, pid := range pids {
		pidsChan <- pid
	}
	close(pidsChan)

	// wait until they are all done
	wg.Wait()
}
