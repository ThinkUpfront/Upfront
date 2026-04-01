# Java Conventions

Sources: Google Java Style Guide, Effective Java (Bloch), Spring conventions

## Commands

```
./gradlew build
./gradlew test
./gradlew check          # includes spotbugs, checkstyle
mvn clean verify         # Maven alternative
```

## Project Structure

```
src/
  main/
    java/com/example/myapp/
      controller/        — REST controllers
      service/           — Business logic
      repository/        — Data access (JPA/JDBC)
      model/             — Domain entities
      dto/               — Data transfer objects
      config/            — Spring configuration
      exception/         — Custom exceptions
  test/
    java/com/example/myapp/
      controller/        — Controller tests (MockMvc)
      service/           — Service unit tests
      integration/       — Integration tests
```

## Classes

```java
// Good — immutable value objects (Java 16+)
public record User(String id, String name, String email) {}

// Good — sealed interfaces for closed hierarchies (Java 17+)
public sealed interface Shape permits Circle, Rectangle, Triangle {}
public record Circle(double radius) implements Shape {}
public record Rectangle(double width, double height) implements Shape {}

// Good — builder for complex construction
User user = User.builder()
    .name("Alice")
    .email("alice@example.com")
    .build();

// Bad — mutable POJOs with getters/setters for everything
public class User {
    private String name;
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
}
```

## Error Handling

```java
// Good — specific exceptions with context
public class UserNotFoundException extends RuntimeException {
    private final String userId;
    public UserNotFoundException(String userId) {
        super("User not found: " + userId);
        this.userId = userId;
    }
    public String getUserId() { return userId; }
}

// Good — use Optional for nullable returns
public Optional<User> findById(String id) {
    return Optional.ofNullable(store.get(id));
}

// Bad — returning null
public User findById(String id) {
    return store.get(id);  // caller must remember to null-check
}

// Bad — catch Exception
try { ... } catch (Exception e) { log.error("failed", e); }

// Bad — swallowing exceptions
try { ... } catch (IOException e) { /* ignore */ }
```

## Dependency Injection

```java
// Good — constructor injection
@Service
public class OrderService {
    private final UserRepository userRepo;
    private final PaymentGateway payments;

    public OrderService(UserRepository userRepo, PaymentGateway payments) {
        this.userRepo = userRepo;
        this.payments = payments;
    }
}

// Bad — field injection
@Service
public class OrderService {
    @Autowired
    private UserRepository userRepo;  // untestable without Spring
}
```

## Testing

```java
// Good — descriptive test names
@Test
void shouldRejectOrderWhenInsufficientStock() {
    var product = new Product("SKU-1", 0);
    var order = new Order(product, 5);

    var result = service.place(order);

    assertThat(result.isRejected()).isTrue();
    assertThat(result.reason()).isEqualTo("Insufficient stock");
}

// Good — AssertJ for fluent assertions
assertThat(users)
    .hasSize(3)
    .extracting(User::name)
    .containsExactly("Alice", "Bob", "Charlie");

// Good — MockMvc for controller tests
mockMvc.perform(get("/api/users/{id}", "123"))
    .andExpect(status().isOk())
    .andExpect(jsonPath("$.name").value("Alice"));

// Good — @SpringBootTest only for integration tests
// Use plain unit tests with mocks for service logic

// Bad — testing private methods
// Bad — @SpringBootTest for every test (slow)
```

## Streams

```java
// Good — streams for transformations
List<String> names = users.stream()
    .filter(u -> u.isActive())
    .map(User::name)
    .sorted()
    .toList();

// Bad — stream for simple iteration
users.stream().forEach(u -> process(u));
// Use: users.forEach(u -> process(u));
// Or:  for (var user : users) { process(user); }
```

## Style

- Classes: `PascalCase`
- Methods/fields: `camelCase`
- Constants: `UPPER_SNAKE_CASE`
- Packages: all lowercase, no underscores
- One top-level class per file
- Use `var` for local variables when type is obvious (Java 10+)
- Use text blocks for multiline strings (Java 15+)
- Prefer `List.of()`, `Map.of()` over `Arrays.asList()`, `Collections.unmodifiable*()`
- No wildcard imports: `import java.util.*`
- Annotations order: `@Override` before `@Deprecated` before `@SuppressWarnings`
