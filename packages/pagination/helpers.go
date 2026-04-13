package pagination

import (
	"net/url"
	"strconv"
	"strings"
)

func intToString(i int) string {
	return strconv.Itoa(i)
}

// encodeQuery encodes url.Values using %20 for spaces instead of +,
// matching Laravel's behavior.
func encodeQuery(v url.Values) string {
	encoded := v.Encode()

	return strings.ReplaceAll(encoded, "+", "%20")
}
