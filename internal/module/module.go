package module

type Module interface {
	Run()
}

type Constructor func(*Wrap) (Module, error)
