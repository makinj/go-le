package app

import (
	"time"
)

type Configurer interface {
	GetName() string
}

type app struct {
	name   string
	killed bool
}

//New will create a new app
func New(c Configurer) (a *app) {
	a = &app{
		name:   c.GetName(),
		killed: false,
	}

	return a
}

func (a *app) Run() chan Error {
	a.killed = false

	errchan := make(chan Error, 1)

	go a.loop(errchan)
	return errchan
}

func (a *app) loop(errchan chan Error) {
	defer close(errchan)

	for cnt := 0; cnt < 3 && (a.killed == false); cnt++ {
		time.Sleep(1 * 1000 * 1000 * 1000)
		//TODO add app logic here
	}

	if a.killed {
		return
	} else {
		errchan <- ERR_TIMEOUT
		return
	}
}

func (a *app) GetName() string {
	return a.name
}

func (a *app) Kill() {
	a.killed = true
}
