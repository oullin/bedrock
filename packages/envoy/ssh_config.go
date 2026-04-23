package envoy

import (
	"strings"
)

// SSHConfigFile represents the Host sections from an OpenSSH config file.
type SSHConfigFile struct {
	groups []map[string]string
}

// ParseSSHConfigString parses Host sections from OpenSSH config text.
func ParseSSHConfigString(config string) SSHConfigFile {
	var groups []map[string]string

	var current map[string]string

	for _, raw := range strings.Split(config, "\n") {
		line := strings.TrimSpace(raw)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := splitSSHConfigLine(line)

		if !ok {
			continue
		}

		if key == "match" {
			current = nil

			continue
		}

		if key == "host" {
			current = map[string]string{"host": value}
			groups = append(groups, current)

			continue
		}

		if current == nil {
			continue
		}

		current[key] = value
	}

	return SSHConfigFile{groups: groups}
}

// Groups returns the parsed Host sections.
func (c SSHConfigFile) Groups() []map[string]string {
	groups := make([]map[string]string, len(c.groups))

	for i, group := range c.groups {
		copyGroup := make(map[string]string, len(group))

		for key, value := range group {
			copyGroup[key] = value
		}

		groups[i] = copyGroup
	}

	return groups
}

// FindConfiguredHost returns the Host alias matching target or target's HostName.
func (c SSHConfigFile) FindConfiguredHost(target string) (string, bool) {
	user, host := splitUserHost(target)

	for _, group := range c.groups {
		if user != "" && group["user"] != "" && group["user"] != user {
			continue
		}

		for _, alias := range strings.Fields(group["host"]) {
			if alias == host || group["hostname"] == host {
				return alias, true
			}
		}
	}

	return "", false
}

func splitSSHConfigLine(line string) (string, string, bool) {
	var key, value string

	if idx := strings.Index(line, "="); idx >= 0 {
		key = strings.TrimSpace(line[:idx])
		value = strings.TrimSpace(line[idx+1:])
	} else {
		fields := strings.Fields(line)

		if len(fields) < 2 {
			return "", "", false
		}

		key = fields[0]
		value = strings.TrimSpace(line[len(fields[0]):])
	}

	key = strings.ToLower(strings.TrimSpace(key))
	value = unquoteSSHArgument(strings.TrimSpace(value))

	if key == "" {
		return "", "", false
	}

	return key, value, true
}

func unquoteSSHArgument(value string) string {
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		return value[1 : len(value)-1]
	}

	return value
}

func splitUserHost(target string) (string, string) {
	if idx := strings.Index(target, "@"); idx >= 0 {
		return target[:idx], target[idx+1:]
	}

	return "", target
}
