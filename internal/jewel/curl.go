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

const CURL_KIND = "curl"

var supportedCurlOptions = []string{
	"-k (insecure)",
	"-v (verbose)",
	"-s (silent)",
	"-L (follow redirects)",
	"-i (include headers)",
	"-I (HEAD only)",
	"done",
}

var curlOptionFlag = map[string]string{
	"-k (insecure)":         "-k",
	"-v (verbose)":          "-v",
	"-s (silent)":           "-s",
	"-L (follow redirects)": "-L",
	"-i (include headers)":  "-i",
	"-I (HEAD only)":        "-I",
}

type Curl struct {
	baseJewel
	Url      string   `json:"url"`
	Method   string   `json:"method"`
	Options  []string `json:"options"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

func (c *Curl) GetEmoji() string { return "🌐" }

func (c *Curl) ToJson() (json.RawMessage, error) { return json.Marshal(c) } //nolint:gosec

func ParseCurl(data json.RawMessage) (*Curl, error) {
	var curl Curl
	if err := json.Unmarshal(data, &curl); err != nil {
		return nil, err
	}
	return &curl, nil
}

func selectCurlOptions() []string {
	var selected []string
	for {
		choice := common.SelectField("Add option (select 'done' when finished):", supportedCurlOptions)
		if choice == "done" {
			break
		}
		flag := curlOptionFlag[choice]
		alreadyAdded := false
		for _, o := range selected {
			if o == flag {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			selected = append(selected, flag)
			fmt.Printf("Added: %s\n", flag)
		} else {
			fmt.Printf("%s already added\n", flag)
		}
	}
	return selected
}

func CreateCurl(name string, description string) (*Curl, error) {
	url := common.ReadField("Insert URL: ")
	method := common.SelectField("Select HTTP method:", []string{"GET", "POST", "PUT", "PATCH", "DELETE"})
	options := selectCurlOptions()
	username := common.ReadField("Insert username (leave empty to skip): ")
	var password string
	if username != "" {
		password = common.ReadSecret("Insert password: ")
	}
	return &Curl{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        CURL_KIND,
			Description: description,
		},
		Url:      url,
		Method:   method,
		Options:  options,
		Username: username,
		Password: password,
	}, nil
}

func init() {
	factory.RegisterJewelCreator(CURL_KIND, func(name string, description string) (factory.Jewel, error) {
		return CreateCurl(name, description)
	})
	factory.RegisterJewelParser(CURL_KIND, func(data json.RawMessage) (factory.Jewel, error) {
		return ParseCurl(data)
	})
	factory.RegisterJewelHelp(CURL_KIND, factory.JewelHelp{
		Emoji:    "🌐",
		Name:     "curl",
		Short:    "stores a curl request with method, options and credentials",
		Behavior: "executes the curl command",
		Operations: map[string]string{
			"add":   "saves a new curl request (URL, method, options, auth) into the open chest",
			"ls":    "lists all curl entries in the open chest",
			"edit":  "edits name, description, URL, method, options, username or password",
			"rm":    "removes a curl entry from the open chest",
			"print": "shows the curl command with masked password",
			"copy":  "copies the full curl command to the clipboard",
		},
	})
}

func (c *Curl) Edit() error {
	field := common.SelectField("which field do you want to edit?", []string{"name", "description", "url", "method", "options", "username", "password"})
	switch field {
	case "name":
		c.Name = common.ReadField("Insert new name: ")
	case "description":
		c.Description = common.ReadField("Insert new description: ")
	case "url":
		c.Url = common.ReadField("Insert new URL: ")
	case "method":
		c.Method = common.SelectField("Select HTTP method:", []string{"GET", "POST", "PUT", "PATCH", "DELETE"})
	case "options":
		c.Options = selectCurlOptions()
	case "username":
		c.Username = common.ReadField("Insert new username: ")
	case "password":
		c.Password = common.ReadSecret("Insert new password: ")
	}
	return nil
}

func (c *Curl) buildArgs(maskPassword bool) []string {
	args := []string{"-X", c.Method}
	args = append(args, c.Options...)
	if c.Username != "" {
		pass := c.Password
		if maskPassword {
			pass = strings.Repeat("*", len(c.Password))
		}
		args = append(args, "-u", fmt.Sprintf("%s:%s", c.Username, pass))
	}
	args = append(args, c.Url)
	return args
}

func (c *Curl) Print() {
	args := c.buildArgs(true)
	fmt.Printf("curl %s\n", strings.Join(args, " "))
}

func (c *Curl) Copy() {
	args := c.buildArgs(false)
	common.WriteToClipboard("curl " + strings.Join(args, " "))
}

func (c *Curl) Use() {
	args := c.buildArgs(false)
	cmd := exec.Command("curl", args...) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "curl failed: %v\n", err)
	}
}
