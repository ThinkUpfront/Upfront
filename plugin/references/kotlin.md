# Kotlin Conventions

Sources: Kotlin Official Style Guide, Android Kotlin Guides, Spring Kotlin support

## Commands

```
./gradlew build
./gradlew test
./gradlew detekt          # static analysis
./gradlew ktlintCheck     # formatting check
```

## Project Structure

```
src/
  main/kotlin/com/example/myapp/
    controller/        — REST controllers
    service/           — Business logic
    repository/        — Data access
    model/             — Domain types (data classes, sealed classes)
    dto/               — Request/response DTOs
    config/            — Spring/DI configuration
    extension/         — Extension functions
  test/kotlin/com/example/myapp/
    controller/        — Controller tests
    service/           — Service unit tests
```

## Data Classes

```kotlin
// Good — data class for value objects
data class User(
    val id: String,
    val name: String,
    val email: String,
)

// Good — copy for immutable updates
val updated = user.copy(name = "New Name")

// Good — destructuring
val (id, name, email) = user

// Bad — mutable data class
data class User(
    var id: String,    // var = mutable, avoid
    var name: String,
)
```

## Null Safety

```kotlin
// Good — leverage the type system
fun findUser(id: String): User? {
    return repository.findById(id)
}

// Good — safe calls and elvis operator
val city = user?.address?.city ?: "Unknown"

// Good — let for null-safe operations
user?.let { sendWelcomeEmail(it) }

// Good — require/check for preconditions
fun process(order: Order) {
    require(order.items.isNotEmpty()) { "Order must have items" }
    check(order.status == Status.PENDING) { "Order must be pending" }
}

// Bad — !! operator (defeats null safety)
val name = user!!.name

// Bad — platform types without null check
val length = javaObject.getString().length  // might be null from Java
```

## Sealed Classes

```kotlin
// Good — sealed for closed hierarchies
sealed class Result<out T> {
    data class Success<T>(val value: T) : Result<T>()
    data class Failure(val error: Throwable) : Result<Nothing>()
}

// Good — exhaustive when
fun handle(result: Result<User>) = when (result) {
    is Result.Success -> showUser(result.value)
    is Result.Failure -> showError(result.error)
}

// Good — sealed interface for state machines
sealed interface ConnectionState {
    data object Disconnected : ConnectionState
    data class Connecting(val attempt: Int) : ConnectionState
    data class Connected(val sessionId: String) : ConnectionState
}
```

## Extension Functions

```kotlin
// Good — extend existing types for domain clarity
fun String.toSlug(): String =
    lowercase().replace(Regex("[^a-z0-9]+"), "-").trim('-')

fun Instant.isExpired(): Boolean =
    this < Instant.now()

// Good — scoped extensions (inside a class, not polluting global)
class OrderService {
    private fun Order.totalWithTax(): BigDecimal =
        total * (BigDecimal.ONE + taxRate)
}

// Bad — extension that accesses private state of another class
// Bad — replacing simple functions with extensions for no reason
```

## Coroutines

```kotlin
// Good — structured concurrency
suspend fun fetchDashboard(): Dashboard = coroutineScope {
    val user = async { userService.getCurrent() }
    val orders = async { orderService.getRecent() }
    Dashboard(user.await(), orders.await())
}

// Good — use withContext for thread switching
suspend fun readFile(path: Path): String =
    withContext(Dispatchers.IO) {
        path.readText()
    }

// Good — Flow for reactive streams
fun observeOrders(): Flow<Order> = flow {
    while (true) {
        emit(orderRepository.getLatest())
        delay(5.seconds)
    }
}

// Bad — GlobalScope (leaks coroutines)
GlobalScope.launch { doSomething() }

// Bad — runBlocking in production code (blocks the thread)
val result = runBlocking { fetchData() }
```

## Dependency Injection (Spring)

```kotlin
// Good — constructor injection (Spring auto-detects)
@Service
class OrderService(
    private val repo: OrderRepository,
    private val payments: PaymentGateway,
    private val logger: Logger = LoggerFactory.getLogger(OrderService::class.java),
) {
    suspend fun place(request: OrderRequest): Order { ... }
}

// Bad — lateinit var injection
@Service
class OrderService {
    @Autowired
    lateinit var repo: OrderRepository  // nullable in disguise
}
```

## Testing

```kotlin
// Good — descriptive names with backticks
@Test
fun `should reject order when stock is insufficient`() {
    val product = Product(sku = "SKU-1", stock = 0)
    val service = OrderService(FakeRepository(product))

    val result = service.place(OrderRequest("SKU-1", quantity = 5))

    assertThat(result).isInstanceOf(Result.Failure::class.java)
    assertThat((result as Result.Failure).error.message)
        .contains("Insufficient stock")
}

// Good — coroutine tests
@Test
fun `should fetch user concurrently`() = runTest {
    val service = UserService(FakeApi())
    val user = service.fetchUser("123")
    assertThat(user.name).isEqualTo("Alice")
}

// Good — test data with copy
private val baseUser = User(id = "1", name = "Alice", email = "a@test.com")

@Test
fun `should reject empty name`() {
    val user = baseUser.copy(name = "")
    assertThat(validate(user)).isFalse()
}

// Bad — mocking data classes (just construct them)
// Bad — PowerMock or reflection-based mocking
```

## Scope Functions

```kotlin
// Good — apply for object configuration
val client = HttpClient().apply {
    timeout = 30.seconds
    retries = 3
    baseUrl = "https://api.example.com"
}

// Good — let for null-safe transformations
val length = name?.let { it.trim().length }

// Good — also for side effects
val user = createUser(request).also {
    logger.info("Created user ${it.id}")
}

// Bad — nested scope functions (unreadable)
user?.let { u ->
    u.address?.let { a ->
        a.city?.let { c ->
            // deeply nested, just use if/when
        }
    }
}
```

## Style

- Classes/interfaces: `PascalCase` — `OrderService`
- Functions/properties: `camelCase` — `placeOrder`
- Constants: `UPPER_SNAKE_CASE` or `PascalCase` for objects
- Packages: all lowercase — `com.example.myapp`
- Files: `PascalCase` matching class name, or `camelCase` for top-level functions
- Test names: backtick descriptive — `` `should do X when Y` ``
- Prefer `val` over `var` (immutable by default)
- Prefer expression body for single-expression functions: `fun double(n: Int) = n * 2`
- Trailing commas in parameter lists and collections
- Named arguments when meaning isn't obvious: `createUser(name = "Alice", admin = false)`
- Use `object` for singletons, not companion object factories
- Prefer `when` over `if-else` chains
