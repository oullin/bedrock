// Package http hosts brain's HTTP layer. Routes mirror laravel-brain's
// routes/brain.php under the `_brain` prefix:
//
//	GET  /_brain/api/manifest                     manifest JSON
//	GET  /_brain/api/graph                        full graph JSON
//	GET  /_brain/api/source?path=...              read a source file
//	POST /_brain/api/scan                         re-run the analyzer pipeline
//	GET  /_brain/api/context                      AI context export
//	POST /_brain/api/generate-rules               write editor rules files
//	POST /_brain/api/stress-test                  enqueue a load test
//	GET  /_brain/api/stress-test/{jobID}          poll a load test
//	GET  /_brain/{any}                            SPA shell
//
// Asset serving lives in cmd/brain so the embed.FS scope stays at the
// binary's root.
package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bedrock/services/brain/api/analysis"
	"github.com/bedrock/services/brain/api/graph"
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

// NewServer wires a server with the canonical analyzer pipeline.
func NewServer(target string, html []byte, assetDir string) *Server {
	return &Server{
		Target:   target,
		Analyzer: analysis.NewDefaultProjectAnalyzer(),
		HTMLBody: html,
		AssetDir: assetDir,
	}
}

// Routes returns a configured handler. We deliberately use net/http directly
// only in this file (and cmd/brain) — analyzers must stay off the HTTP path.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /_brain/api/manifest", s.handleManifest)
	mux.HandleFunc("GET /_brain/api/graph", s.handleGraph)
	mux.HandleFunc("GET /_brain/api/source", s.handleSource)
	mux.HandleFunc("POST /_brain/api/scan", s.handleScan)
	mux.HandleFunc("GET /_brain/api/context", s.handleContext)
	mux.HandleFunc("POST /_brain/api/generate-rules", s.handleGenerateRules)
	mux.HandleFunc("POST /_brain/api/stress-test", s.handleStressTestEnqueue)
	mux.HandleFunc("GET /_brain/api/stress-test/{jobID}", s.handleStressTestPoll)
	if s.AssetDir != "" {
		mux.Handle("GET /_brain/assets/", http.StripPrefix("/_brain/assets/",
			http.FileServer(http.Dir(s.AssetDir))))
	}
	mux.HandleFunc("GET /_brain/", s.handleSPA)
	mux.HandleFunc("GET /", s.handleSPA)
	return mux
}

// EnsureScanned makes sure a graph is loaded before any read endpoint runs.
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
	// disallow escapes out of target
	abs, err := filepath.Abs(filepath.Join(s.Target, p))
	if err != nil || !strings.HasPrefix(abs, s.Target) {
		writeErr(w, http.StatusBadRequest, errors.New("path outside target"))
		return
	}
	body, err := os.ReadFile(abs)
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
	// Phase 12 will replace this with packages/ai/sdk + packages/ai/boost.
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error": "AI context export wired in phase 12",
	})
}

func (s *Server) handleGenerateRules(w http.ResponseWriter, r *http.Request) {
	// Phase 12 wires this to packages/ai/boost.
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error": "rules generator wired in phase 12",
	})
}

func (s *Server) handleStressTestEnqueue(w http.ResponseWriter, r *http.Request) {
	// Phase 11 wires this to the stress runner.
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error": "stress runner wired in phase 11",
	})
}

func (s *Server) handleStressTestPoll(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error": "stress runner wired in phase 11",
	})
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func writeErr(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

// EscapePath is exposed for tests that build URLs around tricky source paths.
func EscapePath(p string) string { return url.PathEscape(p) }
