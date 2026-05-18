package console

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	bcons "github.com/bedrock/packages/console"
	"github.com/bedrock/services/brain/api/graph"
)

func TestCommandsReturnsExpectedThree(t *testing.T) {
	cmds := Commands()

	if len(cmds) != 3 {
		t.Fatalf("len(Commands) = %d, want 3", len(cmds))
	}

	want := map[string]bool{
		"brain:scan":           false,
		"brain:export-context": false,
		"brain:generate-rules": false,
	}

	for _, c := range cmds {
		if _, ok := want[c.Name]; ok {
			want[c.Name] = true

			continue
		}

		t.Errorf("unexpected command %q", c.Name)
	}

	for name, found := range want {
		if !found {
			t.Errorf("missing command %q", name)
		}
	}
}

func TestCommandsHaveRunHandlers(t *testing.T) {
	for _, c := range Commands() {
		if c.Run == nil {
			t.Errorf("command %q has nil Run", c.Name)
		}
	}
}

func TestResolvePathsDefaults(t *testing.T) {
	in := bcons.NewInput()
	target, output := resolvePaths(in)

	if !filepath.IsAbs(target) {
		t.Errorf("target = %q, want absolute", target)
	}

	wantOutput := filepath.Join(target, "storage", "brain")

	if output != wantOutput {
		t.Errorf("output = %q, want %q", output, wantOutput)
	}
}

func TestResolvePathsHonorsExplicitTarget(t *testing.T) {
	in := bcons.NewInput()
	in.Options = map[string]string{"target": "."}

	target, output := resolvePaths(in)

	if !filepath.IsAbs(target) {
		t.Errorf("target = %q, want absolute", target)
	}

	if !strings.HasPrefix(output, target) {
		t.Errorf("output %q not under target %q", output, target)
	}
}

func TestResolvePathsHonorsExplicitOutput(t *testing.T) {
	dir := t.TempDir()
	in := bcons.NewInput()
	in.Options = map[string]string{"target": dir, "output": filepath.Join(dir, "custom")}

	_, output := resolvePaths(in)

	if output != filepath.Join(dir, "custom") {
		t.Errorf("output = %q, want override honored", output)
	}
}

func TestCountRoutes(t *testing.T) {
	g := graph.NewGraph("svc")

	if got := countRoutes(g); got != 0 {
		t.Errorf("empty graph routes = %d, want 0", got)
	}

	g.AddNode(graph.NewNode("r:1", graph.NodeTypeRoute, ""))
	g.AddNode(graph.NewNode("r:2", graph.NodeTypeRoute, ""))
	g.AddNode(graph.NewNode("c:1", graph.NodeTypeController, ""))

	if got := countRoutes(g); got != 2 {
		t.Errorf("routes = %d, want 2", got)
	}
}

func TestWriteGraphWritesBothFiles(t *testing.T) {
	dir := t.TempDir()
	g := graph.NewGraph("svc")
	g.AddNode(graph.NewNode("r:1", graph.NodeTypeRoute, ""))
	g.Stamp()

	if err := writeGraph(dir, g); err != nil {
		t.Fatalf("writeGraph: %v", err)
	}

	for _, name := range []string{".graph-all.json", ".graph-manifest.json"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}

	body, err := os.ReadFile(filepath.Join(dir, ".graph-manifest.json"))

	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}

	var m graph.Manifest

	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}

	if m.TotalNodes != 1 || m.TotalRoutes != 1 {
		t.Errorf("manifest = %+v, want TotalNodes=1 TotalRoutes=1", m)
	}
}

func TestMustReturnsCommandOnSuccess(t *testing.T) {
	cmd := &bcons.Command{Name: "test:cmd"}

	if got := must(cmd, nil); got != cmd {
		t.Errorf("must should return its first arg when err is nil")
	}
}

func TestMustPanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("must did not panic on error")
		}
	}()

	must(nil, errSentinel{})
}

type errSentinel struct{}

func (errSentinel) Error() string { return "boom" }
