package jewel

import (
	"chest/internal/common"
	"chest/internal/factory"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TODO: DA TESTARE

const OPENSHIFT_KIND = "oc"

type Openshift struct {
	baseJewel
	ApiUrl   string `json:"api_url"`
	ApiToken string `json:"api_token"`
}

func (o *Openshift) GetEmoji() string { return "♦️" }

func (o *Openshift) ToJson() (json.RawMessage, error) { return json.Marshal(o) }

func ParseOpenshift(data json.RawMessage) (*Openshift, error) {
	var oc Openshift
	if err := json.Unmarshal(data, &oc); err != nil {
		return nil, err
	}
	return &oc, nil
}

func CreateOpenshift(name string, description string) (*Openshift, error) {
	apiUrl := common.ReadField("Insert API URL: ")
	apiToken := common.ReadSecret("Insert API token: ")
	return &Openshift{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        OPENSHIFT_KIND,
			Description: description,
		},
		ApiUrl:   apiUrl,
		ApiToken: apiToken,
	}, nil
}

func init() {
	factory.RegisterJewelCreator(OPENSHIFT_KIND, func(name string, description string) (factory.Jewel, error) {
		return CreateOpenshift(name, description)
	})
	factory.RegisterJewelParser(OPENSHIFT_KIND, func(data json.RawMessage) (factory.Jewel, error) {
		return ParseOpenshift(data)
	})
	factory.RegisterJewelHelp(OPENSHIFT_KIND, factory.JewelHelp{
		Emoji:    "♦️ ",
		Name:     "oc",
		Short:    "stores OpenShift API URL and token for oc login",
		Behavior: "runs oc login with the stored credentials",
		Operations: map[string]string{
			"add":   "saves a new OpenShift credential into the open chest",
			"ls":    "lists all OpenShift entries in the open chest",
			"edit":  "edits name, description, API URL or token",
			"rm":    "removes an OpenShift entry from the open chest",
			"print": "shows oc login command with masked token",
			"copy":  "copies the full oc login command to the clipboard",
		},
	})
}

func (o *Openshift) Edit() error {
	field := common.SelectField("which field do you want to edit?", []string{"name", "description", "api_url", "api_token"})
	switch field {
	case "name":
		o.Name = common.ReadField("Insert new name: ")
	case "description":
		o.Description = common.ReadField("Insert new description: ")
	case "api_url":
		o.ApiUrl = common.ReadField("Insert new API URL: ")
	case "api_token":
		o.ApiToken = common.ReadSecret("Insert new API token: ")
	}
	return nil
}

func (o *Openshift) loginCommand() string {
	return fmt.Sprintf("oc login %s --token=%s", o.ApiUrl, o.ApiToken)
}

func (o *Openshift) Print() {
	masked := strings.Repeat("*", len(o.ApiToken))
	fmt.Printf("oc login %s --token=%s\n", o.ApiUrl, masked)
}

func (o *Openshift) Copy() {
	common.WriteToClipboard(o.loginCommand())
}

func (o *Openshift) Use() {
	cmd := exec.Command("oc", "login", o.ApiUrl, "--token="+o.ApiToken)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "oc login failed: %v\n", err)
	}
}
