package module

type Manifest struct {
	Name      string
	NewModule Constructor
}

func NewManifest(name string, cons Constructor) *Manifest {
	return &Manifest{
		Name:      name,
		NewModule: cons,
	}
}
