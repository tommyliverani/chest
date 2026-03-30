package jewel

import (
	"chest/internal/common"
	"chest/internal/factory"
	"encoding/json"
	"fmt"
	"strings"
)

const USERNAME_PASSWORD_KIND = "up"

type UsernameAndPassword struct {
	baseJewel
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *UsernameAndPassword) GetEmoji() string { return "💎" }

func (u *UsernameAndPassword) ToJson() (json.RawMessage, error) { return json.Marshal(u) } //nolint:gosec

func ParseUsernameAndPassword(data json.RawMessage) (*UsernameAndPassword, error) {
	var up UsernameAndPassword
	if err := json.Unmarshal(data, &up); err != nil {
		return nil, err
	}
	return &up, nil
}

func CreateUsernameAndPassword(name string, description string) (*UsernameAndPassword, error) {
	username := common.ReadField("Insert username: ")
	password := common.ReadSecret("Insert password: ")
	return &UsernameAndPassword{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        USERNAME_PASSWORD_KIND,
			Description: description,
		},
		Username: username,
		Password: password,
	}, nil
}

func init() {
	factory.RegisterJewelCreator(USERNAME_PASSWORD_KIND, func(name string, description string) (factory.Jewel, error) {
		return CreateUsernameAndPassword(name, description)
	})
	factory.RegisterJewelParser(USERNAME_PASSWORD_KIND, func(data json.RawMessage) (factory.Jewel, error) {
		return ParseUsernameAndPassword(data)
	})
	factory.RegisterJewelHelp(USERNAME_PASSWORD_KIND, factory.JewelHelp{
		Emoji:    "💎",
		Name:     "psw",
		Short:    "stores a username and password pair",
		Behavior: "copies the password to the clipboard",
		Operations: map[string]string{
			"add":   "saves a new username/password into the open chest",
			"ls":    "lists all username/password entries in the open chest",
			"edit":  "edits name, description, username or password",
			"rm":    "removes an entry from the open chest",
			"print": "shows username and masked password",
			"copy":  "copies the password to the clipboard",
		},
	})
}

func (u *UsernameAndPassword) Edit() error {
	field := common.SelectField("which field do you want to edit?", []string{"name", "description", "username", "password"})
	switch field {
	case "name":
		u.Name = common.ReadField("Insert new name: ")
	case "description":
		u.Description = common.ReadField("Insert new description: ")
	case "username":
		u.Username = common.ReadField("Insert new username: ")
	case "password":
		u.Password = common.ReadSecret("Insert new password: ")
	}
	return nil
}

func (u *UsernameAndPassword) Print() {
	masked := strings.Repeat("*", len(u.Password))
	fmt.Printf("Username: %s\nPassword: %s\n", u.Username, masked)
}

func (u *UsernameAndPassword) Copy() {
	common.WriteToClipboard(u.Password)
}

func (u *UsernameAndPassword) Use() {
	u.Copy()
}
