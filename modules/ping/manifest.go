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
	base.BaseConfigurer
	GetPongId() string
}

type Config struct {
	base.BaseConfig
	Name   string `json:"name"`
	PongId string `json:"pongid"`
}

func (c Config) GetName() string {
	return c.Name
}

func (c Config) GetPongId() string {
	return c.PongId
}

type Ping struct {
	base.Module
	PongId string
}

func NewPing(wrap *module.Wrap) (*Ping, error) {
	ctmp, err := wrap.GetModuleConfigurer()
	if err != nil {
		return nil, err
	}

	c, ok := ctmp.(Configurer)
	if !ok {
		fmt.Println(ctmp)
		return nil, fmt.Errorf("Configurer does not implement ping config interface")
	}

	pongid := c.GetPongId()

	b, err := base.MakeModule(wrap)
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
			go func() {
				err := p.ping()
				if err != nil {
					p.AddError(err)
				}

			}()
		}

	}
	return
}

type Ponger interface {
	Pong()
}

func (p *Ping) ping() error {
	repo := p.Wrap.GetRepo()

	wchan := repo.ResolveWrap(p.PongId)

	var ponger Ponger

	w := <-wchan
	ptmp := w.GetModule()
	var ok bool
	ponger, ok = ptmp.(Ponger)
	if !ok {
		return fmt.Errorf("%s does not implement Ponger interface", p.PongId)
	}

	fmt.Printf("Ping\n")
	go ponger.Pong()
	return nil
}
