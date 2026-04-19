package crm

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bedrock/packages/inertia/flash"
	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/services/demo/inertia/api/internal/database"
)

// Container contains the host application integrations required by the CRM package.
type Container struct {
	DB          *sql.DB
	RequireAuth func(http.Handler) http.Handler
	Render      func(http.ResponseWriter, *http.Request, string, protocol.Props)
	Redirect    func(http.ResponseWriter, *http.Request, string)
	RouteURL    func(string, map[string]string) string
	SetFlash    func(http.ResponseWriter, flash.Message) error
	CurrentUser func(*http.Request) *database.User
}

// Validate checks that all required dependencies are set.
func (c Container) Validate() error {
	var errs []error

	if c.DB == nil {
		errs = append(errs, errors.New("crm: DB must not be nil"))
	}

	if c.RequireAuth == nil {
		errs = append(errs, errors.New("crm: RequireAuth must not be nil"))
	}

	if c.Render == nil {
		errs = append(errs, errors.New("crm: Render must not be nil"))
	}

	if c.Redirect == nil {
		errs = append(errs, errors.New("crm: Redirect must not be nil"))
	}

	if c.RouteURL == nil {
		errs = append(errs, errors.New("crm: RouteURL must not be nil"))
	}

	if c.SetFlash == nil {
		errs = append(errs, errors.New("crm: SetFlash must not be nil"))
	}

	if c.CurrentUser == nil {
		errs = append(errs, errors.New("crm: CurrentUser must not be nil"))
	}

	return errors.Join(errs...)
}
