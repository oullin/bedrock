package notifications

import cn "github.com/bedrock/packages/contracts/notifications"

// FormatNotifiables normalises various notifiable inputs into a consistent
// slice of Notifiable.
func FormatNotifiables(notifiables any) []cn.Notifiable {
	switch v := notifiables.(type) {
	case cn.Notifiable:
		return []cn.Notifiable{v}
	case []cn.Notifiable:
		return v
	case []any:
		result := make([]cn.Notifiable, 0, len(v))

		for _, item := range v {
			if n, ok := item.(cn.Notifiable); ok {
				result = append(result, n)
			}
		}

		return result
	default:
		return nil
	}
}
