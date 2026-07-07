#!/usr/bin/env bash
# Turns bash.md's "-e is narrower than it sounds" claims into assertions
# instead of prose someone has to trust. Each check runs in a subshell (parens)
# so one exemption's `set -e` state can't leak into the next.
set -uo pipefail

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

# 1. -e never fires inside an if condition.
out="$( (set -e; if false; then :; fi; echo reached) )"
assert "if-condition exemption" "$out" "reached"

# 2. A command on the left of && is exempt.
out="$( (set -e; false && echo x; echo reached) )"
assert "&&-left exemption" "$out" "reached"

# 3. A failing command inside a function is exempt when the call is tested.
out="$( (set -e; f() { false; echo reached; }; f || true) )"
assert "function call under || exemption" "$out" "reached"

# 4. The same function called plainly aborts before its own echo.
out="$( (set -e; f() { false; echo unreached; }; f; echo after) 2>/dev/null; echo "status=$?" )"
assert "bare call is NOT exempt" "$out" "status=1"

exit "$fail"
