package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	// database connection string
	ConnStr string
}

var Path = "./config.json"
var Config = new(config)

func init() {
	// defaults
	Config.ConnStr = "user=xonstat host=localhost dbname=xonstatdb sslmode=disable"

	// set the config file path via environment variable
	if ecp := os.Getenv("BADGE_CONFIG"); ecp != "" {
		Path = ecp
	}

	file, err := os.Open(Path)
	if err != nil {
		if len(Path) > 1 {
			fmt.Printf("Error: could not read config file %s.\n", Path)
		}
		return
	}
	decoder := json.NewDecoder(file)

	// overwrite in-mem config with new values
	err = decoder.Decode(Config)
	if err != nil {
		fmt.Printf("Error decoding file %s\n%s\n", Path, err)
	}
}
