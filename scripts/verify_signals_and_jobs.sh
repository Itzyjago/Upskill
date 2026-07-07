#!/usr/bin/env bash
# Verifies linux.md's trap/background-job claims against whatever shell this
# actually runs on -- and says so, since that shell is MSYS/git-bash on
# Windows here, not real Linux. See notes/linux.md for what that does and
# doesn't cover.
set -uo pipefail

echo "uname -s: $(uname -s)   (MINGW64_NT.../MSYS = git-bash on Windows, not real Linux)"

fail=0
assert() {
    local desc="$1" got="$2" want="$3"
    if [[ "$got" != "$want" ]]; then
        echo "FAIL: $desc (got '$got', want '$want')"
        fail=1
    else
        echo "ok: $desc"
    fi
}

# 1. trap fires on TERM sent to a background job.
marker="$(mktemp)"
(
    trap 'echo caught >> '"$marker"'; exit 0' TERM
    sleep 5
) &
child=$!
sleep 0.3
kill -TERM "$child" 2>/dev/null
wait "$child" 2>/dev/null
assert "trap caught SIGTERM" "$(cat "$marker" 2>/dev/null)" "caught"
rm -f "$marker"

# 2. background job + wait propagates the real exit status, not just "done".
( exit 7 ) &
bg_pid=$!
wait "$bg_pid"
assert "wait propagates background exit status" "$?" "7"

# 3. \$! tracks the most recently backgrounded job's PID, not the shell's.
sleep 1 &
job_pid=$!
assert "\$! is a real pid, not the shell's" "$([[ "$job_pid" != "$$" ]] && echo distinct)" "distinct"
wait "$job_pid" 2>/dev/null

exit "$fail"
