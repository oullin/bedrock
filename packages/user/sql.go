package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	auth "github.com/gollin/packages/auth"
)

// SQLRepository stores users in SQL.
type SQLRepository struct {
	db       *sql.DB
	table    string
	password auth.PasswordHasher
}

// NewSQLRepository creates a new SQL-backed user repository.
func NewSQLRepository(db *sql.DB, hasher auth.PasswordHasher, table string) (*SQLRepository, error) {
	password, err := auth.EnsureHasher(hasher)

	if err != nil {
		return nil, err
	}

	if table == "" {
		table = "users"
	}

	return &SQLRepository{
		db:       db,
		table:    table,
		password: password,
	}, nil
}

// Create inserts a user.
func (r *SQLRepository) Create(ctx context.Context, user auth.Authenticatable) error {
	record, ok := user.(*User)

	if !ok {
		return fmt.Errorf("user: unsupported type %T", user)
	}

	payload, err := json.Marshal(record.TwoFactorRecoveryCodes)

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `INSERT INTO `+r.table+` (id, name, email, password_hash, api_token, remember_token, email_verified_at, two_factor_secret, two_factor_recovery_codes, two_factor_confirmed_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		record.ID,
		record.Name,
		normalizeEmail(record.Email),
		record.PasswordHash,
		record.APIToken,
		record.RememberToken,
		record.EmailVerifiedAt,
		record.TwoFactorSecret,
		string(payload),
		record.TwoFactorConfirmedAt,
		record.CreatedAt,
		record.UpdatedAt,
	)

	return err
}

// Update replaces a stored user.
func (r *SQLRepository) Update(ctx context.Context, user auth.Authenticatable) error {
	record, ok := user.(*User)

	if !ok {
		return fmt.Errorf("user: unsupported type %T", user)
	}

	payload, err := json.Marshal(record.TwoFactorRecoveryCodes)

	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, `UPDATE `+r.table+` SET name = ?, email = ?, password_hash = ?, api_token = ?, remember_token = ?, email_verified_at = ?, two_factor_secret = ?, two_factor_recovery_codes = ?, two_factor_confirmed_at = ?, updated_at = ? WHERE id = ?`,
		record.Name,
		normalizeEmail(record.Email),
		record.PasswordHash,
		record.APIToken,
		record.RememberToken,
		record.EmailVerifiedAt,
		record.TwoFactorSecret,
		string(payload),
		record.TwoFactorConfirmedAt,
		time.Now().UTC(),
		record.ID,
	)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return auth.ErrUserNotFound
	}

	return nil
}

// DeleteByID deletes a user by id.
func (r *SQLRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM `+r.table+` WHERE id = ?`, id)

	return err
}

// FindByEmail retrieves a user by email.
func (r *SQLRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	return r.findOne(ctx, `SELECT id, name, email, password_hash, api_token, remember_token, email_verified_at, two_factor_secret, two_factor_recovery_codes, two_factor_confirmed_at, created_at, updated_at FROM `+r.table+` WHERE email = ?`, normalizeEmail(email))
}

// RetrieveByID retrieves a user by id.
func (r *SQLRepository) RetrieveByID(ctx context.Context, id string) (auth.Authenticatable, error) {
	return r.findOne(ctx, `SELECT id, name, email, password_hash, api_token, remember_token, email_verified_at, two_factor_secret, two_factor_recovery_codes, two_factor_confirmed_at, created_at, updated_at FROM `+r.table+` WHERE id = ?`, id)
}

// RetrieveByToken retrieves a user by remember token.
func (r *SQLRepository) RetrieveByToken(ctx context.Context, id string, token string) (auth.Authenticatable, error) {
	return r.findOne(ctx, `SELECT id, name, email, password_hash, api_token, remember_token, email_verified_at, two_factor_secret, two_factor_recovery_codes, two_factor_confirmed_at, created_at, updated_at FROM `+r.table+` WHERE id = ? AND remember_token = ?`, id, token)
}

// RetrieveByCredentials retrieves a user by auth credentials.
func (r *SQLRepository) RetrieveByCredentials(ctx context.Context, credentials map[string]string) (auth.Authenticatable, error) {
	if email, ok := credentials["email"]; ok {
		return r.FindByEmail(ctx, email)
	}

	if token, ok := credentials["api_token"]; ok {
		return r.findOne(ctx, `SELECT id, name, email, password_hash, api_token, remember_token, email_verified_at, two_factor_secret, two_factor_recovery_codes, two_factor_confirmed_at, created_at, updated_at FROM `+r.table+` WHERE api_token = ?`, token)
	}

	return nil, auth.ErrUserNotFound
}

// UpdateRememberToken persists a remember token.
func (r *SQLRepository) UpdateRememberToken(ctx context.Context, user auth.Authenticatable, token string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE `+r.table+` SET remember_token = ?, updated_at = ? WHERE id = ?`, token, time.Now().UTC(), user.GetAuthIdentifier())

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return auth.ErrUserNotFound
	}

	return nil
}

// ValidateCredentials compares a provided password against the stored hash.
func (r *SQLRepository) ValidateCredentials(ctx context.Context, user auth.Authenticatable, credentials map[string]string) (bool, error) {
	password, ok := credentials["password"]

	if !ok {
		return true, nil
	}

	return r.password.Check(ctx, password, user.GetAuthPassword(), nil)
}

// RehashPasswordIfRequired updates a password hash when needed.
func (r *SQLRepository) RehashPasswordIfRequired(ctx context.Context, user auth.Authenticatable, credentials map[string]string, force bool) error {
	password, ok := credentials["password"]

	if !ok || strings.TrimSpace(password) == "" {
		return nil
	}

	if !force && !r.password.NeedsRehash(user.GetAuthPassword(), nil) {
		return nil
	}

	hash, err := r.password.Hash(ctx, password)

	if err != nil {
		return err
	}

	user.SetAuthPassword(hash)

	return r.Update(ctx, user)
}

func (r *SQLRepository) findOne(ctx context.Context, query string, args ...any) (*User, error) {
	row := r.db.QueryRowContext(ctx, query, args...)

	var record User

	var emailVerifiedAt sql.NullTime

	var twoFactorConfirmedAt sql.NullTime

	var recoveryCodes string

	if err := row.Scan(
		&record.ID,
		&record.Name,
		&record.Email,
		&record.PasswordHash,
		&record.APIToken,
		&record.RememberToken,
		&emailVerifiedAt,
		&record.TwoFactorSecret,
		&recoveryCodes,
		&twoFactorConfirmedAt,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}

		return nil, err
	}

	if emailVerifiedAt.Valid {
		value := emailVerifiedAt.Time
		record.EmailVerifiedAt = &value
	}

	if twoFactorConfirmedAt.Valid {
		value := twoFactorConfirmedAt.Time
		record.TwoFactorConfirmedAt = &value
	}

	if strings.TrimSpace(recoveryCodes) != "" {
		if err := json.Unmarshal([]byte(recoveryCodes), &record.TwoFactorRecoveryCodes); err != nil {
			return nil, err
		}
	}

	return &record, nil
}
