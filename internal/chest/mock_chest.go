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
	Jewels []json.RawMessage `json:"jewels"` // jewelName -> Jewel json
}

// func init() {
// 	factory.RegisterChestCreator(MOCK, CreateMockChest)
// 	factory.RegisterChestParser(MOCK, ParseMockChest)
// }

func CreateMockChest(name string, description string) (factory.Chest, error) {
	return &MockChest{
		baseChest: baseChest{
			Id:          common.GenerateChestID(),
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

func (m *MockChest) GetKeyJewelKind() string {
	return jewel.KEY
}

func ParseMockChest(data json.RawMessage) (factory.Chest, error) {
	var mc MockChest
	err := json.Unmarshal(data, &mc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MockChest: %w", err)
	}
	return &mc, nil
}

func (b *MockChest) UpdateJewel(jewelName string, jewelToUpdate factory.Jewel, keyJewel factory.Jewel) error {
	jsonJewel, err := jewelToUpdate.ToJson()
	if err != nil {
		return err
	}
	for i, raw := range b.Jewels {
		name, err := common.GetNameFromJson(raw)
		if err != nil {
			return err
		}
		if name == jewelName {
			b.Jewels[i] = json.RawMessage(jsonJewel)
			return nil
		}
	}
	return fmt.Errorf("jewel with name '%s' not found in chest", jewelName)
}

func (b *MockChest) ToJson() (json.RawMessage, error) { return json.Marshal(b) }

func (m *MockChest) GetEmoji() string { return "📦" }

func (m *MockChest) Delete() error { return common.DeleteChestJsonById(m.GetName()) }

func (m *MockChest) Open(jewel factory.Jewel) error {
	//TODO: in the real chest we would check if the provided key jewel is correct and only return the content if it is
	return nil
}

func (m *MockChest) Edit(keyJewel factory.Jewel) error {
	chestField := common.SelectField("which field do you want to edit?", []string{"description", "name"})
	if chestField == "name" {
		newName := common.ReadField("Insert new name: ")
		m.Name = newName
	}
	if chestField == "description" {
		newDescription := common.ReadField("Insert new description: ")

		m.Description = newDescription
	}
	return nil
}

func (m *MockChest) GetJewels(keyJewel factory.Jewel) ([]factory.Jewel, error) {
	// for the mock chest, we ignore the key jewel and return all jewels
	var jewels []factory.Jewel
	for _, raw := range m.Jewels {
		j, err := factory.ParseJewel(raw)
		if err != nil {
			return nil, err
		}
		jewels = append(jewels, j)
	}
	return jewels, nil
}

func (m *MockChest) AddJewel(jewelToAdd factory.Jewel, keyJewel factory.Jewel) error {
	jewelBytes, err := jewelToAdd.ToJson()
	if err != nil {
		return fmt.Errorf("failed to marshal jewel: %w", err)
	}
	m.Jewels = append(m.Jewels, json.RawMessage(jewelBytes))
	return nil
}

func (m *MockChest) RemoveJewel(jewel factory.Jewel, keyJewel factory.Jewel) error {
	jewelName := jewel.GetName()
	for i, raw := range m.Jewels {
		name, err := common.GetNameFromJson(raw)
		if err != nil {
			return err
		}
		if name == jewelName {
			m.Jewels = append(m.Jewels[:i], m.Jewels[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("jewel with name '%s' not found in chest", jewelName)
}
