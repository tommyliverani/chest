package common

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func WriteToClipboard(text string) {
	candidates := [][]string{
		{"xclip", "-selection", "clipboard"},
		{"xsel", "--clipboard", "--input"},
		{"wl-copy"},
	}
	for _, args := range candidates {
		if _, err := exec.LookPath(args[0]); err != nil {
			continue
		}
		cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return
		}
	}
	fmt.Fprintf(os.Stderr, "clipboard: install xclip, xsel, or wl-copy\n")
	_ = os.Stderr
}
