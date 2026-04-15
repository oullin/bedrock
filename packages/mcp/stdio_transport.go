package mcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
)

// StdioTransport implements Transport over standard input/output. It reads
// newline-delimited JSON-RPC messages from an io.Reader and writes responses
// to an io.Writer.
type StdioTransport struct {
	in      io.Reader
	out     io.Writer
	handler func(ctx context.Context, message, sessionID string) (string, error)
}

// NewStdioTransport creates a StdioTransport backed by os.Stdin and os.Stdout.
func NewStdioTransport() *StdioTransport {
	return &StdioTransport{in: os.Stdin, out: os.Stdout}
}

// NewStdioTransportWithIO creates a StdioTransport with custom reader/writer.
// Useful in tests.
func NewStdioTransportWithIO(in io.Reader, out io.Writer) *StdioTransport {
	return &StdioTransport{in: in, out: out}
}

// OnReceive registers the handler that processes each incoming JSON-RPC message.
func (t *StdioTransport) OnReceive(h func(ctx context.Context, message, sessionID string) (string, error)) {
	t.handler = h
}

// Send writes a JSON-RPC message to the output writer followed by a newline.
func (t *StdioTransport) Send(_ context.Context, message, _ string) error {
	_, err := fmt.Fprintln(t.out, message)
	return err
}

// Run reads lines from the input reader until ctx is cancelled or EOF, passing
// each line to the registered handler and writing the response to the output.
func (t *StdioTransport) Run(ctx context.Context) error {
	scanner := bufio.NewScanner(t.in)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return err
			}
			return nil // EOF
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		if t.handler == nil {
			continue
		}

		resp, err := t.handler(ctx, line, "")
		if err != nil {
			// Write a JSON-RPC internal error response.
			errResp := ErrorResponse(nil, CodeInternalError, err.Error())
			b, _ := errResp.ToJSON()
			fmt.Fprintln(t.out, string(b)) //nolint:errcheck
			continue
		}

		fmt.Fprintln(t.out, resp) //nolint:errcheck
	}
}

// SessionID returns an empty string; stdio connections have no session ID.
func (t *StdioTransport) SessionID() string { return "" }
