package failed

import "time"

// SetNow overrides the clock a provider uses when stamping failed_at
// values. It is exported so package _test files can port PHPUnit tests
// that call Carbon::setTestNow; production code should not depend on
// it.
func SetNow(p any, now func() time.Time) {
	switch v := p.(type) {
	case *DatabaseFailedJobProvider:
		v.now = now
	case *DatabaseUuidFailedJobProvider:
		v.now = now
	case *FileFailedJobProvider:
		v.now = now
	case *DynamoDbFailedJobProvider:
		v.now = now
	}
}
