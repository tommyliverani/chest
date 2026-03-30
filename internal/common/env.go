package common

import (
	"fmt"
	"os"
	"path/filepath"
)

var defaultSessionDir = fmt.Sprintf("/run/user/%d", os.Getuid())

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0750)
}

func init() {
	if err := ensureDir(GetChestHome()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func GetChestHome() string {
	homeDir, err := os.UserHomeDir()
	CheckWithMsg("Failed to get user home directory", err)
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
