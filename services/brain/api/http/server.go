// Package http hosts brain's HTTP layer. Routes mirror upstream-brain's
// routes/brain.php under the `_request_cycle` prefix:
//
//	GET  /_request_cycle/api/manifest                     manifest JSON
//	GET  /_request_cycle/api/graph                        full graph JSON
//	GET  /_request_cycle/api/source?path=...              read a source file
//	POST /_request_cycle/api/scan                         re-run the analyzer pipeline
//	GET  /_request_cycle/api/context                      AI context export
//	POST /_request_cycle/api/generate-rules               write editor rules files
//	POST /_request_cycle/api/stress-test                  enqueue a load test
//	GET  /_request_cycle/api/stress-test/{jobID}          poll a load test
//	GET  /_request_cycle/{any}                            SPA shell
//
// Dispatch runs through packages/httpx/routingx so handler results and errors
// flow through the same primitives as services/demo and packages/billing.
// Asset serving lives in cmd/brain so the embed.FS scope stays at the
// binary's root.
package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bedrock/packages/filesystem"
	"github.com/bedrock/packages/httpx"
	"github.com/bedrock/packages/httpx/routingx"
	"github.com/bedrock/packages/routing"
	"github.com/bedrock/services/brain/api/ai"
	"github.com/bedrock/services/brain/api/analysis"
	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/stress"
)

// Server holds the most recently scanned graph plus the analyzer to run on
// demand from POST /api/scan.
type Server struct {
	Target   string
	Analyzer *analysis.ProjectAnalyzer
	HTMLBody []byte
	AssetDir string

	mu       sync.RWMutex
	graph    *graph.Graph
	manifest graph.Manifest
}

type generateRulesRequest struct {
	Force bool `json:"force"`
}

type stressRequest struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	Concurrency int    `json:"concurrency"`
	Requests    int    `json:"requests"`
	TimeoutMs   int    `json:"timeoutMs"`
}

const (
	routePrefix = "/_request_cycle"
	assetPrefix = routePrefix + "/assets/"
)

func NewServer(target string, html []byte, assetDir string) *Server {
	return &Server{
		Target:   target,
		Analyzer: analysis.NewDefaultProjectAnalyzer(),
		HTMLBody: html,
		AssetDir: assetDir,
	}
}

func (s *Server) Routes() http.Handler {
	router := routing.NewRouter(nil, nil)

	router.Group(map[string]any{"prefix": routePrefix}, func(r *routing.Router) {
		r.Get("/api/manifest", s.handleManifest)
		r.Get("/api/graph", s.handleGraph)
		r.Get("/api/source", s.handleSource)
		r.Post("/api/scan", s.handleScan)
		r.Get("/api/context", s.handleContext)
		r.Post("/api/generate-rules", s.handleGenerateRules)
		r.Post("/api/stress-test", s.handleStressTestEnqueue)
		r.Get("/api/stress-test/{jobID}", s.handleStressTestPoll)
	})

	router.Fallback(s.handleSPA)

	return s.withAssetServer(routingx.NewHandler(router))
}

func (s *Server) withAssetServer(next http.Handler) http.Handler {
	if s.AssetDir == "" {
		return next
	}

	fileServer := http.StripPrefix(assetPrefix, http.FileServer(http.Dir(s.AssetDir)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, assetPrefix) {
			fileServer.ServeHTTP(w, r)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) EnsureScanned() error {
	s.mu.RLock()
	loaded := s.graph != nil
	s.mu.RUnlock()

	if loaded {
		return nil
	}

	return s.rescan(context.Background())
}

func (s *Server) rescan(_ context.Context) error {
	g, err := s.Analyzer.AnalyzeTarget(s.Target)

	if err != nil {
		return err
	}

	s.mu.Lock()
	s.graph = g
	s.manifest = graph.Manifest{
		Project:     g.Meta.Project,
		AnalyzedAt:  g.Meta.AnalyzedAt,
		TotalNodes:  g.Meta.NodeCount,
		TotalEdges:  g.Meta.EdgeCount,
		TotalRoutes: countRoutes(g),
		Tabs:        []graph.TabEntry{},
	}
	s.mu.Unlock()

	return nil
}

func countRoutes(g *graph.Graph) int {
	c := 0

	for _, n := range g.Nodes {
		if n.Type == graph.NodeTypeRoute {
			c++
		}
	}

	return c
}

func (s *Server) handleManifest(w http.ResponseWriter, r *http.Request) {
	if err := s.EnsureScanned(); err != nil {
		writeErr(w, http.StatusInternalServerError, err)

		return
	}

	s.mu.RLock()
	m := s.manifest
	s.mu.RUnlock()
	writeJSON(w, http.StatusOK, m)
}

func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
	if err := s.EnsureScanned(); err != nil {
		writeErr(w, http.StatusInternalServerError, err)

		return
	}

	s.mu.RLock()
	g := s.graph
	s.mu.RUnlock()
	writeJSON(w, http.StatusOK, g)
}

func (s *Server) handleSource(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("path")

	if p == "" {
		writeErr(w, http.StatusBadRequest, errors.New("path query param required"))

		return
	}

	abs, err := filepath.Abs(filepath.Join(s.Target, p))

	if err != nil || !strings.HasPrefix(abs, s.Target) {
		writeErr(w, http.StatusBadRequest, errors.New("path outside target"))

		return
	}

	body, err := filesystem.New().Get(abs)

	if err != nil {
		writeErr(w, http.StatusNotFound, err)

		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"path": p, "body": string(body)})
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	if err := s.rescan(r.Context()); err != nil {
		writeErr(w, http.StatusInternalServerError, err)

		return
	}

	s.handleManifest(w, r)
}

func (s *Server) handleContext(w http.ResponseWriter, r *http.Request) {
	if err := s.EnsureScanned(); err != nil {
		writeErr(w, http.StatusInternalServerError, err)

		return
	}

	s.mu.RLock()
	g := s.graph
	s.mu.RUnlock()

	includeData := r.URL.Query().Get("data") == "1"
	body := ai.RenderMarkdown(g, ai.ContextOptions{IncludeData: includeData})

	switch r.URL.Query().Get("format") {
	case "json":
		writeJSON(w, http.StatusOK, map[string]any{
			"markdown":      body,
			"tokenEstimate": ai.EstimateTokens(body),
		})
	default:
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		_, _ = w.Write([]byte(body))
	}
}

func (s *Server) handleGenerateRules(w http.ResponseWriter, r *http.Request) {
	if err := s.EnsureScanned(); err != nil {
		writeErr(w, http.StatusInternalServerError, err)

		return
	}

	var req generateRulesRequest

	if r.ContentLength > 0 {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	s.mu.RLock()
	g := s.graph
	s.mu.RUnlock()
	body := ai.RenderMarkdown(g, ai.ContextOptions{IncludeData: false})
	written, err := ai.GenerateRules(s.Target, body, req.Force)

	if err != nil {
		if err == ai.ErrConflict {
			writeErr(w, http.StatusConflict, err)

			return
		}

		writeErr(w, http.StatusInternalServerError, err)

		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"written": written})
}

func (s *Server) handleStressTestEnqueue(w http.ResponseWriter, r *http.Request) {
	var req stressRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err)

		return
	}

	if req.URL == "" {
		writeErr(w, http.StatusBadRequest, errors.New("url is required"))

		return
	}
	// Refuse non-localhost targets — protects against accidental load
	// generation against shared infra.
	if !isLocalURL(req.URL) {
		writeErr(w, http.StatusForbidden,
			errors.New("brain only stresses localhost URLs; pass an explicit allowlist via config in a later release"))

		return
	}

	cfg := stress.Config{
		URL: req.URL, Method: req.Method,
		Concurrency: req.Concurrency, Requests: req.Requests,
	}

	if req.TimeoutMs > 0 {
		cfg.Timeout = time.Duration(req.TimeoutMs) * time.Millisecond
	}

	result, err := stress.Run(r.Context(), cfg)

	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)

		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleStressTestPoll(w http.ResponseWriter, r *http.Request) {
	// The runner is synchronous in phase 11; polling is not yet useful.
	// Reserved for an async/queued mode in a future release.
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error": "stress runs are synchronous in this release; poll endpoint reserved",
	})
}

func isLocalURL(raw string) bool {
	u, err := url.Parse(raw)

	if err != nil {
		return false
	}

	host := u.Hostname()

	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

func (s *Server) handleSPA(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if len(s.HTMLBody) == 0 {
		_, _ = w.Write([]byte(`<!doctype html><meta charset=utf8><title>brain</title>
<h1>brain</h1><p>SPA build missing. Run <code>pnpm --filter brain-app build</code>.</p>`))

		return
	}

	_, _ = w.Write(s.HTMLBody)
}

// writeJSON renders v as JSON through httpx so brain shares the same response
// primitives as services/demo and packages/billing.
func writeJSON(w http.ResponseWriter, status int, v any) {
	_ = httpx.NewJsonResponse(w, v, status, httpx.JsonOptions{Indent: true}).Send()
}

func writeErr(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

// EscapePath is exposed for tests that build URLs around tricky source paths.
func EscapePath(p string) string { return url.PathEscape(p) }
