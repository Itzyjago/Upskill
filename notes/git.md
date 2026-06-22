# Git notes

## Undoing things
- Unstage a file: `git restore --staged <file>`
- Discard working changes: `git restore <file>`
- Amend the last commit: `git commit --amend`
- Undo last commit, keep changes: `git reset --soft HEAD~1`

## Branching
```bash
git switch -c feature/login   # create + switch
git switch main
git merge feature/login
```

## Rebase vs merge
- `merge` keeps history as it happened (extra merge commits).
- `rebase` replays your commits on top for a linear history.
- Rule: rebase local-only work, never rebase shared/pushed branches.

## Recovery
- `git reflog` shows where HEAD has been — almost nothing is truly lost.
- `git restore --source=<sha> <file>` pulls one file from an old commit.

## Inspecting
- `git log --oneline --graph --all` — the whole picture.
- `git blame <file>` — who changed each line and when.
