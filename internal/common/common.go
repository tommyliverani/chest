package common

//ok

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func CheckWithMsg(msg string, err error) {
	if err != nil {
		if msg != "" {
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m", msg)
		}
		fmt.Fprintf(os.Stderr, "\033[31m%v\033[0m\n", err)
		os.Exit(1)
	}
}

func PrintErrorAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", msg)
	}
	os.Exit(1)
}

func Check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m%v\033[0m\n", err)
		os.Exit(1)
	}
}

func GetExistingChestIds() ([]string, error) {
	entries, err := os.ReadDir(GetChestHome())
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			id := strings.TrimSuffix(entry.Name(), ".json")
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func GetChestJsonById(chestId string) (json.RawMessage, error) {
	chestPath := filepath.Join(GetChestHome(), chestId+".json")
	chestFile, errRead := os.ReadFile(chestPath)
	if errRead != nil {
		return nil, errRead
	}
	return json.RawMessage(chestFile), nil
}

func DeleteChestJsonById(chestId string) error {
	if strings.ContainsAny(chestId, "/\\") || chestId == ".." {
		return fmt.Errorf("invalid chest ID: %q", chestId)
	}
	if err := os.Remove(filepath.Join(GetChestHome(), chestId) + ".json"); err != nil {
		return err
	}
	return nil
}

func GetJsonSessions() (json.RawMessage, error) {
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
	if err := os.WriteFile(getChestSessionFilePath(), newSession, 0600); err != nil {
		return err
	}
	return nil
}

func CreateOrUpdateJsonChestFile(chestId string, jsonData json.RawMessage) (string, error) {
	chestPath := filepath.Join(GetChestHome(), chestId+".json")
	if err := os.WriteFile(chestPath, jsonData, 0600); err != nil {
		return "", err
	}
	return chestPath, nil
}

func GenerateChestID() string {
	now := time.Now().UnixNano()
	return strconv.FormatInt(now, 36)
}

func GetKindFromJson(data json.RawMessage) (string, error) {
	var helper struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(data, &helper); err != nil {
		return "", err
	}
	return helper.Kind, nil
}

func GetNameFromJson(data json.RawMessage) (string, error) {
	var helper struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &helper); err != nil {
		return "", err
	}
	return helper.Name, nil
}
