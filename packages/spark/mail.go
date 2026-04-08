package spark

import "fmt"

// CustomPlanInquiryMail is the notification sent to the contact address when
// a user submits a custom plan inquiry.
type CustomPlanInquiryMail struct {
	Name    string
	Email   string
	Company string
	Message string
}

// Subject returns the email subject line.
func (m *CustomPlanInquiryMail) Subject() string {
	return fmt.Sprintf("New custom plan enquiry from %s", m.Name)
}

// ReplyTo returns the inquirer's email for easy reply.
func (m *CustomPlanInquiryMail) ReplyTo() string { return m.Email }

// Body returns the email body text.
func (m *CustomPlanInquiryMail) Body() string {
	return fmt.Sprintf("Name: %s\nEmail: %s\nCompany: %s\nMessage: %s",
		m.Name, m.Email, m.Company, m.Message)
}

// CustomPlanInquiryConfirmationMail is the confirmation sent to the inquirer.
type CustomPlanInquiryConfirmationMail struct {
	Name string
}

// Subject returns the email subject line.
func (m *CustomPlanInquiryConfirmationMail) Subject() string {
	return "We've received your enquiry"
}

// ReplyTo returns empty (no reply-to for confirmation).
func (m *CustomPlanInquiryConfirmationMail) ReplyTo() string { return "" }

// Body returns the email body text.
func (m *CustomPlanInquiryConfirmationMail) Body() string {
	return fmt.Sprintf("Hi %s,\n\nThank you for your enquiry. We'll be in touch shortly.", m.Name)
}

// Ensure both implement MailMessage.
var (
	_ MailMessage = (*CustomPlanInquiryMail)(nil)
	_ MailMessage = (*CustomPlanInquiryConfirmationMail)(nil)
)
