# Rust Conventions

Sources: Rust API Guidelines, rust-clippy, Rust by Example

## Commands

```
cargo build
cargo test
cargo clippy -- -D warnings
cargo fmt --check
cargo doc --no-deps
```

## Project Structure

```
src/
  main.rs            — Binary entrypoint (thin, calls lib)
  lib.rs             — Library root, public API
  config.rs          — Configuration types
  error.rs           — Error types
  models/            — Domain types
  handlers/          — Request handlers (Axum/Actix)
  db/                — Database access
tests/
  integration/       — Integration tests
benches/             — Benchmarks
```

## Error Handling

```rust
// Good — thiserror for library errors
#[derive(Debug, thiserror::Error)]
pub enum AppError {
    #[error("user {0} not found")]
    NotFound(String),
    #[error("database error")]
    Db(#[from] sqlx::Error),
    #[error("invalid input: {0}")]
    Validation(String),
}

// Good — anyhow for application code
fn main() -> anyhow::Result<()> {
    let config = load_config().context("failed to load config")?;
    run(config)?;
    Ok(())
}

// Good — ? operator with context
let file = File::open(path)
    .with_context(|| format!("failed to open {}", path.display()))?;

// Bad — unwrap in production code
let config = load_config().unwrap();

// Bad — generic error strings without context
Err("something went wrong".into())
```

## Ownership

```rust
// Good — borrow when you don't need ownership
fn greet(name: &str) {
    println!("Hello, {name}");
}

// Good — take ownership when you need to store it
struct Cache {
    entries: Vec<Entry>,
}

impl Cache {
    fn add(&mut self, entry: Entry) {  // takes ownership
        self.entries.push(entry);
    }
}

// Good — clone explicitly when needed, don't hide it
let name = user.name.clone();

// Bad — &String instead of &str
fn greet(name: &String) { ... }

// Bad — unnecessary cloning to satisfy borrow checker
// (usually means the design needs rethinking)
```

## Traits

```rust
// Good — implement standard traits (Debug, Display, Clone, PartialEq)
#[derive(Debug, Clone, PartialEq)]
pub struct Config {
    pub host: String,
    pub port: u16,
}

impl fmt::Display for Config {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}:{}", self.host, self.port)
    }
}

// Good — small, focused traits
trait Validate {
    fn validate(&self) -> Result<(), ValidationError>;
}

// Bad — god trait with 15 methods
// Bad — missing Debug on public types
```

## Enums

```rust
// Good — enums for state machines
enum ConnectionState {
    Disconnected,
    Connecting { attempt: u32 },
    Connected { session_id: String },
    Failed { error: String, retries: u32 },
}

// Good — exhaustive matching
match state {
    ConnectionState::Disconnected => reconnect(),
    ConnectionState::Connecting { attempt } if attempt > 3 => fail(),
    ConnectionState::Connecting { .. } => wait(),
    ConnectionState::Connected { session_id } => use_session(session_id),
    ConnectionState::Failed { retries, .. } if retries > 0 => retry(),
    ConnectionState::Failed { error, .. } => log_error(error),
}
```

## Testing

```rust
#[cfg(test)]
mod tests {
    use super::*;

    // Good — descriptive test names
    #[test]
    fn parse_config_rejects_negative_port() {
        let input = r#"{"host": "localhost", "port": -1}"#;
        let result = Config::parse(input);
        assert!(result.is_err());
        assert!(result.unwrap_err().to_string().contains("port"));
    }

    // Good — test helpers in the test module
    fn sample_config() -> Config {
        Config { host: "localhost".into(), port: 8080 }
    }

    // Good — proptest for property-based testing
    proptest! {
        #[test]
        fn roundtrip_serialization(config: Config) {
            let json = serde_json::to_string(&config).unwrap();
            let parsed: Config = serde_json::from_str(&json).unwrap();
            assert_eq!(config, parsed);
        }
    }
}

// Integration tests in tests/ directory
// test binary per file, access only public API
```

## Style

- Types: `PascalCase` — `HttpClient`
- Functions/methods: `snake_case` — `fetch_user`
- Constants: `UPPER_SNAKE_CASE` — `MAX_RETRIES`
- Lifetimes: short lowercase — `'a`, `'de`
- Crates: `kebab-case` in Cargo.toml, `snake_case` in code
- No `get_` prefix on getters: `fn name(&self) -> &str`
- Builder pattern for complex structs: `Config::builder().host("localhost").build()`
- Prefer `impl Trait` in argument position over generics for simple cases
- Clippy at `warn` level minimum, `deny` in CI
- `#[must_use]` on functions where ignoring the return value is likely a bug
