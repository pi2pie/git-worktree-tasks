# wtree-tasks-functions

# wtask <task-name> [base-branch]
wtask() {
  if [[ $# -lt 1 ]]; then
    echo "Usage: wtask <task-name> [base-branch]"
    return 1
  fi

  local TASK=$1
  local BASE=${2:-main}    # default base branch
  local REPO=$(basename "$(git rev-parse --show-toplevel)")
  local WT_NAME="${REPO}_${TASK}"
  local WT_PATH="../${WT_NAME}"

  # If folder already exists, don't overwrite
  if [[ -d "$WT_PATH" ]]; then
    echo "❗ Worktree path '$WT_PATH' already exists!"
    return 1
  fi

  # Construct branch name (safe: replace spaces with -)
  local BRANCH=$(echo "$TASK" | sed 's/[^a-zA-Z0-9_/-]/-/g')

  # Create worktree and new branch
  git worktree add -b "$BRANCH" "$WT_PATH" "$BASE"
  if [[ $? -ne 0 ]]; then
    echo "❗ Failed to add worktree"
    return 1
  fi

  cd "$WT_PATH" || return
  echo "Worktree ready: $WT_PATH (branch: $BRANCH)"
}

# Usage: wt_finish <task-name> [target-branch]
wt_finish() {
  if [[ $# -lt 1 ]]; then
    echo "Usage: wt_finish <task-name> [target-branch]"
    return 1
  fi

  local TASK=$1
  local TARGET=${2:-main}
  local REPO=$(basename "$(git rev-parse --show-toplevel)")
  local BRANCH=$(echo "$TASK" | sed 's/[^a-zA-Z0-9_/-]/-/g')
  local WT_PATH="../${REPO}_${TASK}"

  # 1) Switch to target branch
  git checkout "$TARGET" || { echo "❗ Failed to checkout $TARGET"; return 1; }

  # 2) Merge task branch
  git merge "$BRANCH" || { echo "❗ Merge conflict; resolve first"; return 1; }

  # 3) Remove worktree
  git worktree remove "$WT_PATH" || { echo "❗ Failed to remove worktree"; }

  # Optional: prune leftover metadata
  git worktree prune

  # 4) Delete branch locally
  git branch -d "$BRANCH"

  echo "Merged $BRANCH into $TARGET; cleaned up worktree and branch."
}