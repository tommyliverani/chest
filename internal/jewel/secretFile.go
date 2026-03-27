package jewel

import (
	"bufio"
	"chest/internal/common"
	"chest/internal/factory"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const SECRET_FILE_KIND = "file"

type SecretFile struct {
	baseJewel
	FileName string `json:"file_name"`
	Content  string `json:"content"` // base64-encoded file content
}

func (s *SecretFile) GetEmoji() string { return "📜" }

func (s *SecretFile) ToJson() (json.RawMessage, error) { return json.Marshal(s) }

func ParseSecretFile(data json.RawMessage) (*SecretFile, error) {
	var sf SecretFile
	if err := json.Unmarshal(data, &sf); err != nil {
		return nil, err
	}
	return &sf, nil
}

func CreateSecretFile(name string, description string) (*SecretFile, error) {
	fileName := common.ReadField("Insert file name (e.g. myfile.txt): ")
	fmt.Println("Paste the file content (press Ctrl+D when done):")
	content := readFileContent()
	return &SecretFile{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        SECRET_FILE_KIND,
			Description: description,
		},
		FileName: fileName,
		Content:  base64.StdEncoding.EncodeToString([]byte(content)),
	}, nil
}

func init() {
	factory.RegisterJewelCreator(SECRET_FILE_KIND, func(name string, description string) (factory.Jewel, error) {
		return CreateSecretFile(name, description)
	})
	factory.RegisterJewelParser(SECRET_FILE_KIND, func(data json.RawMessage) (factory.Jewel, error) {
		return ParseSecretFile(data)
	})
	factory.RegisterJewelHelp(SECRET_FILE_KIND, factory.JewelHelp{
		Emoji:    "📜",
		Name:     "file",
		Short:    "stores a secret file (content pasted via Ctrl+D)",
		Behavior: "saves the file to ~/Desktop or a path chosen by the user",
		Operations: map[string]string{
			"add":   "saves a new file (paste content, Ctrl+D to finish) into the open chest",
			"ls":    "lists all file entries in the open chest",
			"edit":  "edits name, description or file content",
			"rm":    "removes a file entry from the open chest",
			"print": "shows the file content after confirmation",
			"copy":  "saves the file to a path chosen by the user",
		},
	})
}

func (s *SecretFile) Edit() error {
	field := common.SelectField("which field do you want to edit?", []string{"name", "description", "file"})
	switch field {
	case "name":
		s.Name = common.ReadField("Insert new name: ")
	case "description":
		s.Description = common.ReadField("Insert new description: ")
	case "file":
		s.FileName = common.ReadField("Insert new file name: ")
		fmt.Println("Paste the new file content (press Ctrl+D when done):")
		content := readFileContent()
		s.Content = base64.StdEncoding.EncodeToString([]byte(content))
	}
	return nil
}

func (s *SecretFile) Print() {
	confirm := common.SelectField(fmt.Sprintf("Are you sure you want to print the content of '%s'?", s.FileName), []string{"No", "Yes"})
	if confirm != "Yes" {
		return
	}
	data, err := base64.StdEncoding.DecodeString(s.Content)
	common.Check(err)
	fmt.Printf("--- %s ---\n%s\n--- end ---\n", s.FileName, string(data))
}

// Copy copies the file to the clipboard area, suggesting ~/Desktop/<filename> as default.
func (s *SecretFile) Copy() {
	destination := expandHome(common.ReadField(fmt.Sprintf("Save '%s' to path: ", s.FileName)))
	if destination == "" {
		destination = filepath.Join(".", s.FileName)
	}
	data, err := base64.StdEncoding.DecodeString(s.Content)
	common.Check(err)
	common.Check(os.WriteFile(destination, data, 0600))
	fmt.Printf("Saved to %s\n", destination)
}

// Use saves the file to the default path (~/Desktop) or a path chosen by the user.
func (s *SecretFile) Use() {
	home, _ := os.UserHomeDir()
	defaultPath := filepath.Join(home, "Desktop", s.FileName)
	confirm := common.SelectField(fmt.Sprintf("Save '%s' to %s?", s.FileName, defaultPath), []string{"Yes", "No, choose path"})
	var destination string
	if confirm == "Yes" {
		destination = defaultPath
	} else {
		destination = expandHome(common.ReadField("Insert destination path: "))
		if destination == "" {
			destination = defaultPath
		}
	}
	data, err := base64.StdEncoding.DecodeString(s.Content)
	common.Check(err)
	common.Check(os.WriteFile(destination, data, 0600))
	fmt.Printf("Saved to %s\n", destination)
}

func readFileContent() string {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
