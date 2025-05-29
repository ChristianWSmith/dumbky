package model

type modelImpl struct {
}

type Model interface {
}

func NewModel() Model {
	return &modelImpl{}
}
