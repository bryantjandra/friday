package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	APIKey string
	DBPath string
}

func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		APIKey: os.Getenv("ANTHROPIC_API_KEY"),

		/* Unix convention for storing config data for a program
		   Dot folder is used to not clutter home directory
		*/
		DBPath: filepath.Join(home, ".friday", "friday.db"),
	}

	return cfg, nil

}
