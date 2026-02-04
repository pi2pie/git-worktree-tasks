package cli

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/pi2pie/git-worktree-tasks/internal/git"
	"github.com/pi2pie/git-worktree-tasks/internal/worktree"
	"github.com/pi2pie/git-worktree-tasks/ui"
	"github.com/spf13/cobra"
)

type statusOptions struct {
	output string
	target string
	task   string
	branch string
	abs    bool
	grid   bool
	strict bool
}

type statusRow struct {
	Task         string `json:"task"`
	Branch       string `json:"branch"`
	Path         string `json:"path"`
	ModifiedTime string `json:"modified_time"`
	Base         string `json:"base"`
	Target       string `json:"target"`
	LastCommit   string `json:"last_commit"`
	Dirty        bool   `json:"dirty"`
	Ahead        int    `json:"ahead"`
	Behind       int    `json:"behind"`
}

func newStatusCommand() *cobra.Command {
	opts := &statusOptions{output: "table"}
	cmd := &cobra.Command{
		Use:   "status [task]",
		Short: "Show detailed worktree status",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			runner := defaultRunner()
			mode := modeClassic
			var codexHome string
			var codexWorktrees string
			if cfg, ok := configFromContext(cmd.Context()); ok {
				mode = cfg.Mode
				if mode == modeCodex {
					var err error
					codexHome, err = codexHomeDir()
					if err != nil {
						return err
					}
					codexWorktrees = codexWorktreesRoot(codexHome)
				} else {
					if home, err := codexHomeDir(); err == nil {
						codexHome = home
						codexWorktrees = codexWorktreesRoot(codexHome)
					}
				}
				if !cmd.Flags().Changed("output") {
					opts.output = cfg.Status.Output
				}
				if !flagChangedAny(cmd, "absolute-path", "abs") {
					opts.abs = cfg.Status.AbsolutePath
				}
				if !cmd.Flags().Changed("grid") {
					opts.grid = cfg.Status.Grid
				}
				if !cmd.Flags().Changed("strict") {
					opts.strict = cfg.Status.Strict
				}
			}
			repoRoot, err := repoRoot(ctx, runner)
			if err != nil {
				return err
			}
			repo, err := git.RepoBaseName(ctx, runner)
			if err != nil {
				return err
			}
			if len(args) == 1 && opts.task != "" {
				return fmt.Errorf("use either --task or [task], not both")
			}
			var query string
			if len(args) == 1 {
				if mode == modeCodex {
					query = strings.TrimSpace(args[0])
					if query == "" {
						return fmt.Errorf("task query cannot be empty")
					}
				} else {
					query, err = normalizeTaskQuery(args[0])
					if err != nil {
						return err
					}
				}
			}
			if opts.task != "" {
				if mode == modeCodex {
					query = strings.TrimSpace(opts.task)
					if query == "" {
						return fmt.Errorf("task query cannot be empty")
					}
				} else {
					query, err = normalizeTaskQuery(opts.task)
					if err != nil {
						return err
					}
					opts.strict = true
				}
				if mode != modeCodex {
					opts.strict = true
				}
			}

			target := opts.target
			if target == "" {
				current, err := git.CurrentBranch(ctx, runner)
				if err != nil {
					if errors.Is(err, git.ErrNoCommits) {
						current, err = git.SymbolicRefShort(ctx, runner, "HEAD")
					}
					if err != nil {
						return err
					}
				}
				target = current
			}

			worktrees, err := worktree.List(ctx, runner, repoRoot)
			if err != nil {
				return err
			}

			rows := make([]statusRow, 0, len(worktrees))
			for _, wt := range worktrees {
				branch := strings.TrimPrefix(wt.Branch, "refs/heads/")
				task := "-"
				var wtAbs string
				if mode == modeCodex {
					var err error
					wtAbs, err = worktree.NormalizePath(repoRoot, wt.Path)
					if err != nil {
						return err
					}
					opaqueID, _, ok := codexWorktreeInfo(codexWorktrees, wtAbs)
					if !ok {
						continue
					}
					task = opaqueID
					if branch == "" {
						branch = "detached"
					}
					if query != "" && !matchesTask(task, query, opts.strict) {
						continue
					}
				} else {
					if codexWorktrees != "" {
						var err error
						wtAbs, err = worktree.NormalizePath(repoRoot, wt.Path)
						if err != nil {
							return err
						}
						if _, _, ok := codexWorktreeInfo(codexWorktrees, wtAbs); ok {
							continue
						}
					}
					task, _ = worktree.TaskFromPath(repo, wt.Path)
					if task == "" {
						task = "-"
					}
					if query != "" && !matchesTask(task, query, opts.strict) {
						continue
					}
				}
				if opts.branch != "" && branch != opts.branch {
					continue
				}

				statusInfo, err := worktree.Status(ctx, runner, wt.Path, target)
				if err != nil {
					return err
				}
				if wtAbs == "" {
					var err error
					wtAbs, err = worktree.NormalizePath(repoRoot, wt.Path)
					if err != nil {
						return err
					}
				}
				modified := ""
				info, err := os.Stat(wtAbs)
				if err != nil {
					if !os.IsNotExist(err) {
						return fmt.Errorf("stat worktree %s: %w", wtAbs, err)
					}
				} else {
					modified = info.ModTime().UTC().Format(time.RFC3339)
				}

				rows = append(rows, statusRow{
					Task:         task,
					Branch:       branch,
					Path:         displayPathForMode(repoRoot, wt.Path, opts.abs, mode, codexHome),
					ModifiedTime: modified,
					Base:         statusInfo.Base,
					Target:       target,
					LastCommit:   statusInfo.LastCommit,
					Dirty:        statusInfo.Dirty,
					Ahead:        statusInfo.Ahead,
					Behind:       statusInfo.Behind,
				})
				if query != "" && !opts.strict {
					break
				}
			}

			if mode != modeCodex && len(rows) == 0 {
				fallbackBranch := opts.branch
				if fallbackBranch == "" {
					fallbackBranch = query
				}
				path, ok, err := fallbackPathForBranch(ctx, runner, repoRoot, fallbackBranch)
				if err != nil {
					return err
				}
				if ok {
					statusInfo, err := worktree.Status(ctx, runner, path, target)
					if err != nil {
						return err
					}
					branch, err := git.CurrentBranchAt(ctx, runner, path)
					if err != nil {
						return err
					}
					modified := ""
					info, err := os.Stat(path)
					if err != nil {
						if !os.IsNotExist(err) {
							return fmt.Errorf("stat worktree %s: %w", path, err)
						}
					} else {
						modified = info.ModTime().UTC().Format(time.RFC3339)
					}
					rows = append(rows, statusRow{
						Task:         "-",
						Branch:       branch,
						Path:         displayPath(repoRoot, path, opts.abs),
						ModifiedTime: modified,
						Base:         statusInfo.Base,
						Target:       target,
						LastCommit:   statusInfo.LastCommit,
						Dirty:        statusInfo.Dirty,
						Ahead:        statusInfo.Ahead,
						Behind:       statusInfo.Behind,
					})
				}
			}

			return renderStatus(cmd, opts.output, rows, opts.grid)
		},
	}

	cmd.Flags().StringVarP(&opts.output, "output", "o", opts.output, "output format: table, json, or csv")
	cmd.Flags().StringVar(&opts.target, "target", "", "target branch for ahead/behind comparison")
	cmd.Flags().StringVar(&opts.task, "task", "", "filter by task name")
	cmd.Flags().StringVar(&opts.branch, "branch", "", "filter by branch name")
	cmd.Flags().BoolVar(&opts.abs, "absolute-path", false, "show absolute paths instead of relative")
	cmd.Flags().BoolVar(&opts.abs, "abs", false, "alias for --absolute-path")
	cmd.Flags().BoolVar(&opts.grid, "grid", false, "render table with grid borders")
	cmd.Flags().BoolVar(&opts.strict, "strict", false, "require exact task match (after trimming and slugifying)")

	return cmd
}

func renderStatus(cmd *cobra.Command, format string, rows []statusRow, grid bool) error {
	switch format {
	case "table":
		columns := []tableColumn{
			{Header: "TASK", MinWidth: 6},
			{Header: "BRANCH", MinWidth: 10, Flexible: true, Truncate: true, Style: func(value string) lipgloss.Style { return ui.AccentStyle }},
			{Header: "PATH", MinWidth: 16, Flexible: true, Truncate: true},
			{Header: "MODIFIED", MinWidth: 10, Flexible: true, Truncate: true, Style: func(value string) lipgloss.Style { return ui.MutedStyle }},
			{Header: "BASE", MinWidth: 8, Flexible: true, Truncate: true, Style: func(value string) lipgloss.Style { return ui.MutedStyle }},
			{Header: "TARGET", MinWidth: 8, Flexible: true, Truncate: true, Style: func(value string) lipgloss.Style { return ui.MutedStyle }},
			{Header: "LAST_COMMIT", MinWidth: 12, MaxWidth: 24, Flexible: true, Truncate: true, Style: func(value string) lipgloss.Style { return ui.MutedStyle }},
			{Header: "DIRTY", MinWidth: 5, Style: func(value string) lipgloss.Style {
				if value == "true" {
					return ui.WarningStyle
				}
				return ui.SuccessStyle
			}},
			{Header: "AHEAD", MinWidth: 5, Style: func(value string) lipgloss.Style {
				if value != "0" {
					return ui.WarningStyle
				}
				return ui.MutedStyle
			}},
			{Header: "BEHIND", MinWidth: 6, Style: func(value string) lipgloss.Style {
				if value != "0" {
					return ui.ErrorStyle
				}
				return ui.MutedStyle
			}},
		}
		tableRows := make([][]string, 0, len(rows))
		for _, row := range rows {
			tableRows = append(tableRows, []string{
				row.Task,
				row.Branch,
				row.Path,
				row.ModifiedTime,
				row.Base,
				row.Target,
				row.LastCommit,
				strconv.FormatBool(row.Dirty),
				strconv.Itoa(row.Ahead),
				strconv.Itoa(row.Behind),
			})
		}
		renderTable(cmd, columns, tableRows, grid)
		return nil
	case "json":
		payload, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(payload)); err != nil {
			return err
		}
		return nil
	case "csv":
		writer := csv.NewWriter(cmd.OutOrStdout())
		if err := writer.Write([]string{
			"task", "branch", "path", "modified_time", "base", "target", "last_commit", "dirty", "ahead", "behind",
		}); err != nil {
			return err
		}
		for _, row := range rows {
			if err := writer.Write([]string{
				row.Task,
				row.Branch,
				row.Path,
				row.ModifiedTime,
				row.Base,
				row.Target,
				row.LastCommit,
				strconv.FormatBool(row.Dirty),
				strconv.Itoa(row.Ahead),
				strconv.Itoa(row.Behind),
			}); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}
