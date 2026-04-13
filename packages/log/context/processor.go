package context

import (
	"github.com/bedrock/packages/log"
)

// ContextLogProcessor enriches log records with data from a context Repository.
type ContextLogProcessor struct {
	repo *Repository
}

var _ log.Processor = (*ContextLogProcessor)(nil)

// NewContextLogProcessor creates a processor that merges repository data into
// log record extras.
func NewContextLogProcessor(repo *Repository) *ContextLogProcessor {
	return &ContextLogProcessor{repo: repo}
}

// Process merges the repository's public data into the record's Extra field.
func (p *ContextLogProcessor) Process(record log.Record) log.Record {
	data := p.repo.All()

	if len(data) == 0 {
		return record
	}

	if record.Extra == nil {
		record.Extra = make(map[string]any, len(data))
	}

	for k, v := range data {
		record.Extra[k] = v
	}

	return record
}
