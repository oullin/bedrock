package hashing

import "errors"

var (
	errBcryptAlgorithm   = errors.New("hashing: this password does not use the bcrypt algorithm")
	errArgon2iAlgorithm  = errors.New("hashing: this password does not use the argon2i algorithm")
	errArgon2idAlgorithm = errors.New("hashing: this password does not use the argon2id algorithm")
)
