package tcellutil

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// ColourCodeByAge converts a time to a colour scale
// to visually represent the passing of time
func ColourCodeByAge(t time.Time) tcell.Color {
	d := time.Since(t)
	day := 24 * time.Hour
	switch {
	case d < day:
		return tcell.Color220
	case d < day*7:
		return tcell.Color214
	case d < day*15:
		return tcell.Color208
	case d < day*30:
		return tcell.Color202
	default:
		return tcell.Color160
	}
}
