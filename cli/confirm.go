package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/pi2pie/git-worktree-tasks/ui"
)

func confirmPrompt(in io.Reader, out io.Writer, message string) (bool, error) {
	prompt := ui.PromptStyle.Render(message)
	confirmation := fmt.Sprintf("%s %s %s",
		ui.MutedStyle.Render("Type"),
		ui.AccentStyle.Render("'yes'"),
		ui.MutedStyle.Render("to confirm:"),
	)
	if _, err := fmt.Fprintf(out, "%s %s ", prompt, confirmation); err != nil {
		return false, err
	}
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(line), "yes"), nil
}
