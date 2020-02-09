package router

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("router", module.Constructor(NewModule))
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewRouter(wrap)
	return m, err
}

type Configurer interface {
	base.Configurer
}

type Config struct {
	base.Config
	Rules []Rule `json:"Rules"`
}

type Rule struct {
}

type Module struct {
	base.Module
	Rules []Rule
}

func NewRouter(w *module.Wrap) (*Module, error) {
	c := &Config{}
	err := w.MapModuleConfigurer(c)
	if err != nil {
		return nil, err
	}

	b, err := base.MakeModule(w)
	if err != nil {
		return nil, err
	}

	p := &Module{
		Module: b,
		Rules:  c.Rules,
	}
	return p, nil
}

func (p *Module) Receive(val interface{}) {
	spew.Dump(val)
	spew.Dump(p.Rules)
	return
}
