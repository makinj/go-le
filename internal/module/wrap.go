package module

import (
	"fmt"
	"sync"

	"github.com/makinj/go-le/internal/lifecycle"
	"github.com/mitchellh/mapstructure"
)

type WrapConfig struct {
	Type   string
	Config map[string]interface{}
	Id     string
}

type WrapConfigurer interface {
	GetType() string
	GetConfig() map[string]interface{}
	GetId() string
}

func (wc WrapConfig) GetType() string {
	return wc.Type
}

func (wc WrapConfig) GetId() string {
	return wc.Id
}

func (wc WrapConfig) GetConfig() map[string]interface{} {
	return wc.Config
}

type Wrapper interface {
}

type Wrap struct {
	*Controller
	id       string
	manifest *Manifest
	modconf  map[string]interface{}
	modulemu sync.Mutex
	module   Module
	handle   *lifecycle.Handle
}

func NewWrap(cont *Controller, id string, man *Manifest, mconf map[string]interface{}) (*Wrap, error) {

	if man == nil {
		return nil, fmt.Errorf("[%s]\tCannot find module type definition", id)
	}

	handle, err := lifecycle.NewHandle()
	if err != nil {
		return nil, fmt.Errorf("[%s]\tError making lifecycle handle: %s", id, err)
	}

	wrap := &Wrap{
		Controller: cont,
		id:         id,
		manifest:   man,
		modconf:    mconf,
		handle:     handle,
	}

	m, err := man.NewModule(wrap)
	if err != nil {
		return nil, fmt.Errorf("[%s]\tError making module: %s", id, err)
	}
	wrap.module = m

	return wrap, nil
}

func (w *Wrap) Start() {
	w.handle.ShouldStart()
	go w.loop()
	return
}

func (w *Wrap) Running() bool {
	return w.handle.GetIsRunning()
}

func (w *Wrap) loop() {
	defer w.handle.Stopped()
	w.handle.Started()

	for w.handle.GetShouldRun() {
		if w.manifest == nil {
			//fmt.Println("No Manifest found for module.")
			w.handle.AddError(fmt.Errorf("[%s]\tCannot find module type definition", w.id))
			return
		}

		m, err := w.manifest.NewModule(w)
		if err != nil {
			w.handle.AddError(err)
			return
		}

		w.module = m
		m.Start()
		for m.GetIsRunning() && w.handle.GetShouldRun() {
			select {
			case err, ok := <-(m.GetErrChan()):
				if err != nil && ok {
					w.handle.AddError(err)
				}
				if !ok {
					m.Stop()
				}
			case <-w.handle.ShouldRunChan:
				m.Stop()
			}
		}

		for m.GetIsRunning() {
			select {
			case err, ok := <-(m.GetErrChan()):
				if err != nil && ok {
					w.handle.AddError(err)
				}
				if !ok {
					break
				}
			case <-(m.GetIsRunningChan()):
			}
		}
	}
	return
}

func (w *Wrap) Stop() {
	w.handle.ShouldStop()
	return
}

func (w *Wrap) MapModuleConfigurer(conf interface{}) error {
	err := mapstructure.Decode(w.modconf, conf)
	if err != nil {
		return fmt.Errorf("Failed to map config for Wrap(%s): %s", w.id, err)
	}
	return nil
}

func (w *Wrap) GetErrChan() chan error {
	return w.handle.GetErrChan()
}

func (w *Wrap) GetIsRunningChan() chan interface{} {
	return w.handle.IsRunningChan
}

func (w *Wrap) GetId() string {
	return w.id
}

func (w *Wrap) GetController() *Controller {
	return w.Controller
}

func (w *Wrap) GetModule() Module {

	return w.module
}
