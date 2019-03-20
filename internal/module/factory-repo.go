package module

import (
	"fmt"
)

//
type FactoryRepo struct {
	manifests map[string]*Manifest
}

//NewFactoryRepo will create a new FactoryRepo
func NewFactoryRepo() (fr *FactoryRepo, err error) {

	fr = &FactoryRepo{
		manifests: make(map[string]*Manifest),
	}

	return fr, nil
}

func (fr *FactoryRepo) Register(m *Manifest) error {
	//TKTK Lock here
	//TKTK defer Unlock here

	name := m.Name

	//check existence
	_, exists := fr.manifests[name]
	if exists {
		return fmt.Errorf("Type already registered with name: %s", name)
	}

	fr.manifests[name] = m

	return nil
}

func (fr *FactoryRepo) GetManifest(mtype string) (*Manifest, error) {
	man, ok := fr.manifests[mtype]
	if !ok {
		return nil, fmt.Errorf("Module type=%s is not registered", mtype)
	}
	/*
		f := func(wrap *Wrap) (Module, error) {
			mod, err := man.Load(wrap)
			if err != nil {
				return nil, fmt.Errorf("Unable to create module: %s", err)
			}
			return mod, nil
		}
	*/
	return man, nil
}
