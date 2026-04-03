package foundation

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/support/crypto"
)

// VerificationService sends and validates signed email verification links.
type VerificationService struct {
	Config auth.Config
	Users  auth.UserRepository
	Signer auth.LinkSigner
	Mailer auth.Mailer
	Clock  auth.Clock
}

// Send sends an email verification notification.
func (s *VerificationService) Send(ctx context.Context, user auth.Authenticatable) error {
	if user == nil {
		return auth.ErrUserNotFound
	}

	verifiable, ok := user.(auth.MustVerifyEmail)

	if !ok {
		return fmt.Errorf("foundation: user does not support email verification")
	}

	profile, ok := user.(auth.UserProfile)

	if !ok {
		return fmt.Errorf("foundation: user does not expose profile email")
	}

	expiresAt := s.Clock.Now().Add(s.Config.VerificationTTL)
	emailHash := crypto.EmailHash(normalize(profile.GetEmail()))
	signature, err := s.Signer.Sign(ctx, "email-verification", []string{user.GetAuthIdentifier(), emailHash}, expiresAt.Unix())

	if err != nil {
		return fmt.Errorf("sign verification link: %w", err)
	}

	base := strings.TrimRight(s.Config.BaseURL, "/")
	link := fmt.Sprintf("/email/verify/%s/%s?expires=%d&signature=%s", user.GetAuthIdentifier(), emailHash, expiresAt.Unix(), signature)

	if base != "" {
		link = base + link
	}

	if err := s.Mailer.Send(ctx, auth.MailMessage{
		To:      profile.GetEmail(),
		Subject: "Verify your email address",
		Body:    fmt.Sprintf("Verify your email address by visiting %s", link),
		Metadata: map[string]string{
			"link":      link,
			"expiresAt": strconv.FormatInt(expiresAt.Unix(), 10),
			"email":     verifiable.GetEmailForVerification(),
		},
	}); err != nil {
		return fmt.Errorf("send verification email: %w", err)
	}

	return nil
}

// Verify validates and fulfills an email verification link.
func (s *VerificationService) Verify(ctx context.Context, userID string, emailHash string, expiresAt int64, signature string) (auth.Authenticatable, error) {
	user, err := s.Users.RetrieveByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	verifiable, ok := user.(auth.MustVerifyEmail)

	if !ok {
		return nil, fmt.Errorf("foundation: user does not support email verification")
	}

	profile, ok := user.(auth.UserProfile)

	if !ok {
		return nil, fmt.Errorf("foundation: user does not expose profile email")
	}

	expectedHash := crypto.EmailHash(normalize(profile.GetEmail()))

	if expectedHash != emailHash {
		return nil, auth.ErrEmailVerificationInvalid
	}

	if err := s.Signer.Verify(ctx, "email-verification", []string{userID, emailHash}, expiresAt, signature); err != nil {
		return nil, err
	}

	verifiable.MarkEmailAsVerified(s.Clock.Now())

	if err := s.Users.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func normalize(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}
