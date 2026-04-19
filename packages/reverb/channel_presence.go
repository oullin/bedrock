package reverb

import (
	"context"
	"encoding/json"
	"sync"

	contractsReverb "github.com/bedrock/packages/contracts/reverb"
)

// presenceMember holds the identity data for a single presence subscriber.
type presenceMember struct {
	UserID   string
	UserInfo any
}

// PresenceChannel extends the base channel with member presence tracking.
// Each member is identified by user_id; multiple connections from the same
// user_id count as a single presence member.
type PresenceChannel struct {
	channel // base channel (not PrivateChannel, to avoid double mutex)
	app     *App
	mu      sync.RWMutex
	members map[string]presenceMember // socketID → presenceMember
}

var _ contractsReverb.Channel = (*PresenceChannel)(nil)
var _ contractsReverb.PresenceChanneler = (*PresenceChannel)(nil)

// NewPresenceChannel constructs a PresenceChannel.
func NewPresenceChannel(name string, app *App) *PresenceChannel {
	return &PresenceChannel{
		channel: channel{
			name:  name,
			app:   app,
			conns: make(map[string]contractsReverb.Connection),
		},
		app:     app,
		members: make(map[string]presenceMember),
	}
}

// Subscribe verifies auth, registers the member, and notifies existing
// subscribers if this is a new user_id.
func (ch *PresenceChannel) Subscribe(ctx context.Context, conn contractsReverb.Connection, auth, channelData string) error {
	if !VerifyChannelAuth(ch.app.Secret(), conn.SocketID(), ch.name, channelData, auth) {
		return ErrUnauthorized
	}

	var cd struct {
		UserID   string `json:"user_id"`
		UserInfo any    `json:"user_info"`
	}

	_ = json.Unmarshal([]byte(channelData), &cd)

	ch.mu.Lock()

	// Check whether this user_id is already present before this socket joins.
	isNewUser := true

	for _, m := range ch.members {
		if m.UserID == cd.UserID {
			isNewUser = false

			break
		}
	}

	ch.members[conn.SocketID()] = presenceMember{UserID: cd.UserID, UserInfo: cd.UserInfo}

	// Also add to base channel conns.
	ch.channel.mu.Lock()
	ch.channel.conns[conn.SocketID()] = conn
	ch.channel.mu.Unlock()

	ch.mu.Unlock()

	// Notify existing subscribers when a new user_id joins.
	if isNewUser {
		data := map[string]any{"user_id": cd.UserID, "user_info": cd.UserInfo}
		dataBytes, _ := json.Marshal(data)
		socketID := conn.SocketID()
		addedEvent := contractsReverb.Event{
			Event:   "pusher_internal:member_added",
			Data:    string(dataBytes),
			Channel: ch.name,
		}

		_ = ch.channel.broadcast(ctx, addedEvent, &socketID)
	}

	// Build subscription_succeeded with presence data.
	presenceData := ch.subscriptionData()
	presencePayload, err := json.Marshal(map[string]any{"presence": presenceData})

	if err != nil {
		return err
	}

	msg, err := json.Marshal(map[string]string{
		"event":   "pusher_internal:subscription_succeeded",
		"data":    string(presencePayload),
		"channel": ch.name,
	})

	if err != nil {
		return err
	}

	return conn.Send(ctx, msg)
}

// Unsubscribe removes the connection from the channel and broadcasts
// member_removed if no other connection from that user_id remains.
func (ch *PresenceChannel) Unsubscribe(ctx context.Context, conn contractsReverb.Connection) {
	ch.mu.Lock()

	member, exists := ch.members[conn.SocketID()]
	delete(ch.members, conn.SocketID())

	// Remove from base channel conns.
	ch.channel.mu.Lock()
	delete(ch.channel.conns, conn.SocketID())
	ch.channel.mu.Unlock()

	// Count remaining connections for this user_id.
	remaining := 0

	if exists {
		for _, m := range ch.members {
			if m.UserID == member.UserID {
				remaining++
			}
		}
	}

	ch.mu.Unlock()

	if exists && remaining == 0 {
		data := map[string]any{"user_id": member.UserID}
		dataBytes, _ := json.Marshal(data)
		removedEvent := contractsReverb.Event{
			Event:   "pusher_internal:member_removed",
			Data:    string(dataBytes),
			Channel: ch.name,
		}

		_ = ch.channel.broadcastAll(ctx, removedEvent)
	}
}

// Broadcast sends the event to all subscribers except the sender.
func (ch *PresenceChannel) Broadcast(ctx context.Context, event contractsReverb.Event, except *string) error {
	return ch.channel.broadcast(ctx, event, except)
}

// BroadcastToAll sends the event to every subscriber.
func (ch *PresenceChannel) BroadcastToAll(ctx context.Context, event contractsReverb.Event) error {
	return ch.channel.broadcastAll(ctx, event)
}

// Members returns a deduplicated map of user_id → user_info for all subscribers.
func (ch *PresenceChannel) Members() map[string]any {
	ch.mu.RLock()

	defer ch.mu.RUnlock()

	result := make(map[string]any)

	for _, m := range ch.members {
		result[m.UserID] = m.UserInfo
	}

	return result
}

// MemberCount returns the number of unique users subscribed.
func (ch *PresenceChannel) MemberCount() int {
	return len(ch.Members())
}

// MemberIDs returns the unique user IDs of all subscribers.
func (ch *PresenceChannel) MemberIDs() []string {
	members := ch.Members()
	ids := make([]string, 0, len(members))

	for id := range members {
		ids = append(ids, id)
	}

	return ids
}

// subscriptionData builds the PresenceMemberData payload for subscription_succeeded.
func (ch *PresenceChannel) subscriptionData() PresenceMemberData {
	members := ch.Members()
	ids := ch.MemberIDs()

	return PresenceMemberData{
		Count: len(ids),
		IDs:   ids,
		Hash:  members,
	}
}
