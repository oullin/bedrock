package analysis

import (
	"testing"

	"github.com/bedrock/services/brain/api/graph"
	"github.com/bedrock/services/brain/api/parser"
)

func TestSecurityAnalyzer_FindsInjectionAndSecrets(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, dir, "go.mod", "module example.com/fixture\n\ngo 1.26\n")
	mustWrite(t, dir, "bad.go", `package fixture

import (
	"fmt"
	"os"
	"os/exec"
)

type db struct{}
func (d *db) Query(q string, args ...any) any { return nil }

func bad(d *db, userID string) {
	d.Query("SELECT * FROM users WHERE id = " + userID)
	d.Query(fmt.Sprintf("SELECT %s FROM t", userID))
}

func runCmd(name string) {
	exec.Command(name)
}

func secrets() {
	os.Setenv("DB_PASSWORD", "hunter2")
	os.Setenv("STRIPE_API_KEY", "sk_live_abc")
}
`)
	proj, err := parser.Load(dir)

	if err != nil {
		t.Fatal(err)
	}

	g := graph.NewGraph("fixture")
	sa := &SecurityAnalyzer{}

	if err := sa.Analyze(&Context{Project: proj, Graph: g}); err != nil {
		t.Fatal(err)
	}

	counts := map[string]int{}

	for _, i := range sa.Issues {
		counts[i.Rule]++
	}

	if counts["sql_injection"] != 2 {
		t.Errorf("sql_injection count = %d, want 2 (%v)", counts["sql_injection"], sa.Issues)
	}

	if counts["command_injection"] != 1 {
		t.Errorf("command_injection count = %d, want 1", counts["command_injection"])
	}

	if counts["hardcoded_secret"] != 2 {
		t.Errorf("hardcoded_secret count = %d, want 2", counts["hardcoded_secret"])
	}
}
