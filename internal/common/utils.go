package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetChestFilePaths() []string {
	entries, err := os.ReadDir(ChestBasePath)
	if err != nil {
		fmt.Printf("failed to read directory %s: %v\n", ChestBasePath, err)
		return []string{}
	}
	var jsonFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			fullPath := filepath.Join(ChestBasePath, entry.Name())
			jsonFiles = append(jsonFiles, fullPath)
		}
	}
	return jsonFiles
}

func GetChestNames() []string {
	entries, err := os.ReadDir(ChestBasePath)
	if err != nil {
		fmt.Printf("failed to read directory %s: %v\n", ChestBasePath, err)
		return []string{}
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			name := strings.TrimSuffix(entry.Name(), ".json")
			names = append(names, name)
		}
	}
	return names
}

func DeleteChestFile(chestname string) error {
	return os.Remove(filepath.Join(ChestBasePath, chestname) + ".json")
}
