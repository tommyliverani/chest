package common

import (
	"os"
	"path/filepath"
)

var homeDir, _ = os.UserHomeDir()
var ChestBasePath = getEnv("CHEST_PATHS", filepath.Join(homeDir, ".chest"))
var ChestPaths = GetChestFilePaths()

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
