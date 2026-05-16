package analysis

import (
	"go/ast"
	"strings"

	"github.com/bedrock/services/brain/api/graph"
	"golang.org/x/tools/go/packages"
)

// SecurityIssue is a structured finding attached to a service_provider-shaped
// "security:findings" pseudo-node. Field names match upstream-brain.
type SecurityIssue struct {
	Rule     string `json:"rule"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	File     string `json:"file"`
}

// SecurityAnalyzer raises lightweight findings for known-bad Go patterns:
//
//   - String-concatenated SQL passed to *.Query/Exec/QueryRow.
//   - fmt.Sprintf-built SQL passed to *.Query/Exec/QueryRow.
//   - exec.Command with a non-literal first arg (potential command injection).
//   - os.Setenv with hard-coded secrets-looking values (keys named password,
//     secret, token, key — checked by the static value being non-empty).
//
// Findings are advisory. Phase 7 is the seed; phase 7b/c (future) can add
// more rules without changing this analyzer's contract.
type SecurityAnalyzer struct {
	Issues []SecurityIssue
}

func (s *SecurityAnalyzer) Name() string { return "security" }

func (s *SecurityAnalyzer) Analyze(ctx *Context) error {
	return ctx.Project.EachFile(func(_ *packages.Package, file *ast.File, _ string) error {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			s.checkSQLInjection(ctx, call)
			s.checkCommandInjection(ctx, call)
			s.checkHardcodedSecret(ctx, call)
			return true
		})
		return nil
	})
}

func (s *SecurityAnalyzer) raise(ctx *Context, rule, severity, msg string, pos ast.Node) {
	issue := SecurityIssue{
		Rule:     rule,
		Severity: severity,
		Message:  msg,
		File:     ctx.Project.Position(pos),
	}
	s.Issues = append(s.Issues, issue)
	id := "security:" + rule + ":" + issue.File
	node := graph.NewNode(id, graph.NodeTypeServiceProvider, rule).
		Set("kind", "security_finding").
		Set("severity", severity).
		Set("message", msg).
		Set("file", issue.File)
	ctx.Graph.AddNode(node)
}

func (s *SecurityAnalyzer) checkSQLInjection(ctx *Context, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	switch sel.Sel.Name {
	case "Query", "Exec", "QueryRow", "QueryContext", "ExecContext", "QueryRowContext":
	default:
		return
	}
	if len(call.Args) == 0 {
		return
	}
	first := call.Args[0]
	// Skip context.Context (first arg of *Context methods).
	if ident, ok := first.(*ast.Ident); ok && ident.Name == "ctx" && len(call.Args) > 1 {
		first = call.Args[1]
	}
	if isStringConcat(first) {
		s.raise(ctx, "sql_injection", "high",
			"SQL built via string concatenation; use parameter placeholders", call)
		return
	}
	if isSprintfCall(first) {
		s.raise(ctx, "sql_injection", "high",
			"SQL built via fmt.Sprintf; use parameter placeholders", call)
	}
}

func (s *SecurityAnalyzer) checkCommandInjection(ctx *Context, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	if sel.Sel.Name != "Command" && sel.Sel.Name != "CommandContext" {
		return
	}
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok || pkgIdent.Name != "exec" {
		return
	}
	argIdx := 0
	if sel.Sel.Name == "CommandContext" {
		argIdx = 1
	}
	if argIdx >= len(call.Args) {
		return
	}
	if _, isLit := call.Args[argIdx].(*ast.BasicLit); !isLit {
		s.raise(ctx, "command_injection", "medium",
			"exec.Command called with a non-literal program path", call)
	}
}

func (s *SecurityAnalyzer) checkHardcodedSecret(ctx *Context, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Setenv" {
		return
	}
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok || pkgIdent.Name != "os" {
		return
	}
	if len(call.Args) < 2 {
		return
	}
	name, _ := stringLit(call.Args[0])
	val, _ := stringLit(call.Args[1])
	if name == "" || val == "" {
		return
	}
	lower := strings.ToLower(name)
	hot := strings.Contains(lower, "secret") ||
		strings.Contains(lower, "password") ||
		strings.Contains(lower, "token") ||
		strings.Contains(lower, "key") ||
		strings.Contains(lower, "apikey")
	if hot {
		s.raise(ctx, "hardcoded_secret", "high",
			"os.Setenv stores a literal value into "+name+"; load from a secret manager instead", call)
	}
}

func isStringConcat(e ast.Expr) bool {
	be, ok := e.(*ast.BinaryExpr)
	if !ok || be.Op.String() != "+" {
		return false
	}
	// at least one operand is a string literal — heuristic for SQL string-building
	if _, ok := be.X.(*ast.BasicLit); ok {
		return true
	}
	if _, ok := be.Y.(*ast.BasicLit); ok {
		return true
	}
	// or one operand is itself a concat
	return isStringConcat(be.X) || isStringConcat(be.Y)
}

func isSprintfCall(e ast.Expr) bool {
	call, ok := e.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok || pkgIdent.Name != "fmt" {
		return false
	}
	return sel.Sel.Name == "Sprintf"
}
