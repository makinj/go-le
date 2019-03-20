package app

import "github.com/makinj/go-le/internal/module"

//Config represents the information that can be used to configure a go-le application
type Config struct {
	Name    string
	Modules []module.WrapConfig
}

func (c Config) GetName() string {
	return c.Name
}

func (c Config) GetModules() []module.WrapConfig {
	return c.Modules
}
