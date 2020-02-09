package mapper

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("mapper", module.Constructor(NewModule))
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
	Mapping    map[string]interface{} `json:"Mapping"`
	ReceiverId string                 `json:"ReceiverId"`
	Default    interface{}            `json:"Default"`
	FromKey    string                 `json:"FromKey"`
	ToKey      string                 `json:"ToKey"`
}

type Module struct {
	base.Module
	Mapping    map[string]interface{}
	ReceiverId string
	Default    interface{}
	FromKey    string
	ToKey      string
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
		Module:     b,
		Mapping:    c.Mapping,
		Default:    c.Default,
		ReceiverId: c.ReceiverId,
		FromKey:    c.FromKey,
		ToKey:      c.ToKey,
	}
	return p, nil
}

type message map[string]interface{}

type Receiver interface {
	Receive(interface{})
}

func (p *Module) Receive(val interface{}) {
	spew.Dump(val)
	msg, ok := val.(map[string]interface{})
	if !ok {
		p.AddError(fmt.Errorf("message does not implement message interface"))
		return
	}

	fromvaltmp, ok := msg[p.FromKey]
	if !ok {
		p.AddError(fmt.Errorf("message does not have key for mapping"))
		return
	}

	fromval, ok := fromvaltmp.(string)
	if !ok {
		p.AddError(fmt.Errorf("message value isn't a string and can't be mapped"))
		return
	}

	toval, ok := p.Mapping[fromval]
	if !ok {
		toval = p.Default
	}

	tomsg := make(map[string]interface{})
	tomsg[p.ToKey] = toval
	spew.Dump(tomsg)

	repo := p.Wrap.GetRepo()

	wchan := repo.ResolveWrap(p.ReceiverId)
	log.Println("waiting for module")
	w := <-wchan
	log.Println("getting module")
	rtmp := w.GetModule()
	log.Println("got module")
	r, ok := rtmp.(Receiver)
	if !ok {
		p.AddError(fmt.Errorf("%s does not implement Receiver interface", p.ReceiverId))
		return
	}

	fmt.Printf("Sending\n")
	go r.Receive(tomsg)

	//spew.Dump(p)
	return
}
