package hashing

// Driver identifies a hashing algorithm.
type Driver string

const (
	DriverBcrypt   Driver = "bcrypt"
	DriverArgon2i  Driver = "argon2i"
	DriverArgon2id Driver = "argon2id"
)
