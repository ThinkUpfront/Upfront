# .NET / C# Conventions

Sources: Microsoft .NET Guidelines, Framework Design Guidelines, xUnit patterns

## Commands

```
dotnet build
dotnet test
dotnet format --verify-no-changes
dotnet run --project src/MyApp
```

## Project Structure

```
src/
  MyApp/
    Controllers/       — API controllers
    Services/          — Business logic
    Models/            — Domain entities
    DTOs/              — Data transfer objects
    Data/              — DbContext, migrations, repositories
    Extensions/        — Extension method classes
    Program.cs         — Host configuration
tests/
  MyApp.Tests/
    Unit/              — Fast, isolated tests
    Integration/       — Tests hitting real dependencies
```

## Classes

```csharp
// Good — records for immutable data (C# 9+)
public record User(string Id, string Name, string Email);

// Good — primary constructors (C# 12+)
public class OrderService(IOrderRepository repo, IPaymentGateway payments)
{
    public async Task<Order> PlaceAsync(OrderRequest request) { ... }
}

// Good — sealed unless designed for inheritance
public sealed class EmailValidator { ... }

// Bad — mutable DTOs with public setters everywhere
public class User
{
    public string Name { get; set; }  // anyone can mutate
}

// Bad — service locator pattern
var service = ServiceLocator.Get<IOrderService>();
```

## Dependency Injection

```csharp
// Good — constructor injection (automatic with primary constructors)
public class OrderService(IOrderRepository repo, ILogger<OrderService> logger)
{
    public async Task<Order> GetAsync(string id)
    {
        logger.LogInformation("Fetching order {OrderId}", id);
        return await repo.GetByIdAsync(id)
            ?? throw new OrderNotFoundException(id);
    }
}

// Good — register in Program.cs
builder.Services.AddScoped<IOrderRepository, SqlOrderRepository>();
builder.Services.AddScoped<OrderService>();

// Bad — new-ing up dependencies
public class OrderService
{
    private readonly SqlOrderRepository _repo = new();  // untestable
}
```

## Error Handling

```csharp
// Good — custom exceptions with context
public class OrderNotFoundException : Exception
{
    public string OrderId { get; }
    public OrderNotFoundException(string orderId)
        : base($"Order {orderId} not found") => OrderId = orderId;
}

// Good — Result pattern for expected failures
public record Result<T>
{
    public T? Value { get; init; }
    public string? Error { get; init; }
    public bool IsSuccess => Error is null;

    public static Result<T> Ok(T value) => new() { Value = value };
    public static Result<T> Fail(string error) => new() { Error = error };
}

// Bad — catch Exception
try { ... } catch (Exception ex) { _logger.LogError(ex, "failed"); }

// Bad — throwing for flow control
try { var user = GetUser(id); }
catch (UserNotFoundException) { return NotFound(); }
// Use: TryGetUser or return null
```

## Async

```csharp
// Good — async all the way down
public async Task<User> GetUserAsync(string id)
{
    var entity = await _context.Users.FindAsync(id);
    return entity is null ? throw new UserNotFoundException(id) : Map(entity);
}

// Good — CancellationToken on all async public APIs
public async Task<List<Order>> ListAsync(CancellationToken ct = default)
{
    return await _context.Orders.ToListAsync(ct);
}

// Bad — .Result or .Wait() (deadlock risk)
var user = GetUserAsync(id).Result;

// Bad — async void (except event handlers)
async void ProcessOrder() { ... }
```

## Testing

```csharp
// Good — Arrange-Act-Assert with xUnit
[Fact]
public async Task PlaceOrder_RejectsWhenOutOfStock()
{
    // Arrange
    var repo = new FakeOrderRepository();
    var service = new OrderService(repo, NullLogger<OrderService>.Instance);
    var request = new OrderRequest("SKU-1", Quantity: 5);

    // Act
    var result = await service.PlaceAsync(request);

    // Assert
    Assert.False(result.IsSuccess);
    Assert.Equal("Insufficient stock", result.Error);
}

// Good — Theory for parameterized tests
[Theory]
[InlineData("", false)]
[InlineData("bad", false)]
[InlineData("user@example.com", true)]
public void ValidateEmail(string email, bool expected)
{
    Assert.Equal(expected, EmailValidator.IsValid(email));
}

// Good — test doubles via interfaces, not mocking frameworks
public class FakeOrderRepository : IOrderRepository
{
    private readonly List<Order> _orders = [];
    public Task<Order?> GetByIdAsync(string id)
        => Task.FromResult(_orders.FirstOrDefault(o => o.Id == id));
}

// Bad — mocking everything with Moq
// Bad — [SetUp] and [TearDown] (use constructor/IDisposable)
```

## LINQ

```csharp
// Good — method syntax for simple queries
var activeUsers = users
    .Where(u => u.IsActive)
    .OrderBy(u => u.Name)
    .ToList();

// Good — query syntax for joins
var results = from o in orders
              join u in users on o.UserId equals u.Id
              select new { o.Total, u.Name };

// Bad — LINQ for simple iteration
users.ToList().ForEach(u => Process(u));
// Use: foreach (var user in users) { Process(user); }
```

## Style

- Classes/records/interfaces: `PascalCase`
- Methods/properties: `PascalCase`
- Private fields: `_camelCase` with underscore prefix
- Local variables/parameters: `camelCase`
- Constants: `PascalCase` (not UPPER_SNAKE)
- Interfaces: `I` prefix — `IOrderRepository`
- Async methods: `Async` suffix — `GetUserAsync`
- Namespaces: match folder structure
- File-scoped namespaces (C# 10+): `namespace MyApp.Services;`
- `var` for local variables when type is obvious
- Null-conditional: `user?.Name` over null checks
- Collection expressions (C# 12+): `List<int> nums = [1, 2, 3];`
- Pattern matching over type checks: `if (shape is Circle { Radius: > 0 } c)`
