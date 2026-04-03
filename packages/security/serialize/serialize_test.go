package serialize

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

// TestBasicTypes tests serialization of basic PHP types
func TestBasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"null", nil, "N;"},
		{"bool true", true, "b:1;"},
		{"bool false", false, "b:0;"},
		{"int positive", 42, "i:42;"},
		{"int negative", -42, "i:-42;"},
		{"int zero", 0, "i:0;"},
		{"float", 3.14, "d:3.14;"},
		{"string empty", "", `s:0:"";`},
		{"string simple", "hello", `s:5:"hello";`},
		{"string with spaces", "hello world", `s:11:"hello world";`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}

			// Test round-trip
			unmarshalled, err := Unmarshal(result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// Special handling for nil
			if tt.input == nil {
				if unmarshalled != nil {
					t.Errorf("Expected nil, got %v", unmarshalled)
				}
			}
		})
	}
}

// TestSpecialFloats tests NaN and Infinity handling
func TestSpecialFloats(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"NaN", math.NaN(), "d:NAN;"},
		{"positive infinity", math.Inf(1), "d:INF;"},
		{"negative infinity", math.Inf(-1), "d:-INF;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}

			// Test unmarshalling
			unmarshalled, err := Unmarshal(result)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			f, ok := unmarshalled.(float64)
			if !ok {
				t.Fatalf("Expected float64, got %T", unmarshalled)
			}

			if math.IsNaN(tt.input) {
				if !math.IsNaN(f) {
					t.Errorf("Expected NaN, got %v", f)
				}
			} else if math.IsInf(tt.input, 1) {
				if !math.IsInf(f, 1) {
					t.Errorf("Expected +Inf, got %v", f)
				}
			} else if math.IsInf(tt.input, -1) {
				if !math.IsInf(f, -1) {
					t.Errorf("Expected -Inf, got %v", f)
				}
			}
		})
	}
}

// TestArrays tests array/slice serialization
func TestArrays(t *testing.T) {
	// Indexed array
	arr := []interface{}{"a", "b", "c"}
	result, err := Marshal(arr)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := `a:3:{i:0;s:1:"a";i:1;s:1:"b";i:2;s:1:"c";}`
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Round-trip
	unmarshalled, err := Unmarshal(result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	slice, ok := unmarshalled.([]interface{})
	if !ok {
		t.Fatalf("Expected slice, got %T", unmarshalled)
	}

	if len(slice) != 3 {
		t.Errorf("Expected length 3, got %d", len(slice))
	}
}

// TestMaps tests map serialization
func TestMaps(t *testing.T) {
	m := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	result, err := Marshal(m)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Check it's a valid serialization (order might vary)
	if !strings.HasPrefix(result, "a:2:{") {
		t.Errorf("Expected to start with 'a:2:{', got %q", result)
	}

	// Round-trip
	unmarshalled, err := Unmarshal(result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	resultMap, ok := unmarshalled.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", unmarshalled)
	}

	if resultMap["name"] != "John" {
		t.Errorf("Expected name=John, got %v", resultMap["name"])
	}

	if resultMap["age"].(int64) != 30 {
		t.Errorf("Expected age=30, got %v", resultMap["age"])
	}
}

// TestPHPObject tests PHP object serialization
func TestPHPObject(t *testing.T) {
	obj := PHPObject{
		ClassName: "User",
		Properties: map[string]interface{}{
			"id":   123,
			"name": "John",
		},
	}

	result, err := MarshalObject(obj)
	if err != nil {
		t.Fatalf("MarshalObject failed: %v", err)
	}

	// Check format
	if !strings.HasPrefix(result, `O:4:"User":2:{`) {
		t.Errorf("Unexpected format: %q", result)
	}

	// Round-trip
	unmarshalled, err := Unmarshal(result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	phpObj, ok := unmarshalled.(PHPObject)
	if !ok {
		t.Fatalf("Expected PHPObject, got %T", unmarshalled)
	}

	if phpObj.ClassName != "User" {
		t.Errorf("Expected class User, got %s", phpObj.ClassName)
	}

	if len(phpObj.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(phpObj.Properties))
	}
}

// TestMaxDepth tests depth limiting
func TestMaxDepth(t *testing.T) {
	// Create nested structure
	deep := []interface{}{
		[]interface{}{
			[]interface{}{
				[]interface{}{"bottom"},
			},
		},
	}

	// Should succeed with depth 5
	_, err := Marshal(deep, WithMaxDepth(5))
	if err != nil {
		t.Errorf("Unexpected error with depth 5: %v", err)
	}

	// Should fail with depth 2
	_, err = Marshal(deep, WithMaxDepth(2))
	if err == nil {
		t.Error("Expected error with depth 2, got nil")
	}

	// Should fail with depth 2
	_, err = Unmarshal(`a:1:{i:0;a:1:{i:0;a:1:{i:0;i:1;}}}`, WithMaxDepth(2))
	if err == nil {
		t.Fatalf("Expected error at depth=2 with maxDepth=2, got nil")
	}

	// still reject deeper nesting
	_, err = Unmarshal(`a:1:{i:0;a:1:{i:0;a:1:{i:0;a:1:{i:0;i:1;}}}}`, WithMaxDepth(2))
	if err == nil {
		t.Fatalf("Expected error for exceeding max depth=2, but got nil")
	}
}

// TestAllowedClasses tests class filtering
func TestAllowedClasses(t *testing.T) {
	data := `O:4:"User":1:{s:2:"id";i:123;}`

	// Should succeed with allowed classes
	result, err := Unmarshal(data, WithAllowedClasses([]string{"User", "Admin"}))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	obj, ok := result.(PHPObject)
	if !ok || obj.ClassName != "User" {
		t.Error("Expected User object")
	}

	// Should fail with different allowed classes
	_, err = Unmarshal(data, WithAllowedClasses([]string{"Admin"}))
	if err == nil {
		t.Error("Expected error for disallowed class")
	}

	// Should fail with no classes allowed
	_, err = Unmarshal(data, WithAllowedClasses(nil))
	if err == nil {
		t.Error("Expected error when all classes disabled")
	}
}

// TestUTF8Strings tests UTF-8 string handling
func TestUTF8Strings(t *testing.T) {
	tests := []string{
		"Hello 世界",
		"Émojis: 😀🎉",
		"Русский текст",
		"العربية",
	}

	for _, str := range tests {
		t.Run(str, func(t *testing.T) {
			data, err := Marshal(str)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			result, err := Unmarshal(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result != str {
				t.Errorf("Expected %q, got %q", str, result)
			}
		})
	}
}

// TestEmptyCollections tests empty arrays and maps
func TestEmptyCollections(t *testing.T) {
	// Empty slice
	var emptySlice []interface{}
	data, _ := Marshal(emptySlice)
	result, _ := Unmarshal(data)

	// Should be empty map (PHP arrays with no elements)
	if m, ok := result.(map[string]interface{}); !ok || len(m) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}

	// Empty map
	emptyMap := map[string]interface{}{}
	data, _ = Marshal(emptyMap)
	result, _ = Unmarshal(data)

	if m, ok := result.(map[string]interface{}); !ok || len(m) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}
}

// TestErrorMessages tests that errors include position info
func TestErrorMessages(t *testing.T) {
	badData := "i:123" // Missing semicolon

	_, err := Unmarshal(badData)
	if err == nil {
		t.Fatal("Expected error for malformed data")
	}

	if !strings.Contains(err.Error(), "position") {
		t.Errorf("Error should include position info: %v", err)
	}
}

// TestNestedStructures tests complex nested data
func TestNestedStructures(t *testing.T) {
	complexData := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "John",
				"tags": []interface{}{"admin", "user"},
			},
			map[string]interface{}{
				"id":   2,
				"name": "Jane",
				"tags": []interface{}{"user"},
			},
		},
		"count": 2,
		"meta": map[string]interface{}{
			"version":   "1.0",
			"timestamp": 1234567890,
		},
	}

	data, err := Marshal(complexData)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	result, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", result)
	}

	if resultMap["count"].(int64) != 2 {
		t.Errorf("Expected count=2, got %v", resultMap["count"])
	}
}

// TestIsValidMarshaled tests validation function
func TestIsValidMarshaled(t *testing.T) {
	validData := []string{
		"N;",
		"b:1;",
		"i:42;",
		`s:5:"hello";`,
		`a:2:{i:0;s:1:"a";i:1;s:1:"b";}`,
	}

	for _, data := range validData {
		if !IsValidMarshaled(data) {
			t.Errorf("Expected %q to be valid", data)
		}
	}

	invalidData := []string{
		"x:1;",      // Invalid type
		"i:42",      // Missing semicolon
		`s:5:"hi";`, // Wrong length
		`a:1:{`,     // Incomplete
	}

	for _, data := range invalidData {
		if IsValidMarshaled(data) {
			t.Errorf("Expected %q to be invalid", data)
		}
	}
}

// TestMustMarshal tests panic behavior
func TestMustMarshal(t *testing.T) {
	// Should not panic for valid data
	result := MustMarshal(42)
	if result != "i:42;" {
		t.Errorf("Expected i:42;, got %s", result)
	}
}

// TestMustUnmarshal tests panic behavior for unmarshal
func TestMustUnmarshal(t *testing.T) {
	// Should not panic for valid data
	result := MustUnmarshal("i:42;")
	if result.(int64) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}

	// Should panic for invalid data
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid data")
		}
	}()

	MustUnmarshal("invalid") // Should panic
}

// TestIntegerKeys tests maps with integer keys
func TestIntegerKeys(t *testing.T) {
	// PHP allows mixed integer and string keys
	phpData := `a:3:{i:0;s:5:"first";i:1;s:6:"second";s:4:"name";s:4:"test";}`

	result, err := Unmarshal(phpData)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", result)
	}

	// Integer keys become string keys in the map
	if m["0"] != "first" {
		t.Errorf("Expected 'first' at key '0', got %v", m["0"])
	}

	if m["name"] != "test" {
		t.Errorf("Expected 'test' at key 'name', got %v", m["name"])
	}
}

// TestPHPPrivateProperties tests handling of PHP private/protected properties
func TestPHPPrivateProperties(t *testing.T) {
	// PHP private property format: \x00ClassName\x00propertyName
	// PHP protected property format: \x00*\x00propertyName
	phpData := "O:4:\"User\":2:{s:8:\"\x00User\x00id\";i:123;s:11:\"\x00*\x00password\";s:6:\"secret\";}"

	result, err := Unmarshal(phpData)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	obj, ok := result.(PHPObject)
	if !ok {
		t.Fatalf("Expected PHPObject, got %T", result)
	}

	// Property names should have visibility markers stripped
	if _, hasId := obj.Properties["id"]; !hasId {
		t.Error("Expected 'id' property (stripped of \\x00User\\x00)")
	}

	if _, hasPass := obj.Properties["password"]; !hasPass {
		t.Error("Expected 'password' property (stripped of \\x00*\\x00)")
	}
}

// TestBinaryData tests handling of binary strings
func TestBinaryData(t *testing.T) {
	// Binary data with null bytes
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE}
	binaryString := string(binaryData)

	data, err := Marshal(binaryString)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	result, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	resultStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	if resultStr != binaryString {
		t.Error("Binary data corrupted during round-trip")
	}
}

// TestLargeNumbers tests handling of large integers
func TestLargeNumbers(t *testing.T) {
	tests := []struct {
		name  string
		value int64
	}{
		{"max int64", math.MaxInt64},
		{"min int64", math.MinInt64},
		{"zero", 0},
		{"large positive", 9223372036854775806},
		{"large negative", -9223372036854775807},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			result, err := Unmarshal(data)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if result.(int64) != tt.value {
				t.Errorf("Expected %d, got %v", tt.value, result)
			}
		})
	}
}

// TestUintOverflow tests uint to int conversion with strict mode
func TestUintOverflow(t *testing.T) {
	var largeUint uint64 = math.MaxUint64

	// Should fail in strict mode (default)
	_, err := Marshal(largeUint)
	if err == nil {
		t.Error("Expected error for uint overflow in strict mode")
	}

	// Should succeed with non-strict mode
	_, err = Marshal(largeUint, WithStrictPHP(false))
	if err != nil {
		t.Errorf("Unexpected error in non-strict mode: %v", err)
	}
}

// TestSequentialArrayDetection tests proper slice vs map detection
func TestSequentialArrayDetection(t *testing.T) {
	// Sequential from 0
	sequential := `a:3:{i:0;s:1:"a";i:1;s:1:"b";i:2;s:1:"c";}`
	result, _ := Unmarshal(sequential)
	if _, ok := result.([]interface{}); !ok {
		t.Error("Expected slice for sequential array")
	}

	// Non-sequential (gap)
	nonSequential := `a:3:{i:0;s:1:"a";i:2;s:1:"b";i:3;s:1:"c";}`
	result, _ = Unmarshal(nonSequential)
	if _, ok := result.(map[string]interface{}); !ok {
		t.Error("Expected map for non-sequential array")
	}

	// Sequential but not starting from 0
	notFromZero := `a:3:{i:1;s:1:"a";i:2;s:1:"b";i:3;s:1:"c";}`
	result, _ = Unmarshal(notFromZero)
	if _, ok := result.(map[string]interface{}); !ok {
		t.Error("Expected map for array not starting from 0")
	}
}

// TestGetStringLength tests the helper function
func TestGetStringLength(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"", 0},
		{"Hello 世界", 12}, // UTF-8: 3 bytes per Chinese character
		{"😀", 4},         // UTF-8: 4 bytes for emoji
	}

	for _, tt := range tests {
		result := GetStringLength(tt.input)
		if result != tt.expected {
			t.Errorf("GetStringLength(%q) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

// TestIsUTF8 tests UTF-8 validation
func TestIsUTF8(t *testing.T) {
	validUTF8 := []string{
		"hello",
		"世界",
		"😀🎉",
		"",
	}

	for _, str := range validUTF8 {
		if !IsUTF8(str) {
			t.Errorf("Expected %q to be valid UTF-8", str)
		}
	}

	// Invalid UTF-8 sequence
	invalidUTF8 := string([]byte{0xFF, 0xFE, 0xFD})
	if IsUTF8(invalidUTF8) {
		t.Error("Expected invalid UTF-8 to be detected")
	}
}

// TestNegativeValues tests proper handling of negative numbers
func TestNegativeValues(t *testing.T) {
	// Negative string length should error
	badData := `s:-5:"test";`
	_, err := Unmarshal(badData)
	if err == nil {
		t.Error("Expected error for negative string length")
	}

	// Negative array count should error
	badData = `a:-1:{}`
	_, err = Unmarshal(badData)
	if err == nil {
		t.Error("Expected error for negative array count")
	}
}

// TestCombinedOptions tests using multiple options together
func TestCombinedOptions(t *testing.T) {
	data := map[string]interface{}{
		"nested": map[string]interface{}{
			"deep": map[string]interface{}{
				"value": "test",
			},
		},
		"zebra": 1,
		"apple": 2,
	}

	// Combine multiple options
	result, err := Marshal(data, WithMaxDepth(10), WithStrictPHP(true))
	if err != nil {
		t.Fatalf("Marshal with combined options failed: %v", err)
	}

	// Unmarshal with options
	unmarshalled, err := Unmarshal(result, WithMaxDepth(10))
	if err != nil {
		t.Fatalf("Unmarshal with combined options failed: %v", err)
	}

	if unmarshalled == nil {
		t.Error("Unmarshal returned nil")
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("empty string key", func(t *testing.T) {
		m := map[string]interface{}{
			"": "empty key",
		}
		data, err := Marshal(m)
		if err != nil {
			t.Errorf("Failed to marshal map with empty key: %v", err)
		}
		result, err := Unmarshal(data)
		if err != nil {
			t.Errorf("Failed to unmarshal: %v", err)
		}
		if resultMap, ok := result.(map[string]interface{}); ok {
			if resultMap[""] != "empty key" {
				t.Error("Empty key not preserved")
			}
		}
	})

	t.Run("single element array", func(t *testing.T) {
		arr := []interface{}{42}
		data, err := Marshal(arr)
		if err != nil {
			t.Errorf("Failed to marshal: %v", err)
		}
		result, err := Unmarshal(data)
		if err != nil {
			t.Errorf("Failed to unmarshal: %v", err)
		}
		if slice, ok := result.([]interface{}); !ok || len(slice) != 1 {
			t.Error("Single element array not handled correctly")
		}
	})

	t.Run("zero float", func(t *testing.T) {
		data, _ := Marshal(0.0)
		result, _ := Unmarshal(data)
		if result.(float64) != 0.0 {
			t.Error("Zero float not handled correctly")
		}
	})
}

// TestRealWorldPHPData tests with actual PHP serialized data
func TestRealWorldPHPData(t *testing.T) {
	// Sample data that might come from a real PHP application
	realWorldExamples := []string{
		// Simple user object
		`O:4:"User":3:{s:2:"id";i:1;s:4:"name";s:4:"John";s:5:"email";s:13:"john@test.com";}`,

		// Array of users
		`a:2:{i:0;O:4:"User":2:{s:2:"id";i:1;s:4:"name";s:4:"John";}i:1;O:4:"User":2:{s:2:"id";i:2;s:4:"name";s:4:"Jane";}}`,

		// Complex nested structure
		`a:3:{s:5:"users";a:1:{i:0;a:2:{s:2:"id";i:1;s:4:"name";s:4:"John";}}s:5:"total";i:1;s:4:"page";i:1;}`,

		// Session data
		`a:4:{s:7:"user_id";i:123;s:8:"username";s:8:"john_doe";s:9:"logged_in";b:1;s:10:"last_visit";i:1609459200;}`,
	}

	for i, phpData := range realWorldExamples {
		t.Run(fmt.Sprintf("example_%d", i+1), func(t *testing.T) {
			result, err := Unmarshal(phpData)
			if err != nil {
				t.Errorf("Failed to unmarshal real-world PHP data: %v", err)
			}
			if result == nil {
				t.Error("Unmarshal returned nil")
			}
		})
	}
}

// BenchmarkMarshalSimple benchmarks simple value marshaling
func BenchmarkMarshalSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Marshal(42)
	}
}

// BenchmarkMarshalArray benchmarks array marshaling
func BenchmarkMarshalArray(b *testing.B) {
	arr := make([]interface{}, 100)
	for i := range arr {
		arr[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Marshal(arr)
	}
}

// BenchmarkMarshalMap benchmarks map marshaling
func BenchmarkMarshalMap(b *testing.B) {
	m := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		m[fmt.Sprintf("key_%d", i)] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Marshal(m)
	}
}

// BenchmarkUnmarshalSimple benchmarks simple value unmarshalling
func BenchmarkUnmarshalSimple(b *testing.B) {
	data := "i:42;"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Unmarshal(data)
	}
}

// BenchmarkUnmarshalArray benchmarks array unmarshalling
func BenchmarkUnmarshalArray(b *testing.B) {
	arr := make([]interface{}, 100)
	for i := range arr {
		arr[i] = i
	}
	data, _ := Marshal(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Unmarshal(data)
	}
}

// BenchmarkUnmarshalMap benchmarks map unmarshalling
func BenchmarkUnmarshalMap(b *testing.B) {
	m := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		m[fmt.Sprintf("key_%d", i)] = i
	}
	data, _ := Marshal(m)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Unmarshal(data)
	}
}

// BenchmarkDeterministicMap benchmarks map serialization
func BenchmarkMapSerialization(b *testing.B) {
	m := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		m[fmt.Sprintf("key_%d", i)] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Marshal(m)
	}
}
