package drivers

import (
	"fmt"
	"strconv"
)

func toString(id uint64) string {
	return fmt.Sprintf("%d", id)
}

func parseStatInt(stats map[string]string, key string) int64 {
	v, ok := stats[key]

	if !ok {
		return 0
	}

	n, err := strconv.ParseInt(v, 10, 64)

	if err != nil {
		return 0
	}

	return n
}
