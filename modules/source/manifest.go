package source

import (
	"fmt"
	"log"
	"time"

	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("source", module.Constructor(NewModule))
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewSource(wrap)
	return m, err
}

type Configurer interface {
	base.Configurer
	GetReceiverId() string
	GetMessage() string
}

type Message interface{}

type Config struct {
	base.Config
	ReceiverId string `json:"ReceiverId"`
	Message    Message
}

func (c Config) GetReceiverId() string {
	return c.ReceiverId
}

func (c Config) GetMessage() Message {
	return c.Message
}

type Module struct {
	base.Module
	ReceiverId string
	Message    Message
}

func NewSource(w *module.Wrap) (*Module, error) {
	c := &Config{}
	err := w.MapModuleConfigurer(c)
	if err != nil {
		return nil, err
	}

	rid := c.GetReceiverId()
	msg := c.GetMessage()

	b, err := base.MakeModule(w)
	if err != nil {
		return nil, err
	}

	p := &Module{
		Module:     b,
		ReceiverId: rid,
		Message:    msg,
	}
	p.Module.Loop = p.loop
	return p, nil
}

func (p *Module) loop() {
	defer p.Module.Stopped()

	ticker := time.NewTicker(1000000000)

	fmt.Printf("Module module loop started\n")
	for p.Module.GetShouldRun() {
		select {
		case <-p.GetShouldRunChan():
		case <-ticker.C:
			go func() {
				err := p.send()
				if err != nil {
					p.AddError(err)
				}

			}()
		}

	}
	return
}

type Receiver interface {
	Receive(interface{})
}

func (p *Module) send() error {
	repo := p.Wrap.GetRepo()
	wchan := repo.ResolveWrap(p.ReceiverId)

	log.Println("waiting for chan")
	w := <-wchan
	rtmp := w.GetModule()
	var ok bool

	log.Println("got module")
	r, ok := rtmp.(Receiver)
	if !ok {
		return fmt.Errorf("%s does not implement Receiver interface", p.ReceiverId)
	}

	fmt.Printf("Sending\n")
	go r.Receive(p.Message)
	return nil
}
