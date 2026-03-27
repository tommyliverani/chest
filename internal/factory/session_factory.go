package factory

import (
	"chest/internal/common"
	"encoding/json"
	"maps"
	"slices"
)

//ok

func getSessions() SessionMap {
	sessionData, err := common.GetJsonSessions()
	common.CheckWithMsg("Failed to read sessions", err)
	var sessions SessionMap
	common.CheckWithMsg("Failed to unmarshal sessions", json.Unmarshal(sessionData, &sessions))
	return sessions
}

func saveSessions(sessions SessionMap) {
	sessionData, err := json.MarshalIndent(sessions, "", "  ")
	common.CheckWithMsg("Failed to marshal sessions", err)
	common.CheckWithMsg("Failed to update sessions", common.UpdateExistingJsonSession(sessionData))
}

func GetKeyJewelFromSession(chestId string) (Jewel, bool) {
	session, isOpen := getSessions()[chestId]
	if !isOpen {
		return nil, false
	}
	jewel, err := ParseJewel(session.KeyJewel)
	if err != nil {
		return nil, false
	}
	return jewel, true
}

func IsOpen(chestId string) bool {
	_, isOpen := getSessions()[chestId]
	return isOpen
}

func StoreSession(chestId string, jewel Jewel) {
	sessions := getSessions()
	if sessions == nil {
		sessions = make(SessionMap)
	}
	jsonJewel, err := jewel.ToJson()
	common.CheckWithMsg("Failed to convert jewel to JSON", err)
	sessions[chestId] = Session{KeyJewel: jsonJewel}
	saveSessions(sessions)
}

func DeleteSession(chestId string) {
	sessions := getSessions()
	if sessions == nil {
		return
	}
	delete(sessions, chestId)
	saveSessions(sessions)
}

func GetOpenChestIds() []string {
	return slices.Collect(maps.Keys(getSessions()))
}
