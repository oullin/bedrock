package process

// Command describes a process command.
//
// When Shell is set, the command is executed through the platform shell. When
// Shell is empty, Name and Args are executed directly.
type Command struct {
	Name  string
	Args  []string
	Shell string
}

// String returns the display form used by fakes and assertions.
func (c Command) String() string {
	if c.Shell != "" {
		return c.Shell
	}

	if len(c.Args) == 0 {
		return c.Name
	}

	out := c.Name

	for _, arg := range c.Args {
		out += " " + arg
	}

	return out
}
