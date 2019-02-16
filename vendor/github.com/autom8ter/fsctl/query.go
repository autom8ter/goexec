package fsctl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompt prompts user for input with default value.
func (f *Fs) Prompt(key, question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("string | " + question)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	text = strings.TrimRight(text, "`")
	text = strings.TrimLeft(text, "`")
	if strings.Contains(text, "?") {
		newtext := strings.Split(text, "?")
		text = newtext[0]
	}
	f.Set(key, text)
	return text
}
