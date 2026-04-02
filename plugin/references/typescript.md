# TypeScript Conventions

Sources: Google TypeScript Style Guide, TypeScript handbook, ts-reset

## Types

```ts
// Good — interface for objects, type for unions/intersections
interface User {
  id: string;
  name: string;
  role: 'admin' | 'member';
}

type Result<T> = { ok: true; value: T } | { ok: false; error: Error };

// Bad — any
function process(data: any): any { ... }

// Bad — type assertion to bypass checks
const user = data as User;  // only when you truly know better than TS

// Good — type narrowing
function isUser(value: unknown): value is User {
  return typeof value === 'object' && value !== null && 'id' in value;
}
```

## Functions

```ts
// Good — explicit return types on public APIs
function parseConfig(raw: string): Config { ... }

// Good — inferred return types for simple internal functions
const double = (n: number) => n * 2;

// Good — use readonly for arrays/objects you don't mutate
function sum(items: readonly number[]): number {
  return items.reduce((a, b) => a + b, 0);
}

// Bad — optional params before required params
function create(name?: string, id: string) { ... }
```

## Nullability

```ts
// Good — explicit null checks
function getUser(id: string): User | null {
  const user = db.find(id);
  return user ?? null;
}

// Good — optional chaining
const city = user?.address?.city;

// Bad — non-null assertion to hide problems
const name = user!.name;  // only when truly guaranteed
```

## Enums vs Unions

```ts
// Prefer — union types (tree-shakeable, no runtime cost)
type Status = 'active' | 'inactive' | 'pending';

// Acceptable — const enum for numeric values
const enum Direction { Up, Down, Left, Right }

// Avoid — regular enum (generates runtime object)
enum Color { Red, Blue }
```

## Generics

```ts
// Good — constrained generics
function first<T>(items: readonly T[]): T | undefined {
  return items[0];
}

// Good — meaningful generic names for complex signatures
function merge<TBase, TOverride>(base: TBase, override: TOverride): TBase & TOverride { ... }

// Bad — unconstrained, single-letter in complex cases
function process<A, B, C>(a: A, b: B): C { ... }
```

## Error Handling

```ts
// Good — typed errors
class NotFoundError extends Error {
  constructor(public readonly resource: string, public readonly id: string) {
    super(`${resource} ${id} not found`);
    this.name = 'NotFoundError';
  }
}

// Good — Result type instead of throwing
type Result<T, E = Error> = { ok: true; value: T } | { ok: false; error: E };

// Bad — throwing strings
throw 'something went wrong';

// Bad — catch without type narrowing
catch (e) { console.log(e.message); }  // e is unknown
```

## Imports

```ts
// Good — type-only imports (stripped at compile time)
import type { User } from './types';

// Good — named imports
import { formatDate, parseDate } from './utils/date';

// Bad — namespace imports for large modules
import * as utils from './utils';

// Bad — relative paths climbing more than 2 levels
import { foo } from '../../../shared/utils';  // use path aliases
```

## Style

- No `I` prefix on interfaces: `User` not `IUser`
- No `Enum` suffix: `Status` not `StatusEnum`
- Prefer `unknown` over `any`
- Prefer `undefined` over `null` for optional values (match TS conventions)
- Use `satisfies` for type-safe object literals: `const config = { ... } satisfies Config`
- Strict mode always: `"strict": true` in tsconfig
