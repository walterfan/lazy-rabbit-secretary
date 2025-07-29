package task

type TaskHandler interface {
	Execute(params string) error
}

var taskHandlers = make(map[string]TaskHandler)

func RegisterTaskHandler(name string, taskHandler TaskHandler) {
	taskHandlers[name] = taskHandler
}
