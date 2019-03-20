package module

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

type Wrap struct {
	id         string
	controller *Controller
	manifest   *Manifest
	modconf    map[string]interface{}
	errchan    chan error
	running    chan struct{}
	module     Module
}

func NewWrap(cont *Controller, id string, man *Manifest, mconf map[string]interface{}) (*Wrap, error) {
	return &Wrap{
		id:         id,
		manifest:   man,
		modconf:    mconf,
		controller: cont,
	}, nil
}

func (w *Wrap) Start() {
	w.errchan = make(chan error)
	w.running = make(chan struct{})

	go w.loop()
	return
}

func (w *Wrap) Running() bool {
	select {
	case <-w.running:
		return false
	default:
		return true
	}
}

func (w *Wrap) loop() {
	defer close(w.running)
	defer close(w.errchan)

	for w.Running() {
		moderrs := make(chan error)

		m, err := w.manifest.NewModule(w)
		if err != nil {
			w.errchan <- err
			return
		}

		go m.Run()
		w.module = m

		select {
		case err, ok := <-moderrs:
			if ok {
				w.errchan <- err
			}
		case <-w.running:
		}
		//TODO add something to kill the module
		close(moderrs)

	}
	return
}

func (w *Wrap) Stop() {
	close(w.running)
	close(w.errchan)
	return
}

func (w Wrap) GetModuleConfigurer() (interface{}, error) {
	return w.manifest.MapConfig(w.modconf)
}
