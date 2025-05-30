package model

import (
	"dumbky/internal/state"

	"fyne.io/fyne/v2/data/binding"
)

type modelImpl struct {
	enabled binding.Bool
	key     binding.String
	value   binding.String
}

type Bindings struct {
	Enabled    binding.Bool
	Key, Value binding.String
}

type Model interface {
	GetBindings() Bindings
	ToState() state.KeyValueState
	LoadState(kvs state.KeyValueState)
	IsEnabled() bool
	GetKey() string
	GetValue() string
	SetEnabledListener(handler func())
}

func NewModel() Model {
	enabled := binding.NewBool()
	enabled.Set(true)
	return &modelImpl{
		enabled: enabled,
		key:     binding.NewString(),
		value:   binding.NewString(),
	}
}

func (m *modelImpl) ToState() state.KeyValueState {
	enabled, _ := m.enabled.Get()
	key, _ := m.key.Get()
	value, _ := m.value.Get()
	return state.KeyValueState{
		Enabled: enabled,
		Key:     key,
		Value:   value,
	}
}

func (m *modelImpl) GetBindings() Bindings {
	return Bindings{
		Enabled: m.enabled,
		Key:     m.key,
		Value:   m.value,
	}
}

func (m *modelImpl) LoadState(kvs state.KeyValueState) {
	m.enabled.Set(kvs.Enabled)
	m.key.Set(kvs.Key)
	m.value.Set(kvs.Value)
}

func (m *modelImpl) IsEnabled() bool {
	enabled, _ := m.enabled.Get()
	return enabled
}

func (m *modelImpl) GetKey() string {
	key, _ := m.key.Get()
	return key
}

func (m *modelImpl) GetValue() string {
	value, _ := m.value.Get()
	return value
}

func (m *modelImpl) SetEnabledListener(handler func()) {
	m.enabled.AddListener(binding.NewDataListener(handler))
}
