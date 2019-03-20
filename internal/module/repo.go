package module

import "fmt"

//
type Repo struct {
	wrappers map[string]*Wrap
}

//New will create a new factory
func NewRepo() (r *Repo, err error) {

	r = &Repo{
		wrappers: make(map[string]*Wrap),
	}

	return r, nil
}

func (r *Repo) Register(w *Wrap) error {
	//TKTK Lock here
	//TKTK defer Unlock here

	id := w.id

	//check existence
	_, exists := r.wrappers[id]
	if exists {
		return fmt.Errorf("Wrap already registered with id: %s", id)
	}

	r.wrappers[id] = w

	return nil
}

func (r *Repo) GetWraps() map[string]*Wrap {
	return r.wrappers
}
