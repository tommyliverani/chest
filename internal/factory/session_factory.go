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
