package app

import (
	"fmt"

	"github.com/makinj/go-le/internal/lifecycle"
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
	controller *module.Controller
	handle     *lifecycle.Handle
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

	handle, err := lifecycle.NewHandle()
	if err != nil {
		return nil, err
	}

	a = &App{
		name:       c.GetName(),
		controller: cont,
		handle:     handle,
	}

	go a.handleErrs()

	return a, nil
}

//Run starts the application loop goroutine and provides the error channel for the running application
func (a *App) Start() {
	a.handle.ShouldStart()
	go a.controller.StartModules()
	a.handle.Started()
}

func (a *App) handleErrs() {
	for err := range a.controller.GetErrChan() {
		if err != nil {
			a.handle.AddError(fmt.Errorf("Controller received error: %s\n", err))
		}
	}

	return
}

//Kill instructs the main goroutine to kill the application
func (a *App) Stop() {
	a.handle.ShouldStop()
	a.controller.StopModules().Wait()
	a.handle.Stopped()
	return
}

//GetName returns the app name
func (a *App) GetName() string {
	return a.name
}

func (a *App) GetErrChan() chan error {
	return a.handle.GetErrChan()
}

func (a *App) GetIsRunningChan() chan interface{} {
	return a.handle.IsRunningChan
}

func (a *App) GetIsRunning() bool {
	return a.handle.GetIsRunning()
}

func (a *App) GetShouldRunChan() chan interface{} {
	return a.handle.ShouldRunChan
}

func (a *App) GetShouldRun() bool {
	return a.handle.GetShouldRun()
}
