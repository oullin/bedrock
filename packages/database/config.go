package database

import (
	"net/url"
	"strings"
)

// ConnectionConfig holds the configuration for a single database connection.
type ConnectionConfig struct {
	Driver   string         `json:"driver"`
	Host     string         `json:"host"`
	Port     int            `json:"port"`
	Database string         `json:"database"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Charset  string         `json:"charset"`
	Prefix   string         `json:"prefix"`
	Schema   string         `json:"schema"`
	SSLMode  string         `json:"sslmode"`
	Options  map[string]any `json:"options"`
	URL      string         `json:"url"`
}

// ParseDatabaseURL parses a database connection URL into a ConnectionConfig.
// It supports formats like:
//
//	mysql://user:pass@host:3306/dbname
//	postgres://user:pass@host:5432/dbname?sslmode=disable
//	sqlite:///path/to/database.db
func ParseDatabaseURL(rawURL string) (*ConnectionConfig, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	cfg := &ConnectionConfig{
		Driver:  u.Scheme,
		Host:    u.Hostname(),
		Options: make(map[string]any),
	}

	if p := u.Port(); p != "" {
		port := 0
		for _, c := range p {
			port = port*10 + int(c-'0')
		}
		cfg.Port = port
	}

	if u.User != nil {
		cfg.Username = u.User.Username()
		cfg.Password, _ = u.User.Password()
	}

	cfg.Database = strings.TrimPrefix(u.Path, "/")

	for k, v := range u.Query() {
		if len(v) > 0 {
			cfg.Options[k] = v[0]
		}
	}

	if sslmode, ok := cfg.Options["sslmode"]; ok {
		cfg.SSLMode = sslmode.(string)
	}

	return cfg, nil
}
