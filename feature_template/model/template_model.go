package model

type modelImpl struct {
}

type Bindings struct {
}

type Model interface {
	GetBindings() Bindings
}

func NewModel() Model {
	m := &modelImpl{}
	return m
}

func (m *modelImpl) GetBindings() Bindings {
	return Bindings{}
}
