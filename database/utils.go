package database

import (
	"strings"

	"github.com/raphael-p/datashard/logger"
)

func LazyInit(err error) bool {
	if err != nil && strings.Contains(err.Error(), "no such table") {
		logger.Warning("no database found, intialising")
		Initialise()
		return true
	}
	return false
}
