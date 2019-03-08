package app

//Config represents the information that can be used to configure a go-le application
type Config struct {
	Name string
}

func (c Config) GetName() string {
	return c.Name
}
