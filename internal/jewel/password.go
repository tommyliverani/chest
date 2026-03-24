package jewel

import (
	"chest/internal/common"
	"chest/internal/factory"
	"encoding/json"
	"fmt"
)

type Password struct {
	baseJewel
	Password string `json:"password"`
}

const PASSWORD = "psw"

func (p *Password) GetEmoji() string {
	return "🔑"
}

func (j *Password) ToJson() (json.RawMessage, error) { return json.Marshal(j) }

func ParsePassword(data json.RawMessage) (*Password, error) {
	var p Password
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, fmt.Errorf("error while parsing password: %w", err)
	}
	return &p, nil
}

func CreatePassword(name string, description string) (*Password, error) {
	password, err := common.ReadSecret("Insert password: ")
	if err != nil {
		return nil, err
	}
	return &Password{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        PASSWORD,
			Description: description,
		},
		Password: password,
	}, nil
}

func init() {
	factory.RegisterJewelCreator(PASSWORD, func(name string, description string) (factory.Jewel, error) {
		return CreatePassword(name, description)
	})
	factory.RegisterJewelParser(PASSWORD, func(data json.RawMessage) (factory.Jewel, error) {
		return ParsePassword(data)
	})
}
