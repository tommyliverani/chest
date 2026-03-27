package common

import (
	"fmt"
	"os"

	"golang.design/x/clipboard"
)

func WriteToClipboard(text string) {
	if err := clipboard.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "clipboard init failed: %v\n", err)
		return
	}
	clipboard.Write(clipboard.FmtText, []byte(text))
}
