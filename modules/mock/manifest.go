package mock

import (
	"fmt"

	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("mock", module.Constructor(NewModule), &Config{})
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewMock(wrap)
	return m, err
}

type Configurer interface {
	base.Configurer
}

type Config struct {
	base.Config
}

type Mock struct {
	base.Module
}

func (c Config) GetName() string {
	return c.Name
}

func NewMock(wrap *module.Wrap) (*Mock, error) {
	ctmp, err := wrap.GetModuleConfigurer()
	if err != nil {
		return nil, err
	}

	c, ok := ctmp.(Configurer)
	if !ok {
		return nil, fmt.Errorf("Configurer does not implement module interface")
	}

	b, err := base.MakeModule(wrap, c)
	if err != nil {
		return nil, err
	}
	m := &Mock{
		b,
	}
	m.Module.Loop = m.loop
	return m, nil
}

func (m *Mock) loop() {
	defer m.Module.Stopped()
	fmt.Printf("Mock module loop started\n")
	for m.Module.GetShouldRun() {
		<-m.GetShouldRunChan()
	}
	return
}
