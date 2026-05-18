package credentials_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscreds "github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/bedrock/packages/queue/drivers/sqs/credentials"
)

// stubFactory returns a Factory that yields the given provider/error
// pair and increments the call counter so tests can assert on cache
// behaviour.
func stubFactory(provider credentials.Provider, err error, calls *int) credentials.Factory {
	return func(_ context.Context) (credentials.Provider, error) {
		*calls++

		return provider, err
	}
}

func TestRegistryResolveCachesAfterFirstCall(t *testing.T) {
	t.Parallel()

	reg := credentials.NewRegistry()
	stub := awscreds.NewStaticCredentialsProvider("AKIA-cache", "secret", "")
	calls := 0

	reg.Register("prod", stubFactory(stub, nil, &calls))

	if _, err := reg.Resolve(context.Background(), "prod"); err != nil {
		t.Fatalf("first resolve: %v", err)
	}

	got, err := reg.Resolve(context.Background(), "prod")

	if err != nil {
		t.Fatalf("second resolve: %v", err)
	}

	value, err := got.Retrieve(context.Background())

	if err != nil {
		t.Fatalf("retrieve: %v", err)
	}

	if value.AccessKeyID != "AKIA-cache" {
		t.Errorf("AccessKeyID: got %q, want AKIA-cache", value.AccessKeyID)
	}

	if calls != 1 {
		t.Errorf("factory call count: got %d, want 1 (cached after first)", calls)
	}
}

func TestRegistryResolveUnknownProvider(t *testing.T) {
	t.Parallel()

	reg := credentials.NewRegistry()

	_, err := reg.Resolve(context.Background(), "missing")

	if !errors.Is(err, credentials.ErrUnknownProvider) {
		t.Fatalf("error: got %v, want ErrUnknownProvider", err)
	}
}

func TestRegistryRegisterReplacesAndInvalidates(t *testing.T) {
	t.Parallel()

	reg := credentials.NewRegistry()
	first := awscreds.NewStaticCredentialsProvider("k1", "s1", "")
	second := awscreds.NewStaticCredentialsProvider("k2", "s2", "")

	firstCalls := 0
	secondCalls := 0

	reg.Register("svc", stubFactory(first, nil, &firstCalls))

	got, err := reg.Resolve(context.Background(), "svc")

	if err != nil {
		t.Fatalf("resolve first: %v", err)
	}

	value, err := got.Retrieve(context.Background())

	if err != nil || value.AccessKeyID != "k1" {
		t.Fatalf("first resolve key: got %+v err=%v, want k1", value, err)
	}

	reg.Register("svc", stubFactory(second, nil, &secondCalls))

	got, err = reg.Resolve(context.Background(), "svc")

	if err != nil {
		t.Fatalf("resolve second: %v", err)
	}

	value, err = got.Retrieve(context.Background())

	if err != nil || value.AccessKeyID != "k2" {
		t.Errorf("post-replace key: got %+v err=%v, want k2 (cache should have been invalidated)", value, err)
	}

	if secondCalls != 1 {
		t.Errorf("replacement factory invoked %d times, want 1", secondCalls)
	}
}

func TestRegistryNamesIsSorted(t *testing.T) {
	t.Parallel()

	reg := credentials.NewRegistry()

	reg.Register("zeta", credentials.Static("a", "b", ""))
	reg.Register("alpha", credentials.Static("a", "b", ""))
	reg.Register("mike", credentials.Static("a", "b", ""))

	got := reg.Names()
	want := []string{"alpha", "mike", "zeta"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("names: got %v, want %v", got, want)
	}
}

func TestStaticFactoryReturnsExpectedProvider(t *testing.T) {
	t.Parallel()

	factory := credentials.Static("AKIA", "secret", "tok")
	provider, err := factory(context.Background())

	if err != nil {
		t.Fatalf("factory: %v", err)
	}

	value, err := provider.Retrieve(context.Background())

	if err != nil {
		t.Fatalf("retrieve: %v", err)
	}

	if value.AccessKeyID != "AKIA" || value.SecretAccessKey != "secret" || value.SessionToken != "tok" {
		t.Errorf("retrieved credentials wrong: %+v", value)
	}
}

func TestChainPicksFirstSuccess(t *testing.T) {
	t.Parallel()

	failing := func(_ context.Context) (credentials.Provider, error) {
		return nil, errors.New("nope")
	}

	winning := credentials.Static("key", "secret", "")
	chain := credentials.Chain(failing, winning)

	provider, err := chain(context.Background())

	if err != nil {
		t.Fatalf("chain: %v", err)
	}

	value, err := provider.Retrieve(context.Background())

	if err != nil {
		t.Fatalf("retrieve: %v", err)
	}

	if value.AccessKeyID != "key" {
		t.Errorf("chain selected wrong provider: %+v", value)
	}
}

func TestChainBubblesLastError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("doomed")

	chain := credentials.Chain(
		func(_ context.Context) (credentials.Provider, error) { return nil, errors.New("first") },
		func(_ context.Context) (credentials.Provider, error) { return nil, sentinel },
	)

	_, err := chain(context.Background())

	if !errors.Is(err, sentinel) {
		t.Errorf("chain error: got %v, want sentinel", err)
	}
}

// Compile-time guard: credentials.Provider must remain a strict alias
// for aws.CredentialsProvider so callers can pass values across the
// boundary without intermediate conversions.
var _ aws.CredentialsProvider = (credentials.Provider)(nil)
