package authflows

import (
	"fmt"
	"strings"

	configpkg "github.com/gollin/packages/config"
)

// Config controls AuthFlows route registration and behavior.
type Config struct {
	Guard              string
	Passwords          string
	Username           string
	Email              string
	LowercaseUsernames bool
	Home               string
	Prefix             string
	Domain             string
	Middleware         []string
	AuthMiddleware     string
	Views              bool
	Limiters           map[string]string
	Features           []string
	Options            map[string]map[string]any
	Paths              map[string]string
}

// ConfigFromRepository loads AuthFlows config from the shared config repository.
func ConfigFromRepository(repo *configpkg.Repository) (Config, error) {
	guard, err := repo.String("authflows.guard")
	if err != nil {
		return Config{}, err
	}
	passwords, err := repo.String("authflows.passwords")
	if err != nil {
		return Config{}, err
	}
	username, err := repo.String("authflows.username")
	if err != nil {
		return Config{}, err
	}
	email, err := repo.String("authflows.email")
	if err != nil {
		return Config{}, err
	}
	lowercase, err := repo.Bool("authflows.lowercase_usernames")
	if err != nil {
		return Config{}, err
	}
	home, err := repo.String("authflows.home")
	if err != nil {
		return Config{}, err
	}
	prefix, err := repo.String("authflows.prefix")
	if err != nil {
		return Config{}, err
	}
	domain, err := repo.String("authflows.domain")
	if err != nil {
		return Config{}, err
	}
	middleware, err := repo.StringSlice("authflows.middleware")
	if err != nil {
		return Config{}, err
	}
	authMiddleware, err := repo.String("authflows.auth_middleware")
	if err != nil {
		return Config{}, err
	}
	views, err := repo.Bool("authflows.views")
	if err != nil {
		return Config{}, err
	}
	featureList, err := repo.StringSlice("authflows.features")
	if err != nil {
		return Config{}, err
	}
	limitsRaw, err := repo.Map("authflows.limiters")
	if err != nil {
		return Config{}, err
	}
	limiters := make(map[string]string, len(limitsRaw))
	for key, value := range limitsRaw {
		limiters[key] = fmt.Sprint(value)
	}

	options := map[string]map[string]any{}
	if repo.Has("authflows.options") {
		raw, err := repo.Map("authflows.options")
		if err != nil {
			return Config{}, err
		}
		for key, value := range raw {
			child, ok := value.(map[string]any)
			if !ok {
				continue
			}
			options[key] = cloneOptionMap(child)
		}
	}

	paths := map[string]string{}
	if repo.Has("authflows.paths") {
		raw, err := repo.Map("authflows.paths")
		if err != nil {
			return Config{}, err
		}
		for key, value := range raw {
			paths[key] = fmt.Sprint(value)
		}
	}

	return Config{
		Guard:              guard,
		Passwords:          passwords,
		Username:           username,
		Email:              email,
		LowercaseUsernames: lowercase,
		Home:               home,
		Prefix:             strings.TrimSuffix(prefix, "/"),
		Domain:             domain,
		Middleware:         middleware,
		AuthMiddleware:     authMiddleware,
		Views:              views,
		Limiters:           limiters,
		Features:           featureList,
		Options:            options,
		Paths:              paths,
	}, nil
}

func cloneOptionMap(items map[string]any) map[string]any {
	cloned := make(map[string]any, len(items))
	for key, value := range items {
		cloned[key] = value
	}
	return cloned
}
