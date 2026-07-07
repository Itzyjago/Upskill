# scripts

Standalone Python (stdlib only) utilities, run with `python scripts/<name>.py`.
Tests live alongside as `test_<name>.py`; run from inside `scripts/` with
`python -m unittest test_<name>`.

- `parse_access_log.py` — summarizes wordcount's structured JSON access log
  by route (count, error rate, p50/p95 latency). `testdata/` holds a real
  captured sample, not synthesized data — see the file's own docstring for
  how it was captured.
- `regex_lookaround.py` — verifies `notes/regex.md`'s lookaround examples
  against Python's `re` for real.
