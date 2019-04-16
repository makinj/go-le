package barcode

import (
	"fmt"
	"time"

	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/base"
)

var Manifest *module.Manifest

func init() {
	Manifest = module.NewManifest("barcode", module.Constructor(NewModule), &Config{})
}

func NewModule(wrap *module.Wrap) (module.Module, error) {
	m, err := NewBarcode(wrap)
	return m, err
}

type Configurer interface {
	base.Configurer
	GetReceiverId() string
}

type Config struct {
	base.Config
	Name       string `json:"name"`
	ReceiverId string `json:"receiverid"`
}

func (c Config) GetName() string {
	return c.Name
}

func (c Config) GetReceiverId() string {
	return c.ReceiverId
}

type Barcode struct {
	base.Module
	ReceiverId string
}

func NewBarcode(wrap *module.Wrap) (*Barcode, error) {
	ctmp, err := wrap.GetModuleConfigurer()
	if err != nil {
		return nil, err
	}

	c, ok := ctmp.(Configurer)
	if !ok {
		fmt.Println(ctmp)
		return nil, fmt.Errorf("Configurer does not implement barcode config interface")
	}

	rid := c.GetReceiverId()

	b, err := base.MakeModule(wrap)
	if err != nil {
		return nil, err
	}

	p := &Barcode{
		Module:     b,
		ReceiverId: rid,
	}
	p.Module.Loop = p.loop
	return p, nil
}

func (p *Barcode) loop() {
	defer p.Module.Stopped()

	ticker := time.NewTicker(1000000000)

	fmt.Printf("Barcode module loop started\n")
	for p.Module.GetShouldRun() {
		select {
		case <-p.GetShouldRunChan():
		case <-ticker.C:
			go func() {
				err := p.barcode()
				if err != nil {
					p.AddError(err)
				}

			}()
		}

	}
	return
}

func (p *Barcode) barcode() error {
	repo := p.Wrap.GetRepo()

	wchan := repo.ResolveWrap(p.ReceiverId)

	var receiverer Receiverer

	w := <-wchan
	ptmp := w.GetModule()
	var ok bool
	receiverer, ok = ptmp.(Receiverer)
	if !ok {
		return fmt.Errorf("%s does not implement Receiverer interface", p.ReceiverId)
	}

	fmt.Printf("Barcode\n")
	go receiverer.Receiver()
	return nil
}
