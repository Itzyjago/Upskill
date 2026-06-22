# JavaScript notes

## Equality
- Use `===` / `!==`. `==` coerces and surprises you (`0 == ''` is `true`).
- `NaN === NaN` is `false`; test with `Number.isNaN(x)`.

## `var` vs `let` vs `const`
- `var` is function-scoped and hoisted — avoid.
- `let` for reassignment, `const` for everything else.
- `const` freezes the binding, not the value: a `const` array can still `.push()`.

## Async
- A `Promise` is a value that resolves later. `await` unwraps it.
- `await` inside a loop runs serially. For parallel work:
  ```js
  const results = await Promise.all(items.map(fetchOne));
  ```
- Always handle rejection — an unhandled rejection can crash Node.

## Useful array methods
- `map` transforms, `filter` selects, `reduce` folds.
- `find` returns the first match; `some` / `every` return booleans.

## Gotcha
- `typeof null === 'object'` — a historical bug that's here forever.
