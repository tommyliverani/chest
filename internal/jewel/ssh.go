package jewel

import (
	"chest/internal/common"
	"chest/internal/factory"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/term"
)

const SSH_KIND = "ssh"

type SshKey struct {
	baseJewel
	Username string `json:"username"`
	Url      string `json:"url"`
	Password string `json:"password"`
}

func (s *SshKey) GetEmoji() string { return "🧭" }

func (s *SshKey) ToJson() (json.RawMessage, error) { return json.Marshal(s) } //nolint:gosec

func ParseSshKey(data json.RawMessage) (*SshKey, error) {
	var sk SshKey
	if err := json.Unmarshal(data, &sk); err != nil {
		return nil, err
	}
	return &sk, nil
}

func CreateSshKey(name string, description string) (*SshKey, error) {
	username := common.ReadField("Insert username: ")
	url := common.ReadField("Insert host/url: ")
	password := common.ReadSecret("Insert password: ")
	return &SshKey{
		baseJewel: baseJewel{
			Name:        name,
			Kind:        SSH_KIND,
			Description: description,
		},
		Username: username,
		Url:      url,
		Password: password,
	}, nil
}

func init() {
	factory.RegisterJewelCreator(SSH_KIND, func(name string, description string) (factory.Jewel, error) {
		return CreateSshKey(name, description)
	})
	factory.RegisterJewelParser(SSH_KIND, func(data json.RawMessage) (factory.Jewel, error) {
		return ParseSshKey(data)
	})
	factory.RegisterJewelHelp(SSH_KIND, factory.JewelHelp{
		Emoji:    "🧭",
		Name:     "ssh",
		Short:    "stores SSH credentials (username, host, password)",
		Behavior: "opens an interactive SSH session to the stored host",
		Operations: map[string]string{
			"add":   "saves new SSH credentials into the open chest",
			"ls":    "lists all SSH entries in the open chest",
			"edit":  "edits name, description, username, host or password",
			"rm":    "removes an SSH entry from the open chest",
			"print": "shows username and host (password hidden)",
			"copy":  "copies the ssh command (ssh user@host) to the clipboard",
		},
	})
}

func (s *SshKey) Edit() error {
	field := common.SelectField("which field do you want to edit?", []string{"name", "description", "username", "url", "password"})
	switch field {
	case "name":
		s.Name = common.ReadField("Insert new name: ")
	case "description":
		s.Description = common.ReadField("Insert new description: ")
	case "username":
		s.Username = common.ReadField("Insert new username: ")
	case "url":
		s.Url = common.ReadField("Insert new host/url: ")
	case "password":
		s.Password = common.ReadSecret("Insert new password: ")
	}
	return nil
}

func (s *SshKey) Print() {
	confirm := common.SelectField(fmt.Sprintf("Are you sure you want to print the credentials for '%s'?", s.Name), []string{"No", "Yes"})
	if confirm != "Yes" {
		return
	}
	fmt.Printf("Username: %s\nHost:     %s\nPassword: %s\n", s.Username, s.Url, s.Password)
}

func (s *SshKey) Copy() {
	common.WriteToClipboard(fmt.Sprintf("ssh %s@%s", s.Username, s.Url))
}

// Use opens an interactive SSH session using the stored password (no external tools required).
func (s *SshKey) Use() {
	host := s.Url
	if !strings.Contains(host, ":") {
		host += ":22"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot determine home directory: %v\n", err)
		return
	}
	knownHostsFile := filepath.Join(homeDir, ".ssh", "known_hosts")
	hostKeyCallback, err := knownhosts.New(knownHostsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot load known_hosts (%s): %v\nRun 'ssh %s' first to add the host.\n", knownHostsFile, err, s.Url)
		return
	}

	config := &gossh.ClientConfig{
		User: s.Username,
		Auth: []gossh.AuthMethod{
			gossh.Password(s.Password),
		},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := gossh.Dial("tcp", host, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SSH connection failed: %v\n", err)
		return
	}
	defer func() { _ = client.Close() }()

	session, err := client.NewSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create session: %v\n", err)
		return
	}
	defer func() { _ = session.Close() }()

	fd := int(os.Stdin.Fd()) //nolint:gosec
	if term.IsTerminal(fd) {
		oldState, err := term.MakeRaw(fd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to set terminal raw mode: %v\n", err)
			return
		}
		defer func() { _ = term.Restore(fd, oldState) }()

		w, h, err := term.GetSize(fd)
		if err != nil {
			w, h = 80, 24
		}

		if err := session.RequestPty("xterm-256color", h, w, gossh.TerminalModes{
			gossh.ECHO:          1,
			gossh.TTY_OP_ISPEED: 14400,
			gossh.TTY_OP_OSPEED: 14400,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to request PTY: %v\n", err)
			return
		}
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start shell: %v\n", err)
		return
	}

	if err := session.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "session ended: %v\n", err)
	}
}
