package main

import (
	"flag"
	"fmt"
	"github.com/antzucaro/badges/config"
	"github.com/ungerik/go-cairo"
	"log"
	"os"
	"runtime/pprof"
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

	skins := LoadSkins("skins")
	for name := range skins {
		err := os.MkdirAll(fmt.Sprintf("output/%s", name), os.FileMode(0755))
		if err != nil {
			fmt.Println(err)
		}
	}

	// CPU profiling
	f, err := os.Create("cpuprofile.dat")
	if err != nil {
		fmt.Println(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// cairo surface cache
	surfaceCache := make(map[string]*cairo.Surface)

	for _, pid := range pids {
		pd, err := pp.GetPlayerData(pid)
		if err != nil {
			fmt.Println(err)
		}

		if len(pd.Nick) == 0 {
			fmt.Printf("No data for player #%d!\n", pid)
		} else {
			fmt.Printf("Rendering images for player #%d\n", pid)
			for name, skin := range skins {
				skin.Render(pd, fmt.Sprintf("output/%s/%d.png", name, pid), surfaceCache)
			}
		}
	}
}
