package strutil

import (
	"strings"
)

func LineCount(s string) int {
	if len(s) == 0 {
		return 0
	}
	return strings.Count(s, "\n") + 1
}
