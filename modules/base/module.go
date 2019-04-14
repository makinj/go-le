package base

import (
	"fmt"

	"github.com/makinj/go-le/internal/lifecycle"
	"github.com/makinj/go-le/internal/module"
)

type Loop func()

type BaseConfigurer interface {
	GetName() string
}

type BaseConfig struct {
	Name string `json:"name"`
}

type Module struct {
	*lifecycle.Handle
	Name       string
	Wrap       *module.Wrap
	BaseConfig BaseConfigurer
	Loop       Loop
}

func (c BaseConfig) GetName() string {
	return c.Name
}

func MakeModule(w *module.Wrap) (Module, error) {
	ctmp, err := w.GetModuleConfigurer()
	if err != nil {
		return Module{}, err
	}

	c, ok := ctmp.(BaseConfigurer)
	if !ok {
		return Module{}, fmt.Errorf("BaseConfigurer does not implement base module interface")
	}

	n := c.GetName()

	h, err := lifecycle.NewHandle()
	if err != nil {
		return Module{}, err
	}

	fmt.Printf("Creating Module with Name='%s'\n", n)
	m := Module{
		Name:       n,
		Wrap:       w,
		Handle:     h,
		BaseConfig: c,
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
