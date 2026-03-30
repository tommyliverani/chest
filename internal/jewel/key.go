package jewel

import (
	"chest/internal/common"
	"chest/internal/factory"
	"encoding/json"
	"fmt"
)

type Key struct {
	baseJewel
	Key string `json:"key"`
}

const KEY = "key"

func (p *Key) GetEmoji() string {
	return "🗝️ "
}

func (j *Key) ToJson() (json.RawMessage, error) { return json.Marshal(j) }

func ParseKey(data json.RawMessage) (*Key, error) {
	var p Key
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func CreateKey(name string, description string) (*Key, error) {
	key := common.ReadSecret("Insert the key: ")
	return &Key{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        KEY,
			Description: description,
		},
		Key: key,
	}, nil
}

func init() {
	factory.RegisterJewelCreator(KEY, func(name string, description string) (factory.Jewel, error) {
		return CreateKey(name, description)
	})
	factory.RegisterJewelParser(KEY, func(data json.RawMessage) (factory.Jewel, error) {
		return ParseKey(data)
	})
	factory.RegisterJewelHelp(KEY, factory.JewelHelp{
		Emoji:    "🗝️ ",
		Name:     "key",
		Short:    "stores a secret key or password",
		Behavior: "copies the key value to the clipboard",
		Operations: map[string]string{
			"add":   "saves a new secret key into the open chest",
			"ls":    "lists all keys in the open chest",
			"edit":  "edits name, description or key value",
			"rm":    "removes a key from the open chest",
			"print": "shows the key value after confirmation",
			"copy":  "copies the key value to the clipboard",
		},
	})
}

func (m *Key) Edit() error {
	jewelField := common.SelectField("which field do you want to edit?", []string{"name", "description", "key"})
	switch jewelField {
	case "name":
		newName := common.ReadField("Insert new name: ")
		m.Name = newName
	case "description":
		newDescription := common.ReadField("Insert new description: ")
		m.Description = newDescription
	case "key":
		newKey := common.ReadSecret("Insert new key: ")
		m.Key = newKey
	}
	return nil
}

func (p *Key) Print() {
	confirm := common.SelectField(fmt.Sprintf("Are you sure you want to print the key '%s'?", p.Name), []string{"No", "Yes"})
	if confirm != "Yes" {
		return
	}
	fmt.Println(p.Key)
}

func (p *Key) Copy() {
	common.WriteToClipboard(p.Key)
}

func (p *Key) Use() {
	p.Copy()
}
