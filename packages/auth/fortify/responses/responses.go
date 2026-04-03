package responses

import (
	"encoding/json"
	"errors"
	"net/http"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/fortify/contracts"
	authmw "github.com/gollin/packages/auth/middleware"
)

// Registry stores the default Fortify responses.
type Registry struct {
	LoginViewResponse                          contracts.LoginViewResponse
	RegisterViewResponse                       contracts.RegisterViewResponse
	RequestPasswordResetLinkViewResponse       contracts.RequestPasswordResetLinkViewResponse
	ResetPasswordViewResponse                  contracts.ResetPasswordViewResponse
	VerifyEmailViewResponse                    contracts.VerifyEmailViewResponse
	ConfirmPasswordViewResponse                contracts.ConfirmPasswordViewResponse
	TwoFactorChallengeViewResponse             contracts.TwoFactorChallengeViewResponse
	RegisterResponse                           contracts.RegisterResponse
	LoginResponse                              contracts.LoginResponse
	LogoutResponse                             contracts.LogoutResponse
	PasswordResetResponse                      contracts.PasswordResetResponse
	SuccessfulPasswordResetLinkRequestResponse contracts.SuccessfulPasswordResetLinkRequestResponse
	FailedPasswordResetLinkRequestResponse     contracts.FailedPasswordResetLinkRequestResponse
	FailedPasswordResetResponse                contracts.FailedPasswordResetResponse
	FailedPasswordConfirmationResponse         contracts.FailedPasswordConfirmationResponse
	PasswordConfirmedResponse                  contracts.PasswordConfirmedResponse
	VerifyEmailResponse                        contracts.VerifyEmailResponse
	EmailVerificationNotificationSentResponse  contracts.EmailVerificationNotificationSentResponse
	TwoFactorEnabledResponse                   contracts.TwoFactorEnabledResponse
	TwoFactorConfirmedResponse                 contracts.TwoFactorConfirmedResponse
	TwoFactorDisabledResponse                  contracts.TwoFactorDisabledResponse
	RecoveryCodesGeneratedResponse             contracts.RecoveryCodesGeneratedResponse
	TwoFactorLoginResponse                     contracts.TwoFactorLoginResponse
	FailedTwoFactorLoginResponse               contracts.FailedTwoFactorLoginResponse
	ProfileInformationUpdatedResponse          contracts.ProfileInformationUpdatedResponse
	PasswordUpdateResponse                     contracts.PasswordUpdateResponse
	LockoutResponse                            contracts.LockoutResponse
}

// DefaultRegistry creates the default JSON response registry.
func DefaultRegistry() *Registry {
	return &Registry{
		LoginViewResponse:                          JSONResponse{Status: http.StatusOK},
		RegisterViewResponse:                       JSONResponse{Status: http.StatusOK},
		RequestPasswordResetLinkViewResponse:       JSONResponse{Status: http.StatusOK},
		ResetPasswordViewResponse:                  JSONResponse{Status: http.StatusOK},
		VerifyEmailViewResponse:                    JSONResponse{Status: http.StatusOK},
		ConfirmPasswordViewResponse:                JSONResponse{Status: http.StatusOK},
		TwoFactorChallengeViewResponse:             JSONResponse{Status: http.StatusOK},
		RegisterResponse:                           JSONResponse{Status: http.StatusCreated},
		LoginResponse:                              JSONResponse{Status: http.StatusOK},
		LogoutResponse:                             JSONResponse{Status: http.StatusOK},
		PasswordResetResponse:                      JSONResponse{Status: http.StatusOK},
		SuccessfulPasswordResetLinkRequestResponse: JSONResponse{Status: http.StatusAccepted},
		FailedPasswordResetLinkRequestResponse:     ErrorResponse{Status: http.StatusUnprocessableEntity},
		FailedPasswordResetResponse:                ErrorResponse{Status: http.StatusUnprocessableEntity},
		FailedPasswordConfirmationResponse:         ErrorResponse{Status: http.StatusForbidden},
		PasswordConfirmedResponse:                  JSONResponse{Status: http.StatusOK},
		VerifyEmailResponse:                        JSONResponse{Status: http.StatusOK},
		EmailVerificationNotificationSentResponse:  JSONResponse{Status: http.StatusAccepted},
		TwoFactorEnabledResponse:                   JSONResponse{Status: http.StatusCreated},
		TwoFactorConfirmedResponse:                 JSONResponse{Status: http.StatusOK},
		TwoFactorDisabledResponse:                  JSONResponse{Status: http.StatusOK},
		RecoveryCodesGeneratedResponse:             JSONResponse{Status: http.StatusOK},
		TwoFactorLoginResponse:                     JSONResponse{Status: http.StatusOK},
		FailedTwoFactorLoginResponse:               ErrorResponse{Status: http.StatusUnauthorized},
		ProfileInformationUpdatedResponse:          JSONResponse{Status: http.StatusOK},
		PasswordUpdateResponse:                     JSONResponse{Status: http.StatusOK},
		LockoutResponse:                            ErrorResponse{Status: http.StatusTooManyRequests},
	}
}

// JSONResponse writes a JSON success response.
type JSONResponse struct {
	Status int
}

// ToResponse writes the payload.
func (r JSONResponse) ToResponse(w http.ResponseWriter, _ *http.Request, payload any) {
	if err, ok := payload.(error); ok && err != nil {
		authmw.WriteJSONError(w, inferStatus(err, r.Status), err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	if payload == nil {
		payload = map[string]any{}
	}
	_ = json.NewEncoder(w).Encode(payload)
}

// ErrorResponse writes an auth-aware JSON error response.
type ErrorResponse struct {
	Status int
}

// ToResponse writes the error payload.
func (r ErrorResponse) ToResponse(w http.ResponseWriter, _ *http.Request, payload any) {
	err, _ := payload.(error)
	if err == nil {
		err = http.ErrAbortHandler
	}
	authmw.WriteJSONError(w, inferStatus(err, r.Status), err)
}

func inferStatus(err error, fallback int) int {
	var validationErr *auth.ValidationError
	var throttleErr *auth.ThrottleError
	switch {
	case errors.As(err, &validationErr):
		return http.StatusUnprocessableEntity
	case errors.As(err, &throttleErr):
		return http.StatusTooManyRequests
	case errors.Is(err, auth.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, auth.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, auth.ErrInvalidToken):
		return http.StatusUnprocessableEntity
	case errors.Is(err, auth.ErrTokenExpired):
		return http.StatusGone
	case errors.Is(err, auth.ErrTwoFactorInvalid):
		return http.StatusUnauthorized
	case errors.Is(err, auth.ErrTwoFactorRequired):
		return http.StatusConflict
	case errors.Is(err, auth.ErrPasswordConfirmationRequired):
		return http.StatusForbidden
	case errors.Is(err, auth.ErrUserExists):
		return http.StatusConflict
	case errors.Is(err, auth.ErrEmailVerificationInvalid):
		return http.StatusUnprocessableEntity
	default:
		if fallback >= 400 {
			return fallback
		}
		return http.StatusInternalServerError
	}
}
