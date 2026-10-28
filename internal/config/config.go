package config

import (
	"time"

	"github.com/raphael-p/dash/pkg/configreader"
)

type fileConfig struct {
	DashDurationSeconds  uint16 `json:"dash_duration_seconds"`
	DescriptionCharLimit int    `json:"description_char_limit"`
	NameCharLimit        int    `json:"name_char_limit"`
}

type AppConfig struct {
	DashDuration         time.Duration
	DescriptionCharLimit int
	NameCharLimit        int
}

var C *AppConfig = &AppConfig{}

func Initialise(dataDir string) {
	var rawConfig *fileConfig = &fileConfig{}
	configreader.ReadConfigFile(dataDir, rawConfig)
	dashDuration := time.Duration(rawConfig.DashDurationSeconds) * time.Second
	C = &AppConfig{dashDuration, rawConfig.DescriptionCharLimit, rawConfig.NameCharLimit}
}
