package app

import (
	"fmt"

	"github.com/makinj/go-le/internal/module"
	"github.com/makinj/go-le/modules/mock"
)

//Configurer interfaces provide go-le apps with the information required to configure the various subcomponents
type Configurer interface {
	GetName() string
	GetModules() []module.WrapConfig
}

//App represents the state of an application
type App struct {
	config     Configurer
	name       string
	running    chan struct{}
	controller *module.Controller
}

//New will create a new app
func New(c Configurer) (a *App, err error) {
	cont, err := module.NewController()
	if err != nil {
		return nil, err
	}

	//TKTK load plugins from config here
	cont.RegisterManifest(mock.Manifest)
	if err != nil {
		return nil, fmt.Errorf("Error registering manifest with controller: %s", err)
	}

	for _, mconf := range c.GetModules() {
		err := cont.RegisterModule(mconf)
		if err != nil {
			return nil, fmt.Errorf("Error registering module with controller: %s", err)
		}

	}

	a = &App{
		name:       c.GetName(),
		controller: cont,
	}

	return a, nil
}

//Run starts the application loop goroutine and provides the error channel for the running application
func (a *App) Run() chan Error {
	a.running = make(chan struct{})

	errchan := make(chan Error, 1)
	go a.handleErrs(errchan)

	go a.controller.StartModules()

	return errchan
}

func (a *App) handleErrs(errchan chan Error) {
	defer close(errchan)

	select {
	case <-a.running:
		break
	case err := <-(a.controller.GetErrChan()):
		if err != nil {
			errchan <- fmt.Errorf("Controller received error: %s\n", err)
		} else {
			errchan <- fmt.Errorf("Controller stopped running\n")
			a.Kill()
		}
	}

	return
}

//Kill instructs the main goroutine to kill the application
func (a *App) Kill() {
	close(a.running)
	return
}

//GetName returns the app name
func (a *App) GetName() string {
	return a.name
}
