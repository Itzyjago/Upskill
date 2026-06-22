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
