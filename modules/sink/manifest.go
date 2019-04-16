package sink

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("sink", module.Constructor(NewModule))
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewSink(wrap)
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

func NewSink(wrap *module.Wrap) (*Module, error) {
	b, err := base.MakeModule(wrap)
	if err != nil {
		return nil, err
	}

	p := &Module{
		Module: b,
	}
	return p, nil
}

func (p *Module) Receive(val interface{}) {
	spew.Dump(val)
	return
}
