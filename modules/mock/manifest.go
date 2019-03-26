package mock

import (
	"fmt"

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
	Name    string
	Wrap    *module.Wrap
	errchan chan error
	running chan struct{}
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
	fmt.Printf("Creating Mock Module with Name='%s'\n", name)
	return &Module{
		Name:    name,
		Wrap:    wrap,
		errchan: make(chan error),
	}, nil
}

func (m *Module) Start() {
	fmt.Printf("Starting Mock Module loop for Name='%s'\n", m.Name)
	m.running = make(chan struct{})
	go func() { m.errchan <- fmt.Errorf("test") }()
	return
}

func (m *Module) Stop() {
	fmt.Printf("Starting Mock Module loop for Name='%s'\n", m.Name)
	close(m.errchan)
	return
}

func (m *Module) GetErrChan() chan error {
	return m.errchan
}

func (m *Module) Running() bool {
	select {
	case <-m.errchan:
		return false
	default:
		return true
	}
}
