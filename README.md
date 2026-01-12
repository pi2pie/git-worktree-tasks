# git-worktree-tasks

A small CLI to manage task-based Git worktrees with predictable naming and cleanup flows.

## Install

Build the binary:

```bash
go build -o git-worktree-tasks ./
```

Run directly:

```bash
go run ./ --help
```

### Global install

Install into your `GOBIN` (or `GOPATH/bin`) without cloning:

```bash
go install github.com/dev-pi2pie/git-worktree-tasks@latest
```

### Local clone

Clone the repo, then install locally:

```bash
git clone <repo-url>
cd git-worktree-tasks
go install
```

## Usage

Create a worktree for a task:

```bash
git-worktree-tasks create "my-task" --base main
```

List worktrees (relative paths by default):

```bash
git-worktree-tasks list
```

Show detailed status:

```bash
git-worktree-tasks status
```

Show absolute paths when needed:

```bash
git-worktree-tasks list --absolute-path
git-worktree-tasks status --absolute-path
```

Finish a task (merge into target and cleanup):

```bash
git-worktree-tasks finish "my-task" --target main
```

Cleanup without merge:

```bash
git-worktree-tasks cleanup "my-task"
```

## Notes

- Default worktree path uses the pattern `../<repo>_<task>`.
- Use `--output json` on list/status for machine-readable output.
