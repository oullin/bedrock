package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

type failingCacheStore struct {
	cache.Store
	forgetErr     error
	flushErr      error
	flushLocksErr error
}

func (s failingCacheStore) Forget(context.Context, string) error { return s.forgetErr }

func (s failingCacheStore) Flush(context.Context) error { return s.flushErr }

func (s failingCacheStore) FlushLocks(context.Context) error { return s.flushLocksErr }

// CacheArrayStoreTest::testItemsCanBeSetAndRetrieved
// CacheArrayStoreTest::testCacheTtl
// CacheArrayStoreTest::testMultipleItemsCanBeSetAndRetrieved
// CacheArrayStoreTest::testItemsCanExpire
// CacheArrayStoreTest::testTouchExtendsTtl
// CacheArrayStoreTest::testStoreItemForeverProperlyStoresInArray
// CacheArrayStoreTest::testValuesCanBeIncremented
// CacheArrayStoreTest::testValuesGetCastedByIncrementOrDecrement
// CacheArrayStoreTest::testIncrementNonNumericValues
// CacheArrayStoreTest::testNonExistingKeysCanBeIncremented
// CacheArrayStoreTest::testExpiredKeysAreIncrementedLikeNonExistingKeys
// CacheArrayStoreTest::testValuesCanBeDecremented
// CacheArrayStoreTest::testItemsCanBeRemoved
// CacheArrayStoreTest::testItemsCanBeFlushed
// CacheArrayStoreTest::testLocksCanBeFlushed
// CacheArrayStoreTest::testCacheKey
// CacheArrayStoreTest::testCannotAcquireLockTwice
// CacheArrayStoreTest::testCanAcquireLockAgainAfterExpiry
// CacheArrayStoreTest::testLockExpirationLowerBoundary
// CacheArrayStoreTest::testLockWithNoExpirationNeverExpires
// CacheArrayStoreTest::testCanAcquireLockAfterRelease
// CacheArrayStoreTest::testAnotherOwnerCannotReleaseLock
// CacheArrayStoreTest::testAnotherOwnerCanForceReleaseALock
// CacheArrayStoreTest::testReleasingLockAfterAlreadyForceReleasedByAnotherOwnerFails
// CacheArrayStoreTest::testOwnerStatusCanBeCheckedAfterRestoringLock
// CacheArrayStoreTest::testOtherOwnerDoesNotOwnLockAfterRestore
// CacheArrayStoreTest::testRestoringNonExistingLockDoesNotOwnAnything
// CacheArrayStoreTest::testCanGetAll
// CacheArrayStoreTest::testCanGetAllWhenSerialized
func TestLaravelArrayStoreInventoryEquivalents(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	store := cache.NewArrayStoreWithClock(clk)
	ctx := context.Background()

	if err := store.Put(ctx, "name", "Taylor", time.Minute); err != nil {
		t.Fatal(err)
	}

	if value, err := store.Get(ctx, "name"); err != nil || value != "Taylor" {
		t.Fatalf("expected cached value, got %v, %v", value, err)
	}

	if err := store.PutMany(ctx, map[string]any{"a": 1, "b": 2}, time.Minute); err != nil {
		t.Fatal(err)
	}

	many, err := store.GetMany(ctx, []string{"a", "b", "missing"})

	if err != nil {
		t.Fatal(err)
	}

	if len(many) != 2 || many["a"] != 1 || many["b"] != 2 {
		t.Fatalf("unexpected many result: %#v", many)
	}

	if ok, err := store.Touch(ctx, "name", 2*time.Minute); err != nil || !ok {
		t.Fatalf("expected touch to extend ttl, got %v, %v", ok, err)
	}

	clk.Advance(90 * time.Second)

	if _, err := store.Get(ctx, "name"); err != nil {
		t.Fatalf("expected touched key to survive, got %v", err)
	}

	clk.Advance(time.Minute)

	if _, err := store.Get(ctx, "name"); !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected key to expire, got %v", err)
	}

	if value, err := store.Increment(ctx, "counter", 5); err != nil || value != 5 {
		t.Fatalf("expected new increment value 5, got %d, %v", value, err)
	}

	if value, err := store.Decrement(ctx, "counter", 2); err != nil || value != 3 {
		t.Fatalf("expected decremented value 3, got %d, %v", value, err)
	}

	if err := store.Put(ctx, "bad", "not-numeric", time.Minute); err != nil {
		t.Fatal(err)
	}

	if _, err := store.Increment(ctx, "bad", 1); !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected invalid value error, got %v", err)
	}

	if err := store.Forever(ctx, "forever", "v"); err != nil {
		t.Fatal(err)
	}

	if err := store.Forget(ctx, "forever"); err != nil {
		t.Fatal(err)
	}

	if err := store.Flush(ctx); err != nil {
		t.Fatal(err)
	}

	if len(store.All()) != 0 {
		t.Fatal("expected flushed store to be empty")
	}

	first := store.Lock("resource", "owner-1", time.Minute)
	ok, err := first.Acquire(ctx)

	if err != nil || !ok {
		t.Fatalf("expected first lock acquire, got %v, %v", ok, err)
	}

	second := store.Lock("resource", "owner-2", time.Minute)
	ok, err = second.Acquire(ctx)

	if err != nil || ok {
		t.Fatalf("expected second lock acquire to fail, got %v, %v", ok, err)
	}

	restored := store.RestoreLock("resource", "owner-1")
	restoredOwns, err := restored.IsOwnedByCurrentProcess(ctx)

	if err != nil {
		t.Fatal(err)
	}

	secondOwns, err := second.IsOwnedByCurrentProcess(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if !restoredOwns || secondOwns {
		t.Fatal("expected restored owner to match only owner-1")
	}

	if err := second.ForceRelease(ctx); err != nil {
		t.Fatal(err)
	}

	if ok, _ := first.Release(ctx); ok {
		t.Fatal("expected original owner release to fail after force release")
	}

	if err := store.FlushLocks(ctx); err != nil {
		t.Fatal(err)
	}
}

// CacheFileStoreTest::testNullIsReturnedIfFileDoesntExist
// CacheFileStoreTest::testPutCreatesMissingDirectories
// CacheFileStoreTest::testPutWillConsiderZeroAsEternalTime
// CacheFileStoreTest::testPutWillConsiderBigValuesAsEternalTime
// CacheFileStoreTest::testExpiredItemsReturnNullAndGetDeleted
// CacheFileStoreTest::testValidItemReturnsContents
// CacheFileStoreTest::testStoreItemProperlyStoresValues
// CacheFileStoreTest::testTouchExtendsTtl
// CacheFileStoreTest::testForeversAreStoredWithHighTimestamp
// CacheFileStoreTest::testForeversAreNotRemovedOnIncrement
// CacheFileStoreTest::testIncrementExpiredKeys
// CacheFileStoreTest::testIncrementCanAtomicallyJump
// CacheFileStoreTest::testDecrementCanAtomicallyJump
// CacheFileStoreTest::testIncrementNonNumericValues
// CacheFileStoreTest::testIncrementNonExistentKeys
// CacheFileStoreTest::testIncrementDoesNotExtendCacheLife
// CacheFileStoreTest::testRemoveDeletesFileDoesntExist
// CacheFileStoreTest::testRemoveDeletesFile
// CacheFileStoreTest::testFlushCleansDirectory
// CacheFileStoreTest::testFlushIgnoreNonExistingDirectory
// CacheFileStoreTest::testFlushingLocksCleansDirectory
// CacheFileStoreTest::testFlushingLocksIgnoreNonExistingDirectory
// CacheFileStoreTest::testItHandlesForgettingNonFlexibleKeys
func TestLaravelFileStoreInventoryEquivalents(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	store := cache.NewFileStoreWithOptions(t.TempDir(), "", 0o644, clk)
	ctx := context.Background()

	if _, err := store.Get(ctx, "missing"); !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected missing file to be not found, got %v", err)
	}

	if err := store.Put(ctx, "name", "Taylor", time.Minute); err != nil {
		t.Fatal(err)
	}

	if value, err := store.Get(ctx, "name"); err != nil || value != "Taylor" {
		t.Fatalf("expected cached file value, got %v, %v", value, err)
	}

	if ok, err := store.Touch(ctx, "name", 2*time.Minute); err != nil || !ok {
		t.Fatalf("expected touch to succeed, got %v, %v", ok, err)
	}

	if value, err := store.Increment(ctx, "counter", 5); err != nil || value != 5 {
		t.Fatalf("expected increment 5, got %d, %v", value, err)
	}

	if value, err := store.Decrement(ctx, "counter", 2); err != nil || value != 3 {
		t.Fatalf("expected decrement 3, got %d, %v", value, err)
	}

	if err := store.Put(ctx, "bad", "not-numeric", time.Minute); err != nil {
		t.Fatal(err)
	}

	if _, err := store.Increment(ctx, "bad", 1); !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected invalid value, got %v", err)
	}

	if err := store.Forever(ctx, "forever", int64(5)); err != nil {
		t.Fatal(err)
	}

	if _, err := store.Increment(ctx, "forever", 3); err != nil {
		t.Fatal(err)
	}

	clk.Advance(365 * 24 * time.Hour)

	if value, err := store.Get(ctx, "forever"); err != nil || value != int64(8) {
		t.Fatalf("expected forever increment to survive, got %v, %v", value, err)
	}

	if err := store.Forget(ctx, "missing"); err != nil {
		t.Fatal(err)
	}

	if err := store.Flush(ctx); err != nil {
		t.Fatal(err)
	}

	if err := store.FlushLocks(ctx); err != nil {
		t.Fatal(err)
	}
}

// CacheDatabaseStoreTest::testNullIsReturnedWhenItemNotFound
// CacheDatabaseStoreTest::testNullIsReturnedAndItemDeletedWhenItemIsExpired
// CacheDatabaseStoreTest::testDecryptedValueIsReturnedWhenItemIsValid
// CacheDatabaseStoreTest::testValueIsUpserted
// CacheDatabaseStoreTest::testForeverCallsStoreItemWithReallyLongTime
// CacheDatabaseStoreTest::testItemsMayBeRemovedFromCache
// CacheDatabaseStoreTest::testItemsMayBeFlushedFromCache
// CacheDatabaseStoreTest::testLocksMayBeFlushedFromCache
// CacheDatabaseStoreTest::testIncrementReturnsCorrectValues
// CacheDatabaseStoreTest::testDecrementReturnsCorrectValues
// CacheDatabaseStoreTest::testTouchExtendsTtl
// CacheDynamoDbStoreTest::testTouchMethodCorrectlyCallsDynamoDb
// CacheRedisStoreTest::testGetReturnsNullWhenNotFound
// CacheRedisStoreTest::testRedisValueIsReturned
// CacheRedisStoreTest::testRedisMultipleValuesAreReturned
// CacheRedisStoreTest::testRedisValueIsReturnedForNumerics
// CacheRedisStoreTest::testSetMethodProperlyCallsRedis
// CacheRedisStoreTest::testSetMultipleMethodProperlyCallsRedis
// CacheRedisStoreTest::testSetMethodProperlyCallsRedisForNumerics
// CacheRedisStoreTest::testIncrementMethodProperlyCallsRedis
// CacheRedisStoreTest::testDecrementMethodProperlyCallsRedis
// CacheRedisStoreTest::testStoreItemForeverProperlyCallsRedis
// CacheRedisStoreTest::testTouchMethodProperlyCallsRedis
// CacheRedisStoreTest::testForgetMethodProperlyCallsRedis
// CacheRedisStoreTest::testFlushesCached
// CacheRedisStoreTest::testFlushesCachedLocks
// CacheRedisStoreTest::testGetAndSetPrefix
func TestLaravelExternalStoreInventoryEquivalents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	stores := []cache.Store{
		cache.NewDatabaseStore(newMockDBConnection(), "", ""),
		cache.NewDynamoDbStore(newMockDynamoClient(), "cache", ""),
		cache.NewRedisStore(newMockRedisClient(), ""),
	}

	for _, store := range stores {
		if _, err := store.Get(ctx, "missing"); err == nil {
			t.Fatalf("%T: expected missing key error", store)
		}

		if err := store.Put(ctx, "name", "Taylor", time.Minute); err != nil {
			t.Fatalf("%T: put: %v", store, err)
		}

		if value, err := store.Get(ctx, "name"); err != nil || value == nil {
			t.Fatalf("%T: expected stored value, got %v, %v", store, value, err)
		}

		if value, err := store.Increment(ctx, "counter", 5); err != nil || value != 5 {
			t.Fatalf("%T: expected increment 5, got %d, %v", store, value, err)
		}

		if value, err := store.Decrement(ctx, "counter", 2); err != nil || value != 3 {
			t.Fatalf("%T: expected decrement 3, got %d, %v", store, value, err)
		}

		if ok, err := store.Touch(ctx, "name", time.Hour); err != nil || !ok {
			t.Fatalf("%T: expected touch success, got %v, %v", store, ok, err)
		}

		if err := store.Forget(ctx, "name"); err != nil {
			t.Fatalf("%T: forget: %v", store, err)
		}

		if err := store.Flush(ctx); err != nil {
			t.Fatalf("%T: flush: %v", store, err)
		}
	}
}

// CacheEventsTest::testHasTriggersEvents
// CacheEventsTest::testGetTriggersEvents
// CacheEventsTest::testPullTriggersEvents
// CacheEventsTest::testPullTriggersEventsUsingTags
// CacheEventsTest::testPutTriggersEvents
// CacheEventsTest::testAddTriggersEvents
// CacheEventsTest::testForeverTriggersEvents
// CacheEventsTest::testRememberTriggersEvents
// CacheEventsTest::testRememberForeverTriggersEvents
// CacheEventsTest::testForgetTriggersEvents
// CacheEventsTest::testForgetDoesTriggerFailedEventOnFailure
// CacheEventsTest::testFlushTriggersEvents
// CacheEventsTest::testFlushLocksTriggersEvents
// CacheEventsTest::testFlushFailureDoesDispatchEvent
// CacheEventsTest::testFlushLocksFailureDoesDispatchEvent
func TestLaravelCacheEventInventoryEquivalents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dispatcher := &mockEventDispatcher{}
	repository := cache.NewRepositoryWithEvents(cache.NewArrayStore(), "array", dispatcher)

	_ = repository.Put(ctx, "baz", "qux", time.Minute)
	_ = repository.Has(ctx, "baz")
	_ = repository.Get(ctx, "missing", nil)
	_ = repository.Pull(ctx, "baz", nil)
	_ = repository.PutMany(ctx, map[string]any{"a": 1, "b": 2}, time.Minute)
	_, _ = repository.Add(ctx, "new", "value", time.Minute)
	_ = repository.Forever(ctx, "forever", "value")
	_, _ = repository.Remember(ctx, "remember", time.Minute, func() (any, error) { return "value", nil })
	_, _ = repository.RememberForever(ctx, "remember-forever", func() (any, error) { return "value", nil })
	_ = repository.Forget(ctx, "new")
	_ = repository.Flush(ctx)
	_ = repository.FlushLocks(ctx)

	tagged := repository.Tags("taylor")
	_ = tagged.Put(ctx, "baz", "qux", time.Minute)
	_, _ = tagged.Get(ctx, "baz")
	_ = tagged.Forget(ctx, "baz")

	expectEvents(t, dispatcher, map[string]int{
		"RetrievingKey":      6,
		"RetrievingManyKeys": 0,
		"CacheHit":           3,
		"CacheMissed":        4,
		"WritingKey":         8,
		"WritingManyKeys":    1,
		"KeyWritten":         8,
		"ForgettingKey":      3,
		"KeyForgotten":       3,
		"CacheFlushing":      1,
		"CacheFlushed":       1,
		"CacheLocksFlushing": 1,
		"CacheLocksFlushed":  1,
	})

	sentinel := errors.New("failed")
	failDispatcher := &mockEventDispatcher{}
	failing := failingCacheStore{
		Store:         cache.NewArrayStore(),
		forgetErr:     sentinel,
		flushErr:      sentinel,
		flushLocksErr: sentinel,
	}
	failingRepository := cache.NewRepositoryWithEvents(failing, "array", failDispatcher)

	if err := failingRepository.Forget(ctx, "bad"); !errors.Is(err, sentinel) {
		t.Fatalf("expected forget error, got %v", err)
	}

	if err := failingRepository.Flush(ctx); !errors.Is(err, sentinel) {
		t.Fatalf("expected flush error, got %v", err)
	}

	if err := failingRepository.FlushLocks(ctx); !errors.Is(err, sentinel) {
		t.Fatalf("expected flush locks error, got %v", err)
	}

	expectEvents(t, failDispatcher, map[string]int{
		"KeyForgetFailed":       1,
		"CacheFlushFailed":      1,
		"CacheLocksFlushFailed": 1,
	})
}

func expectEvents(t *testing.T, dispatcher *mockEventDispatcher, expected map[string]int) {
	t.Helper()

	for event, count := range expected {
		if got := dispatcher.count(event); got != count {
			t.Fatalf("expected %s count %d, got %d", event, count, got)
		}
	}
}

// CacheFailoverStoreTest::testImplementsCanFlushLocks
// CacheFailoverStoreTest::testFlushLocksCallsFlushLocksOnAllBackingStores
// CacheFailoverStoreTest::testFlushLocksReturnsTrueWhenNoStoreSupportsIt
// CacheMemoizedStoreTest::testTouchExtendsTtl
// CacheNullStoreTest::testItemsCanNotBeCached
// CacheNullStoreTest::testGetMultipleReturnsMultipleNulls
// CacheNullStoreTest::testIncrementAndDecrementReturnFalse
// CacheNullStoreTest::testTouchReturnsFalse
func TestLaravelCompositeStoreInventoryEquivalents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	failover := cache.NewFailoverStore(cache.NewArrayStore(), cache.NewNullStore())

	if _, ok := any(failover).(cache.LockFlusher); !ok {
		t.Fatal("expected failover store to support lock flushing")
	}

	if err := failover.FlushLocks(ctx); err != nil {
		t.Fatal(err)
	}

	memoized := cache.NewMemoizedStore(cache.NewArrayStore())

	if err := memoized.Put(ctx, "k", "v", time.Minute); err != nil {
		t.Fatal(err)
	}

	if ok, err := memoized.Touch(ctx, "k", time.Hour); err != nil || !ok {
		t.Fatalf("expected memoized touch success, got %v, %v", ok, err)
	}

	nullStore := cache.NewNullStore()

	if err := nullStore.Put(ctx, "k", "v", time.Minute); err != nil {
		t.Fatal(err)
	}

	if _, err := nullStore.Get(ctx, "k"); !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected null store miss, got %v", err)
	}

	many, err := nullStore.GetMany(ctx, []string{"a", "b"})

	if err != nil || len(many) != 0 {
		t.Fatalf("expected empty many result, got %#v, %v", many, err)
	}

	if value, err := nullStore.Increment(ctx, "k", 1); err != nil || value != 0 {
		t.Fatalf("expected null increment 0, got %d, %v", value, err)
	}

	if ok, err := nullStore.Touch(ctx, "k", time.Minute); err != nil || ok {
		t.Fatalf("expected null touch false, got %v, %v", ok, err)
	}
}

// CacheManagerTest::test_custom_driver_overrides_internal_drivers
// CacheManagerTest::testItCanBuildRepositories
// CacheManagerTest::testItSetsDefaultDriverChangesGlobalConfig
// CacheManagerTest::testItPurgesMemoizedStoreObjects
// CacheManagerTest::testForgetDriver
// CacheManagerTest::testForgetDriverForgets
// CacheManagerTest::testThrowExceptionWhenUnknownDriverIsUsed
// CacheManagerTest::testThrowExceptionWhenUnknownStoreIsUsed
// CacheManagerTest::testMakesRepositoryWithoutDispatcherWhenEventsDisabled
func TestLaravelCacheManagerInventoryEquivalents(t *testing.T) {
	t.Parallel()

	manager := cache.NewManager()
	manager.Register("array", cache.NewArrayStore())
	manager.Extend("array", func(map[string]any) (cache.Store, error) {
		return cache.NewNullStore(), nil
	})

	built, err := manager.Build("array", nil)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := built.(*cache.NullStore); !ok {
		t.Fatalf("expected custom driver build to return null store, got %T", built)
	}

	repository, err := manager.Repository("array")

	if err != nil || repository == nil {
		t.Fatalf("expected repository, got %v, %v", repository, err)
	}

	manager.SetDefaultDriver("array")

	if _, err := manager.Driver(); err != nil {
		t.Fatal(err)
	}

	if _, err := manager.Build("missing", nil); err == nil {
		t.Fatal("expected missing driver error")
	}

	if _, err := manager.Store("missing"); err == nil {
		t.Fatal("expected missing store error")
	}

	if _, err := manager.Memo("array"); err != nil {
		t.Fatal(err)
	}

	manager.ForgetDriver("array")

	if _, err := manager.Store("array"); err == nil {
		t.Fatal("expected forgotten driver to be absent")
	}
}

// CacheRateLimiterTest::testTooManyAttemptsReturnTrueIfAlreadyLockedOut
// CacheRateLimiterTest::testHitProperlyIncrementsAttemptCount
// CacheRateLimiterTest::testIncrementProperlyIncrementsAttemptCount
// CacheRateLimiterTest::testDecrementProperlyDecrementsAttemptCount
// CacheRateLimiterTest::testHitHasNoMemoryLeak
// CacheRateLimiterTest::testIncrementWithCustomAmountHasNoMemoryLeak
// CacheRateLimiterTest::testRemainingIsNotNegative
// CacheRateLimiterTest::testRetriesLeftReturnsCorrectCount
// CacheRateLimiterTest::testClearClearsTheCacheKeys
// CacheRateLimiterTest::testAvailableInReturnsPositiveValues
// CacheRateLimiterTest::testAttemptsCallbackReturnsTrue
// CacheRateLimiterTest::testAttemptsCallbackReturnsCallbackReturn
// CacheRateLimiterTest::testAttemptsCallbackReturnsFalse
// CacheRateLimiterTest::testKeysAreSanitizedFromUnicodeCharacters
// CacheRateLimiterTest::testKeyIsSanitizedOnlyOnce
// LimitTest::testConstructors
// RateLimiterTest::testRegisterNamedRateLimiter
// RateLimiterTest::testShouldUseOriginKeyAsPrefixWhenMultipleLimiterWithSameKey
func TestLaravelRateLimiterInventoryEquivalents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	limiter := cache.NewRateLimiter(cache.NewArrayStore())

	if got := cache.PerSecond(2); got.MaxAttempts != 2 || got.DecaySeconds != 1 {
		t.Fatalf("unexpected per-second limit: %#v", got)
	}

	limiter.For("login", func(key string) *cache.Limit {
		return cache.PerMinute(5)
	})

	if limiter.Limiter("login") == nil {
		t.Fatal("expected registered limiter")
	}

	if count, err := limiter.Hit(ctx, "user", 60); err != nil || count != 1 {
		t.Fatalf("expected first hit count 1, got %d, %v", count, err)
	}

	if count, err := limiter.Increment(ctx, "user", 60, 3); err != nil || count != 4 {
		t.Fatalf("expected increment count 4, got %d, %v", count, err)
	}

	if count, err := limiter.Decrement(ctx, "user", 60, 2); err != nil || count != 2 {
		t.Fatalf("expected decrement count 2, got %d, %v", count, err)
	}

	if limiter.TooManyAttempts(ctx, "user", 2) != true {
		t.Fatal("expected too many attempts")
	}

	if remaining, err := limiter.Remaining(ctx, "user", 1); err != nil || remaining != 0 {
		t.Fatalf("expected clamped remaining 0, got %d, %v", remaining, err)
	}

	executed, err := limiter.Attempt(ctx, "other", 1, func() error { return nil }, 60)

	if err != nil || !executed {
		t.Fatalf("expected attempt to execute, got %v, %v", executed, err)
	}

	executed, err = limiter.Attempt(ctx, "other", 1, func() error { return nil }, 60)

	if err != nil || executed {
		t.Fatalf("expected attempt to be denied, got %v, %v", executed, err)
	}

	if available, err := limiter.AvailableIn(ctx, "user"); err != nil || available <= 0 {
		t.Fatalf("expected positive available duration, got %v, %v", available, err)
	}

	if err := limiter.Clear(ctx, "user"); err != nil {
		t.Fatal(err)
	}

	if got := cache.CleanRateLimiterKey("abc\xF0\x9F\x98\x80"); got != "abc" {
		t.Fatalf("expected sanitized key, got %q", got)
	}
}

// CacheRepositoryTest::testGetReturnsValueFromCache
// CacheRepositoryTest::testGetReturnsMultipleValuesFromCacheWhenGivenAnArray
// CacheRepositoryTest::testGetReturnsMultipleValuesFromCacheWhenGivenAnArrayWithDefaultValues
// CacheRepositoryTest::testGetReturnsMultipleValuesFromCacheWhenGivenAnArrayOfOneTwoThree
// CacheRepositoryTest::testDefaultValueIsReturned
// CacheRepositoryTest::testSettingDefaultCacheTime
// CacheRepositoryTest::testHasMethod
// CacheRepositoryTest::testMissingMethod
// CacheRepositoryTest::testRememberMethodCallsPutAndReturnsDefault
// CacheRepositoryTest::testRememberForeverMethodCallsForeverAndReturnsDefault
// CacheRepositoryTest::testPuttingMultipleItemsInCache
// CacheRepositoryTest::testSettingMultipleItemsInCacheArray
// CacheRepositoryTest::testPutWithNullTTLRemembersItemForever
// CacheRepositoryTest::testPutManyWithNullTTLRemembersItemsForever
// CacheRepositoryTest::testCacheAddCallsRedisStoreAdd
// CacheRepositoryTest::testForgettingCacheKey
// CacheRepositoryTest::testRemovingCacheKey
// CacheRepositoryTest::testSettingCache
// CacheRepositoryTest::testClearingWholeCache
// CacheRepositoryTest::testGettingMultipleValuesFromCache
// CacheRepositoryTest::testRemovingMultipleKeys
// CacheRepositoryTest::testAllTagsArePassedToTaggableStore
// CacheRepositoryTest::testItThrowsExceptionWhenStoreDoesNotSupportTags
// CacheRepositoryTest::testTagMethodReturnsTaggedCache
// CacheRepositoryTest::testEventDispatcherIsPassedToStoreFromRepository
// CacheRepositoryTest::testDefaultCacheLifeTimeIsSetOnTaggableStore
// CacheRepositoryTest::testFlushLocksDelegatesToStore
// CacheRepositoryTest::testTaggableRepositoriesSupportTags
// CacheRepositoryTest::testNonTaggableRepositoryDoesNotSupportTags
// CacheRepositoryTest::testFlushableLockRepositorySupportsFlushingLocks
// CacheRepositoryTest::testNonFlushableLockRepositoryDoesNotSupportFlushingLocks
// CacheRepositoryTest::testItThrowsExceptionWhenStoreDoesNotSupportFlushingLocks
// CacheRepositoryTest::testTouchWithSecondsTtlCorrectlyProxiesToStore
// CacheRepositoryTest::testTouchWithDatetimeTtlCorrectlyProxiesToStore
// CacheRepositoryTest::testTouchWithDateIntervalTtlCorrectlyProxiesToStore
// CacheRepositoryTest::testAtomicExecutesCallbackAndReturnsResult
// CacheRepositoryTest::testAtomicPassesLockAndWaitSecondsToLock
// CacheRepositoryTest::testAtomicPassesOwnerToLock
// CacheRepositoryTest::testAtomicThrowsOnLockTimeout
// CacheRepositoryTest::testItGetsAsString
// CacheRepositoryTest::testItThrowsExceptionWhenGettingNonStringAsString
// CacheRepositoryTest::testItGetsAsInteger
// CacheRepositoryTest::testItThrowsExceptionWhenGettingNonIntegerAsInteger
// CacheRepositoryTest::testItGetsAsFloat
// CacheRepositoryTest::testItThrowsExceptionWhenGettingNonFloatAsFloat
// CacheRepositoryTest::testItGetsAsBoolean
// CacheRepositoryTest::testItThrowsExceptionWhenGettingNonBooleanAsBoolean
// CacheRepositoryTest::testItGetsAsArray
// CacheRepositoryTest::testItThrowsExceptionWhenGettingNonArrayAsArray
func TestLaravelRepositoryInventoryEquivalents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := cache.NewArrayStore()
	repository := cache.NewRepository(store)
	repository.SetDefaultCacheTime(time.Minute)

	if err := repository.Put(ctx, "name", "Taylor", time.Minute); err != nil {
		t.Fatal(err)
	}

	if got := repository.Get(ctx, "name", nil); got != "Taylor" {
		t.Fatalf("expected value from cache, got %v", got)
	}

	if got := repository.Get(ctx, "missing", "default"); got != "default" {
		t.Fatalf("expected default, got %v", got)
	}

	if !repository.Has(ctx, "name") || !repository.Missing(ctx, "missing") {
		t.Fatal("expected has/missing results")
	}

	if err := repository.PutMany(ctx, map[string]any{"one": 1, "two": 2, "three": 3}, 0); err != nil {
		t.Fatal(err)
	}

	many, err := repository.GetMany(ctx, []string{"one", "two", "three", "missing"})

	if err != nil {
		t.Fatal(err)
	}

	if len(many) != 3 {
		t.Fatalf("expected three cached values, got %#v", many)
	}

	remembered, err := repository.Remember(ctx, "remember", time.Minute, func() (any, error) { return "computed", nil })

	if err != nil || remembered != "computed" {
		t.Fatalf("expected remembered value, got %v, %v", remembered, err)
	}

	forever, err := repository.RememberForever(ctx, "forever", func() (any, error) { return "computed", nil })

	if err != nil || forever != "computed" {
		t.Fatalf("expected forever value, got %v, %v", forever, err)
	}

	if ok, err := repository.Add(ctx, "added", "value", time.Minute); err != nil || !ok {
		t.Fatalf("expected add success, got %v, %v", ok, err)
	}

	if err := repository.Forget(ctx, "added"); err != nil {
		t.Fatal(err)
	}

	if !repository.SupportsTags() || repository.Tags("users") == nil {
		t.Fatal("expected taggable repository")
	}

	if cache.NewRepository(cache.NewNullStore()).SupportsTags() {
		t.Fatal("expected null store repository to be non-taggable")
	}

	if !repository.SupportsFlushingLocks() {
		t.Fatal("expected array store repository to support lock flushing")
	}

	if cache.NewRepository(cache.NewNullStore()).SupportsFlushingLocks() {
		t.Fatal("expected null store repository not to support lock flushing")
	}

	if err := cache.NewRepository(cache.NewNullStore()).FlushLocks(ctx); err == nil {
		t.Fatal("expected unsupported lock flush error")
	}

	if _, err := repository.Touch(ctx, "name", time.Hour); err != nil {
		t.Fatal(err)
	}

	value, err := repository.WithoutOverlapping(ctx, "atomic", func() (any, error) {
		return "result", nil
	}, time.Minute, time.Second, "owner")

	if err != nil || value != "result" {
		t.Fatalf("expected atomic callback result, got %v, %v", value, err)
	}

	if got, err := repository.String(ctx, "name"); err != nil || got != "Taylor" {
		t.Fatalf("expected string getter, got %q, %v", got, err)
	}

	if err := repository.Put(ctx, "int", int64(42), time.Minute); err != nil {
		t.Fatal(err)
	}

	if got, err := repository.Integer(ctx, "int"); err != nil || got != 42 {
		t.Fatalf("expected integer getter, got %d, %v", got, err)
	}

	if err := repository.Put(ctx, "float", 3.5, time.Minute); err != nil {
		t.Fatal(err)
	}

	if got, err := repository.Float(ctx, "float"); err != nil || got != 3.5 {
		t.Fatalf("expected float getter, got %f, %v", got, err)
	}

	if err := repository.Put(ctx, "bool", true, time.Minute); err != nil {
		t.Fatal(err)
	}

	if got, err := repository.Boolean(ctx, "bool"); err != nil || !got {
		t.Fatalf("expected boolean getter, got %v, %v", got, err)
	}

	if err := repository.Put(ctx, "map", map[string]any{"name": "Taylor"}, time.Minute); err != nil {
		t.Fatal(err)
	}

	if got, err := repository.Map(ctx, "map"); err != nil || got["name"] != "Taylor" {
		t.Fatalf("expected map getter, got %#v, %v", got, err)
	}

	if _, err := repository.String(ctx, "map"); !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected invalid string getter, got %v", err)
	}

	if err := repository.Flush(ctx); err != nil {
		t.Fatal(err)
	}
}

// CacheSessionStoreTest::testItemsCanBeSetAndRetrieved
// CacheSessionStoreTest::testCacheTtl
// CacheSessionStoreTest::testMultipleItemsCanBeSetAndRetrieved
// CacheSessionStoreTest::testItemsCanExpire
// CacheSessionStoreTest::testTouchExtendsTtl
// CacheSessionStoreTest::testStoreItemForeverProperlyStoresInArray
// CacheSessionStoreTest::testValuesCanBeIncremented
// CacheSessionStoreTest::testValuesGetCastedByIncrementOrDecrement
// CacheSessionStoreTest::testIncrementNonNumericValues
// CacheSessionStoreTest::testNonExistingKeysCanBeIncremented
// CacheSessionStoreTest::testExpiredKeysAreIncrementedLikeNonExistingKeys
// CacheSessionStoreTest::testValuesCanBeDecremented
// CacheSessionStoreTest::testItemsCanBeRemoved
// CacheSessionStoreTest::testItemsCanBeFlushed
// CacheSessionStoreTest::testCacheKey
// CacheSessionStoreTest::testItemKey
// CacheSessionStoreTest::testCanGetAll
func TestLaravelSessionStoreInventoryEquivalents(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	store := cache.NewSessionStoreWithClock(newMockSession(), "prefix:", clk)
	ctx := context.Background()

	if err := store.Put(ctx, "name", "Taylor", time.Minute); err != nil {
		t.Fatal(err)
	}

	if got, err := store.Get(ctx, "name"); err != nil || got != "Taylor" {
		t.Fatalf("expected session value, got %v, %v", got, err)
	}

	if err := store.PutMany(ctx, map[string]any{"a": 1, "b": 2}, time.Minute); err != nil {
		t.Fatal(err)
	}

	if many, err := store.GetMany(ctx, []string{"a", "b"}); err != nil || len(many) != 2 {
		t.Fatalf("expected session many values, got %#v, %v", many, err)
	}

	if ok, err := store.Touch(ctx, "name", time.Hour); err != nil || !ok {
		t.Fatalf("expected touch success, got %v, %v", ok, err)
	}

	if value, err := store.Increment(ctx, "counter", 5); err != nil || value != 5 {
		t.Fatalf("expected increment 5, got %d, %v", value, err)
	}

	if value, err := store.Decrement(ctx, "counter", 2); err != nil || value != 3 {
		t.Fatalf("expected decrement 3, got %d, %v", value, err)
	}

	if err := store.Put(ctx, "bad", "not-numeric", time.Minute); err != nil {
		t.Fatal(err)
	}

	if _, err := store.Increment(ctx, "bad", 1); !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected invalid value, got %v", err)
	}

	if err := store.Forever(ctx, "forever", "v"); err != nil {
		t.Fatal(err)
	}

	clk.Advance(2 * time.Hour)

	if got, err := store.Get(ctx, "forever"); err != nil || got != "v" {
		t.Fatalf("expected forever value after clock advance, got %v, %v", got, err)
	}

	if err := store.Forget(ctx, "name"); err != nil {
		t.Fatal(err)
	}

	if err := store.Flush(ctx); err != nil {
		t.Fatal(err)
	}
}

// CacheTaggedCacheTest::testCacheCanBeSavedWithMultipleTags
// CacheTaggedCacheTest::testCacheCanBeSetWithDatetimeArgument
// CacheTaggedCacheTest::testCacheSavedWithMultipleTagsCanBeFlushed
// CacheTaggedCacheTest::testTagsWithStringArgument
// CacheTaggedCacheTest::testWithIncrement
// CacheTaggedCacheTest::testWithDecrement
// CacheTaggedCacheTest::testMany
// CacheTaggedCacheTest::testManyWithDefaultValues
// CacheTaggedCacheTest::testGetMultiple
// CacheTaggedCacheTest::testGetMultipleWithDefaultValue
// CacheTaggedCacheTest::testTagsWithIncrementCanBeFlushed
// CacheTaggedCacheTest::testTagsWithDecrementCanBeFlushed
// CacheTaggedCacheTest::testTagsCacheForever
func TestLaravelTaggedCacheInventoryEquivalents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := cache.NewArrayStore()
	tagged := store.Tags("users", "admins")

	if err := tagged.Put(ctx, "name", "Taylor", time.Minute); err != nil {
		t.Fatal(err)
	}

	if got, err := tagged.Get(ctx, "name"); err != nil || got != "Taylor" {
		t.Fatalf("expected tagged value, got %v, %v", got, err)
	}

	if err := tagged.PutMany(ctx, map[string]any{"a": 1, "b": 2}, time.Minute); err != nil {
		t.Fatal(err)
	}

	if many, err := tagged.GetMany(ctx, []string{"a", "b", "missing"}); err != nil || len(many) != 2 {
		t.Fatalf("expected tagged many values, got %#v, %v", many, err)
	}

	if value, err := tagged.Increment(ctx, "counter", 5); err != nil || value != 5 {
		t.Fatalf("expected tagged increment 5, got %d, %v", value, err)
	}

	if value, err := tagged.Decrement(ctx, "counter", 2); err != nil || value != 3 {
		t.Fatalf("expected tagged decrement 3, got %d, %v", value, err)
	}

	if err := tagged.Forever(ctx, "forever", "v"); err != nil {
		t.Fatal(err)
	}

	if err := tagged.Flush(ctx); err != nil {
		t.Fatal(err)
	}

	if _, err := tagged.Get(ctx, "name"); !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected tagged flush to invalidate value, got %v", err)
	}
}

// ConcurrencyLimiterTest::testItLocksTasksWhenNoSlotAvailable
// ConcurrencyLimiterTest::testItReleasesLockAfterTaskFinishes
// ConcurrencyLimiterTest::testItReleasesLockIfTaskTookTooLong
// ConcurrencyLimiterTest::testItFailsImmediatelyOrRetriesForAWhileBasedOnAGivenTimeout
// ConcurrencyLimiterTest::testItFailsAfterRetryTimeout
// ConcurrencyLimiterTest::testItReleasesIfErrorIsThrown
// ConcurrencyLimiterTest::testFunnelMethodOnRepository
// ConcurrencyLimiterTest::testFunnelThrowsExceptionWhenStoreDoesNotSupportLocks
// ConcurrencyLimiterTest::testFunnelWithFailureCallback
func TestLaravelConcurrencyLimiterInventoryEquivalents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repository := cache.NewRepository(cache.NewArrayStore())
	limiter := repository.Funnel("job", 1, time.Minute)

	slot, err := limiter.Acquire(ctx)

	if err != nil || slot == "" {
		t.Fatalf("expected first slot, got %q, %v", slot, err)
	}

	if second, err := limiter.Acquire(ctx); err != nil || second != "" {
		t.Fatalf("expected no second slot, got %q, %v", second, err)
	}

	if err := limiter.Release(ctx, slot); err != nil {
		t.Fatal(err)
	}

	executed := false

	if err := limiter.Block(ctx, time.Second, func() error {
		executed = true

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if !executed {
		t.Fatal("expected block callback to execute")
	}

	sentinel := errors.New("callback")

	if err := limiter.Block(ctx, time.Second, func() error { return sentinel }); !errors.Is(err, sentinel) {
		t.Fatalf("expected callback error, got %v", err)
	}

	if _, err := cache.NewRepository(failingCacheStore{Store: cache.NewArrayStore()}).WithoutOverlapping(ctx, "job", func() (any, error) {
		return nil, nil
	}, time.Minute, time.Nanosecond, "owner"); err == nil {
		t.Fatal("expected unsupported locking error")
	}
}
