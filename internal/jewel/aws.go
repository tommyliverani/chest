package jewel

import (
	"chest/internal/common"
	"chest/internal/factory"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// da testare
const AWS_KIND = "aws"

type Aws struct {
	baseJewel
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

func (a *Aws) GetEmoji() string { return "☁️ " }

func (a *Aws) ToJson() (json.RawMessage, error) { return json.Marshal(a) } //nolint:gosec

func ParseAws(data json.RawMessage) (*Aws, error) {
	var aws Aws
	if err := json.Unmarshal(data, &aws); err != nil {
		return nil, err
	}
	return &aws, nil
}

func CreateAws(name string, description string) (*Aws, error) {
	accessKeyId := common.ReadField("Insert Access Key ID: ")
	secretAccessKey := common.ReadSecret("Insert Secret Access Key: ")
	return &Aws{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        AWS_KIND,
			Description: description,
		},
		AccessKeyId:     accessKeyId,
		SecretAccessKey: secretAccessKey,
	}, nil
}

func init() {
	factory.RegisterJewelCreator(AWS_KIND, func(name string, description string) (factory.Jewel, error) {
		return CreateAws(name, description)
	})
	factory.RegisterJewelParser(AWS_KIND, func(data json.RawMessage) (factory.Jewel, error) {
		return ParseAws(data)
	})
	factory.RegisterJewelHelp(AWS_KIND, factory.JewelHelp{
		Emoji:    "☁️ ",
		Name:     "aws",
		Short:    "stores AWS Access Key ID and Secret Access Key",
		Behavior: "writes credentials to ~/.aws/credentials and verifies with sts get-caller-identity",
		Operations: map[string]string{
			"add":   "saves new AWS credentials into the open chest",
			"ls":    "lists all AWS credential entries in the open chest",
			"edit":  "edits name, description, access key ID or secret",
			"rm":    "removes an AWS credential entry from the open chest",
			"print": "shows export AWS_* commands with masked secret",
			"copy":  "copies the full export AWS_* commands to the clipboard",
		},
	})
}

func (a *Aws) Edit() error {
	field := common.SelectField("which field do you want to edit?", []string{"name", "description", "access_key_id", "secret_access_key"})
	switch field {
	case "name":
		a.Name = common.ReadField("Insert new name: ")
	case "description":
		a.Description = common.ReadField("Insert new description: ")
	case "access_key_id":
		a.AccessKeyId = common.ReadField("Insert new Access Key ID: ")
	case "secret_access_key":
		a.SecretAccessKey = common.ReadSecret("Insert new Secret Access Key: ")
	}
	return nil
}

func (a *Aws) exportCommand() string {
	return fmt.Sprintf("export AWS_ACCESS_KEY_ID=%s && export AWS_SECRET_ACCESS_KEY=%s", a.AccessKeyId, a.SecretAccessKey)
}

func (a *Aws) Print() {
	masked := strings.Repeat("*", len(a.SecretAccessKey))
	fmt.Printf("export AWS_ACCESS_KEY_ID=%s && export AWS_SECRET_ACCESS_KEY=%s\n", a.AccessKeyId, masked)
}

func (a *Aws) Copy() {
	common.WriteToClipboard(a.exportCommand())
}

// Use writes the credentials to ~/.aws/credentials under [default] profile
// and runs `aws configure` to confirm the profile is active.
func (a *Aws) Use() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		common.Check(err)
		return
	}
	credPath := filepath.Join(homeDir, ".aws", "credentials")
	if err := os.MkdirAll(filepath.Dir(credPath), 0700); err != nil { //nolint:gosec
		common.Check(err)
		return
	}
	content := fmt.Sprintf("[default]\naws_access_key_id = %s\naws_secret_access_key = %s\n",
		a.AccessKeyId, a.SecretAccessKey)
	if err := os.WriteFile(credPath, []byte(content), 0600); err != nil { //nolint:gosec
		common.Check(err)
		return
	}
	fmt.Printf("Credentials written to %s\n", credPath)
	cmd := exec.Command("aws", "sts", "get-caller-identity")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "aws sts get-caller-identity failed: %v\n", err)
	}
}
