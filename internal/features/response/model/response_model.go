package model

import "fyne.io/fyne/v2/data/binding"

type modelImpl struct {
	status binding.String
	time   binding.String
	body   binding.String
}

type Bindings struct {
	Status, Time, Body binding.String
}

type Model interface {
	GetBindings() Bindings
	SetResponse(status, time, body string)
}

func NewModel() Model {
	return &modelImpl{
		status: binding.NewString(),
		time:   binding.NewString(),
		body:   binding.NewString(),
	}
}

func (m *modelImpl) GetBindings() Bindings {
	return Bindings{
		Status: m.status,
		Time:   m.time,
		Body:   m.body,
	}
}

func (m *modelImpl) SetResponse(status, time, body string) {
	m.status.Set(status)
	m.time.Set(time)
	m.body.Set(body)
}
