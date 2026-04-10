package authflows

import "errors"

// AuthFlows is the composition root for all authentication features.
// It holds references to contracts and action implementations
// provided by the consuming application.
type AuthFlows struct {
	config    Config
	guard     Guard
	provider  UserProvider
	hasher    PasswordHasher
	broker    PasswordBroker
	verifier  EmailVerifier
	events    EventDispatcher
	limiter   RateLimiter
	responder Responder

	createUser    CreatesNewUsers
	authenticator AuthenticatesUsers
	updateProfile UpdatesUserProfileInformation
	updatePass    UpdatesUserPasswords
	resetPass     ResetsUserPasswords
	confirmPass   ConfirmsPasswords
}

// Config returns the AuthFlows configuration.

// Guard returns the configured guard.

// Provider returns the configured user provider.

// Hasher returns the configured password hasher.

// Broker returns the configured password broker.

// Verifier returns the configured email verifier.

// Events returns the configured event dispatcher.

// Limiter returns the configured rate limiter.

// Responder returns the configured response handler.

// CreateUser returns the user creation action.

// Authenticator returns the custom authenticator, if any.

// UpdateProfile returns the profile update action.

// UpdatePassword returns the password update action.

// ResetPassword returns the password reset action.

// ConfirmPassword returns the password confirmation action.

// Builder constructs a AuthFlows instance with required dependencies.
type Builder struct {
	authflows *AuthFlows
	errors  []error
}

func (f *AuthFlows) Config() Config { return f.config }

func (f *AuthFlows) Guard() Guard { return f.guard }

func (f *AuthFlows) Provider() UserProvider { return f.provider }

func (f *AuthFlows) Hasher() PasswordHasher { return f.hasher }

func (f *AuthFlows) Broker() PasswordBroker { return f.broker }

func (f *AuthFlows) Verifier() EmailVerifier { return f.verifier }

func (f *AuthFlows) Events() EventDispatcher { return f.events }

func (f *AuthFlows) Limiter() RateLimiter { return f.limiter }

func (f *AuthFlows) Responder() Responder { return f.responder }

func (f *AuthFlows) CreateUser() CreatesNewUsers { return f.createUser }

func (f *AuthFlows) Authenticator() AuthenticatesUsers { return f.authenticator }

func (f *AuthFlows) UpdateProfile() UpdatesUserProfileInformation { return f.updateProfile }

func (f *AuthFlows) UpdatePassword() UpdatesUserPasswords { return f.updatePass }

func (f *AuthFlows) ResetPassword() ResetsUserPasswords { return f.resetPass }

func (f *AuthFlows) ConfirmPassword() ConfirmsPasswords { return f.confirmPass }

// NewBuilder creates a new AuthFlows builder.
func NewBuilder() *Builder {
	return &Builder{
		authflows: &AuthFlows{
			config: DefaultConfig(),
		},
	}
}

// WithConfig sets the configuration.
func (b *Builder) WithConfig(config Config) *Builder {
	b.authflows.config = config

	return b
}

// WithGuard sets the authentication guard.
func (b *Builder) WithGuard(guard Guard) *Builder {
	b.authflows.guard = guard

	return b
}

// WithProvider sets the user provider.
func (b *Builder) WithProvider(provider UserProvider) *Builder {
	b.authflows.provider = provider

	return b
}

// WithHasher sets the password hasher.
func (b *Builder) WithHasher(hasher PasswordHasher) *Builder {
	b.authflows.hasher = hasher

	return b
}

// WithBroker sets the password broker.
func (b *Builder) WithBroker(broker PasswordBroker) *Builder {
	b.authflows.broker = broker

	return b
}

// WithVerifier sets the email verifier.
func (b *Builder) WithVerifier(verifier EmailVerifier) *Builder {
	b.authflows.verifier = verifier

	return b
}

// WithEvents sets the event dispatcher.
func (b *Builder) WithEvents(events EventDispatcher) *Builder {
	b.authflows.events = events

	return b
}

// WithLimiter sets the rate limiter.
func (b *Builder) WithLimiter(limiter RateLimiter) *Builder {
	b.authflows.limiter = limiter

	return b
}

// WithResponder sets the response handler.
func (b *Builder) WithResponder(responder Responder) *Builder {
	b.authflows.responder = responder

	return b
}

// WithCreateUser sets the user creation action.
func (b *Builder) WithCreateUser(action CreatesNewUsers) *Builder {
	b.authflows.createUser = action

	return b
}

// WithAuthenticator sets a custom authenticator (optional).
func (b *Builder) WithAuthenticator(action AuthenticatesUsers) *Builder {
	b.authflows.authenticator = action

	return b
}

// WithUpdateProfile sets the profile update action.
func (b *Builder) WithUpdateProfile(action UpdatesUserProfileInformation) *Builder {
	b.authflows.updateProfile = action

	return b
}

// WithUpdatePassword sets the password update action.
func (b *Builder) WithUpdatePassword(action UpdatesUserPasswords) *Builder {
	b.authflows.updatePass = action

	return b
}

// WithResetPassword sets the password reset action.
func (b *Builder) WithResetPassword(action ResetsUserPasswords) *Builder {
	b.authflows.resetPass = action

	return b
}

// WithConfirmPassword sets the password confirmation action.
func (b *Builder) WithConfirmPassword(action ConfirmsPasswords) *Builder {
	b.authflows.confirmPass = action

	return b
}

// Build validates required dependencies and returns the AuthFlows instance.
func (b *Builder) Build() (*AuthFlows, error) {
	if b.authflows.guard == nil {
		b.errors = append(b.errors, errors.New("authflows: guard is required"))
	}

	if b.authflows.provider == nil {
		b.errors = append(b.errors, errors.New("authflows: user provider is required"))
	}

	if b.authflows.hasher == nil {
		b.errors = append(b.errors, errors.New("authflows: password hasher is required"))
	}

	if b.authflows.responder == nil {
		b.errors = append(b.errors, errors.New("authflows: responder is required"))
	}

	if b.authflows.config.Features.Registration && b.authflows.createUser == nil {
		b.errors = append(b.errors, errors.New("authflows: CreatesNewUsers action is required when registration is enabled"))
	}

	if b.authflows.config.Features.UpdateProfileInformation && b.authflows.updateProfile == nil {
		b.errors = append(b.errors, errors.New("authflows: UpdatesUserProfileInformation action is required when profile updates are enabled"))
	}

	if b.authflows.config.Features.UpdatePasswords && b.authflows.updatePass == nil {
		b.errors = append(b.errors, errors.New("authflows: UpdatesUserPasswords action is required when password updates are enabled"))
	}

	if b.authflows.config.Features.ResetPasswords && b.authflows.broker == nil {
		b.errors = append(b.errors, errors.New("authflows: password broker is required when password resets are enabled"))
	}

	if b.authflows.config.Features.ResetPasswords && b.authflows.resetPass == nil {
		b.errors = append(b.errors, errors.New("authflows: ResetsUserPasswords action is required when password resets are enabled"))
	}

	if b.authflows.config.Features.EmailVerification && b.authflows.verifier == nil {
		b.errors = append(b.errors, errors.New("authflows: email verifier is required when email verification is enabled"))
	}

	if b.authflows.config.Features.ConfirmPassword && b.authflows.confirmPass == nil {
		b.errors = append(b.errors, errors.New("authflows: ConfirmsPasswords action is required when password confirmation is enabled"))
	}

	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}

	return b.authflows, nil
}
