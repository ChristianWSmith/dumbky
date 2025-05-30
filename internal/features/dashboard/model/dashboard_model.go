package model

type modelImpl struct {
}

type Bindings struct {
}

type Model interface {
}

func NewModel() Model {
	return &modelImpl{}
}
