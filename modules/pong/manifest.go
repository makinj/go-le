package pong

import (
	"fmt"

	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("pong", module.Constructor(NewModule))
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewPong(wrap)
	return m, err
}

type Configurer interface {
	base.Configurer
}

type Config struct {
	base.Config
}

type Module struct {
	base.Module
}

func NewPong(wrap *module.Wrap) (*Module, error) {
	b, err := base.MakeModule(wrap)
	if err != nil {
		return nil, err
	}

	p := &Module{
		Module: b,
	}
	return p, nil
}

func (p *Module) Pong() {
	fmt.Printf("Pong\n")
	return
}
