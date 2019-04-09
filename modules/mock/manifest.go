package mock

import (
	"fmt"

	"github.com/makinj/go-le/internal/lifecycle"
	"github.com/makinj/go-le/internal/module"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("mock", module.Constructor(NewModule), &Config{})
}

type Configurer interface {
	GetName() string
}

type Config struct {
	Name string
}

type Module struct {
	Name   string
	Wrap   *module.Wrap
	handle *lifecycle.Handle
}

func (c Config) GetName() string {
	return c.Name
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	ctmp, err := wrap.GetModuleConfigurer()
	if err != nil {
		return nil, err
	}

	c, ok := ctmp.(Configurer)
	if !ok {
		return nil, fmt.Errorf("Configurer does not implement module interface")
	}

	name := c.GetName()

	handle, err := lifecycle.NewHandle()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Creating Mock Module with Name='%s'\n", name)
	return &Module{
		Name:   name,
		Wrap:   wrap,
		handle: handle,
	}, nil
}

func (m *Module) Start() {
	fmt.Printf("Starting Mock Module loop for Name='%s'\n", m.Name)
	m.handle.ShouldStart()
	m.handle.Started()
	m.handle.AddError(fmt.Errorf("test1"))
	return
}

func (m *Module) Stop() {
	fmt.Printf("Stopping Mock Module loop for Name='%s'\n", m.Name)
	m.handle.ShouldStop()
	m.handle.AddError(fmt.Errorf("test2"))
	m.handle.Stopped()
	return
}

func (m *Module) GetErrChan() chan error {
	return m.handle.GetErrChan()
}

func (m *Module) Running() bool {
	return m.handle.GetIsRunning()
}
