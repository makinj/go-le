package module

import (
	"fmt"
	"sync"
)

//
type Repo struct {
	mu       sync.Mutex
	wrappers map[string]*Wrap
	waiting  map[string]chan struct{}
}

//New will create a new factory
func NewRepo() (r *Repo, err error) {

	r = &Repo{
		wrappers: make(map[string]*Wrap),
		waiting:  make(map[string]chan struct{}),
	}

	return r, nil
}

func (r *Repo) Register(w *Wrap) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	//TKTK Lock here
	//TKTK defer Unlock here

	id := w.id

	//check existence
	_, exists := r.wrappers[id]
	if exists {
		return fmt.Errorf("Wrap already registered with id: %s", id)
	}

	r.wrappers[id] = w

	//notify dependents
	wait, exists := r.waiting[id]
	if exists && wait != nil {
		close(wait)
	}
	return nil
}

func (r *Repo) GetWraps() map[string]*Wrap {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.wrappers
}

func (r *Repo) ResolveWrap(id string) chan *Wrap {
	res := make(chan *Wrap)
	go r.resolveWrap(res, id)
	return res
}

func (r *Repo) resolveWrap(res chan *Wrap, id string) {

	r.mu.Lock()

	w, exists := r.wrappers[id]

	for ; !exists; w, exists = r.wrappers[id] {
		wait, exists := r.waiting[id]
		if !exists {
			wait = make(chan struct{})
			r.waiting[id] = wait
		}

		r.mu.Unlock()

		<-wait
		r.mu.Lock()

	}
	res <- w

	r.mu.Unlock()
}
