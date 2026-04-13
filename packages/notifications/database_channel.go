package notifications

import (
	"context"
	"fmt"
	"reflect"
	"time"

	cn "github.com/bedrock/packages/contracts/notifications"
)

// DatabaseChannel delivers notifications by persisting them to a database via
// the DatabaseNotificationStore.
type DatabaseChannel struct {
	store DatabaseNotificationStore
}

// NewDatabaseChannel creates a DatabaseChannel backed by the given store.
func NewDatabaseChannel(store DatabaseNotificationStore) *DatabaseChannel {
	return &DatabaseChannel{store: store}
}

// compile-time interface check.
var _ cn.Channel = (*DatabaseChannel)(nil)

// Send persists the notification to the database for the notifiable entity.
func (c *DatabaseChannel) Send(ctx context.Context, notifiable cn.Notifiable, notification any) error {
	record, err := c.buildPayload(ctx, notifiable, notification)

	if err != nil {
		return err
	}

	return c.store.Create(ctx, record)
}

func (c *DatabaseChannel) buildPayload(ctx context.Context, notifiable cn.Notifiable, notification any) (*DatabaseNotification, error) {
	data, err := c.getData(ctx, notifiable, notification)

	if err != nil {
		return nil, err
	}

	id := notificationID(notification)
	nType := c.notificationType(ctx, notifiable, notification)
	now := time.Now()

	return &DatabaseNotification{
		ID:             id,
		Type:           nType,
		NotifiableType: reflect.TypeOf(notifiable).String(),
		NotifiableID:   notifiable.GetKey(),
		Data:           data,
		ReadAt:         nil,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (c *DatabaseChannel) getData(ctx context.Context, notifiable cn.Notifiable, notification any) (map[string]any, error) {
	if db, ok := notification.(DatabaseNotificationData); ok {
		return db.ToDatabase(ctx, notifiable)
	}

	if arr, ok := notification.(ArrayNotification); ok {
		return arr.ToArray(ctx, notifiable)
	}

	return nil, fmt.Errorf("%w: %T", ErrMissingDatabaseData, notification)
}

func (c *DatabaseChannel) notificationType(ctx context.Context, notifiable cn.Notifiable, notification any) string {
	if dt, ok := notification.(HasCustomDatabaseType); ok {
		return dt.DatabaseType(ctx, notifiable)
	}

	t := reflect.TypeOf(notification)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.String()
}

func notificationID(notification any) string {
	if n, ok := notification.(interface{ GetID() string }); ok {
		if id := n.GetID(); id != "" {
			return id
		}
	}

	return generateID()
}
