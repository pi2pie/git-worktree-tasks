package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/pi2pie/git-worktree-tasks/cli"
	"github.com/spf13/cobra/doc"
)

func main() {
	outDir := flag.String("out", "man", "output directory for man pages")
	useName := flag.String("use", "git-worktree-tasks", "root command name for man pages")
	title := flag.String("title", "GIT-WORKTREE-TASKS", "man page title")
	source := flag.String("source", "git-worktree-tasks", "man page source")
	flag.Parse()

	manDir := filepath.Join(*outDir, "man1")
	if err := os.MkdirAll(manDir, 0o755); err != nil {
		log.Fatalf("create man directory: %v", err)
	}

	header := &doc.GenManHeader{
		Title:   *title,
		Section: "1",
		Manual:  "Git Worktree Tasks Manual",
		Source:  *source,
	}

	opts := doc.GenManTreeOptions{
		Header: header,
		Path:   manDir,
	}

	root := cli.RootCommand()
	root.DisableAutoGenTag = true
	root.Use = *useName
	root.Aliases = nil
	if err := doc.GenManTreeFromOpts(root, opts); err != nil {
		log.Fatalf("generate man pages for %s: %v", *useName, err)
	}
}
