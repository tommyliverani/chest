package common

import (
	"fmt"
	"os"
	"path/filepath"
)

var defaultSessionDir = fmt.Sprintf("/run/user/%d", os.Getuid())

func GetChestHome() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("cannot determine home directory: %v", err))
	}
	return getEnv("CHEST_HOME", filepath.Join(homeDir, ".chest"))
}

func getChestSessionFilePath() string {
	chestSessionDir := getEnv("CHEST_SESSION_DIR_PATH", filepath.Join(defaultSessionDir, ".chest_session.json"))
	return getEnv("CHEST_SESSION_FILE_PATH", chestSessionDir)
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
