package model

type ModelImpl struct {
}

type Bindings struct {
}

type Model interface {
	GetBindings() Bindings
}

func NewModel() Model {
	return &ModelImpl{}
}

func (m *ModelImpl) GetBindings() Bindings {
	return Bindings{}
}
