# TypeScript notes

## The type system is structural
- TS cares about *shape*, not name. Two types with the same members are
  compatible even if declared separately ("duck typing").
- `interface` and `type` overlap a lot. Rule of thumb: `interface` for object
  shapes you might extend/merge, `type` for unions, primitives, and tuples.

## Narrowing
```ts
function len(x: string | string[]) {
  if (typeof x === "string") return x.length;   // x narrowed to string
  return x.reduce((n, s) => n + s.length, 0);   // x narrowed to string[]
}
```
- `typeof`, `instanceof`, `in`, and truthiness all narrow.
- A custom type guard returns `x is Foo` to teach the compiler.

## Useful utility types
- `Partial<T>` / `Required<T>` — toggle optionality of every field.
- `Pick<T, K>` / `Omit<T, K>` — keep or drop named keys.
- `Record<K, V>` — an object map with known keys.
- `ReturnType<typeof fn>` — reuse a function's return type without restating it.

## Gotchas
- `any` disables checking and spreads silently — prefer `unknown` and narrow.
- Type assertions (`x as Foo`) are a promise to the compiler, not a conversion;
  a wrong one fails at runtime, not compile time.
- `strictNullChecks` is where most real bugs get caught — keep it on.
