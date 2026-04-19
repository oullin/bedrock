//go:build !cgo

package tools

import "fmt"

func validatePostgresReadOnlySQL(_ string) error {
	return fmt.Errorf("postgres read-only validation requires CGO")
}
