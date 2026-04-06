package hashing

import "golang.org/x/crypto/bcrypt"

func parseInfo(hashedValue string) Info {
	if info, ok := parseBcryptInfo(hashedValue); ok {
		return info
	}

	if info, ok := parseArgonInfo(hashedValue); ok {
		return info
	}

	return Info{Options: map[string]int{}}
}

func parseBcryptInfo(hashedValue string) (Info, bool) {
	cost, err := bcrypt.Cost([]byte(hashedValue))

	if err != nil {
		return Info{}, false
	}

	return Info{
		Algorithm: DriverBcrypt,
		Options: map[string]int{
			"rounds": cost,
		},
	}, true
}

func parseArgonInfo(hashedValue string) (Info, bool) {
	parsed, err := parseArgonHash(hashedValue)

	if err != nil {
		return Info{}, false
	}

	algorithm := parsed.algorithm

	if algorithm == "argon2i" {
		algorithm = DriverArgon
	}

	return Info{
		Algorithm: algorithm,
		Options: map[string]int{
			"memory":  int(parsed.memory),
			"time":    int(parsed.time),
			"threads": int(parsed.threads),
		},
	}, true
}
