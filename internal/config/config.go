package config

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/raphael-p/dash/pkg/configreader"
)

type fileConfig struct {
	DashDurationSeconds  uint16 `json:"dash_duration_seconds"`
	DescriptionCharLimit int    `json:"description_char_limit"`
	NameCharLimit        int    `json:"name_char_limit"`
}

type AppConfig struct {
	DashDuration                     time.Duration
	DescriptionCharLimit             int
	NameCharLimit                    int
	HighlightColour, SecondaryColour tcell.Color
}

var C *AppConfig = &AppConfig{}

func Initialise(dataDir string, highlightColour, secondaryColour tcell.Color) {
	rawConfig := &fileConfig{}
	configreader.ReadConfigFile(dataDir, rawConfig)
	dashDuration := time.Duration(rawConfig.DashDurationSeconds) * time.Second
	C = &AppConfig{dashDuration, rawConfig.DescriptionCharLimit, rawConfig.NameCharLimit, highlightColour, secondaryColour}
}

type InputMode int

const (
	DescriptionInput InputMode = iota
	KeybindInput
	NameInput
	TaskIDInput
)

func GetCharLimit(mode InputMode) (charLimit int) {
	switch mode {
	case DescriptionInput:
		return C.DescriptionCharLimit
	case KeybindInput:
		return 1
	case NameInput:
		return C.NameCharLimit
	case TaskIDInput:
		return 9
	default:
		return 0
	}
}

func GetPlaceholder(mode InputMode) string {
	switch mode {
	case DescriptionInput:
		return "enter a task description"
	case NameInput:
		return "enter a task name"
	case TaskIDInput:
		return "enter a task id"
	default:
		return ""
	}
}
