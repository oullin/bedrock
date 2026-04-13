package log

// Handler processes log records.
type Handler interface {
	Handle(record Record) error
	IsHandling(level Level) bool
	Close() error
}

// ProcessableHandler is an embeddable struct that provides processor
// management for concrete handlers.
type ProcessableHandler struct {
	processors []Processor
}

// AddProcessor appends a processor to the chain.

// GetProcessors returns the registered processors.

// ProcessRecord applies all registered processors to the record.

// FormattableHandler is an embeddable struct that provides formatter
// management for concrete handlers.
type FormattableHandler struct {
	formatter Formatter
}

func (h *ProcessableHandler) AddProcessor(p Processor) {
	h.processors = append(h.processors, p)
}

func (h *ProcessableHandler) GetProcessors() []Processor {
	return h.processors
}

func (h *ProcessableHandler) ProcessRecord(record Record) Record {
	for _, p := range h.processors {
		record = p.Process(record)
	}

	return record
}

// SetFormatter sets the formatter used by this handler.
func (h *FormattableHandler) SetFormatter(f Formatter) {
	h.formatter = f
}

// GetFormatter returns the configured formatter, falling back to a
// LineFormatter if none has been set.
func (h *FormattableHandler) GetFormatter() Formatter {
	if h.formatter == nil {
		h.formatter = NewLineFormatter()
	}

	return h.formatter
}
