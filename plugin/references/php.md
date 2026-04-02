# PHP Conventions

Sources: PHP-FIG PSR standards, Laravel conventions, PHPStan/Psalm

## Commands

```
composer install
php artisan test                # Laravel
./vendor/bin/phpunit            # Generic
./vendor/bin/phpstan analyse    # Static analysis
./vendor/bin/pint               # Laravel formatter
./vendor/bin/php-cs-fixer fix   # Generic formatter
```

## Project Structure (Laravel)

```
app/
  Http/
    Controllers/     — Request handling, thin
    Middleware/       — Request/response filters
    Requests/        — Form request validation
  Models/            — Eloquent models
  Services/          — Business logic
  Actions/           — Single-purpose action classes
  Exceptions/        — Custom exceptions
  Events/            — Domain events
  Listeners/         — Event handlers
database/
  migrations/        — Schema changes
  factories/         — Test data factories
  seeders/           — Sample data
tests/
  Unit/              — Fast, isolated
  Feature/           — HTTP tests with full app
routes/
  web.php            — Web routes
  api.php            — API routes
```

## Controllers

```php
// Good — thin controllers, delegate to services/actions
class OrderController extends Controller
{
    public function store(
        StoreOrderRequest $request,
        PlaceOrderAction $action
    ): JsonResponse {
        $order = $action->execute($request->validated());
        return response()->json($order, 201);
    }
}

// Good — single-action controllers for simple endpoints
class ShowDashboardController extends Controller
{
    public function __invoke(Request $request): View
    {
        return view('dashboard', [
            'stats' => DashboardService::getStats($request->user()),
        ]);
    }
}

// Bad — fat controllers with business logic
class OrderController extends Controller
{
    public function store(Request $request)
    {
        $validated = $request->validate([...]);
        // 50 lines of business logic here
        // email sending here
        // payment processing here
    }
}
```

## Type Safety

```php
// Good — strict types in every file
declare(strict_types=1);

// Good — typed properties and return types (PHP 8.0+)
class User
{
    public function __construct(
        public readonly string $id,
        public readonly string $name,
        public readonly string $email,
    ) {}
}

// Good — union types and null safety (PHP 8.0+)
function findUser(string $id): User|null
{
    return User::find($id);
}

// Good — enums (PHP 8.1+)
enum OrderStatus: string
{
    case Pending = 'pending';
    case Confirmed = 'confirmed';
    case Shipped = 'shipped';
    case Cancelled = 'cancelled';
}

// Bad — no type declarations
function getUser($id) { ... }

// Bad — mixed type as escape hatch
function process(mixed $data): mixed { ... }
```

## Error Handling

```php
// Good — custom exceptions with context
class OrderNotFoundException extends RuntimeException
{
    public function __construct(
        public readonly string $orderId,
    ) {
        parent::__construct("Order {$orderId} not found");
    }
}

// Good — Laravel exception rendering
class InsufficientStockException extends RuntimeException
{
    public function render(): JsonResponse
    {
        return response()->json([
            'error' => $this->getMessage(),
        ], 422);
    }
}

// Bad — generic exceptions
throw new \Exception('something went wrong');

// Bad — swallowing exceptions
try { ... } catch (\Exception $e) { /* ignore */ }
```

## Eloquent

```php
// Good — scopes for reusable queries
class User extends Model
{
    public function scopeActive(Builder $query): Builder
    {
        return $query->where('is_active', true);
    }

    public function scopeCreatedAfter(Builder $query, Carbon $date): Builder
    {
        return $query->where('created_at', '>', $date);
    }
}

// Usage
$users = User::active()->createdAfter(now()->subMonth())->get();

// Good — eager loading to prevent N+1
$orders = Order::with(['user', 'items.product'])->get();

// Bad — N+1 queries
$orders = Order::all();
foreach ($orders as $order) {
    echo $order->user->name;  // query per iteration
}

// Bad — raw queries when Eloquent works
DB::select('SELECT * FROM users WHERE active = 1');
// Use: User::where('active', true)->get();
```

## Testing

```php
// Good — Feature test with full HTTP cycle
public function test_user_can_place_order(): void
{
    $user = User::factory()->create();
    $product = Product::factory()->create(['stock' => 10]);

    $response = $this->actingAs($user)->postJson('/api/orders', [
        'product_id' => $product->id,
        'quantity' => 2,
    ]);

    $response->assertStatus(201)
        ->assertJsonPath('quantity', 2);
    $this->assertDatabaseHas('orders', [
        'user_id' => $user->id,
        'product_id' => $product->id,
    ]);
}

// Good — factories for test data
$user = User::factory()
    ->has(Order::factory()->count(3))
    ->create();

// Good — unit test for service logic
public function test_discount_applies_for_bulk_orders(): void
{
    $calculator = new PriceCalculator();
    $price = $calculator->calculate(quantity: 100, unitPrice: 10.00);
    $this->assertEquals(900.00, $price);  // 10% bulk discount
}

// Bad — testing framework internals
// Bad — no assertions (test just "runs")
// Bad — sharing state between tests
```

## Validation

```php
// Good — Form Request classes
class StoreOrderRequest extends FormRequest
{
    public function rules(): array
    {
        return [
            'product_id' => ['required', 'exists:products,id'],
            'quantity' => ['required', 'integer', 'min:1', 'max:100'],
            'notes' => ['nullable', 'string', 'max:500'],
        ];
    }
}

// Bad — inline validation in controllers
$request->validate([
    'product_id' => 'required|exists:products,id',
    // 20 more rules cluttering the controller
]);
```

## Style

- Classes: `PascalCase` — `OrderService`
- Methods/functions: `camelCase` — `placeOrder`
- Properties: `camelCase` — `$orderTotal`
- Constants: `UPPER_SNAKE_CASE`
- Routes: `kebab-case` — `/api/order-items`
- Database tables: `snake_case`, plural — `order_items`
- Database columns: `snake_case` — `created_at`
- Config keys: `snake_case` — `config('app.debug')`
- One class per file
- `declare(strict_types=1)` in every file
- Named arguments for clarity: `new Money(amount: 100, currency: 'USD')`
- Match expression over switch (PHP 8.0+)
- Readonly properties/classes where immutability is desired (PHP 8.2+)
- PHPStan level 8+ or Psalm level 1-2 for static analysis
