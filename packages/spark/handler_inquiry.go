package spark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// CustomPlanInquiryHandler handles custom plan inquiry submissions.
type CustomPlanInquiryHandler struct {
	mailer Mailer
	config *Config
}

// NewCustomPlanInquiryHandler creates a CustomPlanInquiryHandler.
func NewCustomPlanInquiryHandler(mailer Mailer, config *Config) *CustomPlanInquiryHandler {
	return &CustomPlanInquiryHandler{mailer: mailer, config: config}
}

// Store validates the inquiry and sends notification emails.
func (h *CustomPlanInquiryHandler) Store(w http.ResponseWriter, r *http.Request) {
	var input CustomPlanInquiryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	if errs := input.Validate(); errs != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": errs})

		return
	}

	ctx := r.Context()
	contactEmail := h.resolveContactEmail()

	if contactEmail != "" {
		_ = h.mailer.Send(ctx, contactEmail, &customPlanInquiryMail{
			name:    input.Name,
			email:   input.Email,
			company: input.Company,
			message: input.Message,
		})
	}

	_ = h.mailer.Send(ctx, input.Email, &customPlanInquiryConfirmationMail{name: input.Name})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Thanks %s, we've received your enquiry and will be in touch shortly.", input.Name),
	})
}

func (h *CustomPlanInquiryHandler) resolveContactEmail() string {
	email := strings.TrimSpace(h.config.ContactEmail)

	if email == "" || !strings.Contains(email, "@") || strings.HasSuffix(strings.ToLower(email), "@example.com") {
		fallback := strings.TrimSpace(h.config.ContactFallbackEmail)
		if fallback != "" {
			return fallback
		}
	}

	return email
}

// customPlanInquiryMail is used internally by the handler.
type customPlanInquiryMail struct {
	name, email, company, message string
}

func (m *customPlanInquiryMail) Subject() string {
	return fmt.Sprintf("New custom plan enquiry from %s", m.name)
}

func (m *customPlanInquiryMail) ReplyTo() string { return m.email }

func (m *customPlanInquiryMail) Body() string {
	return fmt.Sprintf("Name: %s\nEmail: %s\nCompany: %s\nMessage: %s",
		m.name, m.email, m.company, m.message)
}

// customPlanInquiryConfirmationMail is sent to the inquirer.
type customPlanInquiryConfirmationMail struct{ name string }

func (m *customPlanInquiryConfirmationMail) Subject() string { return "We've received your enquiry" }
func (m *customPlanInquiryConfirmationMail) ReplyTo() string { return "" }
func (m *customPlanInquiryConfirmationMail) Body() string {
	return fmt.Sprintf("Hi %s,\n\nThank you for your enquiry. We'll be in touch shortly.", m.name)
}

// Ensure both implement MailMessage.
var (
	_ MailMessage = (*customPlanInquiryMail)(nil)
	_ MailMessage = (*customPlanInquiryConfirmationMail)(nil)
)

// resolveContactEmail is exported for the handler but the logic is in the
// unexported method above. This satisfies the context.Context requirement
// but is not used directly.
func resolveContactEmail(_ context.Context, _ *Config) string { return "" }
