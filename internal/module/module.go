package module

type Module interface {
	Start()
	Stop()
	GetIsRunning() bool
	GetErrChan() chan error
	GetIsRunningChan() chan interface{}
}

type Constructor func(*Wrap) (Module, error)
