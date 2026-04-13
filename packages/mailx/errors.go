package mailx

import "errors"

var (
	ErrInvalidDriver   = errors.New("mail: invalid driver")
	ErrInvalidAddress  = errors.New("mail: invalid email address")
	ErrNoRecipients    = errors.New("mail: no recipients specified")
	ErrSendFailed      = errors.New("mail: send failed")
	ErrTransportClosed = errors.New("mail: transport closed")
	ErrNoContent       = errors.New("mail: no content specified")
	ErrNoQueue         = errors.New("mail: no queue configured")
	ErrMailerNotFound  = errors.New("mail: mailer not found in configuration")
)
