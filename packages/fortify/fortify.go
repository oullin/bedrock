package fortify

import "errors"

// Fortify is the composition root for all authentication features.
// It holds references to contracts and action implementations
// provided by the consuming application.
type Fortify struct {
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

// Config returns the Fortify configuration.

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

// Builder constructs a Fortify instance with required dependencies.
type Builder struct {
	fortify *Fortify
	errors  []error
}

func (f *Fortify) Config() Config { return f.config }

func (f *Fortify) Guard() Guard { return f.guard }

func (f *Fortify) Provider() UserProvider { return f.provider }

func (f *Fortify) Hasher() PasswordHasher { return f.hasher }

func (f *Fortify) Broker() PasswordBroker { return f.broker }

func (f *Fortify) Verifier() EmailVerifier { return f.verifier }

func (f *Fortify) Events() EventDispatcher { return f.events }

func (f *Fortify) Limiter() RateLimiter { return f.limiter }

func (f *Fortify) Responder() Responder { return f.responder }

func (f *Fortify) CreateUser() CreatesNewUsers { return f.createUser }

func (f *Fortify) Authenticator() AuthenticatesUsers { return f.authenticator }

func (f *Fortify) UpdateProfile() UpdatesUserProfileInformation { return f.updateProfile }

func (f *Fortify) UpdatePassword() UpdatesUserPasswords { return f.updatePass }

func (f *Fortify) ResetPassword() ResetsUserPasswords { return f.resetPass }

func (f *Fortify) ConfirmPassword() ConfirmsPasswords { return f.confirmPass }

// NewBuilder creates a new Fortify builder.
func NewBuilder() *Builder {
	return &Builder{
		fortify: &Fortify{
			config: DefaultConfig(),
		},
	}
}

// WithConfig sets the configuration.
func (b *Builder) WithConfig(config Config) *Builder {
	b.fortify.config = config

	return b
}

// WithGuard sets the authentication guard.
func (b *Builder) WithGuard(guard Guard) *Builder {
	b.fortify.guard = guard

	return b
}

// WithProvider sets the user provider.
func (b *Builder) WithProvider(provider UserProvider) *Builder {
	b.fortify.provider = provider

	return b
}

// WithHasher sets the password hasher.
func (b *Builder) WithHasher(hasher PasswordHasher) *Builder {
	b.fortify.hasher = hasher

	return b
}

// WithBroker sets the password broker.
func (b *Builder) WithBroker(broker PasswordBroker) *Builder {
	b.fortify.broker = broker

	return b
}

// WithVerifier sets the email verifier.
func (b *Builder) WithVerifier(verifier EmailVerifier) *Builder {
	b.fortify.verifier = verifier

	return b
}

// WithEvents sets the event dispatcher.
func (b *Builder) WithEvents(events EventDispatcher) *Builder {
	b.fortify.events = events

	return b
}

// WithLimiter sets the rate limiter.
func (b *Builder) WithLimiter(limiter RateLimiter) *Builder {
	b.fortify.limiter = limiter

	return b
}

// WithResponder sets the response handler.
func (b *Builder) WithResponder(responder Responder) *Builder {
	b.fortify.responder = responder

	return b
}

// WithCreateUser sets the user creation action.
func (b *Builder) WithCreateUser(action CreatesNewUsers) *Builder {
	b.fortify.createUser = action

	return b
}

// WithAuthenticator sets a custom authenticator (optional).
func (b *Builder) WithAuthenticator(action AuthenticatesUsers) *Builder {
	b.fortify.authenticator = action

	return b
}

// WithUpdateProfile sets the profile update action.
func (b *Builder) WithUpdateProfile(action UpdatesUserProfileInformation) *Builder {
	b.fortify.updateProfile = action

	return b
}

// WithUpdatePassword sets the password update action.
func (b *Builder) WithUpdatePassword(action UpdatesUserPasswords) *Builder {
	b.fortify.updatePass = action

	return b
}

// WithResetPassword sets the password reset action.
func (b *Builder) WithResetPassword(action ResetsUserPasswords) *Builder {
	b.fortify.resetPass = action

	return b
}

// WithConfirmPassword sets the password confirmation action.
func (b *Builder) WithConfirmPassword(action ConfirmsPasswords) *Builder {
	b.fortify.confirmPass = action

	return b
}

// Build validates required dependencies and returns the Fortify instance.
func (b *Builder) Build() (*Fortify, error) {
	if b.fortify.guard == nil {
		b.errors = append(b.errors, errors.New("fortify: guard is required"))
	}

	if b.fortify.provider == nil {
		b.errors = append(b.errors, errors.New("fortify: user provider is required"))
	}

	if b.fortify.hasher == nil {
		b.errors = append(b.errors, errors.New("fortify: password hasher is required"))
	}

	if b.fortify.responder == nil {
		b.errors = append(b.errors, errors.New("fortify: responder is required"))
	}

	if b.fortify.config.Features.Registration && b.fortify.createUser == nil {
		b.errors = append(b.errors, errors.New("fortify: CreatesNewUsers action is required when registration is enabled"))
	}

	if b.fortify.config.Features.UpdateProfileInformation && b.fortify.updateProfile == nil {
		b.errors = append(b.errors, errors.New("fortify: UpdatesUserProfileInformation action is required when profile updates are enabled"))
	}

	if b.fortify.config.Features.UpdatePasswords && b.fortify.updatePass == nil {
		b.errors = append(b.errors, errors.New("fortify: UpdatesUserPasswords action is required when password updates are enabled"))
	}

	if b.fortify.config.Features.ResetPasswords && b.fortify.broker == nil {
		b.errors = append(b.errors, errors.New("fortify: password broker is required when password resets are enabled"))
	}

	if b.fortify.config.Features.ResetPasswords && b.fortify.resetPass == nil {
		b.errors = append(b.errors, errors.New("fortify: ResetsUserPasswords action is required when password resets are enabled"))
	}

	if b.fortify.config.Features.EmailVerification && b.fortify.verifier == nil {
		b.errors = append(b.errors, errors.New("fortify: email verifier is required when email verification is enabled"))
	}

	if b.fortify.config.Features.ConfirmPassword && b.fortify.confirmPass == nil {
		b.errors = append(b.errors, errors.New("fortify: ConfirmsPasswords action is required when password confirmation is enabled"))
	}

	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}

	return b.fortify, nil
}
