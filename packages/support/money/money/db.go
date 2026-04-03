package money

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gollin/packages/support/money/currency"
)

var (
	// dbMoneyValueSeparator is used to join together the Amount and Currency components of money.Money instances
	// allowing them to be stored as strings (via the driver.Valuer interface) and unmarshalled as strings (via
	// the sql.Scanner interface); use SetDBMoneyValueSeparator to change this value in a thread-safe manner.
	dbMoneyValueSeparator   = DefaultDBMoneyValueSeparator
	dbMoneyValueSeparatorMu sync.RWMutex
)

const (
	// DefaultDBMoneyValueSeparator is the default value for DBMoneyValueSeparator; can be used to reset the
	// active separator value
	DefaultDBMoneyValueSeparator = "|"
)

// GetDBMoneyValueSeparator returns the current separator value in a thread-safe manner
func GetDBMoneyValueSeparator() string {
	dbMoneyValueSeparatorMu.RLock()

	defer dbMoneyValueSeparatorMu.RUnlock()

	return dbMoneyValueSeparator
}

// SetDBMoneyValueSeparator sets the separator value in a thread-safe manner
func SetDBMoneyValueSeparator(separator string) error {
	if strings.TrimSpace(separator) == "" {
		return fmt.Errorf("separator [%s] cannot be empty", separator)
	}

	dbMoneyValueSeparatorMu.Lock()

	defer dbMoneyValueSeparatorMu.Unlock()

	dbMoneyValueSeparator = separator

	return nil
}

// Value implements driver.Valuer to serialise a Money instance into a delimited string using the DBMoneyValueSeparator,
// for example, "amount|currency_code"
func (m *Money) Value() (driver.Value, error) {
	if err := ensureMoneyProvided(m); err != nil {
		return nil, err
	}

	curr, err := m.Currency()

	if err != nil {
		return nil, err
	}

	return fmt.Sprintf(
		"%d%s%s",
		m.amount,
		GetDBMoneyValueSeparator(),
		curr.Code,
	), nil
}

// Scan implements sql.Scanner to deserialize a Money instance from a DBMoneyValueSeparator-separated string,
// for example, "amount|currency_code"
func (m *Money) Scan(src interface{}) error {
	if m == nil {
		return fmt.Errorf("cannot scan nil Money")
	}

	var amount Amount

	var parts []string
	curr := &currency.Currency{}
	separator := GetDBMoneyValueSeparator()

	switch v := src.(type) {
	case string:
		parts = strings.Split(v, separator)
	case []byte:
		parts = strings.Split(string(v), separator)
	default:
		return fmt.Errorf("don't know how to scan %T into Money; update your query to return a money.DBMoneyValueSeparator-separated pair of \"amount%scurrency_code\"", src, separator)
	}

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("%#v is not valid to scan into Money; update your query to return a money.DBMoneyValueSeparator-separated pair of \"amount%scurrency_code\"", src, separator)
	}

	if a, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
		amount = a
	} else {
		return fmt.Errorf("scanning %#v into an Amount: %v", parts[0], err)
	}

	if err := curr.DbScan(parts[1]); err != nil {
		return fmt.Errorf("scanning %#v into a Currency: %v", parts[1], err)
	}

	// allocate new Money with the scanned amount and curr
	*m = Money{
		amount:   amount,
		currency: curr,
	}

	return nil
}
