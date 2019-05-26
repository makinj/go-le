package dash

import (
	"fmt"

	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("dash", module.Constructor(Constructor))
}

func Constructor(wrap *module.Wrap) (module.Module, error) {
	m, err := NewModule(wrap)
	return m, err
}

type Configurer interface {
	base.Configurer
	GetReceiverId() string
	GetInterface() string
	GetPort() uint
}

type Config struct {
	base.Config
	Interface  string `json:"Interface"`
	Port       uint   `json:"Port"`
	ReceiverId string `json:"ReceiverId"`
}

func (c Config) GetReceiverId() string {
	return c.ReceiverId
}

func (c Config) GetInterface() string {
	return c.Interface
}

func (c Config) GetPort() uint {
	return c.Port
}

type Module struct {
	base.Module
	Port       uint
	Interface  string
	ReceiverId string
	server     *server
}

func NewModule(w *module.Wrap) (*Module, error) {
	c := &Config{}
	err := w.MapModuleConfigurer(c)
	if err != nil {
		return nil, err
	}

	rid := c.GetReceiverId()

	iface := c.GetInterface()

	p := c.GetPort()

	srv, err := NewServer(iface, p)
	if err != nil {
		return nil, err
	}

	b, err := base.MakeModule(w)
	if err != nil {
		return nil, err
	}

	m := &Module{
		Module:     b,
		Interface:  iface,
		Port:       p,
		ReceiverId: rid,
		server:     srv,
	}
	m.Module.Loop = m.loop
	return m, nil
}

func (m *Module) loop() {
	defer m.Module.Stopped()
	msgs, errs := m.server.Run()
	for m.Module.GetShouldRun() {
		select {
		case <-m.GetShouldRunChan():
		case msg := <-msgs:
			go func() {
				err := m.publishPress(msg)
				if err != nil {
					m.AddError(err)
				}

			}()
		case err := <-errs:
			m.AddError(err)

		}

	}
	return
}

type Receiver interface {
	Receive(interface{})
}

type Press struct {
	ip string
}

func (p *Press) GetIp() string {
	return p.ip
}

func (m *Module) publishPress(val string) error {
	repo := m.Wrap.GetRepo()

	wchan := repo.ResolveWrap(m.ReceiverId)

	var r Receiver

	w := <-wchan
	r, ok := w.GetModule().(Receiver)
	if !ok {
		fmt.Println(w)
		fmt.Println(w.GetModule())
		fmt.Println(r)
		return fmt.Errorf("%s does not implement Receiver interface", m.ReceiverId)
	}

	go r.Receive(&Press{ip: val})
	return nil
}
