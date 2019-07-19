package base

type Configurer interface {
	GetName() string
}

type Config struct {
	Name string `json:"Name"`
}

func (c Config) GetName() string {
	return c.Name
}
