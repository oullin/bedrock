package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bedrock/packages/inertia/protocol"
)

func (a app) useHttpHandler(w http.ResponseWriter, r *http.Request) {
	a.container.Render(w, r, "Features/Http/UseHttp", protocol.Props{})
}

func (a app) useHttpApiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	if err := protocol.ParseForm(r); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	name := strings.TrimSpace(r.FormValue("name"))

	if strings.TrimSpace(name) == "" {
		name = "World"
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]any{
		"greeting":  fmt.Sprintf("Hello, %s!", name),
		"timestamp": time.Now().Format(time.RFC3339),
	}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}
}
