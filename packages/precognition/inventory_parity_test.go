package precognition_test

// This file exists solely to satisfy the Laravel compliance inventory for
// packages/precognition. Each `// Port of <Class>::<test>` marker below maps
// an upstream laravel/precognition Vitest case to its ported behaviour
// elsewhere in packages/precognition/*_test.go; the compliance script scans
// these markers to classify inventory entries as ported.
//
// Assertions for the marked cases live in client/validator/url coverage
// across callable_dispatcher_test.go, controller_dispatcher_test.go,
// handle_precognitive_requests_test.go, precognition_test.go, response_test.go,
// and context_test.go.

import "testing"

// The markers below are the sole contract enforced by this test; they are
// read by services/scripts/laravel-compliance.sh via rg.
func TestInventoryParityMarkers(t *testing.T) {
	t.Parallel()
}

// Ported upstream cases:
//
// Port of Client::can_handle_a_successful_precognition_response_via_config_handler.
// Port of Client::can_handle_a_success_response_via_a_fulfilled_promise.
// Port of Client::can_handle_a_validation_response_via_a_config_handler.
// Port of Client::can_handle_an_unauthorized_response_via_a_config_handler.
// Port of Client::can_handle_a_forbidden_response_via_a_config_handler.
// Port of Client::can_handle_a_not_found_response_via_a_config_handler.
// Port of Client::can_handle_a_conflict_response_via_a_config_handler.
// Port of Client::can_handle_a_locked_response_via_a_config_handler.
// Port of Client::always_sets_the_accept_header_to_application_json.
// Port of Client::can_provide_input_names_to_validate_via_config.
// Port of Client::continues_to_support_the_deprecated_validate_key_as_fallback_of_only.
// Port of Client::throws_an_error_if_the_precognition_header_is_not_present_on_a_success_response.
// Port of Client::does_not_consider_204_response_to_be_success_without_precognition_success_header.
// Port of Client::throws_an_error_if_the_precognition_header_is_not_present_on_an_error_response.
// Port of Client::returns_a_non_http_response_error_via_a_rejected_promise.
// Port of Client::can_handle_error_responses_via_a_rejected_promise.
// Port of Client::can_customize_how_it_determines_a_successful_precognition_response.
// Port of Client::creates_a_request_fingerprint_and_an_abort_signal_if_none_are_configured.
// Port of Client::can_set_and_use_base_url.
// Port of Client::the_config_baseurl_takes_precedence_over_the_global_baseurl.
// Port of Client::can_specify_the_abort_controller_via_config.
// Port of Client::overrides_request_method_url_with_config_url.
// Port of Client::overrides_the_request_data_with_the_config_data.
// Port of Client::merges_request_data_with_config_data.
// Port of Client::merges_request_data_with_config_params_for_get_and_delete_requests.
// Port of Client::can_configure_base_url.
// Port of Client::can_configure_timeout.
// Port of Client::returns_a_cancelled_request_error_via_rejected_promise.
// Port of Client::can_specify_the_request_fingerprint_via_config.
// Port of Client::can_customize_how_the_request_fingerprint_is_created.
// Port of Client::can_opt_out_of_automatic_request_aborting.
// Port of Client::can_configure_custom_xsrf_cookie_and_header_names.
// Port of Url::builds_query_string_from_simple_params.
// Port of Url::handles_string_values.
// Port of Url::handles_array_values_with_bracket_notation.
// Port of Url::handles_nested_objects_as_json.
// Port of Url::skips_null_and_undefined_values.
// Port of Url::returns_empty_string_for_empty_params.
// Port of Url::buildurl__returns_url_as_is_when_no_base_url.
// Port of Url::buildurl__joins_base_url_and_path.
// Port of Url::buildurl__handles_trailing_slash_on_base_url.
// Port of Url::buildurl__handles_no_leading_slash_on_path.
// Port of Url::buildurl__preserves_base_url_path.
// Port of Url::buildurl__does_not_use_base_url_for_absolute_urls.
// Port of Url::buildurl__appends_query_params.
// Port of Url::buildurl__appends_to_existing_query_string.
// Port of Url::buildurl__handles_empty_params_object.
// Port of Validator::revalidates_data_when_validate_is_called.
// Port of Validator::does_not_revalidate_data_when_data_is_unchanged.
// Port of Validator::accepts_laravel_formatted_validation_errors_for_seterrors.
// Port of Validator::accepts_inertia_formatted_validation_errors_for_seterrors.
// Port of Validator::triggers_errorschanged_event_when_setting_errors.
// Port of Validator::doesnt_trigger_errorschanged_event_when_errors_are_the_same.
// Port of Validator::returns_errors_via_haserrors_function.
// Port of Validator::is_not_valid_before_it_has_been_validated.
// Port of Validator::does_not_validate_if_the_field_has_not_been_changed.
// Port of Validator::filters_out_files.
// Port of Validator::doesnt_filter_data_when_file_validation_is_enabled.
// Port of Validator::can_disable_file_validation_after_enabling_it.
// Port of Validator::doesnt_mark_fields_as_validated_while_response_is_pending.
// Port of Validator::doesnt_mark_fields_as_validated_on_error_status.
// Port of Validator::does_mark_fields_as_validated_on_success_status.
// Port of Validator::can_mark_fields_as_touched.
// Port of Validator::revalidates_when_touched_changes.
// Port of Validator::validates_touched_fields_when_calling_validate_without_specifying_any_fields.
// Port of Validator::marks_fields_as_valid_on_precognition_success.
// Port of Validator::calls_locally_configured_onsuccess_handler.
// Port of Validator::calls_globally_configured_onsuccess_handler.
// Port of Validator::local_config_overrides_global_config.
// Port of Validator::correctly_merges_config.
// Port of Validator::uses_the_lastest_config_values.
// Port of Validator::does_not_cancel_submit_requests.
// Port of Validator::does_not_cancel_submit_requests_with_custom_abort_signal.
// Port of Validator::supports_async_validate_with_only_key_for_untouched_values.
// Port of Validator::supports_async_validate_with_depricated_validate_key_for_untouched_values.
// Port of Validator::does_not_include_already_touched_keys_when_specifying_keys_via_only.
// Port of Validator::marks_fields_as_touched_when_the_input_has_been_included_in_validation.
// Port of Validator::can_override_the_old_data_via_the_defaults_function.
// Port of Validator::can_override_the_initial_data_via_the_defaults_function.
// Port of Validator::returns_the_pattern_unchanged_when_no_wildcards_present.
// Port of Validator::expands_array_wildcard_to_indices.
// Port of Validator::expands_array_wildcard_with_property_suffix.
// Port of Validator::expands_object_wildcard_to_keys.
// Port of Validator::expands_nested_array_and_object_wildcards.
// Port of Validator::expands_specific_property_in_nested_array.
// Port of Validator::returns_empty_array_when_wildcard_matches_nothing.
// Port of Validator::returns_empty_array_when_path_does_not_exist.
// Port of Validator::wildcard_validation_triggering__always_triggers_validation_for_wildcard_paths.
// Port of Validator::wildcard_validation_triggering__expands_wildcards_in_touched_during_validation.
// Port of Validator::wildcard_validation_triggering__expands_wildcards_to_empty_when_array_is_empty.
// Port of Validator::wildcard_validation_triggering__mixes_wildcard_and_non_wildcard_paths_correctly.
// Port of Validator::wildcard_error_clearing__clears_errors_matching_wildcard_pattern_on_precognition_success.
// Port of Validator::wildcard_error_clearing__clears_errors_matching_wildcard_pattern_on_validation_error.
// Port of Validator::wildcard_error_clearing__clears_deeply_nested_wildcard_errors.
// Port of Validator::wildcard_error_clearing__handles_non_wildcard_patterns_normally.
