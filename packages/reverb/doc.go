// Reverb speaks the Pusher protocol, making it a drop-in replacement for
// Pusher in any upstream application. Clients connect over WebSocket using the
// standard Pusher client libraries (or packages/echo), authenticate via
// HMAC-SHA256, and receive real-time events broadcast from the server.
//
// # Architecture
//
// The package is composed of these primary concerns:
//
//   - Application management: each App holds credentials and per-app limits.
//   - Connection management: each live WebSocket is tracked as a Conn.
//   - Channel management: subscriptions are tracked per channel per app.
//   - Protocol: Pusher-protocol message parsing and serialisation.
//   - Authentication: HMAC-SHA256 signature verification for private/presence channels.
//   - HTTP API: REST endpoints for triggering events from upstream backends.
//   - Dispatcher: in-process or Redis-backed broadcast routing.
//
// # Channel types
//
// The ChannelBroker (NewChannel) routes by name prefix:
//
//	private-cache-*    → PrivateCacheChannel
//	presence-cache-*   → PresenceCacheChannel
//	cache-*            → CacheChannel
//	private-*          → PrivateChannel
//	presence-*         → PresenceChannel
//	<anything else>    → Channel (public)
//
// # Quick start
//
//	cfg := reverb.DefaultConfig()
//	cfg.Apps = []reverb.AppConfig{{
//	    ID: "app-1", Key: "key-1", Secret: "secret-1",
//	    AllowedOrigins: []string{"*"},
//	}}
//
//	apps := reverb.NewAppManager(cfg.Apps)
//	conns := reverb.NewConnectionManager()
//	channels := reverb.NewChannelManager(apps)
//	dispatcher := reverb.NewSyncDispatcher(channels)
//
//	ws := reverb.NewServer(cfg, apps, conns, channels, dispatcher)
//	api := reverb.NewHTTPHandler(apps, channels, dispatcher)
//
//	mux := http.NewServeMux()
//	mux.Handle("/app/", ws)
//	mux.Handle("/apps/", api)
//
//	http.ListenAndServe(":8080", mux)
package reverb
