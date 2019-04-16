package barcode

import (
	"fmt"

	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("barcode", module.Constructor(Constructor))
}

func Constructor(wrap *module.Wrap) (module.Module, error) {
	m, err := NewModule(wrap)
	return m, err
}

type Configurer interface {
	base.Configurer
	GetReceiverId() string
	GetFileName() string
}

type Config struct {
	base.Config
	FileName   string `json:"FileName"`
	ReceiverId string `json:"ReceiverId"`
}

func (c Config) GetReceiverId() string {
	return c.ReceiverId
}

func (c Config) GetFileName() string {
	return c.FileName
}

type Module struct {
	base.Module
	FileName   string
	ReceiverId string
	device     *device
}

func NewModule(w *module.Wrap) (*Module, error) {
	c := &Config{}
	err := w.MapModuleConfigurer(c)
	if err != nil {
		return nil, err
	}

	rid := c.GetReceiverId()

	fn := c.GetFileName()

	dev, err := NewDevice(fn)
	if err != nil {
		return nil, err
	}

	b, err := base.MakeModule(w)
	if err != nil {
		return nil, err
	}

	m := &Module{
		Module:     b,
		ReceiverId: rid,
		device:     dev,
	}
	m.Module.Loop = m.loop
	return m, nil
}

func (m *Module) loop() {
	defer m.Module.Stopped()
	msgs, errs := m.device.Run()
	for m.Module.GetShouldRun() {
		select {
		case <-m.GetShouldRunChan():
		case msg := <-msgs:
			go func() {
				err := m.publishBarcode(msg)
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

type Barcode struct {
	value string
}

func (m *Module) publishBarcode(val string) error {
	repo := m.Wrap.GetRepo()

	wchan := repo.ResolveWrap(m.ReceiverId)

	var r Receiver

	w := <-wchan
	rtmp := w.GetModule()
	var ok bool
	r, ok = rtmp.(Receiver)
	if !ok {
		return fmt.Errorf("%s does not implement Receiver interface", m.ReceiverId)
	}

	go r.Receive(&Barcode{value: val})
	return nil
}
