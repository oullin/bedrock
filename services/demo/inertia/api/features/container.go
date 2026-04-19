package features

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
)

// Container contains the host application integrations required by the features package.
type Container struct {
	DB          *sql.DB
	RequireAuth func(http.Handler) http.Handler
	Render      func(http.ResponseWriter, *http.Request, string, protocol.Props)
	Redirect    func(http.ResponseWriter, *http.Request, string)
	Location    func(http.ResponseWriter, *http.Request, string)
	RouteURL    func(string, map[string]string) string
	SetFlash    func(http.ResponseWriter, flash.Message) error
}

// Validate checks that all required dependencies are set.
func (c Container) Validate() error {
	var errs []error

	if c.DB == nil {
		errs = append(errs, errors.New("features: DB must not be nil"))
	}

	if c.RequireAuth == nil {
		errs = append(errs, errors.New("features: RequireAuth must not be nil"))
	}

	if c.Render == nil {
		errs = append(errs, errors.New("features: Render must not be nil"))
	}

	if c.Redirect == nil {
		errs = append(errs, errors.New("features: Redirect must not be nil"))
	}

	if c.Location == nil {
		errs = append(errs, errors.New("features: Location must not be nil"))
	}

	if c.RouteURL == nil {
		errs = append(errs, errors.New("features: RouteURL must not be nil"))
	}

	if c.SetFlash == nil {
		errs = append(errs, errors.New("features: SetFlash must not be nil"))
	}

	return errors.Join(errs...)
}
