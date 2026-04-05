package support

import "time"

var inspiringQuotes = []string{
	"Bedrock keeps its footing by building one sound layer at a time.",
	"Discipline turns rough stone into a lasting foundation.",
	"Small, deliberate steps outlast dramatic rewrites.",
}

// InspiringQuote returns a deterministic quote for console output and tests.
func InspiringQuote() string {
	day := time.Now().UTC().YearDay()

	return inspiringQuotes[day%len(inspiringQuotes)]
}
