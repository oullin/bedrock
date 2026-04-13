package mail

// Envelope describes the addressing and metadata of an email message.
type Envelope struct {
	From     Address
	To       []Address
	CC       []Address
	BCC      []Address
	ReplyTo  []Address
	Subject  string
	Tags     []string
	Metadata map[string]string
	Using    []func(*Message)
}
