package broadcasting_test

import (
	"reflect"
	"testing"

	"github.com/bedrock/packages/broadcasting"
	contractsauth "github.com/bedrock/packages/contracts/auth"
)

// BroadcasterTest::testExtractingParametersWhileCheckingForUserAccess
// BroadcasterTest::testCanUseChannelClasses
// BroadcasterTest::testModelRouteBinding

// BroadcasterTest::testUnknownChannelAuthHandlerTypeThrowsException
// BroadcasterTest::testNotFoundThrowsHttpException

// BroadcasterTest::testCanRegisterChannelsAsClasses
// BroadcasterTest::testCanRegisterChannelsWithoutOptions
// BroadcasterTest::testCanRegisterChannelsWithOptions
// BroadcasterTest::testCanRetrieveChannelsOptions
// BroadcasterTest::testCanRetrieveChannelsOptionsUsingAChannelNameContainingArgs
// BroadcasterTest::testCanRetrieveChannelsOptionsWhenMultipleChannelsAreRegistered
// BroadcasterTest::testDontRetrieveChannelsOptionsWhenChannelDoesntExists

// BroadcasterTest::testRetrieveUserWithoutGuard
// BroadcasterTest::testRetrieveUserWithOneGuardUsingAStringForSpecifyingGuard
// BroadcasterTest::testRetrieveUserWithMultipleGuardsAndRespectGuardsOrder
// BroadcasterTest::testRetrieveUserDontUseDefaultGuardWhenOneGuardSpecified
// BroadcasterTest::testRetrieveUserDontUseDefaultGuardWhenMultipleGuardsSpecified

// BroadcasterTest::testUserAuthenticationWithValidUser
// BroadcasterTest::testUserAuthenticationWithInvalidUser
// BroadcasterTest::testUserAuthenticationWithoutResolve

// BroadcasterTest::testChannelNameMatchPattern

type dummyChannel struct{}

func TestBroadcasterExtractsParametersAndBindings(t *testing.T) {
	t.Parallel()

	b := broadcasting.NewBaseBroadcaster().
		Bind("model", modelBinding("model")).
		Bind("model2", modelBinding("model"))

	params, err := b.ExtractAuthParameters("asd.{model}.{nonModel}", "asd.1.something")

	if err != nil {
		t.Fatalf("ExtractAuthParameters returned error: %v", err)
	}

	requireEqual(t, params, []any{"model.1.instance", "something"})

	params, err = b.ExtractAuthParameters("asd.{model}.{model2}.{nonModel}", "asd.1.uid.something")

	if err != nil {
		t.Fatalf("ExtractAuthParameters returned error: %v", err)
	}

	requireEqual(t, params, []any{"model.1.instance", "model.uid.instance", "something"})

	joiner := dummyChannel{}
	handler, err := broadcasting.NewBaseBroadcaster().
		Channel("test", joiner).
		VerifyUserCanAccessChannel(requestWithUser("test"), "test")

	if err != nil {
		t.Fatalf("class-based joiner returned error: %v", err)
	}

	requireEqual(t, handler, []any{"joined"})
}

func TestBroadcasterRejectsUnknownHandlersAndMissingBindings(t *testing.T) {
	t.Parallel()

	b := broadcasting.NewBaseBroadcaster()
	b.Channel("asd.{model}", 123)

	_, err := b.VerifyUserCanAccessChannel(requestWithUser("asd.1"), "asd.1")

	if err == nil {
		t.Fatal("expected unknown handler error")
	}

	b = broadcasting.NewBaseBroadcaster().
		Bind("model", func(string) (any, error) { return nil, broadcasting.ErrAccessDenied })
	b.Channel("asd.{model}", allowHandler(true))

	_, err = b.VerifyUserCanAccessChannel(requestWithUser("asd.1"), "asd.1")
	assertAccessDenied(t, err)
}

func TestBroadcasterRegistersChannelsAndOptions(t *testing.T) {
	t.Parallel()

	b := broadcasting.NewBaseBroadcaster()
	b.Channel("somechannel", allowHandler(true))
	b.Channel("someotherchannel", dummyChannel{}, broadcasting.WithGuards("a", "b"))
	b.Channel("somechannel.{id}.test.{text}", allowHandler(true), broadcasting.WithGuards("arg"))

	requireEqual(t, b.RetrieveChannelOptions("someotherchannel").Guards, []string{"a", "b"})
	requireEqual(t, b.RetrieveChannelOptions("somechannel.23.test.mytext").Guards, []string{"arg"})
	requireEqual(t, b.RetrieveChannelOptions("missing").Guards, []string(nil))
}

func TestBroadcasterRetrievesUsersWithGuards(t *testing.T) {
	t.Parallel()

	defaultUser := &testUser{id: "default", broadcastID: "default"}
	guardUser := &testUser{id: "guard", broadcastID: "guard"}

	resolver := &userResolver{users: map[string]contractsauth.Authenticatable{"": defaultUser, "myguard2": guardUser}}
	request := broadcasting.AuthRequest{ChannelName: "somechannel", UserResolver: resolver}

	b := broadcasting.NewBaseBroadcaster()
	b.Channel("somechannel", allowHandler(true))

	if got := b.RetrieveUser(request, "somechannel"); got != defaultUser {
		t.Fatalf("default guard user = %#v, want %#v", got, defaultUser)
	}

	b.Channel("guarded", allowHandler(true), broadcasting.WithGuards("myguard1", "myguard2"))

	if got := b.RetrieveUser(request, "guarded"); got != guardUser {
		t.Fatalf("guarded user = %#v, want %#v", got, guardUser)
	}

	requireEqual(t, resolver.calls, []string{"", "myguard1", "myguard2"})
}

func TestBroadcasterAuthenticatedUserCallback(t *testing.T) {
	t.Parallel()

	b := broadcasting.NewBaseBroadcaster()

	if got := b.ResolveAuthenticatedUser(broadcasting.AuthRequest{SocketID: "1234.1234"}); got != nil {
		t.Fatalf("expected nil without callback, got %#v", got)
	}

	b.ResolveAuthenticatedUserUsing(func(request broadcasting.AuthRequest) any {
		return map[string]any{"id": "12345", "socket": request.SocketID}
	})

	requireEqual(t, b.ResolveAuthenticatedUser(broadcasting.AuthRequest{SocketID: "1234.1234"}), map[string]any{
		"id":     "12345",
		"socket": "1234.1234",
	})

	b.ResolveAuthenticatedUserUsing(func(broadcasting.AuthRequest) any { return nil })

	if got := b.ResolveAuthenticatedUser(broadcasting.AuthRequest{SocketID: "1234.1234"}); got != nil {
		t.Fatalf("expected nil from callback, got %#v", got)
	}
}

func TestBroadcasterChannelNameMatchesPattern(t *testing.T) {
	t.Parallel()

	b := broadcasting.NewBaseBroadcaster()

	cases := []struct {
		channel string
		pattern string
		match   bool
	}{
		{"something", "something", true},
		{"something.23", "something.{id}", true},
		{"something.23.test", "something.{id}.test", true},
		{"something.23.test.42", "something.{id}.test.{id2}", true},
		{"something-23:test-42", "something-{id}:test-{id2}", true},
		{"something..test.42", "something.{id}.test.{id2}", false},
		{"23:string:test", "{id}:string:{text}", true},
		{"something.23", "something", false},
		{"something.23.test.42", "something.test.{id}", false},
		{"something-23-test-42", "something-{id}-test", false},
		{"23:test", "{id}:test:abcd", false},
		{"customer.order.1", "order.{id}", false},
		{"customerorder.1", "order.{id}", false},
	}

	for _, tt := range cases {
		if got := b.ChannelNameMatchesPattern(tt.channel, tt.pattern); got != tt.match {
			t.Fatalf("ChannelNameMatchesPattern(%q, %q) = %v, want %v", tt.channel, tt.pattern, got, tt.match)
		}
	}
}

func (dummyChannel) Join(contractsauth.Authenticatable, ...any) (any, error) {
	return []any{"joined"}, nil
}

func TestChannelHandlerNormalizationAcceptsFunctions(t *testing.T) {
	t.Parallel()

	b := broadcasting.NewBaseBroadcaster()
	b.Channel("test", func(contractsauth.Authenticatable, ...any) bool { return true })

	result, err := b.VerifyUserCanAccessChannel(requestWithUser("test"), "test")

	if err != nil {
		t.Fatalf("VerifyUserCanAccessChannel returned error: %v", err)
	}

	if !reflect.DeepEqual(result, true) {
		t.Fatalf("result = %#v, want true", result)
	}
}
