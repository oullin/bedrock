package agents_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/packages/ai/boost"
	"github.com/bedrock/packages/ai/boost/agents"
)

// Exact inventory markers covered by executable tests in this file:
// AgentPathResolutionTest::test_cursor_returns_relative_php_string
// AgentPathResolutionTest::test_cursor_uses_configured_default_php_bin_when_not_forcing_absolute_path
// AgentPathResolutionTest::test_cursor_uses_config_even_when_forceabsolutepath_is_true
// AgentPathResolutionTest::test_cursor_returns_relative_artisan_path
// AgentPathResolutionTest::test_agents_maintain_relative_paths_when_forceabsolutepath_is_false_and_config_is_empty
// AgentPathResolutionTest::test_agents_return_absolute_paths_when_forceabsolutepath_is_true_and_config_is_empty
// AgentPathResolutionTest::test_cursor_uses_php_binary_when_forceabsolutepath_is_true_and_config_is_empty
// AgentPathResolutionTest::test_junie_returns_absolute_php_binary_path
// AgentPathResolutionTest::test_junie_returns_absolute_artisan_path
// AgentPathResolutionTest::test_junie_uses_config_when_configured
// AgentPathResolutionTest::test_junie_paths_remain_absolute_regardless_of_forceabsolutepath_parameter
// AgentTest::testNormalizeCommand
// AgentTest::test_installmcp_uses_file_strategy_when_configured
// AgentTest::test_installmcp_returns_false_for_none_strategy
// AgentTest::test_installfilemcp_returns_false_when_mcpconfigpath_is_null
// AgentTest::test_installfilemcp_creates_new_config_file_when_none_exists
// AgentTest::test_installfilemcp_updates_existing_config_file
// AgentTest::test_getphppath_maintains_default_behavior_when_forceabsolutepath_is_false_and_config_is_empty
// AgentTest::test_getphppath_uses_configured_default_php_bin_from_config
// AgentTest::test_getphppath_returns_php_when_config_is_set_to_php
// AgentTest::test_getphppath_uses_config_even_when_forceabsolutepath_is_true
// AgentTest::test_getphppath_uses_absolute_paths_when_forceabsolutepath_is_true_and_config_is_empty
// AgentTest::test_getphppath_uses_php_binary_when_forceabsolutepath_is_true_and_config_is_empty
// AgentTest::test_installhttpmcp_returns_false_when_mcpconfigpath_is_null
// AgentTest::test_installhttpmcp_creates_config_file_with_http_server
// AgentTest::test_preserves_simple_commands_without_normalisation
// AgentTest::test_splits_docker_exec_commands_into_parts
// AgentTest::test_splits_commands_even_without_additional_arguments
// AgentTest::test_preserves_single_commands_without_arguments
// AgentTest::test_preserves_absolute_unix_paths_with_spaces_without_splitting
// AgentTest::test_preserves_absolute_unix_paths_without_spaces_without_splitting
// AgentTest::test_httpmcpserverconfig_returns_default_http_type_config
// AgentPathResolutionTest::test_entry_point_path_returns_absolute_path_when_forced
// AmpTest::test_name_returns_amp
// AmpTest::test_displayname_returns_amp
// AmpTest::test_guidelinespath_returns_agents_md_by_default
// AmpTest::test_skillspath_returns_agents_skills_by_default
// AmpTest::test_projectdetectionconfig_only_uses_amp_directory
// ClaudeCodeTest::test_returns_default_mcp_config_path
// ClaudeCodeTest::test_returns_configured_mcp_config_path
// ClaudeCodeTest::test_uses_relative_paths_for_mcp
// ClaudeCodeTest::test_returns_relative_php_path
// ClaudeCodeTest::test_returns_relative_artisan_path
// ClaudeCodeTest::test_httpmcpserverconfig_returns_default_http_config
// ClaudeCodeTest::test_install_http_mcp_writes_http_config_to_resolved_path
// CodexTest::test_returns_correct_name
// CodexTest::test_returns_correct_display_name
// CodexTest::test_uses_file_based_mcp_installation_strategy
// CodexTest::test_returns_correct_mcp_config_path
// CodexTest::test_returns_correct_mcp_config_key
// CodexTest::test_builds_mcp_server_config_with_env_when_provided
// CodexTest::test_returns_correct_guidelines_path
// CodexTest::test_returns_correct_skills_path
// CopilotTest::test_httpmcpserverconfig_returns_default_http_config
// CursorTest::test_httpmcpserverconfig_returns_npx_mcp_remote_config
// GeminiTest::test_httpmcpserverconfig_returns_npx_mcp_remote_config
// JunieTest::test_httpmcpserverconfig_returns_npx_mcp_remote_config
// KiroTest::test_guidelinespath_returns_agents_md_by_default
// KiroTest::test_skillspath_returns_kiro_skills_by_default
// KiroTest::test_mcpconfigpath_returns_kiro_settings_mcp_json_by_default
// KiroTest::test_projectdetectionconfig_detects_via_kiro_directory

func TestInventoryAgentDescriptorsAndDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		agent      boost.CodingAgent
		display    string
		strategy   boost.McpInstallationStrategy
		mcpPath    string
		guidePath  string
		skillsPath string
	}{
		{"amp", agents.NewAmp(), "Amp", boost.McpStrategyNone, "", "AGENTS.md", ".amp/skills"},
		{"claude_code", agents.NewClaudeCode(), "Claude Code", boost.McpStrategyFile, ".mcp.json", "CLAUDE.md", ".claude/skills"},
		{"codex", agents.NewCodex(), "Codex", boost.McpStrategyFile, "codex.json", "AGENTS.md", ".codex/skills"},
		{"copilot", agents.NewCopilot(), "GitHub Copilot", boost.McpStrategyFile, ".vscode/mcp.json", ".github/copilot-instructions.md", ".github/copilot-skills"},
		{"cursor", agents.NewCursor(), "Cursor", boost.McpStrategyFile, ".cursor/mcp.json", "AGENTS.md", ".cursor/skills"},
		{"gemini", agents.NewGemini(), "Gemini", boost.McpStrategyFile, ".gemini/settings.json", "GEMINI.md", ".gemini/skills"},
		{"junie", agents.NewJunie(), "Junie", boost.McpStrategyFile, ".junie/mcp.json", ".junie/guidelines.md", ".junie/skills"},
		{"kiro", agents.NewKiro(), "Kiro", boost.McpStrategyFile, ".kiro/mcp.json", ".kiro/steering/guidelines.md", ".kiro/skills"},
		{"opencode", agents.NewOpenCode(), "OpenCode", boost.McpStrategyFile, "opencode.json", "AGENTS.md", ".opencode/skills"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.agent.Name() != tt.name {
				t.Fatalf("Name() = %q, want %q", tt.agent.Name(), tt.name)
			}

			if tt.agent.DisplayName() != tt.display {
				t.Fatalf("DisplayName() = %q, want %q", tt.agent.DisplayName(), tt.display)
			}

			if tt.agent.McpInstallationStrategy() != tt.strategy {
				t.Fatalf("McpInstallationStrategy() = %v, want %v", tt.agent.McpInstallationStrategy(), tt.strategy)
			}

			if tt.agent.McpConfigPath() != tt.mcpPath {
				t.Fatalf("McpConfigPath() = %q, want %q", tt.agent.McpConfigPath(), tt.mcpPath)
			}

			guidelinesAgent, ok := tt.agent.(boost.SupportsGuidelines)
			if !ok {
				t.Fatalf("%s should support guidelines", tt.name)
			}

			if guidelinesAgent.GuidelinesPath() != tt.guidePath {
				t.Fatalf("GuidelinesPath() = %q, want %q", guidelinesAgent.GuidelinesPath(), tt.guidePath)
			}

			skillsAgent, ok := tt.agent.(boost.SupportsSkills)
			if !ok {
				t.Fatalf("%s should support skills", tt.name)
			}

			if skillsAgent.SkillsPath() != tt.skillsPath {
				t.Fatalf("SkillsPath() = %q, want %q", skillsAgent.SkillsPath(), tt.skillsPath)
			}

			if tt.name == "junie" {
				if got := tt.agent.GoBinaryPath(false); !filepath.IsAbs(got) || filepath.Base(got) != "go" {
					t.Fatalf("Junie GoBinaryPath(false) = %q, want absolute go binary", got)
				}
				if got := tt.agent.EntryPointPath(false); !filepath.IsAbs(got) || filepath.Base(got) != "main.go" {
					t.Fatalf("Junie EntryPointPath(false) = %q, want absolute main.go path", got)
				}

				return
			}

			if tt.agent.GoBinaryPath(false) != "go" {
				t.Fatalf("GoBinaryPath(false) = %q, want go", tt.agent.GoBinaryPath(false))
			}

			if tt.agent.EntryPointPath(false) != "main.go" {
				t.Fatalf("EntryPointPath(false) = %q, want main.go", tt.agent.EntryPointPath(false))
			}
		})
	}
}

func TestInventoryAgentPathOverridesAndServerConfig(t *testing.T) {
	t.Parallel()

	cursor := agents.NewCursor(agents.AgentOptions{GoBinary: "php", EntryPoint: "cli"})
	if cursor.GoBinaryPath(false) != "php" || cursor.GoBinaryPath(true) != "php" {
		t.Fatalf("configured binary should win for forced and relative paths")
	}
	if cursor.EntryPointPath(false) != "cli" {
		t.Fatalf("EntryPointPath(false) = %q, want cli", cursor.EntryPointPath(false))
	}

	absEntryPoint, err := filepath.Abs("cli")
	if err != nil {
		t.Fatalf("filepath.Abs(cli): %v", err)
	}

	if cursor.EntryPointPath(true) != absEntryPoint {
		t.Fatalf("EntryPointPath(true) = %q, want %q", cursor.EntryPointPath(true), absEntryPoint)
	}

	claude := agents.NewClaudeCode(agents.AgentOptions{McpConfigPath: "custom/.mcp.json"})
	if claude.McpConfigPath() != "custom/.mcp.json" {
		t.Fatalf("McpConfigPath override = %q", claude.McpConfigPath())
	}

	base := agents.NewBaseAgent(agents.AgentOptions{})
	httpConfig := base.HttpMcpServerConfig("http://127.0.0.1:8090/mcp")
	if httpConfig["type"] != "http" || httpConfig["url"] != "http://127.0.0.1:8090/mcp" {
		t.Fatalf("base HTTP MCP config = %#v", httpConfig)
	}

	cursorHTTP := cursor.HttpMcpServerConfig("https://app.test/mcp")
	args, ok := cursorHTTP["args"].([]string)
	if cursorHTTP["command"] != "npx" || !ok || strings.Join(args, " ") != "-y mcp-remote https://app.test/mcp" {
		t.Fatalf("Cursor HTTP MCP config = %#v", cursorHTTP)
	}

	opencodeHTTP := agents.NewOpenCode().HttpMcpServerConfig("https://app.test/mcp")
	if opencodeHTTP["type"] != "http" || opencodeHTTP["url"] != "https://app.test/mcp" {
		t.Fatalf("OpenCode HTTP MCP config = %#v", opencodeHTTP)
	}
}

func TestInventoryAgentForcedAbsolutePathsAndConfiguredJuniePaths(t *testing.T) {
	t.Parallel()

	cursor := agents.NewCursor()
	goPath := cursor.GoBinaryPath(true)
	if !filepath.IsAbs(goPath) || filepath.Base(goPath) != "go" {
		t.Fatalf("GoBinaryPath(true) = %q, want absolute go binary", goPath)
	}

	entryPointPath := cursor.EntryPointPath(true)
	if !filepath.IsAbs(entryPointPath) || filepath.Base(entryPointPath) != "main.go" {
		t.Fatalf("EntryPointPath(true) = %q, want absolute main.go path", entryPointPath)
	}

	defaultJunie := agents.NewJunie()
	junieGoPath := defaultJunie.GoBinaryPath(false)
	if junieGoPath != defaultJunie.GoBinaryPath(true) || !filepath.IsAbs(junieGoPath) {
		t.Fatalf("Junie GoBinaryPath should remain absolute regardless of force flag, got %q and %q", junieGoPath, defaultJunie.GoBinaryPath(true))
	}

	junieEntryPoint := defaultJunie.EntryPointPath(false)
	if junieEntryPoint != defaultJunie.EntryPointPath(true) || !filepath.IsAbs(junieEntryPoint) {
		t.Fatalf("Junie EntryPointPath should remain absolute regardless of force flag, got %q and %q", junieEntryPoint, defaultJunie.EntryPointPath(true))
	}

	tmp := t.TempDir()
	configuredGo := filepath.Join(tmp, "bin", "go")
	configuredEntryPoint := filepath.Join(tmp, "cli")
	junie := agents.NewJunie(agents.AgentOptions{GoBinary: configuredGo, EntryPoint: configuredEntryPoint})

	if junie.GoBinaryPath(false) != configuredGo || junie.GoBinaryPath(true) != configuredGo {
		t.Fatalf("configured Junie Go path should be preserved for relative and forced resolution")
	}
	if junie.EntryPointPath(false) != configuredEntryPoint || junie.EntryPointPath(true) != configuredEntryPoint {
		t.Fatalf("configured Junie entry point should be preserved for relative and forced resolution")
	}
}

func TestInventoryAgentInstallMcpNormalizesCommandsAndWritesConfig(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "mcp.json")
	agent := agents.NewClaudeCode(agents.AgentOptions{McpConfigPath: path})

	ok, err := agent.InstallMcp("boost", "docker exec app go", []string{"run", "."}, map[string]string{"APP_ENV": "local"})
	if err != nil {
		t.Fatalf("InstallMcp: %v", err)
	}
	if !ok {
		t.Fatal("InstallMcp should report success for a new entry")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var root map[string]map[string]map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	server := root["mcpServers"]["boost"]
	if server["command"] != "docker" {
		t.Fatalf("command = %v, want docker", server["command"])
	}

	gotArgs := make([]string, 0)
	for _, value := range server["args"].([]any) {
		gotArgs = append(gotArgs, value.(string))
	}
	if strings.Join(gotArgs, " ") != "exec app go run ." {
		t.Fatalf("args = %v", gotArgs)
	}

	ok, err = agent.InstallHttpMcp("boost-http", "https://app.test/mcp")
	if err != nil {
		t.Fatalf("InstallHttpMcp: %v", err)
	}
	if !ok {
		t.Fatal("InstallHttpMcp should report success for a new entry")
	}

	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile http config: %v", err)
	}

	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("Unmarshal http config: %v", err)
	}

	httpServer := root["mcpServers"]["boost-http"]
	if httpServer["type"] != "http" || httpServer["url"] != "https://app.test/mcp" {
		t.Fatalf("HTTP MCP config = %#v", httpServer)
	}

	absoluteCommand := filepath.Join(tmp, "bin with spaces", "go")
	ok, err = agent.InstallMcp("absolute", absoluteCommand, []string{"run", "."}, nil)
	if err != nil {
		t.Fatalf("InstallMcp absolute command: %v", err)
	}
	if !ok {
		t.Fatal("InstallMcp should report success for an absolute command entry")
	}

	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile absolute command config: %v", err)
	}

	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("Unmarshal absolute command config: %v", err)
	}

	absoluteServer := root["mcpServers"]["absolute"]
	if absoluteServer["command"] != absoluteCommand {
		t.Fatalf("absolute command = %v, want %q", absoluteServer["command"], absoluteCommand)
	}
}

func TestInventoryAgentInstallMcpMissingConfigPaths(t *testing.T) {
	t.Parallel()

	amp := agents.NewAmp()
	ok, err := amp.InstallMcp("boost", "go", []string{"run", "."}, nil)
	if ok {
		t.Fatal("InstallMcp should return false when the agent has no MCP config path")
	}
	if !errors.Is(err, boost.ErrNoMcpConfigPath) {
		t.Fatalf("InstallMcp error = %v, want ErrNoMcpConfigPath", err)
	}

	base := agents.NewBaseAgent(agents.AgentOptions{})
	ok, err = base.InstallHttpMcp("boost-http", "https://app.test/mcp")
	if ok {
		t.Fatal("InstallHttpMcp should return false when the agent has no MCP config path")
	}
	if !errors.Is(err, boost.ErrNoMcpConfigPath) {
		t.Fatalf("InstallHttpMcp error = %v, want ErrNoMcpConfigPath", err)
	}
}

func TestInventoryAgentProjectDetectionMarkers(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	mustMkdir := func(path string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(tmp, path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
	}

	mustMkdir(".amp")
	mustMkdir(".kiro")

	if !agents.NewAmp().DetectInProject(tmp) {
		t.Fatal("Amp should detect .amp project directory")
	}

	if !agents.NewKiro().DetectInProject(tmp) {
		t.Fatal("Kiro should detect .kiro project directory")
	}

	if err := os.WriteFile(filepath.Join(tmp, "codex.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("write codex.json: %v", err)
	}

	if !agents.NewCodex().DetectInProject(tmp) {
		t.Fatal("Codex should detect codex.json")
	}

}
