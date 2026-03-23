package chest

import (
	. "chest/internal/common"
	"encoding/json"
	"fmt"
)

const MOCK string = "mock"

type MockChest struct {
	baseChest
	Jewels []json.RawMessage `json:"jewels"`
}

func (m *MockChest) GetEmoji() string {
	return "📦"
}

func (m *MockChest) Delete() error {
	return DeleteChestFile(m.GetName())
}

func (m *MockChest) Edit() error {
	fmt.Printf("whitch field do you want to edit?")
	chestField, err := SelectChestField([]string{"name", "description"})
	if err != nil {
		return err
	}
	if chestField == "name" {
		newName, err := ReadChestName()
		if err != nil {
			return err
		}
		m.Name = newName
	}
	if chestField == "description" {
		newDescription, err := ReadChestDescription()
		if err != nil {
			return err
		}
		m.Description = newDescription
	}
	return nil
}

func CreateMockChest(name string, description string) (Chest, error) {
	return &MockChest{
		baseChest: baseChest{
			Name:        name,
			Kind:        MOCK,
			Description: description,
		},
		Jewels: []json.RawMessage{},
	}, nil
}

func ParseMockChest(data json.RawMessage) (Chest, error) {
	var mc MockChest
	err := json.Unmarshal(data, &mc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MockChest: %w", err)
	}
	return &mc, nil
}

func init() {
	RegisterChestCreator(MOCK, CreateMockChest)
	RegisterChestParser(MOCK, ParseMockChest)
}
