package jobqueue

// SupervisorOptions contains worker supervision settings.
type SupervisorOptions struct {
	Queue string
}

// NewSupervisorOptions creates supervisor options with defaults.
func NewSupervisorOptions(queue string) SupervisorOptions {
	if queue == "" {
		queue = "default"
	}

	return SupervisorOptions{Queue: queue}
}
