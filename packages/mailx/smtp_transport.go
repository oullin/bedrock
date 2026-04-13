package mailx

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"
	"path/filepath"

	cmail "github.com/bedrock/packages/contracts/mail"
)

// SMTPTransport delivers email via an SMTP server.
type SMTPTransport struct {
	host       string
	port       int
	username   string
	password   string
	encryption string // "tls", "starttls", or "" for plain
	localName  string // HELO/EHLO hostname
}

// SMTPOption configures an SMTPTransport.
type SMTPOption func(*SMTPTransport)

var _ Transport = (*SMTPTransport)(nil)

// WithEncryption sets the encryption mode ("tls", "starttls", or "").
func WithEncryption(enc string) SMTPOption {
	return func(t *SMTPTransport) {
		t.encryption = enc
	}
}

// WithLocalName sets the HELO/EHLO hostname.
func WithLocalName(name string) SMTPOption {
	return func(t *SMTPTransport) {
		t.localName = name
	}
}

// NewSMTPTransport creates a new SMTP transport.
func NewSMTPTransport(host string, port int, username, password string, opts ...SMTPOption) *SMTPTransport {
	t := &SMTPTransport{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Send delivers the message via SMTP.
func (t *SMTPTransport) Send(ctx context.Context, message *cmail.Message) (*cmail.SentMessage, error) {
	addr := fmt.Sprintf("%s:%d", t.host, t.port)

	var client *smtp.Client

	var err error

	if t.encryption == "tls" {
		conn, dialErr := tls.DialWithDialer(
			&net.Dialer{},
			"tcp",
			addr,
			&tls.Config{ServerName: t.host},
		)

		if dialErr != nil {
			return nil, fmt.Errorf("%w: %v", ErrSendFailed, dialErr)
		}

		client, err = smtp.NewClient(conn, t.host)
	} else {
		client, err = smtp.Dial(addr)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	defer client.Close()

	if t.localName != "" {
		if err = client.Hello(t.localName); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
		}
	}

	if t.encryption == "starttls" {
		if err = client.StartTLS(&tls.Config{ServerName: t.host}); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
		}
	}

	if t.username != "" || t.password != "" {
		auth := smtp.PlainAuth("", t.username, t.password, t.host)

		if err = client.Auth(auth); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
		}
	}

	from := message.GetFrom().Email

	if from == "" {
		return nil, fmt.Errorf("%w: missing from address", ErrSendFailed)
	}

	if err = client.Mail(from); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	recipients := message.AllRecipients()

	if len(recipients) == 0 {
		return nil, ErrNoRecipients
	}

	for _, rcpt := range recipients {
		if err = client.Rcpt(rcpt.Email); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
		}
	}

	wc, err := client.Data()

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	raw, err := buildRawMessage(message)

	if err != nil {
		wc.Close()

		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	if _, err = wc.Write(raw); err != nil {
		wc.Close()

		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	if err = wc.Close(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	if err = client.Quit(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSendFailed, err)
	}

	return &cmail.SentMessage{
		Original:  message,
		MessageID: generateMessageID(),
	}, nil
}

// String returns the transport name.
func (t *SMTPTransport) String() string {
	return "smtp"
}

// buildRawMessage constructs the full RFC 2822 message bytes.
func buildRawMessage(m *cmail.Message) ([]byte, error) {
	var buf bytes.Buffer

	htmlBody := m.GetHTMLBody()
	textBody := m.GetTextBody()
	attachments := m.GetAttachments()
	embeds := m.GetEmbeds()
	hasAttachments := len(attachments) > 0
	hasEmbeds := len(embeds) > 0
	hasHTML := htmlBody != ""
	hasText := textBody != ""

	// Write top-level headers.
	writeHeader(&buf, "From", m.GetFrom().String())
	writeHeader(&buf, "Subject", m.GetSubject())

	if to := m.GetTo(); len(to) > 0 {
		writeHeader(&buf, "To", formatAddresses(to))
	}

	if cc := m.GetCC(); len(cc) > 0 {
		writeHeader(&buf, "Cc", formatAddresses(cc))
	}

	if replyTo := m.GetReplyTo(); len(replyTo) > 0 {
		writeHeader(&buf, "Reply-To", formatAddresses(replyTo))
	}

	if sender := m.GetSender(); sender.Email != "" {
		writeHeader(&buf, "Sender", sender.String())
	}

	if rp := m.GetReturnPath(); rp != "" {
		writeHeader(&buf, "Return-Path", "<"+rp+">")
	}

	if p := m.GetPriority(); p > 0 {
		writeHeader(&buf, "X-Priority", fmt.Sprintf("%d", p))
	}

	for name, values := range m.GetCustomHeaders() {
		for _, v := range values {
			writeHeader(&buf, name, v)
		}
	}

	writeHeader(&buf, "MIME-Version", "1.0")

	switch {
	case hasAttachments:
		return buildMixed(&buf, htmlBody, textBody, attachments, embeds)
	case hasHTML && hasText:
		return buildAlternative(&buf, htmlBody, textBody, embeds)
	case hasHTML:
		if hasEmbeds {
			return buildRelated(&buf, htmlBody, embeds)
		}

		writeHeader(&buf, "Content-Type", "text/html; charset=utf-8")
		writeHeader(&buf, "Content-Transfer-Encoding", "quoted-printable")
		buf.WriteString("\r\n")
		buf.WriteString(htmlBody)

		return buf.Bytes(), nil
	default:
		writeHeader(&buf, "Content-Type", "text/plain; charset=utf-8")
		writeHeader(&buf, "Content-Transfer-Encoding", "quoted-printable")
		buf.WriteString("\r\n")
		buf.WriteString(textBody)

		return buf.Bytes(), nil
	}
}

func buildMixed(buf *bytes.Buffer, html, text string, attachments []*cmail.Attachment, embeds []*cmail.Embed) ([]byte, error) {
	mw := multipart.NewWriter(buf)
	writeHeader(buf, "Content-Type", "multipart/mixed; boundary="+mw.Boundary())
	buf.WriteString("\r\n")

	// Body part.
	if html != "" && text != "" {
		if err := writeAlternativePart(mw, html, text, embeds); err != nil {
			return nil, err
		}
	} else if html != "" {
		if err := writeHTMLPart(mw, html, embeds); err != nil {
			return nil, err
		}
	} else {
		if err := writeTextPart(mw, text); err != nil {
			return nil, err
		}
	}

	// Attachments.
	for _, a := range attachments {
		if err := writeAttachmentPart(mw, a); err != nil {
			return nil, err
		}
	}

	if err := mw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildAlternative(buf *bytes.Buffer, html, text string, embeds []*cmail.Embed) ([]byte, error) {
	mw := multipart.NewWriter(buf)
	writeHeader(buf, "Content-Type", "multipart/alternative; boundary="+mw.Boundary())
	buf.WriteString("\r\n")

	if err := writeTextPart(mw, text); err != nil {
		return nil, err
	}

	if err := writeHTMLPart(mw, html, embeds); err != nil {
		return nil, err
	}

	if err := mw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildRelated(buf *bytes.Buffer, html string, embeds []*cmail.Embed) ([]byte, error) {
	mw := multipart.NewWriter(buf)
	writeHeader(buf, "Content-Type", "multipart/related; boundary="+mw.Boundary())
	buf.WriteString("\r\n")

	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Type", "text/html; charset=utf-8")
	hdr.Set("Content-Transfer-Encoding", "quoted-printable")

	pw, err := mw.CreatePart(hdr)

	if err != nil {
		return nil, err
	}

	pw.Write([]byte(html))

	for _, e := range embeds {
		if err := writeEmbedPart(mw, e); err != nil {
			return nil, err
		}
	}

	if err := mw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeAlternativePart(parent *multipart.Writer, html, text string, embeds []*cmail.Embed) error {
	var buf bytes.Buffer
	alt := multipart.NewWriter(&buf)

	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Type", "multipart/alternative; boundary="+alt.Boundary())

	pw, err := parent.CreatePart(hdr)

	if err != nil {
		return err
	}

	if err := writeTextPart(alt, text); err != nil {
		return err
	}

	if err := writeHTMLPart(alt, html, embeds); err != nil {
		return err
	}

	if err := alt.Close(); err != nil {
		return err
	}

	_, err = pw.Write(buf.Bytes())

	return err
}

func writeTextPart(mw *multipart.Writer, text string) error {
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Type", "text/plain; charset=utf-8")
	hdr.Set("Content-Transfer-Encoding", "quoted-printable")

	pw, err := mw.CreatePart(hdr)

	if err != nil {
		return err
	}

	_, err = pw.Write([]byte(text))

	return err
}

func writeHTMLPart(mw *multipart.Writer, html string, embeds []*cmail.Embed) error {
	if len(embeds) > 0 {
		var buf bytes.Buffer
		rel := multipart.NewWriter(&buf)

		hdr := make(textproto.MIMEHeader)
		hdr.Set("Content-Type", "multipart/related; boundary="+rel.Boundary())

		pw, err := mw.CreatePart(hdr)

		if err != nil {
			return err
		}

		innerHdr := make(textproto.MIMEHeader)
		innerHdr.Set("Content-Type", "text/html; charset=utf-8")
		innerHdr.Set("Content-Transfer-Encoding", "quoted-printable")

		inner, err := rel.CreatePart(innerHdr)

		if err != nil {
			return err
		}

		inner.Write([]byte(html))

		for _, e := range embeds {
			if err := writeEmbedPart(rel, e); err != nil {
				return err
			}
		}

		if err := rel.Close(); err != nil {
			return err
		}

		_, err = pw.Write(buf.Bytes())

		return err
	}

	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Type", "text/html; charset=utf-8")
	hdr.Set("Content-Transfer-Encoding", "quoted-printable")

	pw, err := mw.CreatePart(hdr)

	if err != nil {
		return err
	}

	_, err = pw.Write([]byte(html))

	return err
}

func writeAttachmentPart(mw *multipart.Writer, a *cmail.Attachment) error {
	mimeType := a.Mime

	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(a.Name))

		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
	}

	name := a.Name

	if name == "" {
		name = filepath.Base(a.Path)
	}

	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Type", mimeType+"; name=\""+name+"\"")
	hdr.Set("Content-Disposition", "attachment; filename=\""+name+"\"")
	hdr.Set("Content-Transfer-Encoding", "base64")

	pw, err := mw.CreatePart(hdr)

	if err != nil {
		return err
	}

	if a.Data != nil {
		reader, err := a.Data()

		if err != nil {
			return err
		}

		data, err := io.ReadAll(reader)

		if err != nil {
			return err
		}

		_, err = pw.Write([]byte(base64.StdEncoding.EncodeToString(data)))

		return err
	}

	// For path-based attachments, the actual file reading would happen
	// here. For now we write the path as a placeholder.
	_, err = pw.Write([]byte(base64.StdEncoding.EncodeToString([]byte(a.Path))))

	return err
}

func writeEmbedPart(mw *multipart.Writer, e *cmail.Embed) error {
	mimeType := e.Mime

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Type", mimeType)
	hdr.Set("Content-Transfer-Encoding", "base64")
	hdr.Set("Content-ID", "<"+e.CID+">")
	hdr.Set("Content-Disposition", "inline; filename=\""+e.Name+"\"")

	pw, err := mw.CreatePart(hdr)

	if err != nil {
		return err
	}

	_, err = pw.Write([]byte(base64.StdEncoding.EncodeToString(e.Data)))

	return err
}

func writeHeader(buf *bytes.Buffer, name, value string) {
	buf.WriteString(name)
	buf.WriteString(": ")
	buf.WriteString(value)
	buf.WriteString("\r\n")
}
