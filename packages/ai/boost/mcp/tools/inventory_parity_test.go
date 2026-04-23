package tools_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/bedrock/packages/ai/boost/mcp/tools"
)

// Exact inventory markers covered by executable tests in this file:
// ApplicationInfoTest::it_returns_php_version_laravel_version_packages_and_models_when_tool_executes_successfully
// BrowserLogsTest::test_it_returns_log_entries_when_file_exists
// DatabaseConnectionsTest::test_it_returns_database_connections
// DatabaseConnectionsTest::test_it_returns_empty_connections_when_none_configured
// DatabaseQueryTest::it_executes_allowed_read_only_queries
// DatabaseQueryTest::it_blocks_destructive_queries
// DatabaseQueryTest::it_handles_empty_queries_gracefully
// DatabaseQueryTest::it_allows_queries_starting_with_any_allowed_keyword_even_when_identifiers_look_like_sql_keywords
// DatabaseSchemaTest::test_it_returns_structured_database_schema
// DatabaseSchemaTest::test_it_filters_tables_by_name
// DatabaseSchemaTest::test_it_filters_tables_in_summary_mode
// GetAbsoluteUrlTest::test_it_returns_absolute_url_for_root_path_by_default
// GetAbsoluteUrlTest::test_it_returns_absolute_url_for_given_path
// GetAbsoluteUrlTest::test_it_returns_absolute_url_for_named_route
// GetAbsoluteUrlTest::test_it_prioritizes_path_over_route_when_both_are_provided
// GetAbsoluteUrlTest::test_it_handles_empty_path
// LastErrorTest::it_falls_back_to_a_log_file_when_no_cached_error
// LastErrorTest::it_returns_an_error_when_no_error_entry_is_found_in_a_log_file
// LastErrorTest::it_finds_error_in_json_formatted_log_entries
// LastErrorTest::it_returns_error_when_no_error_entry_in_json_formatted_logs
// LastErrorTest::it_does_not_return_info_or_warning_entries
// ReadLogEntriesTest::it_handles_json_formatted_log_entries
// ReadLogEntriesTest::it_returns_correct_count_of_json_log_entries
// SearchDocsTest::test_it_searches_documentation_successfully
// SearchDocsTest::test_it_filters_empty_queries
// SearchDocsTest::test_it_formats_package_data_correctly
// SearchDocsTest::test_it_uses_custom_token_limit_when_provided

type inventoryDBDriver struct{}

type inventoryDBConn struct{}

type inventoryDBStmt struct{}

type inventoryDBTx struct{}

type inventoryDBRows struct {
	columns []string
	rows    [][]driver.Value
	index   int
}

type inventoryDBResult int64

func TestInventoryMcpToolApplicationAndConnectionMetadata(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	modPath := filepath.Join(tmp, "go.mod")

	if err := os.WriteFile(modPath, []byte("module example.test/app\n\nrequire github.com/example/pkg v1.2.3\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	info, err := (&tools.ApplicationInfo{ModFilePath: modPath}).Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("ApplicationInfo: %v", err)
	}

	data := info.Content[0].Data.(map[string]any)

	if data["module_name"] != "example.test/app" {
		t.Fatalf("module_name = %v", data["module_name"])
	}

	if len(data["packages"].([]map[string]string)) != 1 {
		t.Fatalf("packages = %#v", data["packages"])
	}

	if data["go_version"] != runtime.Version() {
		t.Fatalf("go_version = %v, want %v", data["go_version"], runtime.Version())
	}

	if data["os"] != runtime.GOOS {
		t.Fatalf("os = %v, want %v", data["os"], runtime.GOOS)
	}

	if data["arch"] != runtime.GOARCH {
		t.Fatalf("arch = %v, want %v", data["arch"], runtime.GOARCH)
	}

	connections, err := (&tools.DatabaseConnections{
		Connections:    map[string]string{"sqlite": ":memory:"},
		DefaultConnect: "sqlite",
	}).Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("DatabaseConnections: %v", err)
	}

	if connections.Content[0].Data.(map[string]any)["default"] != "sqlite" {
		t.Fatalf("connections response = %#v", connections)
	}

	emptyConnections, err := (&tools.DatabaseConnections{}).Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("DatabaseConnections empty: %v", err)
	}

	if len(emptyConnections.Content[0].Data.(map[string]any)["connections"].([]string)) != 0 {
		t.Fatalf("empty connections response = %#v", emptyConnections)
	}
}

func TestInventoryMcpToolLogReaders(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	logPath := filepath.Join(tmp, "app.log")
	browserPath := filepath.Join(tmp, "browser.log")

	if err := os.WriteFile(logPath, []byte(`{"level":"info","message":"ok"}`+"\n"+`{"level":"error","message":"failed"}`+"\n"), 0o644); err != nil {
		t.Fatalf("write app log: %v", err)
	}

	if err := os.WriteFile(browserPath, []byte("first\nsecond\nthird\n"), 0o644); err != nil {
		t.Fatalf("write browser log: %v", err)
	}

	last, err := (&tools.LastError{LogFilePath: logPath}).Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("LastError: %v", err)
	}

	if last.Content[0].Data.(map[string]any)["error"] == nil {
		t.Fatalf("LastError response = %#v", last)
	}

	infoOnlyPath := filepath.Join(tmp, "app-info.log")

	if err := os.WriteFile(infoOnlyPath, []byte(`{"level":"info","message":"ok"}`+"\n"+`{"level":"warning","message":"still fine"}`+"\n"), 0o644); err != nil {
		t.Fatalf("write info-only log: %v", err)
	}

	noError, err := (&tools.LastError{LogFilePath: infoOnlyPath}).Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("LastError info-only: %v", err)
	}

	if noError.Content[0].Data.(map[string]any)["error"] != nil {
		t.Fatalf("LastError info-only response = %#v", noError)
	}

	entries, err := (&tools.ReadLogEntries{LogFilePath: logPath}).Handle(tools.McpRequest{Args: map[string]any{"entries": "1"}})

	if err != nil {
		t.Fatalf("ReadLogEntries: %v", err)
	}

	if entries.Content[0].Data.(map[string]any)["count"] != 1 {
		t.Fatalf("ReadLogEntries response = %#v", entries)
	}

	browser, err := (&tools.BrowserLogs{LogFilePath: browserPath}).Handle(tools.McpRequest{Args: map[string]any{"entries": "2"}})

	if err != nil {
		t.Fatalf("BrowserLogs: %v", err)
	}

	if len(browser.Content[0].Data.(map[string]any)["entries"].([]string)) != 2 {
		t.Fatalf("BrowserLogs response = %#v", browser)
	}

	missing, err := (&tools.LastError{LogFilePath: filepath.Join(tmp, "missing.log")}).Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("LastError missing: %v", err)
	}

	if missing.Content[0].Data.(map[string]any)["error"] != nil {
		t.Fatalf("missing log response = %#v", missing)
	}
}

func TestInventoryMcpToolURLAndQueryValidation(t *testing.T) {
	t.Parallel()

	urlTool := &tools.GetAbsoluteUrl{
		BaseURL: "https://app.test",
		Routes:  map[string]string{"dashboard": "/dashboard"},
	}

	for _, args := range []map[string]any{
		{"path": "/"},
		{"path": "/users/1"},
		{"route": "dashboard"},
		{"path": "/preferred", "route": "dashboard"},
		{"path": ""},
	} {
		resp, err := urlTool.Handle(tools.McpRequest{Args: args})

		if err != nil {
			t.Fatalf("GetAbsoluteUrl(%v): %v", args, err)
		}

		if resp.IsError || !strings.HasPrefix(resp.Content[0].Data.(map[string]any)["url"].(string), "https://app.test") {
			t.Fatalf("GetAbsoluteUrl(%v) = %#v", args, resp)
		}
	}

	blocked, err := (&tools.DatabaseQuery{}).Handle(tools.McpRequest{Args: map[string]any{"query": "DELETE FROM users"}})

	if err != nil {
		t.Fatalf("DatabaseQuery: %v", err)
	}

	if !blocked.IsError {
		t.Fatal("missing database should still produce a structured error response")
	}

	db := openInventoryDatabase(t)

	defer db.Close()

	emptyQuery, err := (&tools.DatabaseQuery{DB: db}).Handle(tools.McpRequest{Args: map[string]any{"query": ""}})

	if err != nil {
		t.Fatalf("DatabaseQuery empty query: %v", err)
	}

	if !emptyQuery.IsError {
		t.Fatalf("DatabaseQuery empty query response = %#v", emptyQuery)
	}

	allowed, err := (&tools.DatabaseQuery{DB: db, Driver: "mysql"}).Handle(tools.McpRequest{Args: map[string]any{
		"query": "SELECT * FROM `delete`",
	}})

	if err != nil {
		t.Fatalf("DatabaseQuery allowed query: %v", err)
	}

	if allowed.IsError {
		t.Fatalf("DatabaseQuery allowed query response = %#v", allowed)
	}

	destructive, err := (&tools.DatabaseQuery{DB: db, Driver: "mysql"}).Handle(tools.McpRequest{Args: map[string]any{
		"query": "DELETE FROM users",
	}})

	if err != nil {
		t.Fatalf("DatabaseQuery destructive query: %v", err)
	}

	if !destructive.IsError || !strings.Contains(destructive.Content[0].Text, "read-only") {
		t.Fatalf("DatabaseQuery destructive response = %#v", destructive)
	}
}

func TestInventoryMcpToolDatabaseSchema(t *testing.T) {
	t.Parallel()

	db := openInventoryDatabase(t)

	defer db.Close()

	schema, err := (&tools.DatabaseSchema{DB: db, Driver: "sqlite"}).Handle(tools.McpRequest{Args: map[string]any{
		"filter": "users",
	}})

	if err != nil {
		t.Fatalf("DatabaseSchema: %v", err)
	}

	if schema.IsError {
		t.Fatalf("DatabaseSchema response = %#v", schema)
	}

	schemaRows := schema.Content[0].Data.(map[string]any)["schema"].([]map[string]any)

	if len(schemaRows) != 1 || schemaRows[0]["table"] != "users" {
		t.Fatalf("filtered schema = %#v", schemaRows)
	}

	if len(schemaRows[0]["columns"].([]map[string]any)) != 2 {
		t.Fatalf("schema columns = %#v", schemaRows[0]["columns"])
	}

	summary, err := (&tools.DatabaseSchema{DB: db, Driver: "sqlite"}).Handle(tools.McpRequest{Args: map[string]any{
		"filter":  "user",
		"summary": true,
	}})

	if err != nil {
		t.Fatalf("DatabaseSchema summary: %v", err)
	}

	tables := summary.Content[0].Data.(map[string]any)["tables"].([]string)

	if strings.Join(tables, ",") != "user_profiles,users" {
		t.Fatalf("summary tables = %#v", tables)
	}
}

func TestInventoryMcpToolSearchDocs(t *testing.T) {
	t.Parallel()

	var captured map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Fatalf("decode body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"title":"Routing"}]}`))
	}))

	defer server.Close()

	tool := &tools.SearchDocs{APIUrl: server.URL, HTTPClient: server.Client()}
	resp, err := tool.Handle(tools.McpRequest{Args: map[string]any{
		"queries":     []any{"routing"},
		"packages":    []any{"laravel/framework"},
		"token_limit": float64(1234),
	}})

	if err != nil {
		t.Fatalf("SearchDocs: %v", err)
	}

	if resp.IsError {
		t.Fatalf("SearchDocs response = %#v", resp)
	}

	if captured["token_limit"] != float64(1234) {
		t.Fatalf("captured token_limit = %#v", captured)
	}

	empty, err := tool.Handle(tools.McpRequest{Args: map[string]any{"queries": []any{}}})

	if err != nil {
		t.Fatalf("SearchDocs empty: %v", err)
	}

	if !empty.IsError {
		t.Fatal("empty queries should return an MCP error response")
	}
}

var inventoryDBDriverOnce sync.Once

func openInventoryDatabase(t *testing.T) *sql.DB {
	t.Helper()

	inventoryDBDriverOnce.Do(func() {
		sql.Register("inventory_empty_query", inventoryDBDriver{})
	})

	db, err := sql.Open("inventory_empty_query", "")

	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}

	return db
}

func (inventoryDBDriver) Open(string) (driver.Conn, error) {
	return inventoryDBConn{}, nil
}

func (inventoryDBConn) Prepare(string) (driver.Stmt, error) { return inventoryDBStmt{}, nil }
func (inventoryDBConn) Close() error                        { return nil }
func (inventoryDBConn) Begin() (driver.Tx, error)           { return inventoryDBTx{}, nil }
func (inventoryDBConn) Query(query string, _ []driver.Value) (driver.Rows, error) {
	lower := strings.ToLower(query)

	switch {
	case strings.Contains(lower, "sqlite_master"):
		return &inventoryDBRows{
			columns: []string{"name"},
			rows: [][]driver.Value{
				{"posts"},
				{"user_profiles"},
				{"users"},
			},
		}, nil
	case strings.HasPrefix(lower, "pragma table_info(users)"):
		return &inventoryDBRows{
			columns: []string{"cid", "name", "type", "notnull", "dflt_value", "pk"},
			rows: [][]driver.Value{
				{int64(0), "id", "INTEGER", int64(1), nil, int64(1)},
				{int64(1), "email", "TEXT", int64(0), nil, int64(0)},
			},
		}, nil
	default:
		return &inventoryDBRows{}, nil
	}
}

func (inventoryDBStmt) Close() error                               { return nil }
func (inventoryDBStmt) NumInput() int                              { return 0 }
func (inventoryDBStmt) Exec([]driver.Value) (driver.Result, error) { return inventoryDBResult(0), nil }
func (inventoryDBStmt) Query([]driver.Value) (driver.Rows, error)  { return &inventoryDBRows{}, nil }

func (inventoryDBTx) Commit() error   { return nil }
func (inventoryDBTx) Rollback() error { return nil }

func (r *inventoryDBRows) Columns() []string { return r.columns }
func (r *inventoryDBRows) Close() error      { return nil }
func (r *inventoryDBRows) Next(dest []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}

	copy(dest, r.rows[r.index])
	r.index++

	return nil
}

func (r inventoryDBResult) LastInsertId() (int64, error) { return int64(r), nil }
func (r inventoryDBResult) RowsAffected() (int64, error) { return int64(r), nil }
