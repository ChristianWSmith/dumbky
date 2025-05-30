package model

type modelImpl struct {
}

type Bindings struct{}

type Model interface {
	GetBindings() Bindings
}

func NewModel() Model {
	return &modelImpl{}
}

func (m *modelImpl) GetBindings() Bindings {
	return Bindings{}
}
