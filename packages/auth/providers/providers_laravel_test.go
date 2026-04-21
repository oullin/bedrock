package providers_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/bedrock/packages/auth"
	"github.com/bedrock/packages/auth/providers"
	cauth "github.com/bedrock/packages/contracts/auth"
)

type rowStub struct {
	row map[string]any
	err error
}

type dbStub struct {
	row       rowStub
	query     string
	queryArgs []any
	exec      string
	execArgs  []any
}

type modelQueryStub struct {
	user        cauth.Authenticatable
	credentials map[string]string
	token       string
	updated     bool
}

func (r rowStub) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	target, ok := dest[0].(*map[string]any)

	if !ok {
		return errors.New("expected map destination")
	}

	*target = r.row

	return nil
}

func (d *dbStub) QueryRow(_ context.Context, query string, args ...any) providers.DBRow {
	d.query = query
	d.queryArgs = args

	return d.row
}

func (d *dbStub) Exec(_ context.Context, query string, args ...any) error {
	d.exec = query
	d.execArgs = args

	return nil
}

func (m *modelQueryStub) FindByID(_ context.Context, _ string) (cauth.Authenticatable, error) {
	return m.user, nil
}

func (m *modelQueryStub) FindByToken(_ context.Context, _ string, token string) (cauth.Authenticatable, error) {
	m.token = token

	return m.user, nil
}

func (m *modelQueryStub) FindByCredentials(_ context.Context, credentials map[string]string) (cauth.Authenticatable, error) {
	m.credentials = credentials

	return m.user, nil
}

func (m *modelQueryStub) UpdateToken(_ context.Context, user cauth.Authenticatable, token string) error {
	m.updated = true
	user.SetRememberToken(token)

	return nil
}

// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testRetrieveByIDReturnsUserWhenUserIsFound
// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testRetrieveByIDReturnsNullWhenUserIsNotFound
func TestLaravelAuthDatabaseUserProviderRetrieveByID(t *testing.T) {
	db := &dbStub{row: rowStub{row: map[string]any{"id": "1", "password": "secret"}}}
	provider := providers.NewDatabaseUserProvider(db, "users", auth.NewBcryptHasher(4), func(row map[string]any) cauth.Authenticatable {
		return auth.NewGenericUser(row)
	})

	got, err := provider.RetrieveByID(context.Background(), "1")

	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.GetAuthIdentifier() != "1" {
		t.Fatalf("RetrieveByID returned %#v", got)
	}

	db.row.err = errors.New("not found")
	got, err = provider.RetrieveByID(context.Background(), "2")

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Fatal("RetrieveByID should return nil when the row is missing")
	}
}

// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testRetrieveByTokenReturnsUser
// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testRetrieveTokenWithBadIdentifierReturnsNull
// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testRetrieveByBadTokenReturnsNull
func TestLaravelAuthDatabaseUserProviderRetrieveByToken(t *testing.T) {
	db := &dbStub{row: rowStub{row: map[string]any{"id": "1", "remember_token": "token"}}}
	provider := providers.NewDatabaseUserProvider(db, "users", auth.NewBcryptHasher(4), func(row map[string]any) cauth.Authenticatable {
		return auth.NewGenericUser(row)
	})

	got, err := provider.RetrieveByToken(context.Background(), "1", "token")

	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.GetAuthIdentifier() != "1" {
		t.Fatalf("RetrieveByToken returned %#v", got)
	}

	if db.query != "SELECT * FROM users WHERE id = $1 AND remember_token = $2 LIMIT 1" {
		t.Errorf("query = %q", db.query)
	}

	if !reflect.DeepEqual(db.queryArgs, []any{"1", "token"}) {
		t.Errorf("args = %#v", db.queryArgs)
	}
}

// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testRetrieveByCredentialsReturnsUserWhenUserIsFound
// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testRetrieveByCredentialsReturnsNullWhenUserIsFound
// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testRetrieveByCredentialsWithMultiplyPasswordsReturnsNull
func TestLaravelAuthDatabaseUserProviderRetrieveByCredentialsFiltersPassword(t *testing.T) {
	db := &dbStub{row: rowStub{row: map[string]any{"id": "1", "email": "taylor@example.com"}}}
	provider := providers.NewDatabaseUserProvider(db, "users", auth.NewBcryptHasher(4), func(row map[string]any) cauth.Authenticatable {
		return auth.NewGenericUser(row)
	})

	got, err := provider.RetrieveByCredentials(context.Background(), map[string]string{
		"email":    "taylor@example.com",
		"password": "secret",
	})

	if err != nil {
		t.Fatal(err)
	}

	if got == nil {
		t.Fatal("RetrieveByCredentials should return the mapped user")
	}

	if db.query != "SELECT * FROM users WHERE email = $1 LIMIT 1" {
		t.Errorf("query = %q", db.query)
	}

	if !reflect.DeepEqual(db.queryArgs, []any{"taylor@example.com"}) {
		t.Errorf("args = %#v", db.queryArgs)
	}
}

// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testCredentialValidation
// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testCredentialValidationFails
// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testCredentialValidationFailsGracefullyWithNullPassword
func TestLaravelAuthDatabaseUserProviderValidateCredentials(t *testing.T) {
	hasher := auth.NewBcryptHasher(4)
	hash, err := hasher.Hash(context.Background(), "secret")

	if err != nil {
		t.Fatal(err)
	}

	user := auth.NewGenericUser(map[string]any{"id": "1", "password": hash})
	provider := providers.NewDatabaseUserProvider(&dbStub{}, "users", hasher, nil)

	ok, err := provider.ValidateCredentials(context.Background(), user, map[string]string{"password": "secret"})

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("ValidateCredentials should accept the correct password")
	}

	ok, err = provider.ValidateCredentials(context.Background(), user, map[string]string{"password": "wrong"})

	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Fatal("ValidateCredentials should reject the wrong password")
	}

	ok, err = provider.ValidateCredentials(context.Background(), user, map[string]string{})

	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Fatal("ValidateCredentials should reject an empty password")
	}
}

// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testRehashPasswordIfRequired
// Port of Illuminate\Tests\Auth\AuthDatabaseUserProviderTest::testDontRehashPasswordIfNotRequired
func TestLaravelAuthDatabaseUserProviderRehashPasswordIfRequired(t *testing.T) {
	hasher := auth.NewBcryptHasher(4)
	oldHash, err := hasher.Hash(context.Background(), "secret")

	if err != nil {
		t.Fatal(err)
	}

	db := &dbStub{}
	user := auth.NewGenericUser(map[string]any{"id": "1", "password": oldHash})
	provider := providers.NewDatabaseUserProvider(db, "users", hasher, nil)

	if err := provider.RehashPasswordIfRequired(context.Background(), user, map[string]string{"password": "secret"}, false); err != nil {
		t.Fatal(err)
	}

	if db.exec != "" {
		t.Fatal("RehashPasswordIfRequired should skip when the hash is current")
	}

	if err := provider.RehashPasswordIfRequired(context.Background(), user, map[string]string{"password": "secret"}, true); err != nil {
		t.Fatal(err)
	}

	if db.exec != "UPDATE users SET password = $1 WHERE id = $2" {
		t.Errorf("exec = %q", db.exec)
	}
}

// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testRetrieveByIDReturnsUser
// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testRetrieveByTokenReturnsUser
// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testRetrieveTokenWithBadIdentifierReturnsNull
// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testRetrieveByBadTokenReturnsNull
func TestLaravelAuthEloquentUserProviderRetrievesUsers(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "remember_token": "token"})
	user.SetRememberToken("token")
	model := &modelQueryStub{user: user}
	provider := providers.NewORMUserProvider(model, auth.NewBcryptHasher(4))

	got, err := provider.RetrieveByID(context.Background(), "1")

	if err != nil {
		t.Fatal(err)
	}

	if got != user {
		t.Fatal("RetrieveByID should return the model user")
	}

	got, err = provider.RetrieveByToken(context.Background(), "1", "token")

	if err != nil {
		t.Fatal(err)
	}

	if got != user {
		t.Fatal("RetrieveByToken should return the token-matched user")
	}

	got, err = provider.RetrieveByToken(context.Background(), "1", "wrong")

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Fatal("RetrieveByToken should reject a bad token")
	}
}

// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testRetrievingWithOnlyPasswordCredentialReturnsNull
// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testRetrieveByCredentialsReturnsUser
// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testRetrieveByCredentialsWithMultiplyPasswordsReturnsNull
func TestLaravelAuthEloquentUserProviderRetrieveByCredentialsFiltersPassword(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "email": "taylor@example.com"})
	model := &modelQueryStub{user: user}
	provider := providers.NewORMUserProvider(model, auth.NewBcryptHasher(4))

	got, err := provider.RetrieveByCredentials(context.Background(), map[string]string{
		"email":    "taylor@example.com",
		"password": "secret",
	})

	if err != nil {
		t.Fatal(err)
	}

	if got != user {
		t.Fatal("RetrieveByCredentials should return the model user")
	}

	if !reflect.DeepEqual(model.credentials, map[string]string{"email": "taylor@example.com"}) {
		t.Errorf("credentials = %#v", model.credentials)
	}
}

// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testCredentialValidation
// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testCredentialValidationFailed
// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testCredentialValidationFailsGracefullyWithNullPassword
// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testRehashPasswordIfRequired
// Port of Illuminate\Tests\Auth\AuthEloquentUserProviderTest::testDontRehashPasswordIfNotRequired
func TestLaravelAuthEloquentUserProviderCredentialsAndRehash(t *testing.T) {
	hasher := auth.NewBcryptHasher(4)
	hash, err := hasher.Hash(context.Background(), "secret")

	if err != nil {
		t.Fatal(err)
	}

	user := auth.NewGenericUser(map[string]any{"id": "1", "password": hash})
	provider := providers.NewORMUserProvider(&modelQueryStub{user: user}, hasher)

	ok, err := provider.ValidateCredentials(context.Background(), user, map[string]string{"password": "secret"})

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("ValidateCredentials should accept the correct password")
	}

	ok, err = provider.ValidateCredentials(context.Background(), user, map[string]string{})

	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Fatal("ValidateCredentials should reject an empty password")
	}

	if err := provider.RehashPasswordIfRequired(context.Background(), user, map[string]string{"password": "secret"}, false); err != nil {
		t.Fatal(err)
	}
}
