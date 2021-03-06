package base

import (
	"fmt"

	"github.com/makinj/go-le/internal/lifecycle"
	"github.com/makinj/go-le/internal/module"
)

type Loop func()

type Module struct {
	*lifecycle.Handle
	Name   string
	Wrap   *module.Wrap
	Config Configurer
	Loop   Loop
}

func MakeModule(w *module.Wrap) (Module, error) {
	c := &Config{}
	err := w.MapModuleConfigurer(c)
	if err != nil {
		return Module{}, err
	}

	n := c.GetName()

	h, err := lifecycle.NewHandle()
	if err != nil {
		return Module{}, err
	}

	fmt.Printf("Creating Module with Name='%s'\n", n)
	m := Module{
		Name:   n,
		Wrap:   w,
		Handle: h,
		Config: c,
	}
	m.Loop = m.defaultLoop
	return m, nil
}

func (m *Module) Start() {
	fmt.Printf("Starting Module loop for Name='%s'\n", m.Name)
	m.ShouldStart()
	m.Started()
	if m.Loop != nil {
		go m.Loop()

	} else {
		go m.defaultLoop()
	}
	return
}

func (m *Module) defaultLoop() {
	defer m.Stopped()
	<-m.GetShouldRunChan()
	return
}

func (m *Module) Stop() {
	fmt.Printf("Stopping Module loop for Name='%s'\n", m.Name)
	m.ShouldStop()
	return
}
