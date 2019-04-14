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
	base.BaseConfigurer
}

type Config struct {
	base.BaseConfig
}

type Mock struct {
	base.Module
}

func (c Config) GetName() string {
	return c.Name
}

func NewMock(wrap *module.Wrap) (*Mock, error) {
	b, err := base.MakeModule(wrap)
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
