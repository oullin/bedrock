package money

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/exception"
)

// JSONRawData represents the raw JSON structure for money unmarshalling.
// It uses json.Number for the amount to preserve precision during decoding.
type JSONRawData struct {
	Amount   json.Number `json:"amount"`
	Currency string      `json:"currency"`
}

// JSON handles go's JSON marshalling and unmarshalling of Money values.
type JSON struct {
	mutex     sync.RWMutex
	unmarshal func(*Money, []byte) error
	marshal   func(Money) ([]byte, error)
	currency  func() (*currency.Currency, error)
}

func ensureJSONProvided(j *JSON) error {
	if j == nil {
		return exception.ErrNoJSONParserProvided
	}

	return nil
}

// NewJson creates a new JSON parser with default settings.
func NewJson() *JSON {
	j := &JSON{}

	j.unmarshal = j.defaultUnmarshalJSON
	j.marshal = j.defaultMarshalJSON
	j.currency = j.defaultJSONCurrency

	return j
}

// NewJsonWithParser creates a new JSON parser with custom functions.
func NewJsonWithParser(
	unmarshal func(*Money, []byte) error,
	marshal func(Money) ([]byte, error),
	currency func() (*currency.Currency, error),
) *JSON {
	j := NewJson()

	if unmarshal != nil {
		j.unmarshal = unmarshal
	}

	if marshal != nil {
		j.marshal = marshal
	}

	if currency != nil {
		j.currency = currency
	}

	return j
}

// Marshal marshals the money instance using the configured marshal function.
func (j *JSON) Marshal(m Money) ([]byte, error) {
	if err := ensureJSONProvided(j); err != nil {
		return nil, err
	}

	j.mutex.RLock()

	fn := j.marshal
	j.mutex.RUnlock()

	return fn(m)
}

// Unmarshal unmarshals the payload into the provided money instance.
func (j *JSON) Unmarshal(m *Money, b []byte) error {
	if err := ensureJSONProvided(j); err != nil {
		return err
	}

	j.mutex.RLock()

	fn := j.unmarshal
	j.mutex.RUnlock()

	return fn(m, b)
}

// SetUnmarshal sets the custom unmarshal function.
func (j *JSON) SetUnmarshal(fn func(*Money, []byte) error) error {
	if err := ensureJSONProvided(j); err != nil {
		return err
	}

	if fn == nil {
		return exception.ErrJSONUnmarshalFuncNil
	}

	j.mutex.Lock()

	defer j.mutex.Unlock()

	j.unmarshal = fn

	return nil
}

// SetMarshal sets the custom marshal function.
func (j *JSON) SetMarshal(fn func(Money) ([]byte, error)) error {
	if err := ensureJSONProvided(j); err != nil {
		return err
	}

	if fn == nil {
		return exception.ErrJSONMarshalFuncNil
	}

	j.mutex.Lock()

	defer j.mutex.Unlock()

	j.marshal = fn

	return nil
}

// SetCurrency sets the custom currency function.
func (j *JSON) SetCurrency(fn func() (*currency.Currency, error)) error {
	if err := ensureJSONProvided(j); err != nil {
		return err
	}

	j.mutex.Lock()

	defer j.mutex.Unlock()

	j.currency = fn

	return nil
}

func (j *JSON) defaultMarshalJSON(m Money) ([]byte, error) {
	if err := ensureJSONProvided(j); err != nil {
		return nil, err
	}

	if m == (Money{}) {
		m = *NewManager().Create(0, "")
	}

	amount, err := m.Amount()

	if err != nil {
		return nil, err
	}

	curr, err := m.Currency()

	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}{
		Amount:   amount,
		Currency: curr.Code,
	})
}

func (j *JSON) defaultUnmarshalJSON(m *Money, b []byte) error {
	if err := ensureJSONProvided(j); err != nil {
		return err
	}

	// Custom struct to parse amount as json.Number for precision preservation
	// This avoids float64's 2^53 precision limit
	var raw JSONRawData

	err := json.Unmarshal(b, &raw)

	if err != nil {
		// Check if it's a type error for amount or currency fields
		var typeError *json.UnmarshalTypeError

		if errors.As(err, &typeError) {
			if typeError.Field == "amount" || typeError.Field == "currency" {
				return exception.ErrInvalidJSONUnmarshal
			}
		}

		// json.Number returns a plain error for invalid number literals
		// Check if the error message indicates an invalid number
		if strings.HasPrefix(err.Error(), "json: invalid number literal") {
			return exception.ErrInvalidJSONUnmarshal
		}

		var syntaxError *json.SyntaxError

		if errors.As(err, &syntaxError) {
			return exception.ErrInvalidJSONUnmarshal
		}

		// Return other errors (like SyntaxError) as-is
		return exception.NewErrInvalidJSONUnmarshalFrom(err)
	}

	var amount int64

	if raw.Amount != "" {
		// Convert to int64 without precision loss
		amount, err = strconv.ParseInt(string(raw.Amount), 10, 64)

		if err != nil {
			// If it's a decimal or exponent form, use big.Float for precision preservation
			// This avoids float64's 2^53 precision limit for large amounts
			var bigAmount big.Float
			_, _, parseErr := bigAmount.Parse(string(raw.Amount), 10)

			if parseErr != nil {
				return exception.ErrInvalidJSONUnmarshal
			}

			// Round to the nearest integer, ties away from zero
			// This matches the behaviour: 12.50 -> 13, -12.50 -> -13
			var rounded big.Int

			bigAmount.Int(&rounded) // Truncates toward zero

			// Check if we need to round up (for positive) or down (for negative)
			var frac big.Float

			frac.Sub(&bigAmount, new(big.Float).SetInt(&rounded))

			absHalf := big.NewFloat(0.5)

			if bigAmount.Sign() < 0 {
				// For negative numbers, check if |frac| >= 0.5 (ties away from zero)
				frac.Neg(&frac)

				if frac.Cmp(absHalf) >= 0 {
					rounded.Sub(&rounded, big.NewInt(1))
				}
			} else {
				// For positive numbers, check if frac >= 0.5 (ties away from zero)
				if frac.Cmp(absHalf) >= 0 {
					rounded.Add(&rounded, big.NewInt(1))
				}
			}

			// Check for int64 overflow
			if !rounded.IsInt64() {
				return exception.ErrInvalidJSONUnmarshal
			}

			amount = rounded.Int64()
		}
	}

	j.mutex.RLock()
	currencyFn := j.currency
	j.mutex.RUnlock()

	curr, err := currencyFn()

	if err != nil {
		return err
	}

	var parsed Money

	// If the currency is empty or whitespace, use the default currency
	if strings.TrimSpace(raw.Currency) == "" {
		parsed = Money{amount: amount, currency: curr}
	} else if parsed, err = j.getUnmarshalJSONMoney(amount, raw); err != nil {
		return err
	}

	*m = parsed

	return nil
}

func (j *JSON) getUnmarshalJSONMoney(amount int64, raw JSONRawData) (Money, error) {
	if err := ensureJSONProvided(j); err != nil {
		return Money{}, err
	}

	curr := currency.NewCurrenciesMap().FindByCode(
		strings.TrimSpace(raw.Currency),
	)

	if curr == nil {
		return Money{}, exception.ErrCurrencyNotFound
	}

	return Money{
		amount:   amount,
		currency: curr,
	}, nil
}

func (j *JSON) defaultJSONCurrency() (*currency.Currency, error) {
	if err := ensureJSONProvided(j); err != nil {
		return nil, err
	}

	c := currency.NewManager().GetDefault()

	if c == nil {
		return nil, fmt.Errorf("no default currency found")
	}

	return c, nil
}
