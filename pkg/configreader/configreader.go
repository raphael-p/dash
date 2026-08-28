package configreader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ReadConfigFile[T *Q, Q any](dataDir string, values T) {
	configPath := filepath.Join(dataDir, "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to open config file: %s", err))
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(values); err != nil {
		panic(fmt.Sprintf("could not parse config file: %s", err))
	}
}
