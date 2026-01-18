package main

import (
	"os"

	"github.com/pi2pie/git-worktree-tasks/cli"
)

func main() {
	os.Exit(cli.Execute())
}
