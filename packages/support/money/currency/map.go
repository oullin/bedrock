package currency

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gollin/packages/support/money/exception"
)

// Map holds a reference to the active currency configuration.
// It uses a pointer to the map to ensure inexpensive copying of the struct.
type Map struct {
	fallback *Currency
	dataset  *map[string]*Currency
}

// NewCurrenciesMap returns the singleton instance of Map.
// It initialises the data via getCurrencies() only on the very first call.
// Later calls return the cached instance immediately.
var NewCurrenciesMap = func() func() Map {
	// --- Hidden Persistent Scope ---
	// These variables live forever but are ONLY accessible inside this block.
	var once sync.Once

	var singleton *map[string]*Currency
	// -------------------------------

	return func() Map {
		once.Do(func() {
			// This runs only once, populating the hidden variable
			data := getCurrencies()
			singleton = &data
		})

		// On later calls, 'once' is skipped, but 'singleton' is remembered.
		return Map{
			dataset: singleton,
		}
	}
}()

// NewCurrenciesMapFrom returns the singleton instance, initialising it with inputData
// if it has not been initialized yet.
//
// WARNING: This follows a "first-write-wins" pattern. The inputData is used ONLY
// on the first call. Later calls ignore inputData and return the existing singleton.
var NewCurrenciesMapFrom = createCurrenciesMapFromFactory()

// createCurrenciesMapFromFactory creates a new singleton factory function
// for initialising currency maps with custom data.
//
// This function returns a closure that implements the singleton pattern with thread-safe access.
// The returned function follows a "first-write-wins" pattern where:
// - The first call with valid inputData initialises the singleton
// - Subsequent calls return the existing singleton, ignoring new inputData
// - Thread-safety is ensured via mutex
//
// This factory is extracted to allow test code to create fresh instances
// without duplicating the singleton implementation logic.
func createCurrenciesMapFromFactory() func(inputData *map[string]*Currency) (Map, error) {
	var mu sync.Mutex

	var singleton *map[string]*Currency

	return func(inputData *map[string]*Currency) (Map, error) {
		mu.Lock()

		defer mu.Unlock()

		// Logic Check: If already initialised, return existing immediately.
		// This ensures that calling NewCurrenciesMapFrom(nil) on a 2nd try doesn't error out.
		if singleton != nil {
			return Map{dataset: singleton}, nil
		}

		if inputData == nil || len(*inputData) == 0 {
			return Map{}, fmt.Errorf("data can't be nil or empty when constructing a given Map")
		}

		singleton = inputData

		return Map{
			dataset: singleton,
		}, nil
	}
}

// Get safely retrieves a currency by its code.
// Returns nil if the currency is not found.
func (cm Map) Get(code string) *Currency {
	if cm.dataset == nil {
		return nil
	}

	return (*cm.dataset)[code]
}

// IsEmpty checks if the currency map is empty or nil.
func (cm Map) IsEmpty() bool {
	if cm.dataset == nil {
		return true
	}

	return len(*cm.dataset) == 0
}

// IsNotEmpty checks if the currency map is present and has elements.
func (cm Map) IsNotEmpty() bool {
	return !cm.IsEmpty()
}

// HasInvalidState checks if the map is nil or contains nil pointers.
func (cm Map) HasInvalidState() (bool, error) {
	if cm.dataset == nil {
		return true, exception.ErrNoCurrencyMapDataset
	}

	for code, currencyPtr := range *cm.dataset {
		if currencyPtr == nil {
			return true, fmt.Errorf("currency [%s] is nil", code)
		}
	}

	return false, nil
}

// FindByCode case-insensitively searches for a currency by its code.
func (cm Map) FindByCode(code string) *Currency {
	if cm.dataset == nil {
		return nil
	}

	lookup := cm.dataset

	if result, ok := (*lookup)[strings.ToUpper(code)]; ok {
		return result
	}

	return nil
}

// GetCodes returns a list of all currency codes present in the map.
func (cm Map) GetCodes() *[]string {
	var once sync.Once
	codes := make([]string, 0, len(*cm.dataset))

	once.Do(func() {
		for code := range *cm.dataset {
			codes = append(codes, code)
		}
	})

	return &codes
}
