package module

type Factory func(configuration map[string]interface{}, wrap *Wrap) (Module, error)
