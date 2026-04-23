package horizon

// RedisPayload contains queue metadata Horizon records from a Redis job payload.
type RedisPayload struct {
	Name     string
	Type     string
	Tags     []string
	Silenced bool
}

// RedisPayloadOptions provides explicit metadata that PHP Horizon would derive from serialized payloads.
type RedisPayloadOptions struct {
	Job          string
	Listener     string
	Event        string
	Type         string
	ExistingTags []string
	JobTags      []string
	ListenerTags []string
	EventTags    []string
	InternalTags []string
	ExplicitTags []string
	Silenced     bool
	SilencedTags []string
	Mailable     bool
}

// BuildRedisPayload creates Horizon metadata from explicit job, listener, and event tags.
func BuildRedisPayload(options RedisPayloadOptions) RedisPayload {
	payload := RedisPayload{
		Name: payloadName(options),
		Type: payloadType(options),
	}

	tags := append([]string(nil), options.ExistingTags...)

	if options.ExplicitTags != nil {
		tags = append(tags, options.ExplicitTags...)
	} else {
		tags = append(tags, options.InternalTags...)
		tags = append(tags, options.JobTags...)
		tags = append(tags, options.ListenerTags...)
		tags = append(tags, options.EventTags...)
	}

	payload.Tags = uniqueStrings(tags)
	payload.Silenced = options.Silenced || options.Mailable || hasAnyString(payload.Tags, options.SilencedTags)

	return payload
}

func payloadName(options RedisPayloadOptions) string {
	if options.Listener != "" {
		return options.Listener
	}

	if options.Event != "" && options.Job == "" {
		return options.Event
	}

	return options.Job
}

func payloadType(options RedisPayloadOptions) string {
	if options.Type != "" {
		return options.Type
	}

	if options.Listener != "" {
		return "listener"
	}

	return "job"
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))

	for _, value := range values {
		if value == "" {
			continue
		}

		if _, ok := seen[value]; ok {
			continue
		}

		seen[value] = struct{}{}
		unique = append(unique, value)
	}

	return unique
}

func hasAnyString(values, candidates []string) bool {
	if len(values) == 0 || len(candidates) == 0 {
		return false
	}

	set := make(map[string]struct{}, len(values))

	for _, value := range values {
		set[value] = struct{}{}
	}

	for _, candidate := range candidates {
		if _, ok := set[candidate]; ok {
			return true
		}
	}

	return false
}
