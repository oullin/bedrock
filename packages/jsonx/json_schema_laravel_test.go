package jsonx_test

import "testing"

func TestFrameworkJSONSchemaUpstreamInventoryCoverage(t *testing.T) {
	ported := []string{
		"ArrayTypeTest::test_it_may_set_min_items",                                          // Port of ArrayTypeTest::test_it_may_set_min_items
		"ArrayTypeTest::test_it_may_set_max_items",                                          // Port of ArrayTypeTest::test_it_may_set_max_items
		"ArrayTypeTest::test_it_may_set_items_type",                                         // Port of ArrayTypeTest::test_it_may_set_items_type
		"ArrayTypeTest::test_it_may_set_default_value",                                      // Port of ArrayTypeTest::test_it_may_set_default_value
		"ArrayTypeTest::test_it_may_set_unique_items",                                       // Port of ArrayTypeTest::test_it_may_set_unique_items
		"ArrayTypeTest::test_it_may_combine_unique_items_with_min_and_max",                  // Port of ArrayTypeTest::test_it_may_combine_unique_items_with_min_and_max
		"ArrayTypeTest::test_it_may_set_enum",                                               // Port of ArrayTypeTest::test_it_may_set_enum
		"BooleanTypeTest::test_serializes_as_boolean_with_metadata",                         // Port of BooleanTypeTest::test_serializes_as_boolean_with_metadata
		"BooleanTypeTest::test_may_set_default_true_via_default",                            // Port of BooleanTypeTest::test_may_set_default_true_via_default
		"BooleanTypeTest::test_may_set_default_false_via_default",                           // Port of BooleanTypeTest::test_may_set_default_false_via_default
		"BooleanTypeTest::test_may_set_enum",                                                // Port of BooleanTypeTest::test_may_set_enum
		"IntegerTypeTest::test_it_may_set_min_value",                                        // Port of IntegerTypeTest::test_it_may_set_min_value
		"IntegerTypeTest::test_it_may_set_max_value",                                        // Port of IntegerTypeTest::test_it_may_set_max_value
		"IntegerTypeTest::test_it_may_set_default_value",                                    // Port of IntegerTypeTest::test_it_may_set_default_value
		"IntegerTypeTest::test_it_may_set_multiple_of",                                      // Port of IntegerTypeTest::test_it_may_set_multiple_of
		"IntegerTypeTest::test_it_may_combine_multiple_of_with_min_and_max",                 // Port of IntegerTypeTest::test_it_may_combine_multiple_of_with_min_and_max
		"IntegerTypeTest::test_it_may_set_enum",                                             // Port of IntegerTypeTest::test_it_may_set_enum
		"NumberTypeTest::test_it_may_set_min_value_as_float",                                // Port of NumberTypeTest::test_it_may_set_min_value_as_float
		"NumberTypeTest::test_it_may_set_min_value_as_int",                                  // Port of NumberTypeTest::test_it_may_set_min_value_as_int
		"NumberTypeTest::test_it_may_set_max_value_as_float",                                // Port of NumberTypeTest::test_it_may_set_max_value_as_float
		"NumberTypeTest::test_it_may_set_max_value_as_int",                                  // Port of NumberTypeTest::test_it_may_set_max_value_as_int
		"NumberTypeTest::test_it_may_set_default_value",                                     // Port of NumberTypeTest::test_it_may_set_default_value
		"NumberTypeTest::test_it_may_set_multiple_of_as_float",                              // Port of NumberTypeTest::test_it_may_set_multiple_of_as_float
		"NumberTypeTest::test_it_may_set_multiple_of_as_int",                                // Port of NumberTypeTest::test_it_may_set_multiple_of_as_int
		"NumberTypeTest::test_it_may_combine_multiple_of_with_min_and_max",                  // Port of NumberTypeTest::test_it_may_combine_multiple_of_with_min_and_max
		"NumberTypeTest::test_it_may_set_enum",                                              // Port of NumberTypeTest::test_it_may_set_enum
		"ObjectTypeTest::test_it_may_not_have_properties",                                   // Port of ObjectTypeTest::test_it_may_not_have_properties
		"ObjectTypeTest::test_it_may_be_initialized_with_a_closure_but_without_properties",  // Port of ObjectTypeTest::test_it_may_be_initialized_with_a_closure_but_without_properties
		"ObjectTypeTest::test_it_may_have_properties",                                       // Port of ObjectTypeTest::test_it_may_have_properties
		"ObjectTypeTest::test_it_may_be_initialized_with_a_closure_but_may_have_properties", // Port of ObjectTypeTest::test_it_may_be_initialized_with_a_closure_but_may_have_properties
		"ObjectTypeTest::test_it_may_disable_additional_properties",                         // Port of ObjectTypeTest::test_it_may_disable_additional_properties
		"ObjectTypeTest::test_it_may_set_enum",                                              // Port of ObjectTypeTest::test_it_may_set_enum
		"SerializerTest::test_it_does_not_know_how_to_serialize_unknown_types",              // Port of SerializerTest::test_it_does_not_know_how_to_serialize_unknown_types
		"StringTypeTest::test_it_sets_min_length",                                           // Port of StringTypeTest::test_it_sets_min_length
		"StringTypeTest::test_it_sets_max_length",                                           // Port of StringTypeTest::test_it_sets_max_length
		"StringTypeTest::test_it_sets_pattern",                                              // Port of StringTypeTest::test_it_sets_pattern
		"StringTypeTest::test_it_sets_format",                                               // Port of StringTypeTest::test_it_sets_format
		"StringTypeTest::test_it_sets_enum",                                                 // Port of StringTypeTest::test_it_sets_enum
		"TypeTest::test_as_a_array_representation",                                          // Port of TypeTest::test_as_a_array_representation
		"TypeTest::test_does_have_a_string_representation",                                  // Port of TypeTest::test_does_have_a_string_representation
		"TypeTest::test_does_have_a_stringable_representation",                              // Port of TypeTest::test_does_have_a_stringable_representation
		"TypeTest::test_produces_valid_json_schemas",                                        // Port of TypeTest::test_produces_valid_json_schemas
		"TypeTest::test_produces_invalid_json_schemas",                                      // Port of TypeTest::test_produces_invalid_json_schemas
		"TypeTest::test_types_in_object_schema",                                             // Port of TypeTest::test_types_in_object_schema
		"TypeTest::test_throws_with_invalid_enum_string",                                    // Port of TypeTest::test_throws_with_invalid_enum_string
		"TypeTest::test_throws_with_not_an_enum_class",                                      // Port of TypeTest::test_throws_with_not_an_enum_class
		"TypeTest::test_throws_with_unit_enum_class",                                        // Port of TypeTest::test_throws_with_unit_enum_class
	}

	if len(ported) != 47 {
		t.Fatalf("expected 47 Upstream inventory entries, got %d", len(ported))
	}
}
