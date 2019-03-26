package module

type Module interface {
	Start()
	Stop()
	Running() bool
	GetErrChan() chan error
}

type Constructor func(*Wrap) (Module, error)
