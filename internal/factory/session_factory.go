package factory

import (
	"chest/internal/common"
	"encoding/json"
	"fmt"
)

func getSessions() (SessionMap, error) {
	sessionData, err := common.GetJsonSession()
	if err != nil {
		return nil, fmt.Errorf("error reading session data: %w", err)
	}
	var sessions SessionMap
	if err := json.Unmarshal(sessionData, &sessions); err != nil {
		return nil, fmt.Errorf("error unmarshalling session data: %w", err)
	}
	return sessions, nil
}

func saveSessions(sessions SessionMap) error {
	sessionData, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling session data: %w", err)
	}
	return common.UpdateExistingJsonSession(sessionData)
}

func GetSession(chestName string) (Session, bool, error) {
	sessions, err := getSessions()
	if err != nil {
		return Session{}, false, err
	}
	session, exists := sessions[chestName]
	return session, exists, nil
}

func GetKeyJewel(chestName string) (json.RawMessage, bool, error) {
	session, exists, err := GetSession(chestName)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	return session.Jewel, true, nil
}

func StoreSession(chestName string, jewel json.RawMessage) error {
	sessions, err := getSessions()
	if err != nil {
		return err
	}
	if sessions == nil {
		sessions = make(SessionMap)
	}
	sessions[chestName] = Session{Jewel: jewel}
	return saveSessions(sessions)
}

func DeleteSession(chestName string) error {
	sessions, err := getSessions()
	if err != nil {
		return err
	}
	delete(sessions, chestName)
	return saveSessions(sessions)
}

func IsOpen(chestName string) bool {
	sessions, err := getSessions()
	if err != nil {
		return false
	}
	_, exists := sessions[chestName]
	return exists
}

func GetOpenChestNames() ([]string, error) {
	sessions, err := getSessions()
	if err != nil {
		return nil, err
	}
	openChests := make([]string, 0, len(sessions))
	for chestName := range sessions {
		openChests = append(openChests, chestName)
	}
	return openChests, nil
}

func GetOpenChests() ([]Chest, error) {
	sessions, err := getSessions()
	if err != nil {
		return nil, err
	}
	openChests := make([]Chest, 0, len(sessions))
	for chestName := range sessions {
		chest, err := GetExistingChest(chestName)
		if err != nil {
			return nil, fmt.Errorf("error retrieving chest '%s': %w", chestName, err)
		}
		openChests = append(openChests, chest)
	}
	return openChests, nil
}

func SelectOpenChest(prompt string) (Chest, error) {
	openChests, err := GetOpenChests()
	if err != nil {
		return nil, err
	}
	openChestsStrings := make([]string, len(openChests))
	for i, chest := range openChests {
		openChestsStrings[i] = GetChestString(chest)
	}
	index, _, err := common.SelectFieldWithIndex(prompt, openChestsStrings)
	if err != nil {
		return nil, err
	}
	return openChests[index], nil
}
