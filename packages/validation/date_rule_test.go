package validation_test

import (
	"testing"
	"time"
)

func utcDate(offsetDays int) string {
	return time.Now().UTC().AddDate(0, 0, offsetDays).Format("2006-01-02")
}

// ValidationDateRuleTest::testDefaultDateRule
func TestDateRule_DefaultDateRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"published_at": utcDate(0)},
		map[string]any{"published_at": "date"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": "not-a-date"},
		map[string]any{"published_at": "date"},
	)
	assertFails(t, v2)
}

// ValidationDateRuleTest::testDateFormatRule
func TestDateRule_DateFormatRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"published_at": "2026-04-22"},
		map[string]any{"published_at": "date_format:Y-m-d"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": "22/04/2026"},
		map[string]any{"published_at": "date_format:Y-m-d"},
	)
	assertFails(t, v2)
}

// ValidationDateRuleTest::testAfterTodayRule
func TestDateRule_AfterTodayRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"published_at": utcDate(1)},
		map[string]any{"published_at": "after:today"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": utcDate(0)},
		map[string]any{"published_at": "after:today"},
	)
	assertFails(t, v2)
}

// ValidationDateRuleTest::testBeforeTodayRule
func TestDateRule_BeforeTodayRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"published_at": utcDate(-1)},
		map[string]any{"published_at": "before:today"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": utcDate(0)},
		map[string]any{"published_at": "before:today"},
	)
	assertFails(t, v2)
}

// ValidationDateRuleTest::testAfterSpecificDateRule
func TestDateRule_AfterSpecificDateRule(t *testing.T) {
	t.Parallel()

	base := utcDate(0)

	v := makeValidator(
		map[string]any{"published_at": utcDate(1)},
		map[string]any{"published_at": "after:" + base},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": utcDate(-1)},
		map[string]any{"published_at": "after:" + base},
	)
	assertFails(t, v2)
}

// ValidationDateRuleTest::testBeforeSpecificDateRule
func TestDateRule_BeforeSpecificDateRule(t *testing.T) {
	t.Parallel()

	base := utcDate(0)

	v := makeValidator(
		map[string]any{"published_at": utcDate(-1)},
		map[string]any{"published_at": "before:" + base},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": utcDate(1)},
		map[string]any{"published_at": "before:" + base},
	)
	assertFails(t, v2)
}

// ValidationDateRuleTest::testAfterOrEqualSpecificDateRule
func TestDateRule_AfterOrEqualSpecificDateRule(t *testing.T) {
	t.Parallel()

	base := utcDate(0)

	v := makeValidator(
		map[string]any{"published_at": base},
		map[string]any{"published_at": "after_or_equal:" + base},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": utcDate(-1)},
		map[string]any{"published_at": "after_or_equal:" + base},
	)
	assertFails(t, v2)
}

// ValidationDateRuleTest::testBeforeOrEqualSpecificDateRule
func TestDateRule_BeforeOrEqualSpecificDateRule(t *testing.T) {
	t.Parallel()

	base := utcDate(0)

	v := makeValidator(
		map[string]any{"published_at": base},
		map[string]any{"published_at": "before_or_equal:" + base},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": utcDate(1)},
		map[string]any{"published_at": "before_or_equal:" + base},
	)
	assertFails(t, v2)
}

// ValidationDateRuleTest::testBetweenDatesRule
func TestDateRule_BetweenDatesRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"published_at": utcDate(0)},
		map[string]any{"published_at": "after:" + utcDate(-1) + "|before:" + utcDate(1)},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": utcDate(-1)},
		map[string]any{"published_at": "after:" + utcDate(-1) + "|before:" + utcDate(1)},
	)
	assertFails(t, v2)
}

// ValidationDateRuleTest::testBetweenOrEqualDatesRule
func TestDateRule_BetweenOrEqualDatesRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"published_at": utcDate(0)},
		map[string]any{"published_at": "after_or_equal:" + utcDate(-1) + "|before_or_equal:" + utcDate(1)},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": utcDate(-1)},
		map[string]any{"published_at": "after_or_equal:" + utcDate(-1) + "|before_or_equal:" + utcDate(1)},
	)
	assertPasses(t, v2)
}

// ValidationDateRuleTest::testChainedRules
func TestDateRule_ChainedRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"published_at": utcDate(0)},
		map[string]any{"published_at": "date|after:" + utcDate(-1) + "|before:" + utcDate(1)},
	)
	assertPasses(t, v)
}

// ValidationDateRuleTest::testDateValidation
func TestDateRule_DateValidation(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"published_at": "2026-04-22"},
		map[string]any{"published_at": "date"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"published_at": "22-04-2026"},
		map[string]any{"published_at": "date"},
	)
	assertFails(t, v2)
}
