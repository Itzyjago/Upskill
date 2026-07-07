# Python notes

## Virtual environments
```bash
python -m venv .venv
source .venv/bin/activate   # Windows: .venv\Scripts\activate
pip install -r requirements.txt
```
Keep one venv per project; never `pip install` into the system Python.

## Comprehensions
```python
squares = [n * n for n in range(10)]
evens   = {n for n in range(10) if n % 2 == 0}
lookup  = {k: v for k, v in pairs}
```
Readable beats clever — if it needs two `for` clauses, write a loop.

## f-strings
```python
print(f"{name} scored {score:.1f}%")
```

## Truthiness
- Empty containers (`[]`, `{}`, `""`, `0`, `None`) are falsy.
- Prefer `if not items:` over `if len(items) == 0:`.

## Context managers
```python
with open("data.txt") as f:
    data = f.read()
```
The file closes even if the block raises.

## Gotcha
- Mutable default args are shared across calls. Use `def f(x=None):` then
  `x = x or []` inside.

## Packaging (the roadmap gap: venvs were covered, packaging wasn't)
A venv answers "how do I isolate dependencies." Packaging answers a different
question: "how does someone else `pip install` what I wrote."
- **`pyproject.toml`** is the modern, standardized answer (PEP 621) — one file
  declares metadata, dependencies, and the build backend, replacing the old
  `setup.py`/`setup.cfg` split:
  ```toml
  [project]
  name = "wordcount-report"
  version = "0.1.0"
  dependencies = ["requests>=2.31"]

  [build-system]
  requires = ["setuptools>=68"]
  build-backend = "setuptools.build_meta"
  ```
- **`pip install -e .`** (editable install) — installs the package as a link
  back to the source tree instead of copying it, so local edits show up
  immediately without reinstalling. The everyday command while developing a
  package, as opposed to just running a script inside a venv.
- **Where dependency *pinning* actually happens**: `pyproject.toml`'s
  `dependencies` are ranges (what the package is *compatible* with);
  `requirements.txt` (or a lockfile like `poetry.lock`/`uv.lock`) pins *exact*
  versions for a reproducible install. Same distinction as `^1.2.3` vs a
  committed lockfile in [semver.md](semver.md) — one says what's allowed, the
  other says what's actually installed.
- **`__init__.py`** marks a directory as a package (import-able); its
  presence is also what makes relative imports (`from .utils import foo`)
  resolve within it.

## Stdlib fluency: `pathlib` over string paths
The other roadmap bullet — reaching for the stdlib instead of a habit carried
over from another language:
```python
from pathlib import Path

for f in Path("notes").glob("*.md"):
    print(f.stem, f.stat().st_size)
```
- `Path` objects overload `/` for joining (`Path("a") / "b"`) instead of
  string-concatenating with the OS separator — the same "stop hand-rolling
  what the stdlib already got right" instinct as `os.path.join`, just newer
  and object-oriented instead of a bag of free functions.
- `.glob()` and `.rglob()` (recursive) beat `os.walk()` for anything that's
  fundamentally "find files matching a pattern" — one call instead of a
  triple-nested loop over `(dirpath, dirnames, filenames)`.

## Stdlib fluency, applied: `scripts/parse_access_log.py`
Past isolated snippets — a real script parsing wordcount's structured JSON
logs, stdlib only (`json`, `dataclasses`, `pathlib`, `argparse`,
`statistics`, `collections.defaultdict`). Tested against a *real* captured
log (`scripts/testdata/`), not synthetic fixtures — `unittest` catches a
regression the same way `go test` does elsewhere in this repo.
