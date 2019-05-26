package bridge

import (
	"fmt"

	"github.com/amimof/huego"
	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("hue-bridge", module.Constructor(NewModule))
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewBridge(wrap)
	return m, err
}

type Configurer interface {
	base.Configurer
	GetIp() string
	GetKey() string
}

type Config struct {
	base.Config
	Ip  string
	Key string
}

func (c Config) GetIp() string {
	return c.Ip
}

func (c Config) GetKey() string {
	return c.Key
}

type Module struct {
	base.Module
	Ip  string
	Key string
}

func NewBridge(w *module.Wrap) (*Module, error) {

	c := &Config{}
	err := w.MapModuleConfigurer(c)
	if err != nil {
		return &Module{}, err
	}

	b, err := base.MakeModule(w)
	if err != nil {
		return nil, err
	}

	ip := c.GetIp()

	k := c.GetKey()

	p := &Module{
		Module: b,
		Ip:     ip,
		Key:    k,
	}
	return p, nil
}

type message interface {
	//GetValue() string
}

func (m *Module) Receive(val interface{}) {
	//spew.Dump(val)

	_, ok := val.(message)
	if !ok {
		m.AddError(fmt.Errorf("message does not implement message interface"))
		return
	}

	bridge := huego.New(m.Ip, m.Key)
	lights, err := bridge.GetGroup(8)
	if err != nil {
		m.AddError(err)
		return
	}

	v := lights.IsOn()

	if v {
		err = lights.Off()
	} else {
		err = lights.On()
	}

	if err != nil {
		m.AddError(err)
		return
	}
	return
}
