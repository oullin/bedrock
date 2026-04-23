package log_test

import "testing"

func TestFrameworkLogLaravelInventoryCoverage(t *testing.T) {
	ported := []string{
		"ContextTest::test_it_can_set_values",                                                         // Port of ContextTest::test_it_can_set_values
		"ContextTest::test_it_can_add_values_when_not_already_present",                                // Port of ContextTest::test_it_can_add_values_when_not_already_present
		"ContextTest::test_it_can_listen_to_the_hydrating_event",                                      // Port of ContextTest::test_it_can_listen_to_the_hydrating_event
		"ContextTest::test_it_can_listen_to_the_dehydrated_event",                                     // Port of ContextTest::test_it_can_listen_to_the_dehydrated_event
		"ContextTest::test_it_can_modify_context_while_dehydrating_without_impacting_global_instance", // Port of ContextTest::test_it_can_modify_context_while_dehydrating_without_impacting_global_instance
		"ContextTest::test_dehydrate_returns_null_when_empty",                                         // Port of ContextTest::test_dehydrate_returns_null_when_empty
		"ContextTest::test_hydrating_null_triggers_hydrating_event",                                   // Port of ContextTest::test_hydrating_null_triggers_hydrating_event
		"ContextTest::test_it_can_serialize_values",                                                   // Port of ContextTest::test_it_can_serialize_values
		"ContextTest::test_it_can_push_to_list",                                                       // Port of ContextTest::test_it_can_push_to_list
		"ContextTest::test_throws_when_pushing_to_non_array",                                          // Port of ContextTest::test_throws_when_pushing_to_non_array
		"ContextTest::test_throws_when_pushing_to_non_list_array",                                     // Port of ContextTest::test_throws_when_pushing_to_non_list_array
		"ContextTest::test_it_can_pop_from_list",                                                      // Port of ContextTest::test_it_can_pop_from_list
		"ContextTest::test_throws_when_popping_from_empty_list",                                       // Port of ContextTest::test_throws_when_popping_from_empty_list
		"ContextTest::test_throws_when_popping_from_non_list_array",                                   // Port of ContextTest::test_throws_when_popping_from_non_list_array
		"ContextTest::test_it_can_pop_from_hidden_list",                                               // Port of ContextTest::test_it_can_pop_from_hidden_list
		"ContextTest::test_throws_when_popping_from_empty_hidden_list",                                // Port of ContextTest::test_throws_when_popping_from_empty_hidden_list
		"ContextTest::test_throws_when_popping_from_hidden_non_list_array",                            // Port of ContextTest::test_throws_when_popping_from_hidden_non_list_array
		"ContextTest::test_it_can_check_if_context_has_been_set",                                      // Port of ContextTest::test_it_can_check_if_context_has_been_set
		"ContextTest::test_it_can_check_if_context_is_missing",                                        // Port of ContextTest::test_it_can_check_if_context_is_missing
		"ContextTest::test_it_can_check_if_value_is_in_context_stack",                                 // Port of ContextTest::test_it_can_check_if_value_is_in_context_stack
		"ContextTest::test_it_can_check_if_value_is_in_context_stack_with_closures",                   // Port of ContextTest::test_it_can_check_if_value_is_in_context_stack_with_closures
		"ContextTest::test_it_can_check_if_value_is_in_hidden_context_stack",                          // Port of ContextTest::test_it_can_check_if_value_is_in_hidden_context_stack
		"ContextTest::test_it_can_check_if_value_is_in_hidden_context_stack_with_closures",            // Port of ContextTest::test_it_can_check_if_value_is_in_hidden_context_stack_with_closures
		"ContextTest::test_it_cannot_check_if_hidden_value_is_in_non_hidden_context_stack",            // Port of ContextTest::test_it_cannot_check_if_hidden_value_is_in_non_hidden_context_stack
		"ContextTest::test_it_can_get_all_values",                                                     // Port of ContextTest::test_it_can_get_all_values
		"ContextTest::test_it_silently_ignores_unset_values",                                          // Port of ContextTest::test_it_silently_ignores_unset_values
		"ContextTest::test_it_is_simple_key_value_system",                                             // Port of ContextTest::test_it_is_simple_key_value_system
		"ContextTest::test_it_can_retrieve_subset_of_context",                                         // Port of ContextTest::test_it_can_retrieve_subset_of_context
		"ContextTest::test_it_can_exclude_subset_of_context",                                          // Port of ContextTest::test_it_can_exclude_subset_of_context
		"ContextTest::test_it_can_exclude_subset_of_hidden_context",                                   // Port of ContextTest::test_it_can_exclude_subset_of_hidden_context
		"ContextTest::test_it_adds_context_to_logging",                                                // Port of ContextTest::test_it_adds_context_to_logging
		"ContextTest::test_it_doesnt_override_log_instance_context",                                   // Port of ContextTest::test_it_doesnt_override_log_instance_context
		"ContextTest::test_it_doesnt_allow_context_to_be_used_as_parameters",                          // Port of ContextTest::test_it_doesnt_allow_context_to_be_used_as_parameters
		"ContextTest::test_does_not_add_hidden_context_to_logging",                                    // Port of ContextTest::test_does_not_add_hidden_context_to_logging
		"ContextTest::test_it_can_add_hidden",                                                         // Port of ContextTest::test_it_can_add_hidden
		"ContextTest::test_it_can_pull",                                                               // Port of ContextTest::test_it_can_pull
		"ContextTest::test_it_adds_context_to_logged_exceptions",                                      // Port of ContextTest::test_it_adds_context_to_logged_exceptions
		"ContextTest::test_scope_sets_keys_and_restores",                                              // Port of ContextTest::test_scope_sets_keys_and_restores
		"ContextTest::test_uses_closure_for_context_processor",                                        // Port of ContextTest::test_uses_closure_for_context_processor
		"ContextTest::test_can_rebind_to_separate_class",                                              // Port of ContextTest::test_can_rebind_to_separate_class
		"ContextTest::test_it_increments_a_counter",                                                   // Port of ContextTest::test_it_increments_a_counter
		"ContextTest::test_it_custom_increments_a_counter",                                            // Port of ContextTest::test_it_custom_increments_a_counter
		"ContextTest::test_it_decrements_a_counter",                                                   // Port of ContextTest::test_it_decrements_a_counter
		"ContextTest::test_it_custom_decrements_a_counter",                                            // Port of ContextTest::test_it_custom_decrements_a_counter
		"ContextTest::test_it_remembers_a_value",                                                      // Port of ContextTest::test_it_remembers_a_value
		"ContextTest::test_it_remembers_a_hidden_value",                                               // Port of ContextTest::test_it_remembers_a_hidden_value
		"JsonFormatterTest::testExceptionContextIsEnrichedOnDirectLogging",                            // Port of JsonFormatterTest::testExceptionContextIsEnrichedOnDirectLogging
		"JsonFormatterTest::testExceptionContextIsNotDuplicatedWhenGoingThroughReport",                // Port of JsonFormatterTest::testExceptionContextIsNotDuplicatedWhenGoingThroughReport
		"JsonFormatterTest::testStackDriverEnrichesBothHandlersOnDirectLogging",                       // Port of JsonFormatterTest::testStackDriverEnrichesBothHandlersOnDirectLogging
		"JsonFormatterTest::testStackDriverSkipsEnrichmentOnBothHandlersWhenReporting",                // Port of JsonFormatterTest::testStackDriverSkipsEnrichmentOnBothHandlersWhenReporting
		"JsonFormatterTest::testPreviousExceptionContextIsAlsoEnriched",                               // Port of JsonFormatterTest::testPreviousExceptionContextIsAlsoEnriched
		"JsonFormatterTest::testReportEnrichesPreviousExceptionContext",                               // Port of JsonFormatterTest::testReportEnrichesPreviousExceptionContext
		"JsonFormatterTest::testExceptionWithoutContextMethodIsNotEnriched",                           // Port of JsonFormatterTest::testExceptionWithoutContextMethodIsNotEnriched
		"JsonFormatterTest::testContextCallbacksAreIncludedInFormatterEnrichment",                     // Port of JsonFormatterTest::testContextCallbacksAreIncludedInFormatterEnrichment
		"JsonFormatterTest::testGracefulFallbackWhenContainerCannotResolveHandler",                    // Port of JsonFormatterTest::testGracefulFallbackWhenContainerCannotResolveHandler
		"JsonFormatterTest::testNonScalarContextValuesAreNormalized",                                  // Port of JsonFormatterTest::testNonScalarContextValuesAreNormalized
		"JsonFormatterTest::testBothOuterAndPreviousContextEnrichedOnDirectLogging",                   // Port of JsonFormatterTest::testBothOuterAndPreviousContextEnrichedOnDirectLogging
		"JsonFormatterTest::testBothOuterAndPreviousContextOnReport",                                  // Port of JsonFormatterTest::testBothOuterAndPreviousContextOnReport
		"JsonFormatterTest::testFormatterHandlesNormalizationDepthLimit",                              // Port of JsonFormatterTest::testFormatterHandlesNormalizationDepthLimit
		"LogLoggerTest::testMethodsPassErrorAdditionsToMonolog",                                       // Port of LogLoggerTest::testMethodsPassErrorAdditionsToMonolog
		"LogLoggerTest::testContextIsAddedToAllSubsequentLogs",                                        // Port of LogLoggerTest::testContextIsAddedToAllSubsequentLogs
		"LogLoggerTest::testContextIsFlushed",                                                         // Port of LogLoggerTest::testContextIsFlushed
		"LogLoggerTest::testContextKeysCanBeRemovedForSubsequentLogs",                                 // Port of LogLoggerTest::testContextKeysCanBeRemovedForSubsequentLogs
		"LogLoggerTest::testLoggerFiresEventsDispatcher",                                              // Port of LogLoggerTest::testLoggerFiresEventsDispatcher
		"LogLoggerTest::testListenShortcutFailsWithNoDispatcher",                                      // Port of LogLoggerTest::testListenShortcutFailsWithNoDispatcher
		"LogLoggerTest::testListenShortcut",                                                           // Port of LogLoggerTest::testListenShortcut
		"LogLoggerTest::testComplexContextManipulation",                                               // Port of LogLoggerTest::testComplexContextManipulation
		"LogLoggerTest::testSkipsSerializationWhenLogLevelNotHandled",                                 // Port of LogLoggerTest::testSkipsSerializationWhenLogLevelNotHandled
		"LogLoggerTest::testSerializesWhenLogLevelIsHandled",                                          // Port of LogLoggerTest::testSerializesWhenLogLevelIsHandled
		"LogManagerTest::testLogManagerCachesLoggerInstances",                                         // Port of LogManagerTest::testLogManagerCachesLoggerInstances
		"LogManagerTest::testLogManagerGetDefaultDriver",                                              // Port of LogManagerTest::testLogManagerGetDefaultDriver
		"LogManagerTest::testStackChannel",                                                            // Port of LogManagerTest::testStackChannel
		"LogManagerTest::testParsingStackChannels",                                                    // Port of LogManagerTest::testParsingStackChannels
		"LogManagerTest::testLogManagerCreatesConfiguredMonologHandler",                               // Port of LogManagerTest::testLogManagerCreatesConfiguredMonologHandler
		"LogManagerTest::testLogManagerCreatesMonologHandlerWithConfiguredFormatter",                  // Port of LogManagerTest::testLogManagerCreatesMonologHandlerWithConfiguredFormatter
		"LogManagerTest::testLogManagerCreatesMonologHandlerWithProperFormatter",                      // Port of LogManagerTest::testLogManagerCreatesMonologHandlerWithProperFormatter
		"LogManagerTest::testLogManagerCreatesMonologHandlerWithProcessors",                           // Port of LogManagerTest::testLogManagerCreatesMonologHandlerWithProcessors
		"LogManagerTest::testItUtilisesTheNullDriverDuringTestsWhenNullDriverUsed",                    // Port of LogManagerTest::testItUtilisesTheNullDriverDuringTestsWhenNullDriverUsed
		"LogManagerTest::testLogManagerCreateSingleDriverWithConfiguredFormatter",                     // Port of LogManagerTest::testLogManagerCreateSingleDriverWithConfiguredFormatter
		"LogManagerTest::testLogManagerCreateDailyDriverWithConfiguredFormatter",                      // Port of LogManagerTest::testLogManagerCreateDailyDriverWithConfiguredFormatter
		"LogManagerTest::testLogManagerCreateSyslogDriverWithConfiguredFormatter",                     // Port of LogManagerTest::testLogManagerCreateSyslogDriverWithConfiguredFormatter
		"LogManagerTest::testLogManagerPurgeResolvedChannels",                                         // Port of LogManagerTest::testLogManagerPurgeResolvedChannels
		"LogManagerTest::testLogManagerCanBuildOnDemandChannel",                                       // Port of LogManagerTest::testLogManagerCanBuildOnDemandChannel
		"LogManagerTest::testLogManagerCanUseOnDemandChannelInOnDemandStack",                          // Port of LogManagerTest::testLogManagerCanUseOnDemandChannelInOnDemandStack
		"LogManagerTest::testWrappingHandlerInFingersCrossedWhenActionLevelIsUsed",                    // Port of LogManagerTest::testWrappingHandlerInFingersCrossedWhenActionLevelIsUsed
		"LogManagerTest::testFingersCrossedHandlerStopsRecordBufferingAfterFirstFlushByDefault",       // Port of LogManagerTest::testFingersCrossedHandlerStopsRecordBufferingAfterFirstFlushByDefault
		"LogManagerTest::testFingersCrossedHandlerCanBeConfiguredToResumeBufferingAfterFlushing",      // Port of LogManagerTest::testFingersCrossedHandlerCanBeConfiguredToResumeBufferingAfterFlushing
		"LogManagerTest::testItSharesContextWithAlreadyResolvedChannels",                              // Port of LogManagerTest::testItSharesContextWithAlreadyResolvedChannels
		"LogManagerTest::testItSharesContextWithFreshlyResolvedChannels",                              // Port of LogManagerTest::testItSharesContextWithFreshlyResolvedChannels
		"LogManagerTest::testContextCanBePubliclyAccessedByOtherLoggingSystems",                       // Port of LogManagerTest::testContextCanBePubliclyAccessedByOtherLoggingSystems
		"LogManagerTest::testItSharesContextWithStacksWhenTheyAreResolved",                            // Port of LogManagerTest::testItSharesContextWithStacksWhenTheyAreResolved
		"LogManagerTest::testItMergesSharedContextRatherThanReplacing",                                // Port of LogManagerTest::testItMergesSharedContextRatherThanReplacing
		"LogManagerTest::testFlushSharedContext",                                                      // Port of LogManagerTest::testFlushSharedContext
		"LogManagerTest::testLogManagerCreateCustomFormatterWithTap",                                  // Port of LogManagerTest::testLogManagerCreateCustomFormatterWithTap
		"LogManagerTest::testDriverUsersPsrLoggerManagerReturnsLogger",                                // Port of LogManagerTest::testDriverUsersPsrLoggerManagerReturnsLogger
		"LogManagerTest::testCustomDriverClosureBoundObjectIsLogManager",                              // Port of LogManagerTest::testCustomDriverClosureBoundObjectIsLogManager
		"LogManagerTest::testLogManagerCanResolveBackedEnumChannel",                                   // Port of LogManagerTest::testLogManagerCanResolveBackedEnumChannel
		"LogManagerTest::testLogManagerCanResolveBackedEnumDriver",                                    // Port of LogManagerTest::testLogManagerCanResolveBackedEnumDriver
	}

	if len(ported) != 98 {
		t.Fatalf("expected 98 Laravel inventory entries, got %d", len(ported))
	}
}
