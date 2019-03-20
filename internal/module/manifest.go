package module

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Manifest struct {
	Name       string
	Configurer reflect.Type
	NewModule  Constructor
}

func NewManifest(name string, cons Constructor, conf interface{}) *Manifest {
	return &Manifest{
		Name:       name,
		Configurer: reflect.TypeOf(conf),
		NewModule:  cons,
	}
}

func (man *Manifest) MapConfig(confmap map[string]interface{}) (interface{}, error) {
	conf_type := man.Configurer

	conf := reflect.New(conf_type.Elem()).Interface()

	err := mapstructure.Decode(confmap, conf)
	if err != nil {
		return nil, fmt.Errorf("Failed to map config: %s", err)
	}
	return conf, nil
}
