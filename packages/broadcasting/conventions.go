package broadcasting

import "strings"

// IsGuardedChannel reports whether the Pusher-style channel requires auth.
func IsGuardedChannel(channel string) bool {
	return strings.HasPrefix(channel, "private-") ||
		strings.HasPrefix(channel, "private-encrypted-") ||
		strings.HasPrefix(channel, "presence-")
}

// NormalizeChannelName removes Pusher-style auth prefixes from a channel.
func NormalizeChannelName(channel string) string {
	switch {
	case strings.HasPrefix(channel, "private-encrypted-"):
		return strings.TrimPrefix(channel, "private-encrypted-")
	case strings.HasPrefix(channel, "private-"):
		return strings.TrimPrefix(channel, "private-")
	case strings.HasPrefix(channel, "presence-"):
		return strings.TrimPrefix(channel, "presence-")
	default:
		return channel
	}
}
