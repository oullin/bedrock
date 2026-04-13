package notifications_test

import (
	"context"
	"sync"

	"github.com/bedrock/packages/bus"
	cevents "github.com/bedrock/packages/contracts/events"
	"github.com/bedrock/packages/contracts/mail"
	cn "github.com/bedrock/packages/contracts/notifications"

	"github.com/bedrock/packages/notifications"
)

// --- Mock Notifiable ---

type mockNotifiable struct {
	key    string
	routes map[string]any
	locale string
}

// --- Mock Channel ---

type channelSendCall struct {
	Notifiable   cn.Notifiable
	Notification any
}

type mockChannel struct {
	mu    sync.Mutex
	calls []channelSendCall
	err   error
}

// --- Mock Event Dispatcher ---

type dispatchedEvent struct {
	Event any
}

type mockEventDispatcher struct {
	mu              sync.Mutex
	dispatched      []dispatchedEvent
	listeners       map[string][]cevents.Listener
	untilResponse   any
	untilErr        error
	cancelOnSending bool
}

// --- Mock Bus Dispatcher ---

type busDispatchCall struct {
	Command any
}

type mockBusDispatcher struct {
	mu    sync.Mutex
	calls []busDispatchCall
	err   error
}

// --- Mock Mailer Factory ---

type mailSendCall struct {
	Mailable mail.Mailable
}

type mockMailer struct {
	mu    sync.Mutex
	calls []mailSendCall
	err   error
}

type mockMailerFactory struct {
	mailer *mockMailer
}

// --- Mock Database Notification Store ---

type storeCreateCall struct {
	Notification *notifications.DatabaseNotification
}

type mockDatabaseNotificationStore struct {
	mu          sync.Mutex
	creates     []storeCreateCall
	readMarks   []string
	unreadMarks []string
	records     map[string]*notifications.DatabaseNotification
	err         error
}

// --- Test Notification Types ---

type testNotification struct {
	notifications.Notification
	viaChannels []string
	shouldSend  *bool
}

type testMailNotification struct {
	testNotification
	message *notifications.MailMessage
}

type testDatabaseNotification struct {
	testNotification
	data map[string]any
}

type testBroadcastNotification struct {
	testNotification
	data map[string]any
}

type testQueuedNotification struct {
	testNotification
}

type testArrayNotification struct {
	testNotification
	data map[string]any
}

func newMockNotifiable(key string) *mockNotifiable {
	return &mockNotifiable{key: key, routes: make(map[string]any)}
}

func (n *mockNotifiable) RouteNotificationFor(_ context.Context, channel string) any {
	return n.routes[channel]
}

func (n *mockNotifiable) GetKey() string { return n.key }

func (n *mockNotifiable) PreferredLocale() string { return n.locale }

func (c *mockChannel) Send(_ context.Context, notifiable cn.Notifiable, notification any) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.calls = append(c.calls, channelSendCall{Notifiable: notifiable, Notification: notification})

	return c.err
}

func (c *mockChannel) CallCount() int {
	c.mu.Lock()

	defer c.mu.Unlock()

	return len(c.calls)
}

func newMockEventDispatcher() *mockEventDispatcher {
	return &mockEventDispatcher{
		listeners: make(map[string][]cevents.Listener),
	}
}

func (d *mockEventDispatcher) Listen(events any, listeners ...cevents.Listener) {
	_ = events
	_ = listeners
}

func (d *mockEventDispatcher) HasListeners(event any) bool {
	_ = event

	return false
}

func (d *mockEventDispatcher) HasWildcardListeners(event any) bool {
	_ = event

	return false
}

func (d *mockEventDispatcher) Subscribe(subscriber cevents.Subscriber) {
	_ = subscriber
}

func (d *mockEventDispatcher) Until(_ context.Context, event any) (any, error) {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.dispatched = append(d.dispatched, dispatchedEvent{Event: event})

	if d.cancelOnSending {
		if _, ok := event.(notifications.NotificationSending); ok {
			return false, nil
		}
	}

	return d.untilResponse, d.untilErr
}

func (d *mockEventDispatcher) Dispatch(_ context.Context, event any) ([]any, error) {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.dispatched = append(d.dispatched, dispatchedEvent{Event: event})

	return nil, nil
}

func (d *mockEventDispatcher) Push(_ context.Context, event any) {
	_ = event
}

func (d *mockEventDispatcher) Flush(_ context.Context, event string) error {
	_ = event

	return nil
}

func (d *mockEventDispatcher) Forget(event any) {
	_ = event
}

func (d *mockEventDispatcher) ForgetPushed() {}

func (d *mockEventDispatcher) GetListeners(event any) []cevents.Listener {
	_ = event

	return nil
}

func (d *mockEventDispatcher) DispatchedEvents() []dispatchedEvent {
	d.mu.Lock()

	defer d.mu.Unlock()

	result := make([]dispatchedEvent, len(d.dispatched))
	copy(result, d.dispatched)

	return result
}

func (d *mockBusDispatcher) Dispatch(_ context.Context, command any) (any, error) {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.calls = append(d.calls, busDispatchCall{Command: command})

	return nil, d.err
}

func (d *mockBusDispatcher) DispatchSync(_ context.Context, command any) (any, error) {
	return d.Dispatch(context.Background(), command)
}

func (d *mockBusDispatcher) DispatchNow(_ context.Context, command any) (any, error) {
	return d.Dispatch(context.Background(), command)
}

func (d *mockBusDispatcher) DispatchAfterResponse(_ context.Context, _ any) error { return nil }

func (d *mockBusDispatcher) PipeThrough(_ ...bus.Pipe) bus.Dispatcher { return d }

func (d *mockBusDispatcher) Map(_ any, _ bus.Handler) bus.Dispatcher { return d }

func (d *mockBusDispatcher) HasCommandHandler(_ any) bool { return false }

func (d *mockBusDispatcher) GetCommandHandler(_ any) (bus.Handler, bool) { return nil, false }

func (d *mockBusDispatcher) Chain(_ []any) *bus.PendingChain { return nil }

func (d *mockBusDispatcher) CallCount() int {
	d.mu.Lock()

	defer d.mu.Unlock()

	return len(d.calls)
}

func (m *mockMailer) Raw(_ context.Context, _ string, _ ...func(*mail.Message)) (*mail.SentMessage, error) {
	return &mail.SentMessage{}, nil
}

func (m *mockMailer) Send(_ context.Context, mailable mail.Mailable) (*mail.SentMessage, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.calls = append(m.calls, mailSendCall{Mailable: mailable})

	return &mail.SentMessage{}, m.err
}

func (m *mockMailer) SendNow(ctx context.Context, mailable mail.Mailable) (*mail.SentMessage, error) {
	return m.Send(ctx, mailable)
}

func (m *mockMailer) CallCount() int {
	m.mu.Lock()

	defer m.mu.Unlock()

	return len(m.calls)
}

func (f *mockMailerFactory) Mailer(_ ...string) (mail.Mailer, error) {
	return f.mailer, nil
}

func newMockDatabaseNotificationStore() *mockDatabaseNotificationStore {
	return &mockDatabaseNotificationStore{
		records: make(map[string]*notifications.DatabaseNotification),
	}
}

func (s *mockDatabaseNotificationStore) Create(_ context.Context, n *notifications.DatabaseNotification) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.creates = append(s.creates, storeCreateCall{Notification: n})
	s.records[n.ID] = n

	return s.err
}

func (s *mockDatabaseNotificationStore) Find(_ context.Context, id string) (*notifications.DatabaseNotification, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	return s.records[id], s.err
}

func (s *mockDatabaseNotificationStore) ForNotifiable(_ context.Context, _, _ string) ([]*notifications.DatabaseNotification, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	var result []*notifications.DatabaseNotification

	for _, r := range s.records {
		result = append(result, r)
	}

	return result, s.err
}

func (s *mockDatabaseNotificationStore) ReadForNotifiable(_ context.Context, _, _ string) ([]*notifications.DatabaseNotification, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	var result []*notifications.DatabaseNotification

	for _, r := range s.records {
		if r.Read() {
			result = append(result, r)
		}
	}

	return result, s.err
}

func (s *mockDatabaseNotificationStore) UnreadForNotifiable(_ context.Context, _, _ string) ([]*notifications.DatabaseNotification, error) {
	s.mu.Lock()

	defer s.mu.Unlock()

	var result []*notifications.DatabaseNotification

	for _, r := range s.records {
		if r.Unread() {
			result = append(result, r)
		}
	}

	return result, s.err
}

func (s *mockDatabaseNotificationStore) MarkAsRead(_ context.Context, id string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.readMarks = append(s.readMarks, id)

	if r, ok := s.records[id]; ok {
		r.MarkAsRead()
	}

	return s.err
}

func (s *mockDatabaseNotificationStore) MarkAsUnread(_ context.Context, id string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.unreadMarks = append(s.unreadMarks, id)

	if r, ok := s.records[id]; ok {
		r.MarkAsUnread()
	}

	return s.err
}

func (s *mockDatabaseNotificationStore) Delete(_ context.Context, _ string) error {
	return s.err
}

func (s *mockDatabaseNotificationStore) CreateCount() int {
	s.mu.Lock()

	defer s.mu.Unlock()

	return len(s.creates)
}

func newTestNotification(channels ...string) *testNotification {
	return &testNotification{
		Notification: notifications.NewNotification(),
		viaChannels:  channels,
	}
}

func (n *testNotification) Via(_ context.Context, _ cn.Notifiable) []string {
	return n.viaChannels
}

func (n *testNotification) ShouldSend(_ context.Context, _ cn.Notifiable, _ string) bool {
	if n.shouldSend != nil {
		return *n.shouldSend
	}

	return true
}

func newTestMailNotification() *testMailNotification {
	msg := notifications.NewMailMessage().
		Subject("Test Subject").
		Line("Hello from test")

	return &testMailNotification{
		testNotification: *newTestNotification("mail"),
		message:          msg,
	}
}

func (n *testMailNotification) ToMail(_ context.Context, _ cn.Notifiable) (*notifications.MailMessage, error) {
	return n.message, nil
}

func newTestDatabaseNotification(data map[string]any) *testDatabaseNotification {
	return &testDatabaseNotification{
		testNotification: *newTestNotification("database"),
		data:             data,
	}
}

func (n *testDatabaseNotification) ToDatabase(_ context.Context, _ cn.Notifiable) (map[string]any, error) {
	return n.data, nil
}

func newTestBroadcastNotification(data map[string]any) *testBroadcastNotification {
	return &testBroadcastNotification{
		testNotification: *newTestNotification("broadcast"),
		data:             data,
	}
}

func (n *testBroadcastNotification) ToBroadcast(_ context.Context, _ cn.Notifiable) (*notifications.BroadcastMessage, error) {
	return notifications.NewBroadcastMessage(n.data), nil
}

func newTestQueuedNotification(channels ...string) *testQueuedNotification {
	return &testQueuedNotification{
		testNotification: *newTestNotification(channels...),
	}
}

func (n *testQueuedNotification) ShouldQueue() {}

func newTestArrayNotification(channels []string, data map[string]any) *testArrayNotification {
	return &testArrayNotification{
		testNotification: *newTestNotification(channels...),
		data:             data,
	}
}

func (n *testArrayNotification) ToArray(_ context.Context, _ cn.Notifiable) (map[string]any, error) {
	return n.data, nil
}
