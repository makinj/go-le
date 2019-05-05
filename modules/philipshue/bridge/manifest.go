package bridge

import (
	"fmt"

	"github.com/amimof/huego"
	"github.com/davecgh/go-spew/spew"
	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("bridge", module.Constructor(NewModule))
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewBridge(wrap)
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

func NewBridge(wrap *module.Wrap) (*Module, error) {
	b, err := base.MakeModule(wrap)
	if err != nil {
		return nil, err
	}

	p := &Module{
		Module: b,
	}
	return p, nil
}

type message interface {
	GetValue() string
}

func (m *Module) Receive(val interface{}) {
	spew.Dump(val)

	_, ok := val.(message)
	if !ok {
		m.AddError(fmt.Errorf("message does not implement message interface"))
		return
	}

	bridge := huego.New("192.168.65.142", "")
	fmt.Println("test")
	lights, err := bridge.GetLights()
	fmt.Println("test")
	if err != nil {
		m.AddError(err)
		return
	}
	fmt.Println("test")

	v := lights[0].IsOn()

	for _, l := range lights {
		if v {
			err = l.Off()
		} else {
			err = l.On()
		}

		if err != nil {
			m.AddError(err)
			return
		}
	}
	return
}
