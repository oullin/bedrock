package features

import (
	"fmt"
	"net/http"

	"github.com/bedrock/packages/inertia/routegen"
)

// DefineRoutes registers all feature showcase route metadata (name, method,
// pattern) on the given registry without mounting handlers.
func DefineRoutes(routes *routegen.Registry) {
	routes.Group("features.forms", "/features/forms", func(g *routegen.Group) {
		g.Add("use-form", "GET", "/use-form")
		g.Add("form-component", "GET", "/form-component")
		g.Add("file-uploads", "GET", "/file-uploads")
		g.Add("validation", "GET", "/validation")
		g.Add("httppreview", "GET", "/httppreview")
		g.Add("optimistic-updates", "GET", "/optimistic-updates")
		g.Add("use-form-context", "GET", "/use-form-context")
		g.Add("dotted-keys", "GET", "/dotted-keys")
		g.Add("routegen", "GET", "/routegen")
	})

	routes.Group("features.navigation", "/features/navigation", func(g *routegen.Group) {
		g.Add("links", "GET", "/links")
		g.Add("preserve-state", "GET", "/preserve-state")
		g.Add("preserve-scroll", "GET", "/preserve-scroll")
		g.Add("view-transitions", "GET", "/view-transitions")
		g.Add("history-management", "GET", "/history-management")
		g.Add("async-requests", "GET", "/async-requests")
		g.Add("async-slow", "GET", "/async-slow")
		g.Add("manual-visits", "GET", "/manual-visits")
		g.Add("redirects", "GET", "/redirects")
		g.Add("scroll-management", "GET", "/scroll-management")
		g.Add("instant-visits", "GET", "/instant-visits")
		g.Add("instant-visit-target", "GET", "/instant-visit-target")
		g.Add("url-fragments", "GET", "/url-fragments")
	})

	routes.Group("features.data-loading", "/features/data-loading", func(g *routegen.Group) {
		g.Add("deferred-props", "GET", "/deferred-props")
		g.Add("partial-reloads", "GET", "/partial-reloads")
		g.Add("infinite-scroll", "GET", "/infinite-scroll")
		g.Add("when-visible", "GET", "/when-visible")
		g.Add("polling", "GET", "/polling")
		g.Add("prop-merging", "GET", "/prop-merging")
		g.Add("optional-props", "GET", "/optional-props")
		g.Add("once-props", "GET", "/once-props/{page}")
	})

	routes.Group("features.prefetching", "/features/prefetching", func(g *routegen.Group) {
		g.Add("link-prefetch", "GET", "/link-prefetch")
		g.Add("stale-while-revalidate", "GET", "/stale-while-revalidate")
		g.Add("manual-prefetch", "GET", "/manual-prefetch")
		g.Add("cache-management", "GET", "/cache-management")
	})

	routes.Group("features.state", "/features/state", func(g *routegen.Group) {
		g.Add("remember", "GET", "/remember")
		g.Add("flash-data", "GET", "/flash-data")
		g.Add("shared-props", "GET", "/shared-props")
	})

	routes.Group("features.layouts", "/features/layouts", func(g *routegen.Group) {
		g.Add("persistent-layouts", "GET", "/persistent-layouts")
		g.Add("persistent-layouts-page-2", "GET", "/persistent-layouts/page-2")
		g.Add("nested-layouts", "GET", "/nested-layouts")
		g.Add("head", "GET", "/head")
		g.Add("layout-props", "GET", "/layout-props")
	})

	routes.Group("features.events", "/features/events", func(g *routegen.Group) {
		g.Add("global-events", "GET", "/global-events")
		g.Add("visit-callbacks", "GET", "/visit-callbacks")
		g.Add("progress", "GET", "/progress")
		g.Add("progress-slow", "GET", "/progress/slow")
	})

	routes.Group("features.http", "/features/http", func(g *routegen.Group) {
		g.Add("use-http", "GET", "/use-http")
	})
}

// RegisterRoutes mounts all feature showcase HTTP routes onto the provided mux.
func RegisterRoutes(routes *routegen.Registry, mux *http.ServeMux, container Container) error {
	if err := container.Validate(); err != nil {
		return fmt.Errorf("features: %w", err)
	}

	a := newApp(container)

	auth := func(h http.HandlerFunc) http.Handler {
		return container.RequireAuth(h)
	}

	// Forms
	routes.Handle("features.forms.use-form", auth(a.useFormHandler), mux)
	routes.Handle("features.forms.form-component", auth(a.formComponentHandler), mux)
	routes.Handle("features.forms.file-uploads", auth(a.fileUploadsHandler), mux)
	routes.Handle("features.forms.validation", auth(a.validationHandler), mux)

	mux.Handle("/features/forms/validation/secondary", auth(a.validationSecondaryHandler))

	routes.Handle("features.forms.httppreview", auth(a.precognitionHandler), mux)
	routes.Handle("features.forms.optimistic-updates", auth(a.optimisticUpdatesHandler), mux)

	mux.Handle("POST /features/forms/optimistic-toggle/{id}", auth(a.optimisticToggleHandler))

	routes.Handle("features.forms.use-form-context", auth(a.formContextHandler), mux)
	routes.Handle("features.forms.dotted-keys", auth(a.dottedKeysHandler), mux)
	routes.Handle("features.forms.routegen", auth(a.wayfinderHandler), mux)

	// Navigation
	routes.Handle("features.navigation.links", auth(a.linksHandler), mux)
	routes.Handle("features.navigation.preserve-state", auth(a.preserveStateHandler), mux)
	routes.Handle("features.navigation.preserve-scroll", auth(a.preserveScrollHandler), mux)
	routes.Handle("features.navigation.view-transitions", auth(a.viewTransitionsHandler), mux)
	routes.Handle("features.navigation.history-management", auth(a.historyManagementHandler), mux)
	routes.Handle("features.navigation.async-requests", auth(a.asyncRequestsHandler), mux)
	routes.Handle("features.navigation.async-slow", auth(a.asyncSlowHandler), mux)
	routes.Handle("features.navigation.manual-visits", auth(a.manualVisitsHandler), mux)
	routes.Handle("features.navigation.redirects", auth(a.redirectsHandler), mux)

	mux.Handle("POST /features/navigation/redirects/{action}", auth(a.redirectsActionHandler))

	routes.Handle("features.navigation.scroll-management", auth(a.scrollManagementHandler), mux)
	routes.Handle("features.navigation.instant-visits", auth(a.instantVisitsHandler), mux)
	routes.Handle("features.navigation.instant-visit-target", auth(a.instantVisitTargetHandler), mux)
	routes.Handle("features.navigation.url-fragments", auth(a.urlFragmentsHandler), mux)

	mux.Handle("POST /features/navigation/url-fragments/{action}", auth(a.urlFragmentsActionHandler))

	// Data Loading
	routes.Handle("features.data-loading.deferred-props", auth(a.deferredPropsHandler), mux)
	routes.Handle("features.data-loading.partial-reloads", auth(a.partialReloadsHandler), mux)
	routes.Handle("features.data-loading.infinite-scroll", auth(a.infiniteScrollHandler), mux)
	routes.Handle("features.data-loading.when-visible", auth(a.whenVisibleHandler), mux)
	routes.Handle("features.data-loading.polling", auth(a.pollingHandler), mux)
	routes.Handle("features.data-loading.prop-merging", auth(a.propMergingHandler), mux)
	routes.Handle("features.data-loading.optional-props", auth(a.optionalPropsHandler), mux)
	routes.Handle("features.data-loading.once-props", auth(a.oncePropsHandler), mux)

	// Prefetching
	routes.Handle("features.prefetching.link-prefetch", auth(a.linkPrefetchHandler), mux)
	routes.Handle("features.prefetching.stale-while-revalidate", auth(a.staleWhileRevalidateHandler), mux)
	routes.Handle("features.prefetching.manual-prefetch", auth(a.manualPrefetchHandler), mux)
	routes.Handle("features.prefetching.cache-management", auth(a.cacheManagementHandler), mux)

	// State
	routes.Handle("features.state.remember", auth(a.rememberHandler), mux)
	routes.Handle("features.state.flash-data", auth(a.flashDataHandler), mux)

	mux.Handle("POST /features/state/flash-data/{action}", auth(a.flashDataActionHandler))

	routes.Handle("features.state.shared-props", auth(a.sharedPropsHandler), mux)

	// Layouts
	routes.Handle("features.layouts.persistent-layouts", auth(a.persistentLayoutsHandler), mux)
	routes.Handle("features.layouts.persistent-layouts-page-2", auth(a.persistentLayoutsPage2Handler), mux)
	routes.Handle("features.layouts.nested-layouts", auth(a.nestedLayoutsHandler), mux)
	routes.Handle("features.layouts.head", auth(a.headHandler), mux)
	routes.Handle("features.layouts.layout-props", auth(a.layoutPropsHandler), mux)

	// Events
	routes.Handle("features.events.global-events", auth(a.globalEventsHandler), mux)
	routes.Handle("features.events.visit-callbacks", auth(a.visitCallbacksHandler), mux)
	routes.Handle("features.events.progress", auth(a.progressHandler), mux)
	routes.Handle("features.events.progress-slow", auth(a.progressSlowHandler), mux)

	// HTTP
	routes.Handle("features.http.use-http", auth(a.useHttpHandler), mux)

	mux.Handle("/features/http/use-http/api", auth(a.useHttpApiHandler))

	return nil
}
