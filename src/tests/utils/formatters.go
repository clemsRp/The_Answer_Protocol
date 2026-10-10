package utils

import "fmt"

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
)

func Truncate(s string) string {
	max := 40
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}

func FormatMismatch(Command, wanted, got string) string {
	return fmt.Sprintf("\n\t%s📡 COMMAND :%s %s\n\t%s🎯 WANTED  :%s %s\n\t%s❌ GOT     :%s %s\n",
		colorCyan, colorReset, Command,
		colorGreen, colorReset, wanted,
		colorRed, colorReset, Truncate(got))
}

func FormatInvalidJSON(got string) string {
	return fmt.Sprintf("\n\t%s⚠️  JSON INVALID%s\n\t%s❌ GOT     :%s %s\n",
		colorRed, colorReset, colorRed, colorReset, Truncate(got))
}
