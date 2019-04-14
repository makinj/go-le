package module

import (
	"fmt"
	"sync"
)

type Controller struct {
	FactoryRepo *FactoryRepo
	Repo        *Repo
	ErrChan     chan error
}

func NewController() (*Controller, error) {
	fact, err := NewFactoryRepo()
	if err != nil {
		return nil, err
	}

	repo, err := NewRepo()
	if err != nil {
		return nil, err
	}

	cont := &Controller{
		FactoryRepo: fact,
		Repo:        repo,
		ErrChan:     make(chan error),
	}

	return cont, nil
}

func (c *Controller) RegisterManifest(man *Manifest) error {
	return c.FactoryRepo.Register(man)
}

func (c *Controller) RegisterModule(wconf WrapConfigurer) error {
	mtype := wconf.GetType()

	man, err := c.FactoryRepo.GetManifest(mtype)

	mconf := wconf.GetConfig()

	id := wconf.GetId()

	//Make Wrap
	w, err := NewWrap(c, id, man, mconf)
	if err != nil {
		return err
	}

	go c.HandleWrapErrors(w)

	err = c.Repo.Register(w)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) StartModules() {
	wraps := c.Repo.GetWraps()
	for _, wrap := range wraps {
		wrap.Start()
	}
	return
}

func (c *Controller) StopModules() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wraps := c.Repo.GetWraps()
	for _, wrap := range wraps {
		wg.Add(1)
		fmt.Printf("stopping module: %s\n", wrap.GetId())
		wrap.Stop()
		go func(w *Wrap) {
			<-w.GetIsRunningChan()
			fmt.Printf("stopped module: %s\n", w.GetId())
			wg.Done()
		}(wrap)
	}
	return wg
}

func (c *Controller) GetModulesErrChan() chan error {
	return c.ErrChan
}
func (c *Controller) GetRepo() *Repo {
	return c.Repo
}

func (c *Controller) HandleWrapErrors(w *Wrap) {
	errchan := w.GetErrChan()
	for err := range errchan {
		c.ErrChan <- err
	}
}

//TKTK add AddTask
