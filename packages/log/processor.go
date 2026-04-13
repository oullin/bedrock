package log

// Processor transforms a log record before it reaches the handler output.
type Processor interface {
	Process(record Record) Record
}

// ProcessorFunc is an adapter to allow the use of ordinary functions as
// processors.
type ProcessorFunc func(record Record) Record

// Process calls the underlying function.
func (f ProcessorFunc) Process(record Record) Record {
	return f(record)
}

var _ Processor = ProcessorFunc(nil)
