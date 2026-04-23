package telescope_test

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/telescope"
	"github.com/bedrock/packages/telescope/storage"
	"github.com/bedrock/packages/telescope/watchers"
)

func TestTelescopeComplianceInventory(t *testing.T) {
	t.Run("ClearCommandTest::test_clear_command_will_delete_all_entries", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()

		if err := repo.Store([]*telescope.IncomingEntry{
			telescope.NewEntry(telescope.EntryTypeLog, map[string]any{"message": "one"}),
			telescope.NewEntry(telescope.EntryTypeQuery, map[string]any{"sql": "select 1"}),
		}); err != nil {
			t.Fatal(err)
		}

		if err := repo.Clear(); err != nil {
			t.Fatal(err)
		}

		if got := repo.Count(); got != 0 {
			t.Fatalf("repo count = %d, want 0", got)
		}
	})

	t.Run("PruneCommandTest::test_prune_command_will_clear_old_records", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		old := telescope.NewEntry(telescope.EntryTypeLog, map[string]any{"message": "old"})
		old.RecordedAt = time.Now().Add(-48 * time.Hour)
		recent := telescope.NewEntry(telescope.EntryTypeLog, map[string]any{"message": "recent"})

		if err := repo.Store([]*telescope.IncomingEntry{old, recent}); err != nil {
			t.Fatal(err)
		}

		pruned, err := repo.Prune(time.Now().Add(-24*time.Hour), false)

		if err != nil {
			t.Fatal(err)
		}

		if pruned != 1 || repo.Count() != 1 {
			t.Fatalf("pruned/count = %d/%d, want 1/1", pruned, repo.Count())
		}
	})

	t.Run("PruneCommandTest::test_prune_command_can_vary_hours", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		entry := telescope.NewEntry(telescope.EntryTypeException, map[string]any{"message": "keep"})
		entry.RecordedAt = time.Now().Add(-72 * time.Hour)

		if err := repo.Store([]*telescope.IncomingEntry{entry}); err != nil {
			t.Fatal(err)
		}

		pruned, err := repo.Prune(time.Now().Add(-24*time.Hour), true)

		if err != nil {
			t.Fatal(err)
		}

		if pruned != 0 || repo.Count() != 1 {
			t.Fatalf("pruned/count = %d/%d, want 0/1", pruned, repo.Count())
		}
	})

	t.Run("AvatarTest::test_it_can_generate_avatar_url", func(t *testing.T) {
		result := &telescope.EntryResult{Content: map[string]any{
			"user": map[string]any{"email": "taylor@example.com"},
		}}

		result.GenerateAvatar(nil)

		hash := fmt.Sprintf("%x", md5.Sum([]byte("taylor@example.com")))
		want := "https://www.gravatar.com/avatar/" + hash + "?s=200&d=mm"

		if result.Avatar != want {
			t.Fatalf("avatar = %q, want %q", result.Avatar, want)
		}
	})

	t.Run("AvatarTest::test_it_can_register_custom_avatar_path", func(t *testing.T) {
		result := &telescope.EntryResult{Content: map[string]any{
			"user": map[string]any{"email": "nuno@example.com"},
		}}

		result.GenerateAvatar(func(user map[string]any) string {
			return "/avatars/" + user["email"].(string)
		})

		if result.Avatar != "/avatars/nuno@example.com" {
			t.Fatalf("avatar = %q", result.Avatar)
		}
	})

	t.Run("AvatarTest::test_it_can_read_custom_avatar_path_on_null_email", func(t *testing.T) {
		result := &telescope.EntryResult{Content: map[string]any{
			"user": map[string]any{"name": "No Email"},
		}}

		result.GenerateAvatar(func(user map[string]any) string {
			return "/avatars/fallback"
		})

		if result.Avatar != "/avatars/fallback" {
			t.Fatalf("avatar = %q", result.Avatar)
		}
	})

	t.Run("DatabaseEntriesRepositoryTest::test_find_entry_by_uuid", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		entry := telescope.NewEntry(telescope.EntryTypeLog, map[string]any{"message": "stored"})
		entry.UUID = "entry-1"

		if err := repo.Store([]*telescope.IncomingEntry{entry}); err != nil {
			t.Fatal(err)
		}

		found, err := repo.Find("entry-1")

		if err != nil {
			t.Fatal(err)
		}

		if found.ID != "entry-1" || found.Content["message"] != "stored" {
			t.Fatalf("found entry = %#v", found)
		}
	})

	t.Run("DatabaseEntriesRepositoryTest::test_update", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		entry := telescope.NewEntry(telescope.EntryTypeJob, map[string]any{"status": "pending"})
		entry.UUID = "job-1"

		if err := repo.Store([]*telescope.IncomingEntry{entry}); err != nil {
			t.Fatal(err)
		}

		update := telescope.NewEntryUpdate("job-1", telescope.EntryTypeJob, map[string]any{"status": "processed"}).
			AddTags("processed")

		if err := repo.Update([]*telescope.EntryUpdate{update}); err != nil {
			t.Fatal(err)
		}

		found, err := repo.Find("job-1")

		if err != nil {
			t.Fatal(err)
		}

		if found.Content["status"] != "processed" {
			t.Fatalf("status = %v, want processed", found.Content["status"])
		}

		assertHasTag(t, found.Tags, "processed")
	})

	t.Run("DatabaseEntriesRepositoryTest::test_store_binary_content", func(t *testing.T) {
		repo := storage.NewInMemoryRepository()
		entry := telescope.NewEntry(telescope.EntryTypeRequest, map[string]any{"body": []byte{0, 1, 2}})
		entry.UUID = "binary-1"

		if err := repo.Store([]*telescope.IncomingEntry{entry}); err != nil {
			t.Fatal(err)
		}

		found, err := repo.Find("binary-1")

		if err != nil {
			t.Fatal(err)
		}

		if got := fmt.Sprint(found.Content["body"]); got != "[0 1 2]" {
			t.Fatalf("body = %v", found.Content["body"])
		}
	})

	t.Run("TelescopeTest::test_run_after_recording_callback", func(t *testing.T) {
		scope, _ := newTestScope(t)
		called := false

		scope.AfterRecording(func(entry *telescope.IncomingEntry) {
			called = entry.Type == telescope.EntryTypeLog
		})
		scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{"message": "logged"}))

		if !called {
			t.Fatal("after-recording callback was not called")
		}
	})

	t.Run("TelescopeTest::test_after_recording_callback_can_store_and_flush", func(t *testing.T) {
		scope, repo := newTestScope(t)
		scope.AfterRecording(func(_ *telescope.IncomingEntry) {
			if err := scope.Store(testContext()); err != nil {
				t.Fatalf("store from callback: %v", err)
			}

			scope.Flush()
		})

		scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{"message": "stored"}))

		if repo.Count() != 1 || len(scope.QueuedEntries()) != 0 {
			t.Fatalf("repo/queue = %d/%d, want 1/0", repo.Count(), len(scope.QueuedEntries()))
		}
	})

	t.Run("TelescopeTest::test_run_after_store_callback", func(t *testing.T) {
		scope, _ := newTestScope(t)

		var batch string

		var stored int

		scope.AfterStoring(func(batchID string, entries []*telescope.IncomingEntry) {
			batch = batchID
			stored = len(entries)
		})

		scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{"message": "stored"}))

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		if batch == "" || stored != 1 {
			t.Fatalf("after-store batch/stored = %q/%d", batch, stored)
		}
	})

	t.Run("CacheWatcherTest::test_cache_watcher_registers_missed_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewCacheWatcher(scope, nil)
		w.Missed("missing-key")
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

		if repo.Entries()[0].Content["type"] != "missed" {
			t.Fatalf("type = %v", repo.Entries()[0].Content["type"])
		}
	})

	t.Run("CacheWatcherTest::test_cache_watcher_registers_store_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewCacheWatcher(scope, nil)
		w.Written("cache-key", "value", 60)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

		if repo.Entries()[0].Content["expires"] != int64(60) {
			t.Fatalf("expires = %v", repo.Entries()[0].Content["expires"])
		}
	})

	t.Run("CacheWatcherTest::test_cache_watcher_registers_hit_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewCacheWatcher(scope, nil)
		w.Hit("cache-key", "value")
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

		if repo.Entries()[0].Content["type"] != "hit" {
			t.Fatalf("type = %v", repo.Entries()[0].Content["type"])
		}
	})

	t.Run("CacheWatcherTest::test_cache_watcher_registers_forget_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewCacheWatcher(scope, nil)
		w.Forgotten("cache-key")
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

		if repo.Entries()[0].Content["type"] != "forget" {
			t.Fatalf("type = %v", repo.Entries()[0].Content["type"])
		}
	})

	t.Run("CacheWatcherTest::test_cache_watcher_hides_hidden_values_when_set", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewCacheWatcher(scope, map[string]any{"hidden": []string{"secret"}})
		w.Written("secret", "value", 0)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

		if repo.Entries()[0].Content["value"] != "********" {
			t.Fatalf("value = %v", repo.Entries()[0].Content["value"])
		}
	})

	t.Run("CacheWatcherTest::test_cache_watcher_hides_hidden_values_when_retrieved", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewCacheWatcher(scope, map[string]any{"hidden": []string{"secret"}})
		w.Hit("secret", "value")
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

		if repo.Entries()[0].Content["value"] != "********" {
			t.Fatalf("value = %v", repo.Entries()[0].Content["value"])
		}
	})

	t.Run("CacheWatcherTest::test_cache_watcher_skips_recording_ignored_cache_keys", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewCacheWatcher(scope, nil)
		w.Hit("telescope:recording", true)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 0)
	})

	t.Run("EventWatcherTest::test_event_watcher_registers_any_events", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewEventWatcher(scope, nil)
		w.Record("App\\Events\\OrderShipped", nil, nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeEvent, 1)

		if repo.Entries()[0].Content["name"] != "App\\Events\\OrderShipped" {
			t.Fatalf("name = %v", repo.Entries()[0].Content["name"])
		}
	})

	t.Run("EventWatcherTest::test_event_watcher_stores_payloads", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewEventWatcher(scope, nil)
		w.Record("App\\Events\\OrderShipped", map[string]any{"order_id": 10}, nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeEvent, 1)
		payload := repo.Entries()[0].Content["payload"].(map[string]any)

		if payload["order_id"] != 10 {
			t.Fatalf("payload = %#v", payload)
		}
	})

	t.Run("EventWatcherTest::test_event_watcher_ignore_event", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewEventWatcher(scope, map[string]any{"ignore": []string{"App\\Events\\Ignored"}})
		w.Record("App\\Events\\Ignored", nil, nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeEvent, 0)
	})

	t.Run("ExceptionWatcherTest::test_exception_watcher_register_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewExceptionWatcher(scope, nil)
		w.Record(errors.New("boom"), 0)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeException, 1)

		if repo.Entries()[0].Content["message"] != "boom" {
			t.Fatalf("message = %v", repo.Entries()[0].Content["message"])
		}
	})

	t.Run("ExceptionWatcherTest::test_exception_watcher_register_throwable_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewExceptionWatcher(scope, nil)
		w.RecordRaw("App\\Exceptions\\Custom", "/app/main.go", 9, "custom", nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeException, 1)

		if repo.Entries()[0].Content["class"] != "App\\Exceptions\\Custom" {
			t.Fatalf("class = %v", repo.Entries()[0].Content["class"])
		}
	})

	t.Run("JobWatcherTest::test_job_registers_entry", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewJobWatcher(scope, nil)
		w.Pending(watchers.JobMeta{UUID: "job-1", Type: "SendMail", Connection: "redis", Queue: "mail"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeJob, 1)

		if repo.Entries()[0].Content["name"] != "SendMail" {
			t.Fatalf("name = %v", repo.Entries()[0].Content["name"])
		}
	})

	t.Run("JobWatcherTest::test_job_registers_entry_with_batchId_in_payload", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewJobWatcher(scope, nil)
		w.Pending(watchers.JobMeta{UUID: "job-1", Type: "SendMail", BatchID: "batch-1"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeJob, 1)

		if repo.Entries()[0].BatchID != scope.CurrentBatchID() {
			t.Fatalf("recorded batch = %q", repo.Entries()[0].BatchID)
		}
	})

	t.Run("JobWatcherTest::test_failed_jobs_register_entry", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewJobWatcher(scope, nil)
		w.Pending(watchers.JobMeta{UUID: "job-1", Type: "Import"})

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		w.Failed("job-1", "Import", errors.New("failed"))

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		found, err := repo.Find("job-1")

		if err != nil {
			t.Fatal(err)
		}

		if found.Content["status"] != watchers.JobStatusFailed {
			t.Fatalf("status = %v", found.Content["status"])
		}
	})

	t.Run("QueryWatcherTest::test_query_watcher_registers_database_queries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewQueryWatcher(scope, nil)
		w.Record("select * from users where id = ?", []any{1}, time.Millisecond, "mysql", "", 0)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeQuery, 1)

		if repo.Entries()[0].Content["connection"] != "mysql" {
			t.Fatalf("connection = %v", repo.Entries()[0].Content["connection"])
		}
	})

	t.Run("QueryWatcherTest::test_query_watcher_can_tag_slow_queries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewQueryWatcher(scope, map[string]any{"slow": float64(10)})
		w.Record("select sleep(1)", nil, 20*time.Millisecond, "mysql", "", 0)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeQuery, 1)
		assertHasTag(t, repo.Entries()[0].Tags, "slow")
	})

	t.Run("QueryWatcherTest::test_query_watcher_can_prepare_bindings", func(t *testing.T) {
		got := watchers.ReplaceBindings("select * from users where id = ? and active = ?", []any{7, true})
		want := "select * from users where id = 7 and active = 1"

		if got != want {
			t.Fatalf("sql = %q, want %q", got, want)
		}
	})

	t.Run("QueryWatcherTest::test_query_watcher_can_prepare_named_bindings", func(t *testing.T) {
		got := watchers.ReplaceNamedBindings("select * from users where id = :id", map[string]any{"id": 7})
		want := "select * from users where id = 7"

		if got != want {
			t.Fatalf("sql = %q, want %q", got, want)
		}
	})

	t.Run("QueryWatcherTest::test_query_watcher_can_prepare_bindings_for_nonstandard_connections", func(t *testing.T) {
		got := watchers.ReplaceBindings("select * from `users` where uuid = ?", []any{"abc-123"})
		want := "select * from `users` where uuid = 'abc-123'"

		if got != want {
			t.Fatalf("sql = %q, want %q", got, want)
		}
	})

	t.Run("RequestWatchersTest::test_request_watcher_registers_requests", func(t *testing.T) {
		scope, repo, w := makeRequestTestScope(t, nil)
		doTestRequest(t, w, http.MethodGet, "/users", "", nil)

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		if repo.Count() != 1 {
			t.Fatalf("count = %d, want 1", repo.Count())
		}
	})

	t.Run("RequestWatchersTest::test_request_watcher_hides_password", func(t *testing.T) {
		scope, repo, w := makeRequestTestScope(t, nil)
		doTestRequest(t, w, http.MethodPost, "/login", `{"password":"secret"}`, map[string]string{"Content-Type": "application/json"})

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		payload := repo.Entries()[0].Content["payload"].(map[string]any)

		if payload["password"] != "********" {
			t.Fatalf("password = %v", payload["password"])
		}
	})

	t.Run("RequestWatchersTest::test_request_watcher_hides_authorization", func(t *testing.T) {
		scope, repo, w := makeRequestTestScope(t, nil)
		doTestRequest(t, w, http.MethodGet, "/users", "", map[string]string{"Authorization": "Bearer token"})

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		headers := repo.Entries()[0].Content["headers"].(map[string]string)

		if headers["Authorization"] != "********" {
			t.Fatalf("authorization = %q", headers["Authorization"])
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_registers_succesful_client_request_and_response", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodGet, "https://example.test/users", http.Header{"Accept": []string{"application/json"}}, nil, &watchers.ClientResponse{
			StatusCode:  http.StatusOK,
			Headers:     http.Header{"Content-Type": []string{"application/json"}},
			Body:        []byte(`{"ok":true}`),
			ContentType: "application/json",
		}, 15*time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)

		if repo.Entries()[0].Content["response_status"] != http.StatusOK {
			t.Fatalf("status = %v", repo.Entries()[0].Content["response_status"])
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_registers_redirect_response", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodGet, "https://example.test/login", nil, nil, &watchers.ClientResponse{
			StatusCode: http.StatusFound,
			Headers:    http.Header{"Location": []string{"/dashboard"}},
		}, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)

		if repo.Entries()[0].Content["response_status"] != http.StatusFound {
			t.Fatalf("status = %v", repo.Entries()[0].Content["response_status"])
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_plain_text_response", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodGet, "https://example.test/plain", nil, nil, &watchers.ClientResponse{
			StatusCode:  http.StatusOK,
			Body:        []byte("plain response"),
			ContentType: "text/plain",
		}, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)

		if repo.Entries()[0].Content["response"] != "plain response" {
			t.Fatalf("response = %v", repo.Entries()[0].Content["response"])
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_registers_server_error_response", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodPost, "https://example.test/fail", nil, nil, &watchers.ClientResponse{
			StatusCode:  http.StatusInternalServerError,
			Body:        []byte(`{"error":"failed"}`),
			ContentType: "application/json",
		}, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)

		if repo.Entries()[0].Content["response_status"] != http.StatusInternalServerError {
			t.Fatalf("status = %v", repo.Entries()[0].Content["response_status"])
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_hides_password", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodPost, "https://example.test/login", http.Header{"Content-Type": []string{"application/json"}}, []byte(`{"password":"secret"}`), nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		payload := repo.Entries()[0].Content["payload"].(map[string]any)

		if payload["password"] != "********" {
			t.Fatalf("client request payload = %#v", payload)
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_hides_authorization", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodGet, "https://example.test/users", http.Header{"Authorization": []string{"Bearer token"}}, nil, nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		headers := repo.Entries()[0].Content["headers"].(map[string]string)

		if headers["Authorization"] != "********" {
			t.Fatalf("authorization = %q", headers["Authorization"])
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_hides_php_auth_pw", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodGet, "https://example.test/users", http.Header{"Php-Auth-Pw": []string{"secret"}}, nil, nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		headers := repo.Entries()[0].Content["headers"].(map[string]string)

		if headers["Php-Auth-Pw"] != "********" {
			t.Fatalf("php auth pw = %q", headers["Php-Auth-Pw"])
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_handles_form_request", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodPost, "https://example.test/form", http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}, []byte("name=taylor"), nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		payload := repo.Entries()[0].Content["payload"].(map[string]any)

		if payload["name"] != "taylor" {
			t.Fatalf("payload = %v", repo.Entries()[0].Content["payload"])
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_handles_multipart_request", func(t *testing.T) {
		body, contentType := multipartBody(t, "avatar", "avatar.txt", "hello", map[string]string{"name": "taylor"})
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodPost, "https://example.test/upload", http.Header{"Content-Type": []string{contentType}}, body, nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		payload := repo.Entries()[0].Content["payload"].(map[string]any)

		if payload["name"] != "taylor" || payload["avatar"].(map[string]any)["name"] != "avatar.txt" {
			t.Fatalf("payload = %#v", payload)
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_handles_file_contents_upload", func(t *testing.T) {
		body, contentType := multipartBody(t, "document", "report.txt", "contents", nil)
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodPost, "https://example.test/upload", http.Header{"Content-Type": []string{contentType}}, body, nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		file := repo.Entries()[0].Content["payload"].(map[string]any)["document"].(map[string]any)

		if file["size"] != int64(len("contents")) {
			t.Fatalf("file = %#v", file)
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_handles_file_contents_upload_without_explicit_filename_or_headers", func(t *testing.T) {
		body, contentType := multipartBody(t, "document", "", "contents", nil)
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodPost, "https://example.test/upload", http.Header{"Content-Type": []string{contentType}}, body, nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		document := repo.Entries()[0].Content["payload"].(map[string]any)["document"]

		if document != "contents" {
			t.Fatalf("document = %#v", document)
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_handles_resource_file_upload", func(t *testing.T) {
		body, contentType := multipartBody(t, "resource", "resource.bin", "binary", nil)
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodPost, "https://example.test/upload", http.Header{"Content-Type": []string{contentType}}, body, nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		file := repo.Entries()[0].Content["payload"].(map[string]any)["resource"].(map[string]any)

		if file["name"] != "resource.bin" {
			t.Fatalf("file = %#v", file)
		}
	})

	t.Run("ClientRequestWatcherTest::test_client_request_watcher_handles_resource_file_upload_with_filename_and_headers", func(t *testing.T) {
		body, contentType := multipartBody(t, "resource", "named.bin", "binary", nil)
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodPost, "https://example.test/upload", http.Header{"Content-Type": []string{contentType}, "X-Upload": []string{"1"}}, body, nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		content := repo.Entries()[0].Content
		file := content["payload"].(map[string]any)["resource"].(map[string]any)
		headers := content["headers"].(map[string]string)

		if file["name"] != "named.bin" || headers["X-Upload"] != "1" {
			t.Fatalf("content = %#v", content)
		}
	})

	t.Run("ClientRequestWatcherTest::test_it_stores_and_displays_array_of_request_headers", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewClientRequestWatcher(scope, nil)
		w.Record(http.MethodGet, "https://example.test/users", http.Header{"X-Test": []string{"one", "two"}}, nil, nil, time.Millisecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeClientRequest, 1)
		headers := repo.Entries()[0].Content["headers"].(map[string]string)

		if headers["X-Test"] != "one" {
			t.Fatalf("headers = %#v", headers)
		}
	})

	t.Run("CommandWatcherTest::test_command_watcher_register_entry", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewCommandWatcher(scope, nil)
		w.Record("queue:work", 0, map[string]any{"connection": "redis"}, map[string]any{"once": true})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeCommand, 1)

		if repo.Entries()[0].Content["command"] != "queue:work" {
			t.Fatalf("command = %v", repo.Entries()[0].Content["command"])
		}
	})

	t.Run("DumpWatcherTest::test_dump_watcher_register_entry", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewDumpWatcher(scope, nil)
		w.Record(map[string]any{"name": "Taylor"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeDump, 1)

		if repo.Entries()[0].Content["dump"] == "" {
			t.Fatal("expected dump content")
		}
	})

	t.Run("EventWatcherTest::test_format_listeners", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewEventWatcher(scope, nil)
		w.Record("App\\Events\\OrderShipped", nil, []string{"SendShipmentNotification"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeEvent, 1)
		listeners := repo.Entries()[0].Content["listeners"].([]string)

		if len(listeners) != 1 || listeners[0] != "SendShipmentNotification" {
			t.Fatalf("listeners = %#v", listeners)
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_allowed_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.Record("viewDashboard", true, "user-1", []any{"dashboard"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["result"] != watchers.GateResultAllowed {
			t.Fatalf("result = %v", repo.Entries()[0].Content["result"])
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_denied_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.Record("deletePost", false, "user-1", []any{"post-1"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["result"] != watchers.GateResultDenied {
			t.Fatalf("result = %v", repo.Entries()[0].Content["result"])
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_allowed_guest_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.Record("viewPublic", true, nil, nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if _, ok := repo.Entries()[0].Content["user"]; ok {
			t.Fatalf("guest gate stored user: %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_denied_guest_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.Record("viewPrivate", false, nil, nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["result"] != watchers.GateResultDenied {
			t.Fatalf("result = %v", repo.Entries()[0].Content["result"])
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_allowed_entries_with_message", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.RecordWithMessage("viewDashboard", true, "allowed", "user-1", nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["message"] != "allowed" {
			t.Fatalf("content = %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_denied_entries_with_message", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.RecordWithMessage("deletePost", false, "denied", "user-1", nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["message"] != "denied" {
			t.Fatalf("content = %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_allowed_guest_entries_with_message", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.RecordWithMessage("viewPublic", true, "guest allowed", nil, nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["message"] != "guest allowed" {
			t.Fatalf("content = %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_allowed_policy_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.Record("PostPolicy@view", true, "user-1", []any{"post-1"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["ability"] != "PostPolicy@view" {
			t.Fatalf("content = %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_after_checks", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.Record("after:view", true, "user-1", nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["ability"] != "after:view" {
			t.Fatalf("content = %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_denied_policy_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.Record("PostPolicy@delete", false, "user-1", []any{"post-1"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["result"] != watchers.GateResultDenied {
			t.Fatalf("content = %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_allowed_response_policy_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.RecordWithMessage("PostPolicy@update", true, "policy allowed", "user-1", []any{"post-1"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["message"] != "policy allowed" {
			t.Fatalf("content = %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("GateWatcherTest::test_gate_watcher_registers_denied_response_policy_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewGateWatcher(scope, nil)
		w.RecordWithMessage("PostPolicy@update", false, "policy denied", "user-1", []any{"post-1"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeGate, 1)

		if repo.Entries()[0].Content["message"] != "policy denied" {
			t.Fatalf("content = %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("BatchWatcherTest::test_job_dispatch_registers_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewBatchWatcher(scope, nil)
		w.Record(watchers.BatchDispatch{ID: "batch-1", Name: "Import users", Total: 3})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeBatch, 1)

		if repo.Entries()[0].Content["total"] != 3 {
			t.Fatalf("content = %#v", repo.Entries()[0].Content)
		}
	})

	t.Run("JobWatcherTest::test_it_handles_pushed_jobs", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewJobWatcher(scope, nil)
		w.Pending(watchers.JobMeta{UUID: "job-2", Type: "ImportUsers", Connection: "database", Queue: "imports", MaxTries: 3, Timeout: 60})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeJob, 1)
		content := repo.Entries()[0].Content

		if content["connection"] != "database" || content["queue"] != "imports" {
			t.Fatalf("job content = %#v", content)
		}
	})

	t.Run("LogWatcherTest::test_log_watcher_registers_entry_for_any_level_by_default", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewLogWatcher(scope, nil)
		w.Record("debug", "debugging", nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 1)
	})

	t.Run("LogWatcherTest::test_log_watcher_only_registers_entries_for_the_specified_error_level_priority", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewLogWatcher(scope, map[string]any{"level": "error"})
		w.Record("warning", "skip", nil)
		w.Record("error", "record", nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 1)
	})

	t.Run("LogWatcherTest::test_log_watcher_only_registers_entries_for_the_specified_debug_level_priority", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewLogWatcher(scope, map[string]any{"level": "debug"})
		w.Record("debug", "record", nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 1)
	})

	t.Run("LogWatcherTest::test_log_watcher_registers_entry_with_exception_key", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewLogWatcher(scope, nil)
		w.Record("error", "with exception", map[string]any{"exception": errors.New("boom")})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 0)
	})

	t.Run("LogWatcherTest::test_log_watcher_do_not_registers_entry_when_disabled_on_the_boolean_format", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewLogWatcher(scope, map[string]any{"enabled": false})
		w.Record("error", "disabled", nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 0)
	})

	t.Run("LogWatcherTest::test_log_watcher_do_not_registers_entry_when_disabled_on_the_array_format", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewLogWatcher(scope, map[string]any{"enabled": map[string]any{"log": false}})
		w.Record("error", "disabled", nil)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 0)
	})

	t.Run("LogWatcherTest::test_log_watcher_interpolates_message", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewLogWatcher(scope, nil)
		w.Record("info", "Hello {name}", map[string]any{"name": "Taylor"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 1)

		if repo.Entries()[0].Content["message"] != "Hello Taylor" {
			t.Fatalf("message = %v", repo.Entries()[0].Content["message"])
		}
	})

	t.Run("MailWatcherTest::test_mail_watcher_registers_entry", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewMailWatcher(scope, nil)
		w.Record(watchers.MailMessage{Mailable: "App\\Mail\\Welcome", Subject: "Welcome", To: []string{"taylor@example.com"}})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeMail, 1)

		if repo.Entries()[0].Content["subject"] != "Welcome" {
			t.Fatalf("subject = %v", repo.Entries()[0].Content["subject"])
		}
	})

	t.Run("ModelWatcherTest::test_model_watcher_registers_entry", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewModelWatcher(scope, nil)
		w.Record("App\\Models\\User", watchers.ModelActionCreated, 1)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeModel, 1)

		if repo.Entries()[0].Content["action"] != watchers.ModelActionCreated {
			t.Fatalf("action = %v", repo.Entries()[0].Content["action"])
		}
	})

	t.Run("ModelWatcherTest::test_model_watcher_can_restrict_events", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewModelWatcher(scope, map[string]any{"ignore": []string{"App\\Models\\Audit"}})
		w.Record("App\\Models\\Audit", watchers.ModelActionCreated, 1)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeModel, 0)
	})

	t.Run("ModelWatcherTest::test_model_watcher_registers_hydration_entry", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewModelWatcher(scope, map[string]any{"hydrations": true})
		w.Record("App\\Models\\User", watchers.ModelActionRetrieved, 1)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeModel, 1)
	})

	t.Run("NotificationWatcherTest::test_notification_watcher_registers_entry", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewNotificationWatcher(scope, nil)
		w.Record("App\\Notifications\\InvoicePaid", "mail", "user-1", true)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeNotification, 1)

		if repo.Entries()[0].Content["channel"] != "mail" {
			t.Fatalf("channel = %v", repo.Entries()[0].Content["channel"])
		}
	})

	t.Run("NotificationWatcherTest::test_notification_watcher_registers_array_routes", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewNotificationWatcher(scope, nil)
		w.Record("App\\Notifications\\InvoicePaid", "mail", []string{"a@example.com", "b@example.com"}, false)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeNotification, 1)

		if repo.Entries()[0].Content["notifiable"] == "" {
			t.Fatalf("notifiable = %v", repo.Entries()[0].Content["notifiable"])
		}
	})

	t.Run("RedisWatcherTest::test_redis_watcher_registers_entries", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewRedisWatcher(scope, nil)
		w.Record("SET key value", "default", 250*time.Microsecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeRedis, 1)

		if repo.Entries()[0].Content["connection"] != "default" {
			t.Fatalf("connection = %v", repo.Entries()[0].Content["connection"])
		}
	})

	t.Run("RedisWatcherTest::test_does_not_register_when_redis_unbound", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewRedisWatcher(scope, nil)
		w.Record("MULTI", "default", time.Microsecond)
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeRedis, 0)
	})

	t.Run("RequestWatchersTest::test_request_watcher_registers_404", func(t *testing.T) {
		scope, repo, w := makeRequestTestScope(t, nil)
		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		rr := httptest.NewRecorder()
		w.Middleware(http.NotFoundHandler()).ServeHTTP(rr, req)

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		if repo.Entries()[0].Content["response_status"] != http.StatusNotFound {
			t.Fatalf("status = %v", repo.Entries()[0].Content["response_status"])
		}
	})

	t.Run("RequestWatchersTest::test_request_watcher_hides_php_auth_pw", func(t *testing.T) {
		scope, repo, w := makeRequestTestScope(t, nil)
		doTestRequest(t, w, http.MethodGet, "/users", "", map[string]string{"Php-Auth-Pw": "secret"})

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		headers := repo.Entries()[0].Content["headers"].(map[string]string)

		if headers["Php-Auth-Pw"] != "********" {
			t.Fatalf("php auth pw = %q", headers["Php-Auth-Pw"])
		}
	})

	t.Run("RequestWatchersTest::test_it_stores_and_displays_array_of_request_and_response_headers", func(t *testing.T) {
		scope, repo, w := makeRequestTestScope(t, nil)
		req := httptest.NewRequest(http.MethodGet, "/headers", nil)
		req.Header.Add("X-Test", "one")
		req.Header.Add("X-Test", "two")
		rr := httptest.NewRecorder()
		w.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.Header().Add("X-Response", "first")
			rw.Header().Add("X-Response", "second")
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte("ok"))
		})).ServeHTTP(rr, req)

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		content := repo.Entries()[0].Content

		if content["headers"].(map[string]string)["X-Test"] != "one" {
			t.Fatalf("headers = %#v", content["headers"])
		}

		if content["response_headers"].(map[string]string)["X-Response"] != "first" {
			t.Fatalf("response headers = %#v", content["response_headers"])
		}
	})

	t.Run("RequestWatchersTest::test_request_watcher_plain_text_response", func(t *testing.T) {
		scope, repo, w := makeRequestTestScope(t, nil)
		req := httptest.NewRequest(http.MethodGet, "/plain", nil)
		rr := httptest.NewRecorder()
		w.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.Header().Set("Content-Type", "text/plain")
			_, _ = rw.Write([]byte("plain"))
		})).ServeHTTP(rr, req)

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		if repo.Entries()[0].Content["response"] != "plain" {
			t.Fatalf("response = %v", repo.Entries()[0].Content["response"])
		}
	})

	t.Run("RequestWatchersTest::test_request_watcher_records_plain_text_payload", func(t *testing.T) {
		scope, repo, w := makeRequestTestScope(t, nil)
		doTestRequest(t, w, http.MethodPost, "/plain", "plain body", map[string]string{"Content-Type": "text/plain"})

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		payload := repo.Entries()[0].Content["payload"].(map[string]any)

		if payload["raw"] != "plain body" {
			t.Fatalf("payload = %#v", payload)
		}
	})

	t.Run("RequestWatchersTest::test_request_watcher_handles_file_uploads", func(t *testing.T) {
		body, contentType := multipartBody(t, "avatar", "avatar.txt", "hello", map[string]string{"name": "taylor"})
		scope, repo, w := makeRequestTestScope(t, nil)
		req := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader(body))
		req.Header.Set("Content-Type", contentType)
		rr := httptest.NewRecorder()
		w.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			_, _ = rw.Write([]byte("ok"))
		})).ServeHTTP(rr, req)

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		payload := repo.Entries()[0].Content["payload"].(map[string]any)

		if payload["name"] != "taylor" || payload["avatar"].(map[string]any)["name"] != "avatar.txt" {
			t.Fatalf("payload = %#v", payload)
		}
	})

	t.Run("RequestWatchersTest::test_request_watcher_handles_unlinked_file_uploads", func(t *testing.T) {
		body, contentType := multipartBody(t, "avatar", "", "hello", nil)
		scope, repo, w := makeRequestTestScope(t, nil)
		req := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader(body))
		req.Header.Set("Content-Type", contentType)
		rr := httptest.NewRecorder()
		w.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			_, _ = rw.Write([]byte("ok"))
		})).ServeHTTP(rr, req)

		if err := scope.Store(testContext()); err != nil {
			t.Fatal(err)
		}

		avatar := repo.Entries()[0].Content["payload"].(map[string]any)["avatar"]

		if avatar != "hello" {
			t.Fatalf("avatar = %#v", avatar)
		}
	})

	t.Run("ViewWatcherTest::test_view_watcher_registers_views", func(t *testing.T) {
		scope, repo := newTestScope(t)
		w := watchers.NewViewWatcher(scope, nil)
		w.Record("users.index", []string{"users"}, []string{"App\\View\\Composers\\UsersComposer"})
		storeAndAssertCount(t, scope, repo, telescope.EntryTypeView, 1)

		if repo.Entries()[0].Content["name"] != "users.index" {
			t.Fatalf("name = %v", repo.Entries()[0].Content["name"])
		}
	})
}

func multipartBody(t *testing.T, field string, filename string, contents string, values map[string]string) ([]byte, string) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for key, value := range values {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("write field: %v", err)
		}
	}

	part, err := writer.CreateFormFile(field, filename)

	if err != nil {
		t.Fatalf("create file: %v", err)
	}

	if _, err := part.Write([]byte(contents)); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	return body.Bytes(), writer.FormDataContentType()
}
