package stringpad

import "strings"

// RightPad up to the required space count.
// If the string is longer than the space count, 1 space is added.
func RightPad(str string, spaceCount uint8) string {
	str = strings.TrimSpace(str)
	padCount := int(spaceCount) - len(str)
	if padCount < 1 {
		return str + " "
	} else {
		return str + strings.Repeat(" ", padCount)
	}
}
