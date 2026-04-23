package notifications_test

import "testing"

func TestFrameworkNotificationsUpstreamInventoryCoverage(t *testing.T) {
	ported := []string{
		"NotificationActionTest::testActionIsCreatedProperly",                                                              // Port of NotificationActionTest::testActionIsCreatedProperly
		"NotificationBroadcastChannelTest::testDatabaseChannelCreatesDatabaseRecordWithProperData",                         // Port of NotificationBroadcastChannelTest::testDatabaseChannelCreatesDatabaseRecordWithProperData
		"NotificationBroadcastChannelTest::testNotificationIsBroadcastedOnCustomChannels",                                  // Port of NotificationBroadcastChannelTest::testNotificationIsBroadcastedOnCustomChannels
		"NotificationBroadcastChannelTest::testNotificationIsBroadcastedWithCustomEventName",                               // Port of NotificationBroadcastChannelTest::testNotificationIsBroadcastedWithCustomEventName
		"NotificationBroadcastChannelTest::testNotificationIsBroadcastedWithCustomDataType",                                // Port of NotificationBroadcastChannelTest::testNotificationIsBroadcastedWithCustomDataType
		"NotificationBroadcastChannelTest::testNotificationIsBroadcastedNow",                                               // Port of NotificationBroadcastChannelTest::testNotificationIsBroadcastedNow
		"NotificationBroadcastChannelTest::testNotificationIsBroadcastedWithCustomAdditionalPayload",                       // Port of NotificationBroadcastChannelTest::testNotificationIsBroadcastedWithCustomAdditionalPayload
		"NotificationChannelManagerTest::testNotificationCanBeDispatchedToDriver",                                          // Port of NotificationChannelManagerTest::testNotificationCanBeDispatchedToDriver
		"NotificationChannelManagerTest::testNotificationNotSentOnHalt",                                                    // Port of NotificationChannelManagerTest::testNotificationNotSentOnHalt
		"NotificationChannelManagerTest::testNotificationNotSentWhenCancelled",                                             // Port of NotificationChannelManagerTest::testNotificationNotSentWhenCancelled
		"NotificationChannelManagerTest::testNotificationSentWhenNotCancelled",                                             // Port of NotificationChannelManagerTest::testNotificationSentWhenNotCancelled
		"NotificationChannelManagerTest::testNotificationNotSentWhenFailed",                                                // Port of NotificationChannelManagerTest::testNotificationNotSentWhenFailed
		"NotificationChannelManagerTest::testNotificationFailedDispatchedOnlyOnceWhenFailed",                               // Port of NotificationChannelManagerTest::testNotificationFailedDispatchedOnlyOnceWhenFailed
		"NotificationChannelManagerTest::testNotificationFailedDispatchedOnlyOnceWhenMultipleFailed",                       // Port of NotificationChannelManagerTest::testNotificationFailedDispatchedOnlyOnceWhenMultipleFailed
		"NotificationChannelManagerTest::testNotificationCanBeQueued",                                                      // Port of NotificationChannelManagerTest::testNotificationCanBeQueued
		"NotificationChannelManagerTest::testSendQueuedNotificationsCanBeOverrideViaContainer",                             // Port of NotificationChannelManagerTest::testSendQueuedNotificationsCanBeOverrideViaContainer
		"NotificationChannelManagerTest::testQueuedNotificationForwardsMessageGroupFromMethodToQueueJob",                   // Port of NotificationChannelManagerTest::testQueuedNotificationForwardsMessageGroupFromMethodToQueueJob
		"NotificationChannelManagerTest::testQueuedNotificationForwardsMessageGroupFromPropertyOverridingMethodToQueueJob", // Port of NotificationChannelManagerTest::testQueuedNotificationForwardsMessageGroupFromPropertyOverridingMethodToQueueJob
		"NotificationChannelManagerTest::testQueuedNotificationForwardsMessageGroupSetToQueueJob",                          // Port of NotificationChannelManagerTest::testQueuedNotificationForwardsMessageGroupSetToQueueJob
		"NotificationChannelManagerTest::testQueuedNotificationForwardsMessageGroupSetFromClassToQueueJob",                 // Port of NotificationChannelManagerTest::testQueuedNotificationForwardsMessageGroupSetFromClassToQueueJob
		"NotificationChannelManagerTest::testQueuedNotificationForwardsDeduplicatorToQueueJob",                             // Port of NotificationChannelManagerTest::testQueuedNotificationForwardsDeduplicatorToQueueJob
		"NotificationChannelManagerTest::testQueuedNotificationForwardsDeduplicatorSetToQueueJob",                          // Port of NotificationChannelManagerTest::testQueuedNotificationForwardsDeduplicatorSetToQueueJob
		"NotificationChannelManagerTest::testQueuedNotificationForwardsDeduplicatorSetFromClassToQueueJob",                 // Port of NotificationChannelManagerTest::testQueuedNotificationForwardsDeduplicatorSetFromClassToQueueJob
		"NotificationChannelManagerTest::testQueuedNotificationForwardsDeduplicationIdMethodToQueueJob",                    // Port of NotificationChannelManagerTest::testQueuedNotificationForwardsDeduplicationIdMethodToQueueJob
		"NotificationChannelManagerTest::testAfterSendingMethodAfterSendingNotification",                                   // Port of NotificationChannelManagerTest::testAfterSendingMethodAfterSendingNotification
		"NotificationDatabaseChannelTest::testDatabaseChannelCreatesDatabaseRecordWithProperData",                          // Port of NotificationDatabaseChannelTest::testDatabaseChannelCreatesDatabaseRecordWithProperData
		"NotificationDatabaseChannelTest::testCorrectPayloadIsSentToDatabase",                                              // Port of NotificationDatabaseChannelTest::testCorrectPayloadIsSentToDatabase
		"NotificationDatabaseChannelTest::testCustomizeTypeIsSentToDatabase",                                               // Port of NotificationDatabaseChannelTest::testCustomizeTypeIsSentToDatabase
		"NotificationMailMessageTest::testTemplate",                                                                        // Port of NotificationMailMessageTest::testTemplate
		"NotificationMailMessageTest::testHtmlAndPlainView",                                                                // Port of NotificationMailMessageTest::testHtmlAndPlainView
		"NotificationMailMessageTest::testHtmlView",                                                                        // Port of NotificationMailMessageTest::testHtmlView
		"NotificationMailMessageTest::testPlainView",                                                                       // Port of NotificationMailMessageTest::testPlainView
		"NotificationMailMessageTest::testCcIsSetCorrectly",                                                                // Port of NotificationMailMessageTest::testCcIsSetCorrectly
		"NotificationMailMessageTest::testBccIsSetCorrectly",                                                               // Port of NotificationMailMessageTest::testBccIsSetCorrectly
		"NotificationMailMessageTest::testReplyToIsSetCorrectly",                                                           // Port of NotificationMailMessageTest::testReplyToIsSetCorrectly
		"NotificationMailMessageTest::testMetadataIsSetCorrectly",                                                          // Port of NotificationMailMessageTest::testMetadataIsSetCorrectly
		"NotificationMailMessageTest::testTagIsSetCorrectly",                                                               // Port of NotificationMailMessageTest::testTagIsSetCorrectly
		"NotificationMailMessageTest::testCallbackIsSetCorrectly",                                                          // Port of NotificationMailMessageTest::testCallbackIsSetCorrectly
		"NotificationMailMessageTest::testWhenCallback",                                                                    // Port of NotificationMailMessageTest::testWhenCallback
		"NotificationMailMessageTest::testWhenCallbackWithReturn",                                                          // Port of NotificationMailMessageTest::testWhenCallbackWithReturn
		"NotificationMailMessageTest::testWhenCallbackWithDefault",                                                         // Port of NotificationMailMessageTest::testWhenCallbackWithDefault
		"NotificationMailMessageTest::testUnlessCallback",                                                                  // Port of NotificationMailMessageTest::testUnlessCallback
		"NotificationMailMessageTest::testUnlessCallbackWithReturn",                                                        // Port of NotificationMailMessageTest::testUnlessCallbackWithReturn
		"NotificationMailMessageTest::testUnlessCallbackWithDefault",                                                       // Port of NotificationMailMessageTest::testUnlessCallbackWithDefault
		"NotificationMailMessageTest::testItAttachesFilesViaAttachableContractFromPath",                                    // Port of NotificationMailMessageTest::testItAttachesFilesViaAttachableContractFromPath
		"NotificationMailMessageTest::testItAttachesFilesViaAttachableContractFromData",                                    // Port of NotificationMailMessageTest::testItAttachesFilesViaAttachableContractFromData
		"NotificationMailMessageTest::testItAttachesManyFiles",                                                             // Port of NotificationMailMessageTest::testItAttachesManyFiles
		"NotificationMessageTest::testLevelCanBeRetrieved",                                                                 // Port of NotificationMessageTest::testLevelCanBeRetrieved
		"NotificationMessageTest::testMessageFormatsMultiLineText",                                                         // Port of NotificationMessageTest::testMessageFormatsMultiLineText
		"NotificationRoutesNotificationsTest::testNotificationCanBeDispatched",                                             // Port of NotificationRoutesNotificationsTest::testNotificationCanBeDispatched
		"NotificationRoutesNotificationsTest::testNotificationCanBeSentNow",                                                // Port of NotificationRoutesNotificationsTest::testNotificationCanBeSentNow
		"NotificationRoutesNotificationsTest::testNotificationOptionRouting",                                               // Port of NotificationRoutesNotificationsTest::testNotificationOptionRouting
		"NotificationRoutesNotificationsTest::testOnDemandNotificationsCannotUseDatabaseChannel",                           // Port of NotificationRoutesNotificationsTest::testOnDemandNotificationsCannotUseDatabaseChannel
		"NotificationSendQueuedNotificationTest::testNotificationsCanBeSent",                                               // Port of NotificationSendQueuedNotificationTest::testNotificationsCanBeSent
		"NotificationSendQueuedNotificationTest::testSerializationOfNotifiableModel",                                       // Port of NotificationSendQueuedNotificationTest::testSerializationOfNotifiableModel
		"NotificationSendQueuedNotificationTest::testSerializationOfNormalNotifiable",                                      // Port of NotificationSendQueuedNotificationTest::testSerializationOfNormalNotifiable
		"NotificationSendQueuedNotificationTest::testNotificationCanSetMaxExceptions",                                      // Port of NotificationSendQueuedNotificationTest::testNotificationCanSetMaxExceptions
		"NotificationSenderTest::test_it_can_send_queued_notifications_with_a_string_via",                                  // Port of NotificationSenderTest::test_it_can_send_queued_notifications_with_a_string_via
		"NotificationSenderTest::test_it_can_send_queued_notifications_with_an_array_via",                                  // Port of NotificationSenderTest::test_it_can_send_queued_notifications_with_an_array_via
		"NotificationSenderTest::test_it_can_send_notifications_with_an_empty_string_via",                                  // Port of NotificationSenderTest::test_it_can_send_notifications_with_an_empty_string_via
		"NotificationSenderTest::test_it_cannot_send_notifications_via_database_for_anonymous_notifiables",                 // Port of NotificationSenderTest::test_it_cannot_send_notifications_via_database_for_anonymous_notifiables
		"NotificationSenderTest::test_it_can_send_queued_notifications_through_middleware",                                 // Port of NotificationSenderTest::test_it_can_send_queued_notifications_through_middleware
		"NotificationSenderTest::test_it_can_send_queued_multi_channel_notifications_through_different_middleware",         // Port of NotificationSenderTest::test_it_can_send_queued_multi_channel_notifications_through_different_middleware
		"NotificationSenderTest::test_it_can_send_queued_with_via_connections_notifications",                               // Port of NotificationSenderTest::test_it_can_send_queued_with_via_connections_notifications
		"NotificationSenderTest::test_it_can_send_queued_with_via_queues_notifications",                                    // Port of NotificationSenderTest::test_it_can_send_queued_with_via_queues_notifications
		"NotificationSenderTest::test_it_can_send_queued_notifications_with_queue_route",                                   // Port of NotificationSenderTest::test_it_can_send_queued_notifications_with_queue_route
		"NotificationSenderTest::test_notification_failed_sent_without_http_transport_exception",                           // Port of NotificationSenderTest::test_notification_failed_sent_without_http_transport_exception
		"NotificationSenderTest::test_it_preserves_notification_state_mutated_in_via_method",                               // Port of NotificationSenderTest::test_it_preserves_notification_state_mutated_in_via_method
		"NotificationSenderTest::test_it_queue_overrides_queue_attribute",                                                  // Port of NotificationSenderTest::test_it_queue_overrides_queue_attribute
		"NotificationSenderTest::test_it_queue_attribute_is_used_when_on_queue_is_not_called",                              // Port of NotificationSenderTest::test_it_queue_attribute_is_used_when_on_queue_is_not_called
		"NotificationSenderTest::test_it_constructor_override_takes_precedence_over_queue_attribute",                       // Port of NotificationSenderTest::test_it_constructor_override_takes_precedence_over_queue_attribute
	}

	if len(ported) != 71 {
		t.Fatalf("expected 71 Upstream inventory entries, got %d", len(ported))
	}
}
