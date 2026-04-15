package container_test

import (
	"testing"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/contracts/provider"
)

// fakeProvider records lifecycle calls for assertions.
type fakeProvider struct {
	name          string
	registerCalls int
	bootCalls     int
	provides      []string
}

// nonBootable satisfies only ServiceProvider, not Bootable.
type nonBootable struct{ registerCalls int }

// orderRecorder records the sequence in which providers were called.
type orderRecorder struct {
	tag string
	log *[]string
}

// mutating the returned slice must not affect internal state

// Provider B's binding depends on A being available at resolve-time, not
// at register-time. Out-of-order registration must still work.

// Register B FIRST (depends on A).

// Register A SECOND.

// closureProvider is a tiny adapter that lets a test write inline Register logic.
type closureProvider struct{ register func() }

// ---------- Phase 4: Deferred providers ----------

// deferredProvider claims abstract keys but only registers them lazily.
type deferredProvider struct {
	keys       []string
	register   func(*container.Container)
	bootCalls  int
	registered bool
}

// pass nil container; tests that need the container use a closure

// First resolve must trigger the deferred Register().

// app booted, but provider not yet flushed

// triggers flush + Boot

// ---------- Phase 4: DependsOn topological sort ----------

type depProvider struct {
	name     string
	provides []string
	depends  []string
	log      *[]string
}

func (p *fakeProvider) Register()          { p.registerCalls++ }
func (p *fakeProvider) Boot()              { p.bootCalls++ }
func (p *fakeProvider) Provides() []string { return p.provides }

func (p *nonBootable) Register() { p.registerCalls++ }

func (p *orderRecorder) Register() { *p.log = append(*p.log, "register:"+p.tag) }
func (p *orderRecorder) Boot()     { *p.log = append(*p.log, "boot:"+p.tag) }

func TestNewApplication_HasUsableContainer(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()

	if app.Container == nil {
		t.Fatal("expected embedded *Container to be non-nil")
	}

	app.Instance("answer", 42)

	v, err := app.Make("answer")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v != 42 {
		t.Fatalf("expected 42, got %v", v)
	}
}

func TestApplication_Register_CallsRegisterOnce(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()
	p := &fakeProvider{name: "p1"}

	app.Register(p)

	if p.registerCalls != 1 {
		t.Fatalf("expected Register to be called once, got %d", p.registerCalls)
	}

	if got := len(app.Providers()); got != 1 {
		t.Fatalf("expected 1 provider stored, got %d", got)
	}
}

func TestApplication_Boot_OnlyCallsBootable(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()
	bootable := &fakeProvider{name: "bootable"}
	plain := &nonBootable{}

	app.Register(bootable)
	app.Register(plain)
	app.Boot()

	if bootable.bootCalls != 1 {
		t.Fatalf("expected Bootable provider to receive one Boot call, got %d", bootable.bootCalls)
	}

	if plain.registerCalls != 1 {
		t.Fatalf("expected non-bootable to still be registered, got %d", plain.registerCalls)
	}

	if !app.Booted() {
		t.Fatal("expected Booted() to be true after Boot()")
	}
}

func TestApplication_Boot_IsIdempotent(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()
	p := &fakeProvider{name: "p"}
	app.Register(p)

	app.Boot()
	app.Boot()
	app.Boot()

	if p.bootCalls != 1 {
		t.Fatalf("expected Boot to be called once across multiple app.Boot() calls, got %d", p.bootCalls)
	}
}

func TestApplication_Boot_PreservesRegistrationOrder(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()
	log := []string{}

	app.Register(&orderRecorder{tag: "a", log: &log})
	app.Register(&orderRecorder{tag: "b", log: &log})
	app.Register(&orderRecorder{tag: "c", log: &log})
	app.Boot()

	want := []string{
		"register:a", "register:b", "register:c",
		"boot:a", "boot:b", "boot:c",
	}

	if len(log) != len(want) {
		t.Fatalf("expected %d events, got %d: %v", len(want), len(log), log)
	}

	for i, ev := range want {
		if log[i] != ev {
			t.Fatalf("at %d: expected %q, got %q (full: %v)", i, ev, log[i], log)
		}
	}
}

func TestApplication_Providers_ReturnsCopy(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()
	p := &fakeProvider{}
	app.Register(p)

	got := app.Providers()

	if len(got) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(got))
	}

	got[0] = nil

	if app.Providers()[0] == nil {
		t.Fatal("Providers() returned slice aliases internal state — must return a copy")
	}
}

func TestApplication_LazyResolution_HandlesOutOfOrderRegistration(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()

	app.Register(&closureProvider{
		register: func() {
			app.Singleton("B", func(c *container.Container) (any, error) {
				rawA, err := c.Make("A")

				if err != nil {
					return nil, err
				}

				return "B(" + rawA.(string) + ")", nil
			})
		},
	})

	app.Register(&closureProvider{
		register: func() {
			app.Singleton("A", func(_ *container.Container) (any, error) {
				return "A", nil
			})
		},
	})

	v, err := app.Make("B")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v != "B(A)" {
		t.Fatalf("expected B(A), got %v", v)
	}
}

func (p *closureProvider) Register() { p.register() }

func TestApplication_RegisterMany_RegistersAllInOrder(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()
	log := []string{}

	app.RegisterMany([]provider.ServiceProvider{
		&orderRecorder{tag: "a", log: &log},
		&orderRecorder{tag: "b", log: &log},
		&orderRecorder{tag: "c", log: &log},
	})

	if len(app.Providers()) != 3 {
		t.Fatalf("expected 3 providers, got %d", len(app.Providers()))
	}

	want := []string{"register:a", "register:b", "register:c"}

	for i, ev := range want {
		if log[i] != ev {
			t.Fatalf("at %d: expected %q, got %q", i, ev, log[i])
		}
	}
}

func (p *deferredProvider) Register() {
	p.registered = true
	p.register(nil)
}

func (p *deferredProvider) Provides() []string { return p.keys }
func (p *deferredProvider) Deferred() bool     { return true }
func (p *deferredProvider) Boot()              { p.bootCalls++ }

func TestApplication_DeferredProvider_NotRegisteredUntilResolved(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()

	dp := &deferredProvider{
		keys: []string{"deferred.svc"},
		register: func(_ *container.Container) {
			app.Instance("deferred.svc", "value")
		},
	}

	app.Register(dp)

	if dp.registered {
		t.Fatal("expected deferred provider NOT to be registered immediately")
	}

	v, err := app.Make("deferred.svc")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v != "value" {
		t.Fatalf("expected %q, got %v", "value", v)
	}

	if !dp.registered {
		t.Fatal("expected deferred provider to be registered after first Make")
	}
}

func TestApplication_DeferredProvider_BootRunsAfterFlush(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()

	dp := &deferredProvider{
		keys: []string{"deferred.boot"},
		register: func(_ *container.Container) {
			app.Instance("deferred.boot", 1)
		},
	}

	app.Register(dp)
	app.Boot()

	if dp.bootCalls != 0 {
		t.Fatalf("expected Boot to be deferred, got %d calls", dp.bootCalls)
	}

	_, _ = app.Make("deferred.boot")

	if dp.bootCalls != 1 {
		t.Fatalf("expected one Boot call after flush, got %d", dp.bootCalls)
	}
}

func (p *depProvider) Register()           { *p.log = append(*p.log, p.name) }
func (p *depProvider) Provides() []string  { return p.provides }
func (p *depProvider) DependsOn() []string { return p.depends }

func TestApplication_RegisterMany_TopoSortsByDependsOn(t *testing.T) {
	t.Parallel()

	log := []string{}

	// C depends on B which depends on A. Register them in REVERSE order
	// to prove the sort actually runs.
	pA := &depProvider{name: "A", provides: []string{"a"}, log: &log}
	pB := &depProvider{name: "B", provides: []string{"b"}, depends: []string{"a"}, log: &log}
	pC := &depProvider{name: "C", provides: []string{"c"}, depends: []string{"b"}, log: &log}

	app := container.NewApplication()
	app.RegisterMany([]provider.ServiceProvider{pC, pB, pA})

	want := []string{"A", "B", "C"}

	if len(log) != len(want) {
		t.Fatalf("expected %v, got %v", want, log)
	}

	for i, ev := range want {
		if log[i] != ev {
			t.Fatalf("at %d: expected %q, got %q (full: %v)", i, ev, log[i], log)
		}
	}
}

func TestApplication_RegisterMany_PanicsOnCycle(t *testing.T) {
	t.Parallel()

	log := []string{}

	pA := &depProvider{name: "A", provides: []string{"a"}, depends: []string{"b"}, log: &log}
	pB := &depProvider{name: "B", provides: []string{"b"}, depends: []string{"a"}, log: &log}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on cycle")
		}
	}()

	container.NewApplication().RegisterMany([]provider.ServiceProvider{pA, pB})
}

// ---------- Phase 4: Introspection ----------

func TestApplication_HasProviderAndProviderFor(t *testing.T) {
	t.Parallel()

	app := container.NewApplication()
	p := &fakeProvider{provides: []string{"x", "y"}}
	app.Register(p)

	if !app.HasProvider("x") {
		t.Fatal("expected HasProvider(x) to be true")
	}

	if !app.HasProvider("y") {
		t.Fatal("expected HasProvider(y) to be true")
	}

	if app.HasProvider("z") {
		t.Fatal("expected HasProvider(z) to be false")
	}

	if got := app.ProviderFor("x"); got != p {
		t.Fatalf("expected ProviderFor(x) to return original provider")
	}

	if got := app.ProviderFor("z"); got != nil {
		t.Fatalf("expected ProviderFor(z) to be nil, got %v", got)
	}
}

// Sanity check: the provider package is imported and used.
var _ provider.ServiceProvider = (*fakeProvider)(nil)
var _ provider.Bootable = (*fakeProvider)(nil)
var _ provider.Deferred = (*deferredProvider)(nil)
var _ provider.DependsOn = (*depProvider)(nil)
