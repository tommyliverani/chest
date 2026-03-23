package jewel

import (
	"encoding/json"
	"fmt"
)

type password struct {
	baseJewel
	Password string `json:"password"`
}

const PASSWORD = "pass"

func (p *password) GetEmoji() string {
	return "🔑"
}

func ParsePassword(data []byte) (*password, error) {
	var p password

	// json.Unmarshal riempie la struct p con i dati presenti in data
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, fmt.Errorf("errore nel parsing della password: %w", err)
	}

	return &p, nil
}
