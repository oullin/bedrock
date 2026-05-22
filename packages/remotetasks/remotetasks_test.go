package remotetasks

import (
	"context"
	"errors"
	"testing"
)

func TestPlanRendersCommandsForHosts(t *testing.T) {
	t.Parallel()

	steps, err := Plan(Task{
		Name:     "deploy",
		Hosts:    []string{"web-2", "web-1"},
		Commands: []string{"cd {{ .release }}", "php artisan migrate --force"},
		Variables: map[string]any{
			"release": "/srv/app/current",
		},
	}, nil)

	if err != nil {
		t.Fatalf("Plan returned error: %v", err)
	}

	if len(steps) != 4 {
		t.Fatalf("len(steps) = %d, want 4", len(steps))
	}

	if steps[0].Host != "web-1" {
		t.Fatalf("first host = %q, want web-1", steps[0].Host)
	}

	if steps[0].Command != "cd /srv/app/current" {
		t.Fatalf("first command = %q", steps[0].Command)
	}

	if steps[1].Command != "php artisan migrate --force" {
		t.Fatalf("second command = %q", steps[1].Command)
	}
}

// SSHConfigFileTest::test_it_removes_leading_and_trailing_whitespace_on_each_line
// SSHConfigFileTest::test_it_splits_keywords_and_arguments_by_equal_sign
// SSHConfigFileTest::test_it_splits_keywords_and_arguments_by_whitespace
// SSHConfigFileTest::test_it_lowercases_keywords_not_arguments
// SSHConfigFileTest::test_it_unquotes_arguments_and_reserves_whitespace
// SSHConfigFileTest::test_it_ignores_comments_and_empty_lines_outside_host_section
// SSHConfigFileTest::test_it_ignores_comments_and_empty_lines_inside_host_section
// SSHConfigFileTest::test_it_parses_sections_without_new_lines_between_them
// SSHConfigFileTest::test_it_parses_sections_separated_by_new_lines
// SSHConfigFileTest::test_it_parses_sections_separated_by_match_keyword
func TestSSHConfigParsing(t *testing.T) {
	t.Parallel()

	config := `
# ignored
 Host=Bar
	Hostname baz.com
# ignored in section
Port 1234
IdentityFile="path "to/file""
Match user john
    IdentityFile ~/.ssh/id_rsa
Host qux
Hostname=qux.com

Host zap
Hostname zap.com
`

	groups := ParseSSHConfigString(config).Groups()

	if len(groups) != 3 {
		t.Fatalf("len(groups) = %d, want 3", len(groups))
	}

	if groups[0]["host"] != "Bar" {
		t.Fatalf("host casing = %q", groups[0]["host"])
	}

	if groups[0]["hostname"] != "baz.com" {
		t.Fatalf("hostname = %q", groups[0]["hostname"])
	}

	if groups[0]["port"] != "1234" {
		t.Fatalf("port = %q", groups[0]["port"])
	}

	if groups[0]["identityfile"] != `path "to/file"` {
		t.Fatalf("identityfile = %q", groups[0]["identityfile"])
	}

	if _, ok := groups[1]["identityfile"]; ok {
		t.Fatalf("match section leaked into host group: %#v", groups[1])
	}
}

// SSHConfigFileTest::test_it_finds_a_matching_host_without_user
// SSHConfigFileTest::test_it_finds_a_matching_host_by_hostname
// SSHConfigFileTest::test_it_returns_null_if_there_are_no_matching_hosts
// SSHConfigFileTest::test_it_finds_a_matching_host_if_user_specified_and_matches_config
// SSHConfigFileTest::test_it_returns_null_for_a_matching_host_if_user_specified_and_is_different_from_one_in_config
// SSHConfigFileTest::test_it_returns_valid_host_for_multiple_aliased_hosts
func TestSSHConfigFindConfiguredHost(t *testing.T) {
	t.Parallel()

	config := ParseSSHConfigString(`
Host foo bar baz
Hostname baz.com
User john
Host qux
Hostname qux.com
`)

	if host, ok := config.FindConfiguredHost("bar"); !ok || host != "bar" {
		t.Fatalf("FindConfiguredHost(bar) = %q, %v", host, ok)
	}

	if host, ok := config.FindConfiguredHost("baz.com"); !ok || host != "foo" {
		t.Fatalf("FindConfiguredHost(baz.com) = %q, %v", host, ok)
	}

	if _, ok := config.FindConfiguredHost("none"); ok {
		t.Fatal("FindConfiguredHost(none) matched unexpectedly")
	}

	if host, ok := config.FindConfiguredHost("john@bar"); !ok || host != "bar" {
		t.Fatalf("FindConfiguredHost(john@bar) = %q, %v", host, ok)
	}

	if _, ok := config.FindConfiguredHost("jane@bar"); ok {
		t.Fatal("FindConfiguredHost(jane@bar) matched unexpectedly")
	}
}

func TestRunUsesRunner(t *testing.T) {
	t.Parallel()

	var called bool
	results, err := Run(context.Background(), RunnerFunc(func(ctx context.Context, step Step) (Result, error) {
		called = true

		return Result{Step: step, Stdout: "ok"}, nil
	}), Task{
		Name:     "status",
		Commands: []string{"whoami"},
	}, nil)

	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !called {
		t.Fatal("runner was not called")
	}

	if results[0].Stdout != "ok" {
		t.Fatalf("stdout = %q, want ok", results[0].Stdout)
	}
}

func TestPlanRequiresTaskName(t *testing.T) {
	t.Parallel()

	_, err := Plan(Task{Commands: []string{"echo ok"}}, nil)

	if !errors.Is(err, ErrMissingTaskName) {
		t.Fatalf("Plan error = %v, want ErrMissingTaskName", err)
	}
}
