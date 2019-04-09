package lifecycle

import (
	"sync"
)

type Handle struct {
	mu             sync.Mutex
	errchan        chan error
	errsubscribers []chan error

	IsRunningChan chan interface{}
	ShouldRunChan chan interface{}
}

func NewHandle() (*Handle, error) {
	h := &Handle{
		errchan: make(chan error),
	}
	go h.handleErrs()
	return h, nil
}

func (h *Handle) handleErrs() {
	for err := range h.errchan {
		h.mu.Lock()
		for _, sub := range h.errsubscribers {
			go func() { sub <- err }()
		}
		h.mu.Unlock()
	}
	for _, sub := range h.errsubscribers {
		close(sub)
	}
	return
}

func (h *Handle) AddError(err error) {
	go func() {
		h.errchan <- err
	}()
	return
}

func (h *Handle) GetErrChan() chan error {
	sub := make(chan error)
	h.mu.Lock()
	h.errsubscribers = append(h.errsubscribers, sub)
	h.mu.Unlock()
	return sub
}

func (h *Handle) ShouldStart() {
	h.mu.Lock()
	if h.ShouldRunChan == nil || !isopen(h.ShouldRunChan) {
		h.ShouldRunChan = make(chan interface{})
	}
	h.mu.Unlock()
}

func (h *Handle) ShouldStop() {
	h.mu.Lock()
	if h.ShouldRunChan != nil && isopen(h.ShouldRunChan) {
		close(h.ShouldRunChan)
	}
	h.mu.Unlock()
}

func (h *Handle) Started() {
	h.mu.Lock()
	if !isopen(h.IsRunningChan) {
		h.IsRunningChan = make(chan interface{})
	}
	h.mu.Unlock()
}

func (h *Handle) Stopped() {
	h.mu.Lock()
	if isopen(h.IsRunningChan) && h.IsRunningChan != nil {
		close(h.IsRunningChan)
	}
	h.mu.Unlock()
}

func (h *Handle) GetShouldRun() bool {
	return isopen(h.ShouldRunChan)
}

func (h *Handle) GetIsRunning() bool {
	return isopen(h.IsRunningChan)
}

func isopen(c chan interface{}) bool {
	if c == nil {
		return false
	}
	select {
	case <-c:
		return false
	default:
		return true
	}
}

//TKTK add invoke function
