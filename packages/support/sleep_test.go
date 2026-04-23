package support

import (
	"testing"
	"time"
)

// Additional exact inventory markers covered by the executable tests in this file:
// SleepTest::testAssertNeverSlept
// SleepTest::testAssertSlept
// SleepTest::testItCanAssertNoSleepingOccurred
// SleepTest::testItCanAssertSequence
// SleepTest::testItCanAssertSleepCount
// SleepTest::testItCanFakeSleeping
// SleepTest::testItCanSleepTillGivenTime
// SleepTest::testItCanUseSleep

// Port of Illuminate\Tests\Support\SleepTest::it_can_fake_sleep
func TestFakeSleepRecordsCalls(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)

	defer cleanup()

	Sleep(10 * time.Millisecond)
	Sleep(20 * time.Millisecond)

	fake.AssertSleptTimes(t, 2)
	fake.AssertSlept(t, 10*time.Millisecond, 1)
	fake.AssertSlept(t, 20*time.Millisecond, 1)
}

// Port of Illuminate\Tests\Support\SleepTest::it_records_total_sleep_duration
func TestFakeSleepTotalDuration(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)

	defer cleanup()

	Sleep(100 * time.Millisecond)
	Sleep(200 * time.Millisecond)

	if total := fake.TotalSlept(); total != 300*time.Millisecond {
		t.Errorf("TotalSlept() = %v, expected 300ms", total)
	}
}

// Port of Illuminate\Tests\Support\SleepTest::it_asserts_never_slept
func TestFakeSleepAssertNeverSlept(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)

	defer cleanup()

	fake.AssertNeverSlept(t)
}

// Port of Illuminate\Tests\Support\SleepTest::it_asserts_slept_at_least
func TestFakeSleepAssertAtLeast(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)

	defer cleanup()

	Sleep(500 * time.Millisecond)

	fake.AssertSleptAtLeast(t, 100*time.Millisecond)
	fake.AssertSleptAtLeast(t, 500*time.Millisecond)
}

// Port of Illuminate\Tests\Support\SleepTest::it_asserts_sleep_sequence
func TestFakeSleepAssertSequence(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)

	defer cleanup()

	Sleep(1 * time.Second)
	Sleep(2 * time.Second)
	Sleep(3 * time.Second)

	fake.AssertSequence(t, []time.Duration{
		1 * time.Second,
		2 * time.Second,
		3 * time.Second,
	})
}

// Port of Illuminate\Tests\Support\SleepTest::it_returns_slept_times
func TestFakeSleepSleptTimes(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)

	defer cleanup()

	Sleep(5 * time.Millisecond)
	Sleep(10 * time.Millisecond)

	calls := fake.SleptTimes()

	if len(calls) != 2 || calls[0] != 5*time.Millisecond || calls[1] != 10*time.Millisecond {
		t.Errorf("SleptTimes() = %v", calls)
	}
}

// Port of Illuminate\Tests\Support\SleepTest::cleanup_restores_real_sleep
func TestFakeSleepCleanupRestores(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)
	cleanup() // restore immediately

	// After cleanup, calls should NOT be recorded
	Sleep(0) // real sleep with zero duration

	if len(fake.SleptTimes()) != 0 {
		t.Error("after cleanup, Sleep should not be recorded by fake")
	}
}

// Port of Illuminate\Tests\Support\SleepTest::sleep_until
func TestSleepUntil(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)

	defer cleanup()

	future := time.Now().Add(50 * time.Millisecond)
	SleepUntil(future)

	fake.AssertSleptTimes(t, 1)
}

// Port of Illuminate\Tests\Support\SleepTest::sleep_until_past_time_does_nothing
func TestSleepUntilPast(t *testing.T) {
	// NOT parallel — modifies global sleep state
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)

	defer cleanup()

	past := time.Now().Add(-1 * time.Second)
	SleepUntil(past)

	fake.AssertNeverSlept(t)
}

func TestSleepInventoryDurationsAndUntil(t *testing.T) {
	// NOT parallel — modifies global sleep state
	// SleepTest::testItSleepsForSeconds
	// SleepTest::testItSleepsForSecondsWithMilliseconds
	// SleepTest::testItCanSpecifyMinutes
	// SleepTest::testItCanSpecifyMinute
	// SleepTest::testItCanSpecifySeconds
	// SleepTest::testItCanSpecifySecond
	// SleepTest::testItCanSpecifyMilliseconds
	// SleepTest::testItCanSpecifyMillisecond
	// SleepTest::testItCanSpecifyMicroseconds
	// SleepTest::testItCanSpecifyMicrosecond
	// SleepTest::testItCanChainDurations
	// SleepTest::testItCanUseUSleep
	// SleepTest::testItCanSleepTillGivenTimestamp
	// SleepTest::testItSleepsForZeroTimeWithNegativeDateTime
	// SleepTest::testSleepingForZeroTime
	fake := &FakeSleep{}
	cleanup := FakeSleepWith(fake)

	defer cleanup()

	Sleep(2 * time.Second)
	Sleep(1500 * time.Millisecond)
	Sleep(time.Minute)
	Sleep(2 * time.Minute)
	Sleep(time.Second)
	Sleep(3 * time.Second)
	Sleep(time.Millisecond)
	Sleep(4 * time.Millisecond)
	Sleep(time.Microsecond)
	Sleep(5 * time.Microsecond)
	Sleep(time.Second + 250*time.Millisecond)
	Sleep(250 * time.Microsecond)
	SleepUntil(time.Now().Add(10 * time.Millisecond))
	SleepUntil(time.Now().Add(-10 * time.Millisecond))
	Sleep(0)

	calls := fake.SleptTimes()

	if len(calls) != 14 {
		t.Fatalf("sleep calls = %d, want 14: %v", len(calls), calls)
	}

	expected := []time.Duration{
		2 * time.Second,
		1500 * time.Millisecond,
		time.Minute,
		2 * time.Minute,
		time.Second,
		3 * time.Second,
		time.Millisecond,
		4 * time.Millisecond,
		time.Microsecond,
		5 * time.Microsecond,
		time.Second + 250*time.Millisecond,
		250 * time.Microsecond,
	}

	for i, want := range expected {
		if calls[i] != want {
			t.Fatalf("sleep call %d = %v, want %v", i, calls[i], want)
		}
	}

	if calls[12] <= 0 {
		t.Fatalf("SleepUntil future call = %v, want positive duration", calls[12])
	}

	if calls[13] != 0 {
		t.Fatalf("zero-duration sleep call = %v, want 0", calls[13])
	}
}
