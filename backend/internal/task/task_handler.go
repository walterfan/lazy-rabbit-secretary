package task

type TaskHandler interface {
    Execute() error
}

var taskHandlers = make(map[string]func() TaskHandler)

func RegisterTaskHandler(name string, handlerFunc func() TaskHandler) {
    taskHandlers[name] = handlerFunc
}