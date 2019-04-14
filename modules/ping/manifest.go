package ping

import (
	"fmt"
	"time"

	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("ping", module.Constructor(NewModule), &Config{})
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewPing(wrap)
	return m, err
}

type Configurer interface {
	base.Configurer
	GetPongId() string
}

type Config struct {
	base.Config
	PongId string
}

type Ping struct {
	base.Module
	PongId string
}

func (c Config) GetName() string {
	return c.Name
}

func NewPing(wrap *module.Wrap) (*Ping, error) {
	ctmp, err := wrap.GetModuleConfigurer()
	if err != nil {
		return nil, err
	}

	c, ok := ctmp.(Configurer)
	if !ok {
		return nil, fmt.Errorf("Configurer does not implement module interface")
	}

	pongid := c.GetPongId()

	b, err := base.MakeModule(wrap, c)
	if err != nil {
		return nil, err
	}

	p := &Ping{
		Module: b,
		PongId: pongid,
	}
	p.Module.Loop = p.loop
	return p, nil
}

func (p *Ping) loop() {
	defer p.Module.Stopped()

	ticker := time.NewTicker(1000000000)

	fmt.Printf("Ping module loop started\n")
	for p.Module.GetShouldRun() {
		select {
		case <-p.GetShouldRunChan():
		case <-ticker.C:
			go p.ping()
		}

	}
	return
}

type Ponger interface {
	Pong()
}

func (p *Ping) ping() {
	ptmp, err := p.Wrap.GetRepo().GetModule(p.PongerId)
	if err != nil {
		p.AddError(fmt.Errorf("Cannot find module with Id: %s", p.PongerId))
		return
	}

	jponger, ok := ptmp.(Ponger)
	if !ok {
		p.AddError(fmt.Errorf("%s does not implement Ponger interface", p.PongerId))
		return
	}

	fmt.Printf("Ping\n")
	ponger.Pong()
}
