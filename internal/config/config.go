package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey  string
	DBPath  string
	BaseURL string
	Model   string
}

func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),

		/* Unix convention for storing config data for a program
		   Dot folder is used to not clutter home directory
		*/
		DBPath:  filepath.Join(home, ".friday", "friday.db"),
		BaseURL: "https://dashscope-intl.aliyuncs.com/apps/anthropic",
		Model:   "qwen3.7-flash",
	}

	/* fail early rather than calling the API with a missing API key */
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("DASHSCOPE_API_KEY environment variable is not set.")
	}

	return cfg, nil

}
