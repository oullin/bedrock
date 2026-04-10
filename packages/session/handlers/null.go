package handlers

import "context"

// NullHandler discards all session data. Useful when sessions are not needed.
type NullHandler struct{}

func (h *NullHandler) Open(_ context.Context, _, _ string) error        { return nil }
func (h *NullHandler) Close(_ context.Context) error                    { return nil }
func (h *NullHandler) Read(_ context.Context, _ string) (string, error) { return "", nil }
func (h *NullHandler) Write(_ context.Context, _, _ string) error       { return nil }
func (h *NullHandler) Destroy(_ context.Context, _ string) error        { return nil }
func (h *NullHandler) GC(_ context.Context, _ int) error                { return nil }
