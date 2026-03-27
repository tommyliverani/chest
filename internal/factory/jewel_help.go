package factory

import "fmt"

// JewelHelp contains the short and long description for a jewel kind.
type JewelHelp struct {
	Emoji      string
	Name       string
	Short      string            // one-line description
	Behavior   string            // what "use" does on this jewel
	Operations map[string]string // command -> what it does
}

var jewelHelpRegistry = make(map[string]JewelHelp)

func RegisterJewelHelp(kind string, help JewelHelp) {
	jewelHelpRegistry[kind] = help
}

// ShortHelp returns the one-line description: "  <emoji> <name> - <short>"
func ShortHelp(kind string) string {
	h, ok := jewelHelpRegistry[kind]
	if !ok {
		return fmt.Sprintf("  %s - no description available", kind)
	}
	return fmt.Sprintf("  %s %-6s - %s", h.Emoji, h.Name, h.Short)
}

// LongHelp returns the full multi-line description for a jewel kind.
func LongHelp(kind string) string {
	h, ok := jewelHelpRegistry[kind]
	if !ok {
		return fmt.Sprintf("No help available for '%s'\n", kind)
	}
	out := fmt.Sprintf("%s %s - %s\n", h.Emoji, h.Name, h.Short)
	if h.Behavior != "" {
		out += fmt.Sprintf("Behavior - %s\n", h.Behavior)
	}
	out += "\nUsage:\n"
	for _, cmd := range []string{"add", "ls", "edit", "rm", "print", "copy"} {
		if desc, exists := h.Operations[cmd]; exists {
			out += fmt.Sprintf("  chest %s %-6s  %s\n", kind, cmd, desc)
		}
	}
	return out
}
