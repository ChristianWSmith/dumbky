package keyvalue

import (
	"dumbky/internal/state"

	"fyne.io/fyne/v2/data/binding"
)

type model struct {
	enabled binding.Bool
	key     binding.String
	value   binding.String
}

func newModel() *model {
	enabled := binding.NewBool()
	enabled.Set(true)
	return &model{
		enabled: enabled,
		key:     binding.NewString(),
		value:   binding.NewString(),
	}
}
func (m *model) toState() state.KeyValueState {
	enabled, _ := m.enabled.Get()
	key, _ := m.key.Get()
	value, _ := m.value.Get()
	return state.KeyValueState{
		Enabled: enabled,
		Key:     key,
		Value:   value,
	}
}

func (m *model) loadState(kvs state.KeyValueState) {
	m.enabled.Set(kvs.Enabled)
	m.key.Set(kvs.Key)
	m.value.Set(kvs.Value)
}

func (m *model) isEnabled() bool {
	enabled, _ := m.enabled.Get()
	return enabled
}

func (m *model) getKey() string {
	key, _ := m.key.Get()
	return key
}

func (m *model) getValue() string {
	value, _ := m.value.Get()
	return value
}

func (m *model) setEnabledListener(handler func()) {
	m.enabled.AddListener(binding.NewDataListener(handler))
}
