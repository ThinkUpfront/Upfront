# Python Conventions

Sources: Google Python Style Guide, PEP 8, Hypermodern Python

## Commands

```
python -m pytest -v
python -m pytest --cov=src
python -m mypy src/
ruff check .
ruff format .
```

## Project Structure

```
src/
  mypackage/
    __init__.py
    models.py        — Domain types (dataclasses/Pydantic)
    services.py      — Business logic
    api/             — Route handlers
    db/              — Database access
tests/
  test_models.py     — Mirrors src/ structure
  conftest.py        — Shared fixtures
pyproject.toml       — Single config file for everything
```

## Type Hints

```python
# Good — type hints on all public functions
def fetch_user(user_id: str) -> User | None:
    ...

# Good — use modern syntax (3.10+)
def process(items: list[str]) -> dict[str, int]:
    ...

# Good — TypeVar for generics
from typing import TypeVar
T = TypeVar('T')
def first(items: list[T]) -> T | None:
    return items[0] if items else None

# Bad — no type hints on public API
def fetch_user(user_id):
    ...

# Bad — old-style Optional
from typing import Optional, List, Dict
def process(items: List[str]) -> Optional[str]:
    ...
```

## Error Handling

```python
# Good — specific exceptions
class UserNotFoundError(Exception):
    def __init__(self, user_id: str) -> None:
        self.user_id = user_id
        super().__init__(f"User {user_id} not found")

# Good — catch specific exceptions
try:
    user = repo.get(user_id)
except UserNotFoundError:
    return None

# Bad — bare except
try:
    do_something()
except:
    pass

# Bad — catching Exception for flow control
try:
    value = my_dict[key]
except Exception:
    value = default
# Use: value = my_dict.get(key, default)
```

## Classes

```python
# Good — dataclasses for data containers
from dataclasses import dataclass

@dataclass(frozen=True)
class Point:
    x: float
    y: float

# Good — Pydantic for validation/serialization
from pydantic import BaseModel

class CreateUserRequest(BaseModel):
    name: str
    email: str

# Bad — manual __init__, __repr__, __eq__ for data classes
class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y
    def __repr__(self): ...
    def __eq__(self, other): ...
```

## Testing

```python
# Good — pytest fixtures
@pytest.fixture
def db_session():
    session = create_test_session()
    yield session
    session.rollback()

def test_create_user(db_session):
    user = create_user(db_session, name="Alice")
    assert user.name == "Alice"
    assert user.id is not None

# Good — parametrize for table-driven tests
@pytest.mark.parametrize("input,expected", [
    ("hello", "HELLO"),
    ("world", "WORLD"),
    ("", ""),
])
def test_uppercase(input, expected):
    assert uppercase(input) == expected

# Bad — unittest.TestCase
class TestUser(unittest.TestCase):
    def setUp(self): ...

# Bad — mocking internal modules
@patch('mypackage.services.repository.get')
```

## Async

```python
# Good — async where I/O bound
async def fetch_users(client: httpx.AsyncClient) -> list[User]:
    response = await client.get("/users")
    response.raise_for_status()
    return [User(**u) for u in response.json()]

# Good — gather for concurrent I/O
results = await asyncio.gather(
    fetch_users(client),
    fetch_roles(client),
)

# Bad — async for CPU-bound work (use multiprocessing)
# Bad — mixing sync and async I/O
```

## Style

- Modules/packages: `snake_case` — `user_service.py`
- Classes: `PascalCase`
- Functions/variables: `snake_case`
- Constants: `UPPER_SNAKE_CASE`
- Private: single underscore prefix `_internal_method`
- No abbreviations: `calculate_total` not `calc_tot`
- Imports: stdlib → third-party → local, separated by blank lines
- Use `pathlib.Path` over `os.path`
- Use f-strings over `.format()` or `%`
- Comprehensions over `map`/`filter` for simple cases
