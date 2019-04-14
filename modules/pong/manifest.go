package pong

import (
	"fmt"

	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("pong", module.Constructor(NewModule), &Config{})
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewPong(wrap)
	return m, err
}

type Configurer interface {
	base.BaseConfigurer
}

type Config struct {
	base.BaseConfig
}

type Pong struct {
	base.Module
}

func NewPong(wrap *module.Wrap) (*Pong, error) {
	b, err := base.MakeModule(wrap)
	if err != nil {
		return nil, err
	}

	p := &Pong{
		Module: b,
	}
	return p, nil
}

func (p *Pong) Pong() {
	fmt.Printf("Pong\n")
	return
}
