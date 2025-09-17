package jobs

// JobHandler interface defines the contract for all job handlers
type JobHandler interface {
	Execute(params string) error
}

// JobHandlers registry stores all registered job handlers
var JobHandlers = make(map[string]JobHandler)

// RegisterJobHandler registers a job handler with a given name
func RegisterJobHandler(name string, handler JobHandler) {
	JobHandlers[name] = handler
}
