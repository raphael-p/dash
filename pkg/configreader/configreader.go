package configreader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func getExecDirectory(executableName, fallback string) string {
	if os.Args[0] == executableName {
		// if executed from compiled binary, use its directory
		ex, err := os.Executable()
		if err != nil {
			panic(fmt.Sprintf("failed to locate executable: %s", err))
		}
		return filepath.Dir(ex)
	} else {
		// otherwise, use fallback directory (relative path)
		return fallback
	}
}

func ReadConfigFile[T *Q, Q any](executableName, configRelativePath string, values T) {
	configPath := filepath.Join(getExecDirectory(executableName, configRelativePath), "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to open config file: %s", err))
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(values); err != nil {
		panic(fmt.Sprintf("could not parse config file: %s", err))
	}
}
