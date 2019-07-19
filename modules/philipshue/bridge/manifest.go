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

type controllable interface {
	IsOn() bool
	On() error
	Off() error
}

func (m *Module) Receive(val interface{}) {
	//spew.Dump(val)

	msg, ok := val.(message)
	if !ok {
		m.AddError(fmt.Errorf("message does not implement message interface"))
		return
	}

	bridge := huego.New(m.Ip, m.Key)

	fmt.Println(bridge.GetGroups())
	fmt.Println(bridge.GetLights())

	var target controllable
	var err error

	switch msg {
	default:
		target, err = bridge.GetGroup(8)

	}

	if err != nil {
		m.AddError(err)
		return
	}

	v := target.IsOn()

	if v {
		err := target.Off()
		if err != nil {
			m.AddError(err)
			return
		}
	} else {
		err := target.On()
		if err != nil {
			m.AddError(err)
			return
		}
	}

	return
}
