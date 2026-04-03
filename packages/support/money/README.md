# Money

A Go implementation of Martin Fowler's Money pattern, inspired by [moneyphp/money](https://github.com/moneyphp/money).

[![Go Report Card](https://goreportcard.com/badge/github.com/gollin/packages/support/money)](https://goreportcard.com/report/github.com/gollin/packages/support/money)
[![Go Reference](https://pkg.go.dev/badge/github.com/gollin/packages/support/money.svg)](https://pkg.go.dev/github.com/gollin/packages/support/money)

## 📖 Table of Contents

- [Why Use a Money Library?](#why-use-a-money-library)
- [Installation](#installation)
- [Requirements](#requirements)
- [Features](#features)
- [Packages](#packages)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
  - [Creating Money](#creating-money)
  - [Arithmetic Operations](#arithmetic-operations)
  - [Comparison Operations](#comparison-operations)
  - [Money Allocation](#money-allocation)
  - [Formatting and Display](#formatting-and-display)
  - [Parsing](#parsing)
    - [Decimal Separator Handling](#decimal-separator-handling)
  - [Aggregation](#aggregation)
  - [Currency Exchange](#currency-exchange)
  - [JSON Serialization](#json-serialization)
  - [Database Support](#database-support)
  - [Currency Management](#currency-management)
- [Supported Currencies](#supported-currencies)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)
- [Development](#development)
- [Differences from MoneyPHP](#differences-from-moneyphp)
- [License](#license)
- [Credits](#credits)
- [Contributing](#contributing)
- [Resources](#resources)

---

## 🤔 Why Use a Money Library?

You shouldn't represent monetary values with floats due to precision issues. This library uses integers internally to avoid floating-point arithmetic errors and provides safe money calculations.

```go
// ❌ DON'T DO THIS
price := 0.1 + 0.2  // 0.30000000000000004

// ✅ DO THIS
mm := money.NewManager()
a := mm.Create(10, currency.USD) // 10 cents
b := mm.Create(20, currency.USD) // 20 cents

sum, _ := mm.Add(a, b) // $0.30
```

## 🚀 Installation

```bash
go get github.com/gollin/packages/support/money
```

## 📋 Requirements

- Go `1.26.0` (per `go.mod`)

## ✨ Features

- **Safe Money Arithmetic**: Add, subtract, multiply without floats
- **Currency Support**: ISO 4217 dataset + custom currencies
- **Formatting**: Currency-aware formatting and major-unit conversion
- **Parsing**: Parse human-entered amounts (symbols/codes, separators)
- **Aggregation**: Sum/Min/Max/Avg via an `Aggregator`
- **Currency Exchange**: Convert via an `exchange.Exchange`
- **JSON + DB Integration**: `json.Marshaler`/`Unmarshaler`, `sql.Scanner`/`driver.Valuer`

## 📦 Packages

- `github.com/gollin/packages/support/money/money`: `Money`, `Manager`, `Aggregator`, `Converter`
- `github.com/gollin/packages/support/money/currency`: ISO currency data + `Manager`
- `github.com/gollin/packages/support/money/exchange`: in-memory exchange rates + conversion
- `github.com/gollin/packages/support/money/parser`: parse strings like `$1,234.56`, `EUR 10,50`
- `github.com/gollin/packages/support/money/format`: currency formatter used by `currency.Currency`
- `github.com/gollin/packages/support/money/exception`: sentinel errors used across the module

## 🏁 Quick Start

```go
package main

import (
	"fmt"

	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/money"
)

func main() {
	mm := money.NewManager()

	// Amounts are in minor units (cents, pence, etc.)
	price := mm.Create(10000, currency.USD) // $100.00
	tax := mm.Create(850, currency.USD)     // $8.50

	total, _ := mm.Add(price, tax)
	display, _ := total.Display()
	fmt.Println(display) // $108.50

	// Convenience constructors (use the default Manager)
	eur := money.FromEUR(50000) // €500.00
	eurDisplay, _ := eur.Display()
	fmt.Println(eurDisplay)
}
```

## 🧠 Core Concepts

### Creating Money

```go
mm := money.NewManager()

// Create with amount in minor units (cents, pence, etc.)
usd := mm.Create(10000, currency.USD) // $100.00

// From float (use only for user input conversion; floats are imprecise)
approx := mm.CreateFromFloat(99.99, currency.USD)

// From exact decimal string (recommended for external/user-provided decimals)
exact, _ := mm.CreateFromString("99.99", currency.USD)

// Convenience constructors (use the default Manager)
eur := money.FromEUR(5000)  // €50.00
jpy := money.FromJPY(10000) // ¥10000
```

### Arithmetic Operations

Use `money.Manager` for arithmetic operations.

```go
mm := money.NewManager()
a := mm.Create(10000, currency.USD) // $100.00
b := mm.Create(2500, currency.USD)  // $25.00

// Addition
sum, _ := mm.Add(a, b)
sumDisplay, _ := sum.Display()
fmt.Println(sumDisplay) // $125.00

// Subtraction
diff, _ := mm.Subtract(a, b)
diffDisplay, _ := diff.Display()
fmt.Println(diffDisplay) // $75.00

// Multiplication
product, _ := mm.Multiply(a, 2)
productDisplay, _ := product.Display()
fmt.Println(productDisplay) // $200.00
```

### Comparison Operations

```go
mm := money.NewManager()
a := mm.Create(10000, currency.USD)
b := mm.Create(5000, currency.USD)

equals, _ := a.Equals(b)                    // false
greaterThan, _ := a.GreaterThan(b)          // true
lessThan, _ := a.LessThan(b)                // false
greaterOrEqual, _ := a.GreaterThanOrEqual(b) // true
lessOrEqual, _ := a.LessThanOrEqual(b)      // false

// Compare returns -1, 0, or 1
cmp, _ := a.Compare(b)  // 1 (a > b)

// Check properties
a.IsZero()      // (false, nil)
a.IsPositive()  // (true, nil)
a.IsNegative()  // (false, nil)
```

### Money Allocation

Split money without losing pennies due to rounding.

```go
mm := money.NewManager()
m := mm.Create(10000, currency.USD) // $100.00

// Split evenly
parts, _ := mm.Split(m, 3)
// parts[0]: $33.34
// parts[1]: $33.33
// parts[2]: $33.33

// Allocate by ratios
m2 := mm.Create(10000, currency.USD) // $100.00
parts, _ = mm.Allocate(m2, 50, 30, 20)
// parts[0]: $50.00 (50%)
// parts[1]: $30.00 (30%)
// parts[2]: $20.00 (20%)
```

### Formatting and Display

```go
mm := money.NewManager()
m := mm.Create(123456, currency.USD)

// Display formatted
display, _ := m.Display() // "$1,234.56"

// Get as major units (float)
major, _ := m.AsMajorUnits() // 1234.56

// Access raw values
amount, _ := m.Amount()   // 123456 (int64)
curr, _ := m.Currency()   // *currency.Currency{Code: "USD", ...}
```

### Parsing

```go
import (
	"github.com/gollin/packages/support/money/money"
	"github.com/gollin/packages/support/money/parser"
)

p := parser.NewParser()

// Parse with currency symbol
val, currency, err := p.ParseAmount("$100.50")        // 100.50, "USD"
m := money.NewManager().CreateFromFloat(val, currency)

val, currency, err = p.ParseAmount("€250.75")        // 250.75, "EUR"
val, currency, err = p.ParseAmount("£99.99")         // 99.99, "GBP"

// Parse with currency code
val, currency, err = p.ParseAmount("100.50 USD")     // 100.50, "USD"
val, currency, err = p.ParseAmount("EUR 250.75")     // 250.75, "EUR"

// Parse with thousands separator
val, currency, err = p.ParseAmount("$1,234.56")      // 1234.56, "USD"

// Parse with default currency
val, currency, err = p.ParseAmount("100.50", "USD")  // 100.50, "USD"

// Parse just decimal values (no currency)
amount, err := p.ParseDecimal("1234.56")             // 1234.56
```

#### Decimal Separator Handling

The parser intelligently handles different decimal separator formats:

```go
// Mixed separators - unambiguous
p.ParseAmount("$1,234.56")     // US format: 1234.56 USD
p.ParseAmount("€1.234,56")     // European format: 1234.56 EUR

// Comma-only input - ambiguous
p.ParseAmount("$1,000")        // Treated as 1000.00 USD (thousands separator)
p.ParseAmount("€10,50")        // Treated as 1050.00 EUR (NOT €10.50)

// For European decimal-only format, use ParseAmountWithDecimalComma
p.ParseAmountWithDecimalComma("€10,50")  // Correctly parses as 10.50 EUR
p.ParseDecimalWithComma("10,50")         // Correctly parses as 10.50
```

**Important**: When only commas are present (no dots), the parser treats them as thousands separators by default. To parse European-style decimal commas (e.g., "10,50" as 10.50), use `ParseAmountWithDecimalComma()` or `ParseDecimalWithComma()` methods.

### Aggregation

```go
mm := money.NewManager()
agg := money.NewAggregator(mm)
prices := []*money.Money{
	mm.Create(10000, currency.USD),
	mm.Create(20000, currency.USD),
	mm.Create(30000, currency.USD),
}

// Sum all
total, _ := agg.Sum(prices...)

// Get minimum
min, _ := agg.Min(prices...)

// Get maximum
max, _ := agg.Max(prices...)

// Calculate average
avg, _ := agg.Avg(prices...)
```

### Currency Exchange

```go
import (
	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/exchange"
	"github.com/gollin/packages/support/money/money"
)

currencies := currency.NewManager()
ex := exchange.NewExchange()
_ = ex.AddRate(currency.USD, currency.EUR, 0.85)
_ = ex.AddRate(currency.USD, currency.GBP, 0.73)

// Get exchange rate
rate, _ := ex.GetRate(currency.USD, currency.EUR) // 0.85

converter, _ := money.NewConverter(currencies, ex)

mm := money.NewManager()
usd := mm.Create(10000, currency.USD) // $100.00
eur, _ := converter.Convert(usd, currency.EUR)
eurDisplay, _ := eur.Display()
fmt.Println(eurDisplay) // €85.00

gbp, _ := converter.ConvertWithRate(usd, currency.GBP, 0.73)
gbpDisplay, _ := gbp.Display()
fmt.Println(gbpDisplay) // £73.00
```

### JSON Serialization

```go
type Product struct {
    Name  string       `json:"name"`
    Price *money.Money `json:"price"`
}

// Marshal
product := Product{
    Name:  "Book",
    Price: money.NewManager().Create(2999, currency.USD),
}
json, _ := json.Marshal(product)
// {"name":"Book","price":{"amount":2999,"currency":"USD"}}

// Unmarshal
var p Product
json.Unmarshal([]byte(jsonStr), &p)
```

### Database Support

```go
type Product struct {
    ID    int          `db:"id"`
    Price *money.Money `db:"price"`
}

// Money implements sql.Scanner and driver.Valuer
// Stores as a delimited string: "amount|currency_code" (separator is configurable)
db.Exec("INSERT INTO products (price) VALUES (?)", product.Price)
// Stores as: 2999|USD
```

### Currency Management

```go
import (
	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/money"
)

// Get currency information
cm := currency.NewManager()
c := cm.FindByCode(currency.USD)
c.Code     // "USD"
c.Fraction // 2 (decimal places)
c.Grapheme // "$"

// Get by numeric code (ISO 4217)
c2 := cm.FindByNumericCode("840") // USD

// Add custom currency
cm.AddFrom("BTC", "₿", "$1", ".", ",", "000", 8)
mm := money.NewManager()
bitcoin := mm.Create(100000000, "BTC") // 1.00000000 ₿
```

## 🌍 Supported Currencies

The default `currency.Manager` includes an ISO 4217 dataset (including many historical entries) and can be extended with custom currencies.

The `money` package also includes many `FromXXX(amount int64)` helpers (e.g. `money.FromUSD`) that create `Money` values using the default manager.

## 🚧 Error Handling

Operations that could fail return errors:

```go
import (
	"fmt"

	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/exception"
	"github.com/gollin/packages/support/money/money"
)

mm := money.NewManager()
usd := mm.Create(100, currency.USD)
eur := mm.Create(100, currency.EUR)

// This will return error
_, err := mm.Add(usd, eur)
if err != nil {
    fmt.Println("Cannot add different currencies:", err)
    // Error message: "currencies don't match"
}
```

Common sentinel errors live in `github.com/gollin/packages/support/money/exception` (e.g. `exception.ErrCurrencyMismatch`).

## ✅ Best Practices

1.  **Always use minor units**: Store amounts as integers in the smallest currency unit (cents, pence, etc.)

    ```go
    // ✅ Good
    price := money.NewManager().Create(9999, currency.USD) // $99.99
    ```

2.  **Check errors**: Money operations can fail (currency mismatch, etc.)

    ```go
    mm := money.NewManager()
    result, err := mm.Add(a, b)
    if err != nil {
        // Handle error
    }
    ```

3.  **Use allocation for splits**: Don't manually divide - use Split() or Allocate()

    ```go
    // ✅ Good
    mm := money.NewManager()
    m := mm.Create(10000, currency.USD)
    parts, _ := mm.Split(m, 3)

    // ❌ Loses pennies
    third := mm.Create(10000/3, currency.USD)
    ```

4.  **Prefer exact decimal strings for user input**: `CreateFromString("12.34", "USD")` avoids float rounding surprises.

## 👨‍💻 Development

- Run tests + auto-format: `make tests`
- Run a per-package test summary: `make pretty-test`
- Format only: `make format`

## ⚖️ Differences from MoneyPHP

- Uses `int64` instead of strings for amounts (supports up to ~92 quadrillion)
- Go-idiomatic API (e.g., `money.FromUSD()` helper, and `money.Manager` for operations)
- Database support via `sql.Scanner` and `driver.Valuer`
- Aggregation via `money.Aggregator`
- Simple in-memory exchange system with explicit rate management

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Credits

Inspired by [moneyphp/money](https://github.com/moneyphp/money) and Martin Fowler's Money pattern.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📚 Resources

- [Martin Fowler's Money Pattern](https://martinfowler.com/eaaCatalog/money.html)
- [MoneyPHP Documentation](https://www.moneyphp.org/)
- [ISO 4217 Currency Codes](https://en.wikipedia.org/wiki/ISO_4217)
