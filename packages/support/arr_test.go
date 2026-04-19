package support

import (
	"testing"
)

// Port of Framework\Tests\Support\SupportArrTest::testAdd
func TestArrAdd(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Desk"}
	result := ArrAdd(m, "price", 100)

	if result["price"] != 100 {
		t.Errorf("ArrAdd should add missing key, got %v", result)
	}

	result = ArrAdd(m, "name", "Chair")

	if result["name"] != "Desk" {
		t.Errorf("ArrAdd should not overwrite existing key, got %v", result["name"])
	}
}

func TestArrAddDotNotation(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"user": map[string]any{"name": "Taylor"},
	}

	ArrAdd(m, "user.age", 30)
	user := m["user"].(map[string]any)

	if user["age"] != 30 {
		t.Errorf("ArrAdd with dot notation failed: %v", m)
	}

	ArrAdd(m, "user.name", "Otwell")

	if user["name"] != "Taylor" {
		t.Errorf("ArrAdd should not overwrite existing nested key, got %v", user["name"])
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testGet
func TestArrGet(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"products": map[string]any{
			"desk": map[string]any{
				"price": 100,
			},
		},
	}

	if val := ArrGet(m, "products.desk.price"); val != 100 {
		t.Errorf("ArrGet = %v, want 100", val)
	}

	if val := ArrGet(m, "missing", "default"); val != "default" {
		t.Errorf("ArrGet with default = %v, want default", val)
	}

	if val := ArrGet(m, "missing"); val != nil {
		t.Errorf("ArrGet missing without default = %v, want nil", val)
	}
}

func TestArrGetNestedMissing(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Taylor"}

	if val := ArrGet(m, "name.first", "default"); val != "default" {
		t.Errorf("ArrGet through non-map should return default, got %v", val)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testSet
func TestArrSet(t *testing.T) {
	t.Parallel()

	m := map[string]any{"products": map[string]any{
		"desk": map[string]any{"price": 100},
	}}

	ArrSet(m, "products.desk.price", 200)
	val := ArrGet(m, "products.desk.price")

	if val != 200 {
		t.Errorf("ArrSet = %v, want 200", val)
	}

	ArrSet(m, "products.desk.discount", 10)

	if ArrGet(m, "products.desk.discount") != 10 {
		t.Error("ArrSet should create new nested key")
	}
}

func TestArrSetCreatesIntermediates(t *testing.T) {
	t.Parallel()

	m := make(map[string]any)
	ArrSet(m, "user.profile.name", "Taylor")

	if ArrGet(m, "user.profile.name") != "Taylor" {
		t.Errorf("ArrSet should create intermediate maps: %v", m)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testHas
func TestArrHas(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"products": map[string]any{
			"desk": map[string]any{"price": 100},
		},
	}

	if !ArrHas(m, "products.desk.price") {
		t.Error("ArrHas should find nested key")
	}

	if ArrHas(m, "products.desk.missing") {
		t.Error("ArrHas should not find missing key")
	}

	if ArrHas(m) {
		t.Error("ArrHas with no keys should return false")
	}
}

func TestArrHasMultipleKeys(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Taylor", "age": 30}

	if !ArrHas(m, "name", "age") {
		t.Error("ArrHas should return true when all keys exist")
	}

	if ArrHas(m, "name", "missing") {
		t.Error("ArrHas should return false when any key is missing")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testForget
func TestArrForget(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"products": map[string]any{
			"desk": map[string]any{"price": 100},
		},
		"name": "Taylor",
	}

	ArrForget(m, "products.desk.price")

	desk := m["products"].(map[string]any)["desk"].(map[string]any)

	if _, ok := desk["price"]; ok {
		t.Error("ArrForget should remove nested key")
	}

	ArrForget(m, "name")

	if _, ok := m["name"]; ok {
		t.Error("ArrForget should remove top-level key")
	}
}

func TestArrForgetMultipleKeys(t *testing.T) {
	t.Parallel()

	m := map[string]any{"a": 1, "b": 2, "c": 3}
	ArrForget(m, "a", "c")

	if _, ok := m["a"]; ok {
		t.Error("ArrForget should remove key 'a'")
	}

	if _, ok := m["c"]; ok {
		t.Error("ArrForget should remove key 'c'")
	}

	if m["b"] != 2 {
		t.Error("ArrForget should not affect other keys")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testPull
func TestArrPull(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Desk", "price": 100}

	name := ArrPull(m, "name")

	if name != "Desk" {
		t.Errorf("ArrPull should return value, got %v", name)
	}

	if _, ok := m["name"]; ok {
		t.Error("ArrPull should remove key after retrieval")
	}
}

func TestArrPullDefault(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Desk"}

	val := ArrPull(m, "missing", "default")

	if val != "default" {
		t.Errorf("ArrPull missing key should return default, got %v", val)
	}
}

func TestArrPullDotNotation(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"user": map[string]any{"name": "Taylor", "age": 30},
	}

	name := ArrPull(m, "user.name")

	if name != "Taylor" {
		t.Errorf("ArrPull nested = %v, want Taylor", name)
	}

	user := m["user"].(map[string]any)

	if _, ok := user["name"]; ok {
		t.Error("ArrPull should remove nested key")
	}

	if user["age"] != 30 {
		t.Error("ArrPull should not affect other nested keys")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testDot
func TestArrDot(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"user": map[string]any{
			"name":  "Taylor",
			"email": "taylor@example.com",
		},
		"active": true,
	}

	result := ArrDot(m)

	if result["user.name"] != "Taylor" {
		t.Errorf("ArrDot[user.name] = %v, want Taylor", result["user.name"])
	}

	if result["user.email"] != "taylor@example.com" {
		t.Errorf("ArrDot[user.email] = %v", result["user.email"])
	}

	if result["active"] != true {
		t.Errorf("ArrDot[active] = %v, want true", result["active"])
	}
}

func TestArrDotWithPrepend(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Taylor"}
	result := ArrDot(m, "user")

	if result["user.name"] != "Taylor" {
		t.Errorf("ArrDot with prepend = %v", result)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testExcept
func TestArrExcept(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Taylor", "age": 30, "city": "Little Rock"}
	result := ArrExcept(m, "age", "city")

	if len(result) != 1 || result["name"] != "Taylor" {
		t.Errorf("ArrExcept = %v, want {name: Taylor}", result)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testOnly
func TestArrOnly(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Taylor", "age": 30, "city": "Little Rock"}
	result := ArrOnly(m, "name", "age")

	if len(result) != 2 {
		t.Errorf("ArrOnly length = %d, want 2", len(result))
	}

	if result["name"] != "Taylor" || result["age"] != 30 {
		t.Errorf("ArrOnly = %v", result)
	}
}

func TestArrOnlyMissingKeys(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Taylor"}
	result := ArrOnly(m, "name", "missing")

	if len(result) != 1 || result["name"] != "Taylor" {
		t.Errorf("ArrOnly with missing keys = %v", result)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testDivide
func TestArrDivide(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Taylor"}
	keys, values := ArrDivide(m)

	if len(keys) != 1 || keys[0] != "name" {
		t.Errorf("ArrDivide keys = %v, want [name]", keys)
	}

	if len(values) != 1 || values[0] != "Taylor" {
		t.Errorf("ArrDivide values = %v, want [Taylor]", values)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testPluck
func TestArrPluck(t *testing.T) {
	t.Parallel()

	items := []map[string]any{
		{"name": "Taylor", "email": "taylor@example.com"},
		{"name": "Abigail", "email": "abigail@example.com"},
	}

	names := ArrPluck(items, "name").([]any)

	if len(names) != 2 || names[0] != "Taylor" || names[1] != "Abigail" {
		t.Errorf("ArrPluck = %v", names)
	}
}

func TestArrPluckWithKey(t *testing.T) {
	t.Parallel()

	items := []map[string]any{
		{"id": "1", "name": "Taylor"},
		{"id": "2", "name": "Abigail"},
	}

	result := ArrPluck(items, "name", "id").(map[string]any)

	if result["1"] != "Taylor" || result["2"] != "Abigail" {
		t.Errorf("ArrPluck with key = %v", result)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testSortRecursive
func TestArrSortRecursive(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"users": []any{"Taylor", "Abigail"},
		"name":  "Desk",
	}

	result := ArrSortRecursive(m)

	users := result["users"].([]any)

	if users[0] != "Abigail" || users[1] != "Taylor" {
		t.Errorf("ArrSortRecursive users = %v", users)
	}
}

func TestArrSortRecursiveDescending(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"users": []any{"Abigail", "Taylor"},
	}

	result := ArrSortRecursive(m, true)
	users := result["users"].([]any)

	if users[0] != "Taylor" || users[1] != "Abigail" {
		t.Errorf("ArrSortRecursive descending = %v", users)
	}
}

func TestArrSortRecursiveNested(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"info": map[string]any{
			"tags": []any{"c", "a", "b"},
		},
	}

	result := ArrSortRecursive(m)
	tags := result["info"].(map[string]any)["tags"].([]any)

	if tags[0] != "a" || tags[1] != "b" || tags[2] != "c" {
		t.Errorf("ArrSortRecursive nested = %v", tags)
	}
}
