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

const DOCKER_KIND = "docker"

// da testare

type DockerRegistry struct {
	baseJewel
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (d *DockerRegistry) GetEmoji() string { return "🐳" }

func (d *DockerRegistry) ToJson() (json.RawMessage, error) { return json.Marshal(d) } //nolint:gosec

func ParseDockerRegistry(data json.RawMessage) (*DockerRegistry, error) {
	var dr DockerRegistry
	if err := json.Unmarshal(data, &dr); err != nil {
		return nil, err
	}
	return &dr, nil
}

func CreateDockerRegistry(name string, description string) (*DockerRegistry, error) {
	url := common.ReadField("Insert registry URL: ")
	username := common.ReadField("Insert username: ")
	password := common.ReadSecret("Insert password: ")
	return &DockerRegistry{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        DOCKER_KIND,
			Description: description,
		},
		Url:      url,
		Username: username,
		Password: password,
	}, nil
}

func init() {
	factory.RegisterJewelCreator(DOCKER_KIND, func(name string, description string) (factory.Jewel, error) {
		return CreateDockerRegistry(name, description)
	})
	factory.RegisterJewelParser(DOCKER_KIND, func(data json.RawMessage) (factory.Jewel, error) {
		return ParseDockerRegistry(data)
	})
	factory.RegisterJewelHelp(DOCKER_KIND, factory.JewelHelp{
		Emoji:    "🐳",
		Name:     "docker",
		Short:    "stores Docker registry credentials for docker login",
		Behavior: "runs docker login with the stored credentials",
		Operations: map[string]string{
			"add":   "saves a new Docker registry credential into the open chest",
			"ls":    "lists all Docker registry entries in the open chest",
			"edit":  "edits name, description, URL, username or password",
			"rm":    "removes a Docker registry entry from the open chest",
			"print": "shows docker login command with masked password",
			"copy":  "copies the full docker login command to the clipboard",
		},
	})
}

func (d *DockerRegistry) Edit() error {
	field := common.SelectField("which field do you want to edit?", []string{"name", "description", "url", "username", "password"})
	switch field {
	case "name":
		d.Name = common.ReadField("Insert new name: ")
	case "description":
		d.Description = common.ReadField("Insert new description: ")
	case "url":
		d.Url = common.ReadField("Insert new registry URL: ")
	case "username":
		d.Username = common.ReadField("Insert new username: ")
	case "password":
		d.Password = common.ReadSecret("Insert new password: ")
	}
	return nil
}

func (d *DockerRegistry) loginCommand() string {
	return fmt.Sprintf("docker login %s -u %s -p %s", d.Url, d.Username, d.Password)
}

func (d *DockerRegistry) Print() {
	masked := strings.Repeat("*", len(d.Password))
	fmt.Printf("docker login %s -u %s -p %s\n", d.Url, d.Username, masked)
}

func (d *DockerRegistry) Copy() {
	common.WriteToClipboard(d.loginCommand())
}

func (d *DockerRegistry) Use() {
	cmd := exec.Command("docker", "login", d.Url, "-u", d.Username, "-p", d.Password) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "docker login failed: %v\n", err)
	}
}
