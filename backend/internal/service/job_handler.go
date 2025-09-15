package service

type JobHandler interface {
	Execute(params string) error
}

var JobHandlers = make(map[string]JobHandler)

func RegisterJobHandler(name string, JobHandler JobHandler) {
	JobHandlers[name] = JobHandler
}
