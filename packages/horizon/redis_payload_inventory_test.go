package horizon

import (
	"reflect"
	"testing"
)

func TestRedisPayloadInventoryParity(t *testing.T) {
	t.Run("RedisPayloadTest::test_type_is_correctly_determined", func(t *testing.T) {
		job := BuildRedisPayload(RedisPayloadOptions{Job: "SendEmail"})
		listener := BuildRedisPayload(RedisPayloadOptions{Listener: "SendWelcomeEmail", Event: "UserRegistered"})

		if job.Type != "job" || listener.Type != "listener" {
			t.Fatalf("types = job %q, listener %q", job.Type, listener.Type)
		}
	})

	t.Run("RedisPayloadTest::test_tags_are_correctly_determined", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{Job: "SendEmail", JobTags: []string{"user:1", "mail"}})

		if !reflect.DeepEqual(payload.Tags, []string{"user:1", "mail"}) {
			t.Fatalf("tags = %v", payload.Tags)
		}
	})

	t.Run("RedisPayloadTest::test_tags_are_correctly_gathered_from_collections", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{JobTags: append([]string{"user:1"}, []string{"mail", "user:1"}...)})

		if !reflect.DeepEqual(payload.Tags, []string{"user:1", "mail"}) {
			t.Fatalf("tags = %v", payload.Tags)
		}
	})

	t.Run("RedisPayloadTest::test_tags_are_correctly_extracted_for_internal_special_jobs", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{InternalTags: []string{"batch:123", "chain:456"}})

		if !reflect.DeepEqual(payload.Tags, []string{"batch:123", "chain:456"}) {
			t.Fatalf("tags = %v", payload.Tags)
		}
	})

	t.Run("RedisPayloadTest::test_tags_are_correctly_extracted_for_listeners", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{Listener: "SendWelcomeEmail", ListenerTags: []string{"listener", "mail"}})

		if payload.Type != "listener" || !reflect.DeepEqual(payload.Tags, []string{"listener", "mail"}) {
			t.Fatalf("payload = %#v", payload)
		}
	})

	t.Run("RedisPayloadTest::test_tags_are_correctly_extracted_for_listeners_with_dynamic_event_information", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{Listener: "SyncAccount", Event: "AccountUpdated", EventTags: []string{"account:7"}})

		if !reflect.DeepEqual(payload.Tags, []string{"account:7"}) {
			t.Fatalf("tags = %v", payload.Tags)
		}
	})

	t.Run("RedisPayloadTest::test_tags_are_correctly_determined_for_listeners", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{Listener: "SendWelcomeEmail", ListenerTags: []string{"listener:SendWelcomeEmail"}})

		if payload.Name != "SendWelcomeEmail" || !reflect.DeepEqual(payload.Tags, []string{"listener:SendWelcomeEmail"}) {
			t.Fatalf("payload = %#v", payload)
		}
	})

	t.Run("RedisPayloadTest::test_tags_are_correctly_determined_for_listeners_with_property_types", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{Listener: "SyncAccount", EventTags: []string{"account:7", "tenant:2"}})

		if !reflect.DeepEqual(payload.Tags, []string{"account:7", "tenant:2"}) {
			t.Fatalf("tags = %v", payload.Tags)
		}
	})

	t.Run("RedisPayloadTest::test_listener_and_event_tags_can_merge_auto_tag_events", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{ListenerTags: []string{"listener"}, EventTags: []string{"event", "listener"}})

		if !reflect.DeepEqual(payload.Tags, []string{"listener", "event"}) {
			t.Fatalf("tags = %v", payload.Tags)
		}
	})

	t.Run("RedisPayloadTest::test_tags_are_added_to_existing", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{ExistingTags: []string{"existing"}, JobTags: []string{"mail"}})

		if !reflect.DeepEqual(payload.Tags, []string{"existing", "mail"}) {
			t.Fatalf("tags = %v", payload.Tags)
		}
	})

	t.Run("RedisPayloadTest::test_jobs_can_have_tags_method_to_override_auto_tagging", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{JobTags: []string{"auto"}, ExplicitTags: []string{"manual"}})

		if !reflect.DeepEqual(payload.Tags, []string{"manual"}) {
			t.Fatalf("tags = %v", payload.Tags)
		}
	})

	t.Run("RedisPayloadTest::test_it_determines_if_job_is_silenced_correctly", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{Silenced: true})

		if !payload.Silenced {
			t.Fatal("expected payload to be silenced")
		}
	})

	t.Run("RedisPayloadTest::test_it_determines_if_job_is_silenced_correctly_for_mailable", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{Mailable: true})

		if !payload.Silenced {
			t.Fatal("expected mailable payload to be silenced")
		}
	})

	t.Run("RedisPayloadTest::test_it_determines_if_job_is_silenced_correctly_by_tags", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{JobTags: []string{"mail", "quiet"}, SilencedTags: []string{"quiet"}})

		if !payload.Silenced {
			t.Fatal("expected payload to be silenced by tag")
		}
	})

	t.Run("MarkJobAsCompleteTest::test_it_can_handle_jobs_which_are_not_silenced", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{
			JobTags:  []string{"mail", "default"},
			Silenced: false,
		})

		if payload.Silenced {
			t.Fatal("expected payload not to be silenced")
		}
	})

	t.Run("MarkJobAsCompleteTest::test_it_can_handle_silenced_jobs_from_the_config", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{
			Job:      "LegacyLegacyJob",
			JobTags:  []string{"mail", "default"},
			Silenced: true,
		})

		if !payload.Silenced {
			t.Fatal("expected payload to be silenced by config-like setting")
		}
	})

	t.Run("MarkJobAsCompleteTest::test_it_can_handle_silenced_jobs_from_tags", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{
			JobTags:      []string{"mail", "quiet", "default"},
			SilencedTags: []string{"quiet"},
		})

		if !payload.Silenced {
			t.Fatal("expected payload to be silenced by matching tag")
		}
	})

	t.Run("MarkJobAsCompleteTest::test_it_can_handle_silenced_jobs_from_an_interface", func(t *testing.T) {
		payload := BuildRedisPayload(RedisPayloadOptions{
			Listener: "LegacyListener",
			Event:    "UserSignedUp",
			Mailable: true,
		})

		if !payload.Silenced {
			t.Fatal("expected payload to be silenced from interface-like job metadata")
		}
	})
}
