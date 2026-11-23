package markerfile

import (
	"os"
	"path/filepath"
	"time"
)

func getMarker(directory, name string) string {
	return filepath.Join(directory, name)
}

func DebounceIsActive(directory, name string, debounceDuration time.Duration) bool {
	marker := getMarker(directory, name)
	if info, err := os.Stat(marker); err == nil && time.Since(info.ModTime()) < debounceDuration {
		return true
	}
	if f, err := os.OpenFile(marker, os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
		f.Close()
		_ = os.Chtimes(marker, time.Now(), time.Now())
	}
	return false
}
