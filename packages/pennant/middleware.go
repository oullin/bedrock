package pennant

import (
	"fmt"
	"net/http"
	"strings"
)

// InactiveFeatureResponder handles requests blocked by feature middleware.
type InactiveFeatureResponder func(http.ResponseWriter, *http.Request, []string)

// EnsureFeaturesAreActive adapts Pennant checks to net/http middleware.
type EnsureFeaturesAreActive struct {
	features     *ScopedFeatureInteraction
	whenInactive InactiveFeatureResponder
}

// NewEnsureFeaturesAreActive creates middleware backed by the given scoped
// feature interaction.
func NewEnsureFeaturesAreActive(features *ScopedFeatureInteraction) *EnsureFeaturesAreActive {
	return &EnsureFeaturesAreActive{features: features}
}

// WhenInactive configures a custom response for inactive or undefined
// features. Passing nil restores the default 404 response.
func (m *EnsureFeaturesAreActive) WhenInactive(responder InactiveFeatureResponder) {
	m.whenInactive = responder
}

// Handle wraps next and only calls it when all named features are active.
func (m *EnsureFeaturesAreActive) Handle(next http.Handler, features ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.features != nil && m.features.AllAreActive(r.Context(), features) {
			next.ServeHTTP(w, r)

			return
		}

		if m.whenInactive != nil {
			m.whenInactive(w, r, features)

			return
		}

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
}

// Using returns a compact middleware descriptor suitable for route metadata.
func (m *EnsureFeaturesAreActive) Using(features ...string) string {
	return fmt.Sprintf("pennant.EnsureFeaturesAreActive:%s", strings.Join(features, ","))
}
