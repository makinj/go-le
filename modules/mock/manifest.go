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
	Name string
	Wrap *module.Wrap
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
		Name: name,
	}, nil
}

func (m *Module) Run() {
	fmt.Printf("Starting Mock Module loop for Name='%s'\n", m.Name)
	return
}
