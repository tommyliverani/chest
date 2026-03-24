package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetExistingChestPaths() ([]string, error) {
	entries, err := os.ReadDir(GetChestHome())
	if err != nil {
		return nil, err
	}
	var jsonFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			fullPath := filepath.Join(GetChestHome(), entry.Name())
			jsonFiles = append(jsonFiles, fullPath)
		}
	}
	return jsonFiles, nil
}

func GetExistingChestNames() ([]string, error) {
	entries, err := os.ReadDir(GetChestHome())
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			name := strings.TrimSuffix(entry.Name(), ".json")
			names = append(names, name)
		}
	}
	return names, nil
}

func DeleteExistingChestFile(chestname string) error {
	if strings.ContainsAny(chestname, "/\\") || chestname == ".." {
		return fmt.Errorf("invalid chest name: %q", chestname)
	}
	return os.Remove(filepath.Join(GetChestHome(), chestname) + ".json")
}

func GetJsonSession() (json.RawMessage, error) {
	sessionData, err := os.ReadFile(getChestSessionFilePath())
	if os.IsNotExist(err) {
		return json.RawMessage("{}"), nil
	}
	if err != nil {
		return nil, err
	}
	return json.RawMessage(sessionData), nil
}

func UpdateExistingJsonSession(newSession json.RawMessage) error {
	return os.WriteFile(getChestSessionFilePath(), newSession, 0600)
}
