package console

import (
	"encoding/json"
	"io"
)

func newIndentEncoder(w io.Writer) *json.Encoder {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc
}
