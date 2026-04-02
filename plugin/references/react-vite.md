# React + Vite Conventions

Sources: bulletproof-react, Google TypeScript Style Guide, React Testing Library docs

## Commands

```
npm run dev
npm run build
npm run test
npm run lint
npx tsc --noEmit
```

## Project Structure

```
src/
  components/     — Shared UI components (Button, Modal, etc.)
  features/       — Feature modules, each self-contained
    auth/
      components/ — Feature-specific components
      hooks/      — Feature-specific hooks
      api/        — Feature-specific API calls
      index.ts    — Public API (re-exports)
  hooks/          — Shared hooks
  lib/            — Configured library instances (axios, dayjs)
  types/          — Shared TypeScript types
  utils/          — Pure utility functions
  routes/         — Route definitions
```

## Components

```tsx
// Good — props interface, destructured, named export
interface UserCardProps {
  user: User;
  onSelect: (id: string) => void;
}

export function UserCard({ user, onSelect }: UserCardProps) {
  return (
    <button onClick={() => onSelect(user.id)}>
      {user.name}
    </button>
  );
}

// Bad — default export, inline types, generic names
export default function Card({ data, cb }: any) { ... }
```

## Hooks

```tsx
// Good — starts with "use", returns typed value, handles cleanup
function useDebounce<T>(value: T, delay: number): T {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const timer = setTimeout(() => setDebounced(value), delay);
    return () => clearTimeout(timer);
  }, [value, delay]);
  return debounced;
}

// Bad — doesn't clean up, mutation inside render
```

## State Management

```tsx
// Good — colocate state, lift only when needed
// Local state for UI concerns (open/closed, form input)
const [isOpen, setIsOpen] = useState(false);

// Server state via TanStack Query, not manual useEffect+fetch
const { data, isLoading } = useQuery({
  queryKey: ['users', filters],
  queryFn: () => fetchUsers(filters),
});

// Bad — global state for everything
// Bad — useEffect + useState for data fetching
useEffect(() => {
  fetch('/api/users').then(r => r.json()).then(setUsers);
}, []);
```

## Testing

```tsx
// Good — test behavior, not implementation
test('disables submit when email is invalid', async () => {
  render(<SignupForm />);
  await userEvent.type(screen.getByLabelText('Email'), 'bad');
  expect(screen.getByRole('button', { name: /submit/i })).toBeDisabled();
});

// Good — use role queries, not test IDs
screen.getByRole('button', { name: /save/i });
screen.getByLabelText('Email');

// Bad — testing state directly
expect(component.state.isValid).toBe(false);

// Bad — snapshot tests for dynamic components
expect(tree).toMatchSnapshot();

// Bad — querySelector or data-testid as first choice
document.querySelector('.submit-btn');
```

## Error Handling

```tsx
// Good — error boundaries for UI regions
<ErrorBoundary fallback={<ErrorFallback />}>
  <Dashboard />
</ErrorBoundary>

// Good — form validation with clear messages
if (!email.includes('@')) {
  setError('email', { message: 'Enter a valid email address' });
  return;
}

// Bad — silent catch
try { await save(); } catch {}
```

## Performance

```tsx
// Good — lazy load routes
const Dashboard = lazy(() => import('./features/dashboard'));

// Good — memoize expensive computations
const sorted = useMemo(() => items.sort(compareFn), [items]);

// Bad — premature React.memo on everything
// Bad — inline object/array literals in JSX props (new ref every render)
<Chart options={{ animate: true }} />  // re-renders Chart every time
```

## Style

- Components: `PascalCase` — `UserCard.tsx`
- Hooks: `camelCase` with `use` prefix — `useAuth.ts`
- Utils: `camelCase` — `formatDate.ts`
- Constants: `UPPER_SNAKE_CASE`
- Types/Interfaces: `PascalCase`, no `I` prefix
- Event handlers: `onVerb` for props, `handleVerb` for implementation
- Boolean props: `isX`, `hasX`, `canX`, `shouldX`
- Files: one component per file, filename matches export name
