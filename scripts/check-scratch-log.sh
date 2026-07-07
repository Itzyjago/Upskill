#!/usr/bin/env bash
# Warns if notes/scratch-log.md is growing past a size where it stops being
# skimmable — same "revisit past ~500 lines" rule PROGRESS.md applies to
# itself, applied here as an actual check instead of a promise to remember.
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
log_file="$repo_root/notes/scratch-log.md"
max_lines="${SCRATCH_LOG_MAX_LINES:-40}"

if [[ ! -f "$log_file" ]]; then
    echo "check-scratch-log: $log_file not found" >&2
    exit 1
fi

line_count="$(wc -l < "$log_file" | tr -d ' ')"

if (( line_count > max_lines )); then
    echo "scratch-log.md is $line_count lines (max $max_lines) — clear resolved entries into real notes"
    exit 1
fi

echo "scratch-log.md: $line_count/$max_lines lines, ok"
