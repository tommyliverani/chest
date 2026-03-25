package chest

import (
	"chest/internal/common"
	"chest/internal/factory"
	"chest/internal/jewel"
	"encoding/json"
	"fmt"
)

const MOCK string = "mock"

type MockChest struct {
	baseChest
	Jewels []json.RawMessage `json:"jewels"`
}

func init() {
	factory.RegisterChestCreator(MOCK, CreateMockChest)
	factory.RegisterChestParser(MOCK, ParseMockChest)
}

func CreateMockChest(name string, description string) (factory.Chest, error) {
	return &MockChest{
		baseChest: baseChest{
			Name:        name,
			Kind:        MOCK,
			Description: description,
		},
		Jewels: []json.RawMessage{},
	}, nil
}

func (m *MockChest) Close() error {
	return nil
}

func ParseMockChest(data json.RawMessage) (factory.Chest, error) {
	var mc MockChest
	err := json.Unmarshal(data, &mc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MockChest: %w", err)
	}
	return &mc, nil
}

func (b *MockChest) ToJson() (json.RawMessage, error) { return json.Marshal(b) }

func (m *MockChest) GetEmoji() string { return "📦" }

func (m *MockChest) Delete() error { return common.DeleteExistingChestFile(m.GetName()) }

func (m *MockChest) Open() (factory.Jewel, error) {
	keyJewel, err := factory.CreateJewel(jewel.PASSWORD, m.GetName()+"KeyJewel", "key jewel to open "+m.GetName())
	if err != nil {
		return nil, fmt.Errorf("Error retriving password: %v\n", err)
	}
	//TODO: in the real chest we would check if the provided key jewel is correct and only return the content if it is
	return keyJewel, nil
}

func (m *MockChest) Edit() error {
	chestField, err := common.SelectField("which field do you want to edit?", []string{"description"})
	if err != nil {
		return err
	}
	// if chestField == "name" {
	// 	newName, err := common.ReadField("Insert new name: ")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	m.Name = newName
	// }
	if chestField == "description" {
		newDescription, err := common.ReadField("Insert new description: ")
		if err != nil {
			return err
		}
		m.Description = newDescription
	}
	return nil
}

func (m *MockChest) GetJewels(keyJewel json.RawMessage) ([]factory.Jewel, error) {
	// for the mock chest, we ignore the key jewel and return all jewels
	var jewels []factory.Jewel
	for _, raw := range m.Jewels {
		jewel, err := factory.ParseJewel(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse jewel: %w", err)
		}
		jewels = append(jewels, jewel)
	}
	return jewels, nil

}

func (m *MockChest) AddJewel(jewelToAdd factory.Jewel, keyJewel json.RawMessage) error {
	jewelBytes, err := jewelToAdd.ToJson()
	if err != nil {
		return fmt.Errorf("failed to marshal jewel: %w", err)
	}
	m.Jewels = append(m.Jewels, json.RawMessage(jewelBytes))
	return nil
}

func (m *MockChest) RemoveJewel(jewelName string, keyJewel json.RawMessage) error {
	for i, raw := range m.Jewels {
		var temp struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(raw, &temp); err != nil {
			continue
		}
		if temp.Name == jewelName {
			m.Jewels = append(m.Jewels[:i], m.Jewels[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("jewel with name '%s' not found in chest", jewelName)
}
