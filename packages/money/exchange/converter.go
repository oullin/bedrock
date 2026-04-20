package exchange

import (
	"github.com/bedrock/packages/money/exception"
)

// Converter provides a simple interface for currency conversion with a fixed exchange
type Converter struct {
	exchange *Exchange
}

// NewConverter creates a new Converter instance with the provided Exchange.
func NewConverter(exchange *Exchange) (*Converter, error) {
	if exchange == nil {
		return nil, exception.ErrNoConverterProvided
	}

	return &Converter{
		exchange: exchange,
	}, nil
}

// GetExchange returns the underlying Exchange instance.
func (c *Converter) GetExchange() (*Exchange, error) {
	if c == nil {
		return nil, exception.ErrNoConverterProvided
	}

	if c.exchange == nil {
		return nil, exception.ErrInvalidExchangeRate
	}

	return c.exchange, nil
}
