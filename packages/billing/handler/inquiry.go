package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bedrock/packages/billing"
)

// InquiryHandler handles custom plan inquiry submissions.
type InquiryHandler struct {
	mailer billing.Mailer
	config *billing.Config
}

// NewInquiryHandler creates an InquiryHandler.

// Store validates the inquiry and sends notification emails.

// InquiryInput holds the validated input for a custom plan inquiry.
type InquiryInput struct {
	Name    string
	Email   string
	Company string
	Message string
}

func NewInquiryHandler(mailer billing.Mailer, config *billing.Config) *InquiryHandler {
	return &InquiryHandler{mailer: mailer, config: config}
}

func (h *InquiryHandler) Store(w http.ResponseWriter, r *http.Request) {
	var input InquiryInput

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
		_ = h.mailer.Send(ctx, contactEmail, &billing.CustomPlanInquiryMail{
			Name:    input.Name,
			Email:   input.Email,
			Company: input.Company,
			Message: input.Message,
		})
	}

	_ = h.mailer.Send(ctx, input.Email, &billing.CustomPlanInquiryConfirmationMail{Name: input.Name})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Thanks %s, we've received your enquiry and will be in touch shortly.", input.Name),
	})
}

func (h *InquiryHandler) resolveContactEmail() string {
	email := strings.TrimSpace(h.config.ContactEmail)

	if email == "" || !strings.Contains(email, "@") || strings.HasSuffix(strings.ToLower(email), "@example.com") {
		fallback := strings.TrimSpace(h.config.ContactFallbackEmail)

		if fallback != "" {
			return fallback
		}
	}

	return email
}

// Validate checks that the inquiry input fields are valid.
func (c *InquiryInput) Validate() billing.ValidationErrors {
	var errs billing.ValidationErrors

	if strings.TrimSpace(c.Name) == "" {
		errs = append(errs, billing.ValidationError{Field: "name", Message: "Your name is required."})
	} else if len(c.Name) > 150 {
		errs = append(errs, billing.ValidationError{Field: "name", Message: "Your name must not be greater than 150 characters."})
	}

	if strings.TrimSpace(c.Email) == "" {
		errs = append(errs, billing.ValidationError{Field: "email", Message: "An email address is required."})
	} else if len(c.Email) > 255 || !strings.Contains(c.Email, "@") {
		errs = append(errs, billing.ValidationError{Field: "email", Message: "Please provide a valid email address."})
	}

	if len(c.Company) > 150 {
		errs = append(errs, billing.ValidationError{Field: "company", Message: "Company name must not be greater than 150 characters."})
	}

	if strings.TrimSpace(c.Message) == "" {
		errs = append(errs, billing.ValidationError{Field: "message", Message: "A message is required."})
	} else if len(c.Message) > 2000 {
		errs = append(errs, billing.ValidationError{Field: "message", Message: "A message must not be greater than 2000 characters."})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
