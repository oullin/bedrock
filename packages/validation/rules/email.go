package rules

import (
	"net"
	"regexp"
	"strings"
)

func init_email() {
	Register("Email", validateEmail)
	Register("Url", validateUrl)
	Register("ActiveUrl", validateActiveUrl)
	Register("Ip", validateIp)
	Register("Ipv4", validateIpv4)
	Register("Ipv6", validateIpv6)
	Register("MacAddress", validateMacAddress)
	Register("Uuid", validateUuid)
	Register("Ulid", validateUlid)
}

// emailRe is the basic RFC 5322-compatible email regex used by the native
// validator (without DNS/spoof checks).
var emailRe = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

func validateEmail(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	s = strings.TrimSpace(s)

	if s == "" {
		return false
	}

	return emailRe.MatchString(s)
}

var urlRe = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)

func validateUrl(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	return urlRe.MatchString(strings.TrimSpace(s))
}

// validateActiveUrl checks DNS resolution (A or AAAA records).
func validateActiveUrl(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	// Extract host from URL
	s = strings.TrimSpace(s)

	for _, prefix := range []string{"https://", "http://", "ftp://"} {
		s = strings.TrimPrefix(s, prefix)
	}

	// Strip path
	if idx := strings.IndexByte(s, '/'); idx != -1 {
		s = s[:idx]
	}

	// Strip port
	host, _, err := net.SplitHostPort(s)

	if err != nil {
		host = s
	}

	if host == "" {
		return false
	}

	addrs, err := net.LookupHost(host)

	return err == nil && len(addrs) > 0
}

func validateIp(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	return net.ParseIP(strings.TrimSpace(s)) != nil
}

func validateIpv4(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	ip := net.ParseIP(strings.TrimSpace(s))

	return ip != nil && ip.To4() != nil
}

func validateIpv6(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	ip := net.ParseIP(strings.TrimSpace(s))

	return ip != nil && ip.To4() == nil
}

var macRe = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:\-]){5}[0-9A-Fa-f]{2}$`)

func validateMacAddress(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	return macRe.MatchString(strings.TrimSpace(s))
}

var uuidRe = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func validateUuid(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	return uuidRe.MatchString(strings.TrimSpace(s))
}

// ULID: 26 chars, Crockford base32 alphabet, first char <= 7 (time bits).
var ulidRe = regexp.MustCompile(`(?i)^[0-7][0-9A-HJKMNP-TV-Z]{25}$`)

func validateUlid(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	return ulidRe.MatchString(strings.TrimSpace(s))
}
