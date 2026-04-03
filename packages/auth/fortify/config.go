package fortify

import (
	"fmt"
	"strings"

	configpkg "github.com/gollin/packages/config"
)

// Config controls Fortify route registration and behavior.
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
	Options            map[string]map[string]bool
	Paths              map[string]string
}

// ConfigFromRepository loads Fortify config from the shared config repository.
func ConfigFromRepository(repo *configpkg.Repository) (Config, error) {
	guard, err := repo.String("fortify.guard")
	if err != nil {
		return Config{}, err
	}
	passwords, err := repo.String("fortify.passwords")
	if err != nil {
		return Config{}, err
	}
	username, err := repo.String("fortify.username")
	if err != nil {
		return Config{}, err
	}
	email, err := repo.String("fortify.email")
	if err != nil {
		return Config{}, err
	}
	lowercase, err := repo.Bool("fortify.lowercase_usernames")
	if err != nil {
		return Config{}, err
	}
	home, err := repo.String("fortify.home")
	if err != nil {
		return Config{}, err
	}
	prefix, err := repo.String("fortify.prefix")
	if err != nil {
		return Config{}, err
	}
	domain, err := repo.String("fortify.domain")
	if err != nil {
		return Config{}, err
	}
	middleware, err := repo.StringSlice("fortify.middleware")
	if err != nil {
		return Config{}, err
	}
	authMiddleware, err := repo.String("fortify.auth_middleware")
	if err != nil {
		return Config{}, err
	}
	views, err := repo.Bool("fortify.views")
	if err != nil {
		return Config{}, err
	}
	featureList, err := repo.StringSlice("fortify.features")
	if err != nil {
		return Config{}, err
	}
	limitsRaw, err := repo.Map("fortify.limiters")
	if err != nil {
		return Config{}, err
	}
	limiters := make(map[string]string, len(limitsRaw))
	for key, value := range limitsRaw {
		limiters[key] = fmt.Sprint(value)
	}

	options := map[string]map[string]bool{}
	if repo.Has("fortify.options") {
		raw, err := repo.Map("fortify.options")
		if err != nil {
			return Config{}, err
		}
		for key, value := range raw {
			child, ok := value.(map[string]any)
			if !ok {
				continue
			}
			options[key] = map[string]bool{}
			for option, enabled := range child {
				boolValue, ok := enabled.(bool)
				if ok {
					options[key][option] = boolValue
				}
			}
		}
	}

	paths := map[string]string{}
	if repo.Has("fortify.paths") {
		raw, err := repo.Map("fortify.paths")
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
