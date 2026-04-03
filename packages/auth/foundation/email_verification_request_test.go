package foundation_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/events"
	"github.com/gollin/packages/auth/foundation"
	configpkg "github.com/gollin/packages/config"
	securitycrypto "github.com/gollin/packages/security/crypto"
	"github.com/gollin/packages/user"
)

type verifiedDispatcher struct {
	event *events.Verified
}

func (d *verifiedDispatcher) DispatchVerified(_ context.Context, event events.Verified) error {
	d.event = &event

	return nil
}

func TestEmailVerificationRequestAuthorizeAndFulfill(t *testing.T) {
	t.Parallel()

	_, currentFile, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	configDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "user", "config")
	repo, err := configpkg.NewBuilder(configDir).Build(context.Background())

	if err != nil {
		t.Fatalf("build user config: %v", err)
	}

	cfg, err := user.ConfigFromRepository(repo)

	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	record := user.New(cfg, "user-1", "", "person@example.com", time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC))
	request := foundation.EmailVerificationRequest{
		User:      record,
		RouteID:   record.GetAuthIdentifier(),
		RouteHash: securitycrypto.EmailHash(record.GetEmailForVerification()),
	}

	if !request.Authorize() {
		t.Fatal("expected request to authorize")
	}

	dispatcher := &verifiedDispatcher{}
	now := time.Date(2026, 4, 3, 1, 0, 0, 0, time.UTC)

	if err := request.FulfillAt(context.Background(), now, dispatcher); err != nil {
		t.Fatalf("FulfillAt: %v", err)
	}

	verifiable := any(record).(auth.MustVerifyEmail)

	if !verifiable.HasVerifiedEmail() {
		t.Fatal("expected user to be verified")
	}

	if dispatcher.event == nil || dispatcher.event.User.GetAuthIdentifier() != record.GetAuthIdentifier() {
		t.Fatal("expected verified event to be dispatched")
	}
}
