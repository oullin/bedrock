package support

import (
	"reflect"
	"testing"
)

// Port of Framework\Tests\Support\SupportArrTest::testAccessible
func TestArrAccessible(t *testing.T) {
	t.Parallel()

	if !ArrAccessible(map[string]any{"name": "Taylor"}) {
		t.Fatal("ArrAccessible should accept maps")
	}

	if !ArrAccessible([]int{1, 2, 3}) {
		t.Fatal("ArrAccessible should accept slices")
	}

	if ArrAccessible(42) {
		t.Fatal("ArrAccessible should reject scalar integers")
	}
}

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

// Port of Framework\Tests\Support\SupportArrTest::testItGetsAString
func TestArrString(t *testing.T) {
	t.Parallel()

	got := ArrString(map[string]any{"name": "Taylor"}, "name")

	if got != "Taylor" {
		t.Fatalf("ArrString = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testItGetsAnInteger
func TestArrInteger(t *testing.T) {
	t.Parallel()

	got := ArrInteger(map[string]any{"count": "10"}, "count")

	if got != 10 {
		t.Fatalf("ArrInteger = %d", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testItGetsAFloat
func TestArrFloat(t *testing.T) {
	t.Parallel()

	got := ArrFloat(map[string]any{"price": "10.50"}, "price")

	if got != 10.50 {
		t.Fatalf("ArrFloat = %f", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testItGetsABoolean
func TestArrBoolean(t *testing.T) {
	t.Parallel()

	got := ArrBoolean(map[string]any{"active": "true"}, "active")

	if !got {
		t.Fatal("ArrBoolean should parse true")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testItGetsAnArray
func TestArrArray(t *testing.T) {
	t.Parallel()

	got := ArrArray(map[string]any{"items": []string{"a", "b"}}, "items")

	if !reflect.DeepEqual(got, []any{"a", "b"}) {
		t.Fatalf("ArrArray = %#v", got)
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

// Port of Framework\Tests\Support\SupportArrTest::testHasAllMethod
func TestArrHasAllMethod(t *testing.T) {
	t.Parallel()

	m := map[string]any{"user": map[string]any{"name": "Taylor", "email": "taylor@example.com"}}

	if !ArrHas(m, "user.name", "user.email") {
		t.Fatal("ArrHas should require all keys")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testHasAnyMethod
func TestArrHasAnyMethod(t *testing.T) {
	t.Parallel()

	m := map[string]any{"user": map[string]any{"name": "Taylor"}}

	if !ArrHasAny(m, "user.email", "user.name") {
		t.Fatal("ArrHasAny should return true when one key exists")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testIsAssoc
func TestArrIsAssoc(t *testing.T) {
	t.Parallel()

	if !ArrIsAssoc(map[string]any{"name": "Taylor"}) {
		t.Fatal("ArrIsAssoc should treat string-keyed maps as associative")
	}

	if ArrIsAssoc(map[int]any{0: "a", 1: "b"}) {
		t.Fatal("ArrIsAssoc should reject list-shaped integer maps")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testIsList
func TestArrIsList(t *testing.T) {
	t.Parallel()

	if !ArrIsList([]string{"a", "b"}) {
		t.Fatal("ArrIsList should accept slices")
	}

	if ArrIsList(map[int]any{1: "a"}) {
		t.Fatal("ArrIsList should reject maps missing zero-based keys")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testExists
func TestArrExists(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"products": map[string]any{
			"desk": map[string]any{"price": 100},
		},
	}

	if !ArrExists(m, "products.desk.price") {
		t.Fatal("ArrExists should find nested key")
	}

	if ArrExists(m, "products.desk.missing") {
		t.Fatal("ArrExists should not find missing key")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testWhereNotNull
func TestArrWhereNotNull(t *testing.T) {
	t.Parallel()

	m := map[string]any{"name": "Taylor", "email": nil, "age": 30}
	result := ArrWhereNotNull(m)

	if len(result) != 2 {
		t.Fatalf("ArrWhereNotNull length = %d, want 2", len(result))
	}

	if _, ok := result["email"]; ok {
		t.Fatal("ArrWhereNotNull should remove nil values")
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testExceptValues
func TestArrExceptValues(t *testing.T) {
	t.Parallel()

	m := map[string]any{"a": 1, "b": 2, "c": 1}
	result := ArrExceptValues(m, 1)

	if len(result) != 1 || result["b"] != 2 {
		t.Fatalf("ArrExceptValues = %v", result)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testUndot
func TestArrUndot(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"user.name":  "Taylor",
		"user.email": "taylor@example.com",
		"active":     true,
	}

	result := ArrUndot(m)
	user, ok := result["user"].(map[string]any)

	if !ok {
		t.Fatalf("ArrUndot user = %T", result["user"])
	}

	if user["name"] != "Taylor" || user["email"] != "taylor@example.com" {
		t.Fatalf("ArrUndot user = %v", user)
	}

	if result["active"] != true {
		t.Fatalf("ArrUndot active = %v", result["active"])
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testJoin
func TestArrJoin(t *testing.T) {
	t.Parallel()

	if got := ArrJoin([]string{"a", "b", "c"}, ", ", " and "); got != "a, b and c" {
		t.Fatalf("ArrJoin = %q", got)
	}

	if got := ArrJoin([]string{"a", "b"}, ", ", " and "); got != "a and b" {
		t.Fatalf("ArrJoin two items = %q", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testTake
func TestArrTake(t *testing.T) {
	t.Parallel()

	if got := ArrTake([]int{1, 2, 3, 4}, 2); len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("ArrTake positive = %v", got)
	}

	if got := ArrTake([]int{1, 2, 3, 4}, -2); len(got) != 2 || got[0] != 3 || got[1] != 4 {
		t.Fatalf("ArrTake negative = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testPush
func TestArrPush(t *testing.T) {
	t.Parallel()

	result := ArrPush([]string{"a"}, "b", "c")

	if len(result) != 3 || result[0] != "a" || result[2] != "c" {
		t.Fatalf("ArrPush = %v", result)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testQuery
func TestArrQuery(t *testing.T) {
	t.Parallel()

	got := ArrQuery(map[string]any{
		"name":   "Taylor Otwell",
		"skills": []string{"go", "php"},
	})

	if got != "name=Taylor+Otwell&skills=go&skills=php" && got != "skills=go&skills=php&name=Taylor+Otwell" {
		t.Fatalf("ArrQuery = %q", got)
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

// Port of Framework\Tests\Support\SupportArrTest::testOnlyValues
func TestArrOnlyValues(t *testing.T) {
	t.Parallel()

	got := ArrOnlyValues(map[string]any{"name": "Taylor", "age": 40}, "age", "name")

	if !reflect.DeepEqual(got, []any{40, "Taylor"}) {
		t.Fatalf("ArrOnlyValues = %#v", got)
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

// Port of Framework\Tests\Support\SupportArrTest::testPluckWithArrayValue
func TestArrPluckWithArrayValue(t *testing.T) {
	t.Parallel()

	got := ArrPluck([]map[string]any{
		{"name": []string{"Taylor", "Otwell"}},
	}, "name").([]any)

	if !reflect.DeepEqual(got[0], []string{"Taylor", "Otwell"}) {
		t.Fatalf("ArrPluck array value = %#v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testPluckWithKeys
func TestArrPluckWithKeys(t *testing.T) {
	t.Parallel()

	got := ArrPluck([]map[string]any{
		{"account": map[string]any{"id": "a"}, "name": "Taylor"},
		{"account": map[string]any{"id": "b"}, "name": "Abigail"},
	}, "name", "account.id").(map[string]any)

	if got["a"] != "Taylor" || got["b"] != "Abigail" {
		t.Fatalf("ArrPluck nested keys = %#v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testArrayPluckWithNestedKeys
func TestArrPluckWithNestedKeys(t *testing.T) {
	t.Parallel()

	got := ArrPluck([]map[string]any{
		{"user": map[string]any{"name": "Taylor"}},
		{"user": map[string]any{"name": "Abigail"}},
	}, "user.name").([]any)

	if !reflect.DeepEqual(got, []any{"Taylor", "Abigail"}) {
		t.Fatalf("ArrPluck nested values = %#v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testArrayPluckWithNestedArrays
func TestArrPluckWithNestedArrays(t *testing.T) {
	t.Parallel()

	got := ArrPluck([]map[string]any{
		{"users": map[string]any{"names": []string{"Taylor"}}},
	}, "users.names").([]any)

	if !reflect.DeepEqual(got[0], []string{"Taylor"}) {
		t.Fatalf("ArrPluck nested arrays = %#v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testMap
func TestArrMap(t *testing.T) {
	t.Parallel()

	got := ArrMap([]int{1, 2, 3}, func(item int) int {
		return item * 2
	})

	if !reflect.DeepEqual(got, []int{2, 4, 6}) {
		t.Fatalf("ArrMap = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testMapWithEmptyArray
func TestArrMapWithEmptyArray(t *testing.T) {
	t.Parallel()

	got := ArrMap([]int{}, func(item int) int {
		return item * 2
	})

	if len(got) != 0 {
		t.Fatalf("ArrMap empty = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testMapNullValues
func TestArrMapNullValues(t *testing.T) {
	t.Parallel()

	got := ArrMap([]any{nil, "Taylor"}, func(item any) any {
		if item == nil {
			return "missing"
		}

		return item
	})

	if !reflect.DeepEqual(got, []any{"missing", "Taylor"}) {
		t.Fatalf("ArrMap null values = %#v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testMapWithKeys
func TestArrMapWithKeys(t *testing.T) {
	t.Parallel()

	got := ArrMapWithKeys([]string{"Taylor", "Abigail"}, func(item string) map[string]any {
		return map[string]any{item: len(item)}
	})

	if got["Taylor"] != 6 || got["Abigail"] != 7 {
		t.Fatalf("ArrMapWithKeys = %v", got)
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

// Port of Framework\Tests\Support\SupportArrTest::testSortDesc
func TestArrSortDesc(t *testing.T) {
	t.Parallel()

	got := ArrSortDesc([]int{1, 3, 2})

	if !reflect.DeepEqual(got, []int{3, 2, 1}) {
		t.Fatalf("ArrSortDesc = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testSortByMany
func TestArrSortByMany(t *testing.T) {
	t.Parallel()

	got := ArrSortByMany([]map[string]any{
		{"name": "Taylor", "age": 40},
		{"name": "Abigail", "age": 30},
		{"name": "Taylor", "age": 35},
	}, SortClause{Key: "name", Direction: SortAsc}, SortClause{Key: "age", Direction: SortDesc})

	if got[0]["name"] != "Abigail" || got[1]["age"] != 40 || got[2]["age"] != 35 {
		t.Fatalf("ArrSortByMany = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testKeyBy
func TestArrKeyBy(t *testing.T) {
	t.Parallel()

	got := ArrKeyBy([]map[string]any{
		{"id": 1, "name": "Taylor"},
		{"id": 2, "name": "Abigail"},
	}, "id")

	if got["1"]["name"] != "Taylor" || got["2"]["name"] != "Abigail" {
		t.Fatalf("ArrKeyBy = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testPrependKeysWith
func TestArrPrependKeysWith(t *testing.T) {
	t.Parallel()

	got := ArrPrependKeysWith(map[string]any{"name": "Taylor"}, "user.")

	if got["user.name"] != "Taylor" {
		t.Fatalf("ArrPrependKeysWith = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testSelect
func TestArrSelect(t *testing.T) {
	t.Parallel()

	got := ArrSelect([]map[string]any{
		{"name": "Taylor", "email": "taylor@example.com"},
		{"name": "Abigail", "email": "abigail@example.com"},
	}, "name")

	if len(got) != 2 || got[0]["name"] != "Taylor" || got[0]["email"] != nil {
		t.Fatalf("ArrSelect = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testReject
func TestArrReject(t *testing.T) {
	t.Parallel()

	got := ArrReject([]int{1, 2, 3}, func(value int, _ int) bool {
		return value > 1
	})

	if !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("ArrReject = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testWhereKey
func TestArrWhereKey(t *testing.T) {
	t.Parallel()

	got := ArrWhereKey(map[string]any{"name": "Taylor", "email": "taylor@example.com"}, func(key string) bool {
		return key == "email"
	})

	if len(got) != 1 || got["email"] != "taylor@example.com" {
		t.Fatalf("ArrWhereKey = %v", got)
	}
}

// Port of Framework\Tests\Support\SupportArrTest::testFrom
func TestArrFrom(t *testing.T) {
	t.Parallel()

	got := ArrFrom([]string{"a", "b"})

	if !reflect.DeepEqual(got, []any{"a", "b"}) {
		t.Fatalf("ArrFrom slice = %#v", got)
	}

	got = ArrFrom("a")

	if !reflect.DeepEqual(got, []any{"a"}) {
		t.Fatalf("ArrFrom scalar = %#v", got)
	}
}
