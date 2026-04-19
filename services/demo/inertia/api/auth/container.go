package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
)

// Container contains the host application integrations required by the auth package.
type Container struct {
	DB           *sql.DB
	CryptoKey    []byte
	Render       func(http.ResponseWriter, *http.Request, string, protocol.Props)
	Redirect     func(http.ResponseWriter, *http.Request, string)
	RouteURL     func(string, map[string]string) string
	SetFlash     func(http.ResponseWriter, flash.Message) error
	SecureCookie bool
}

// Validate checks that all required dependencies are set.
func (c Container) Validate() error {
	var errs []error

	if c.DB == nil {
		errs = append(errs, errors.New("auth: DB must not be nil"))
	}

	if len(c.CryptoKey) == 0 {
		errs = append(errs, errors.New("auth: CryptoKey must not be empty"))
	}

	if c.Render == nil {
		errs = append(errs, errors.New("auth: Render must not be nil"))
	}

	if c.Redirect == nil {
		errs = append(errs, errors.New("auth: Redirect must not be nil"))
	}

	if c.RouteURL == nil {
		errs = append(errs, errors.New("auth: RouteURL must not be nil"))
	}

	if c.SetFlash == nil {
		errs = append(errs, errors.New("auth: SetFlash must not be nil"))
	}

	return errors.Join(errs...)
}
