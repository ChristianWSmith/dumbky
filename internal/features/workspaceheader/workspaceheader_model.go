package workspaceheader

import (
	"fyne.io/fyne/v2/data/binding"
)

type model struct {
	requestNameBinding binding.String
}

func newModel() *model {
	return &model{
		requestNameBinding: binding.NewString(),
	}
}

func (m *model) setRequestNameListener(handler func()) {
	m.requestNameBinding.AddListener(binding.NewDataListener(handler))
}

func (m *model) getRequestName() string {
	requestName, _ := m.requestNameBinding.Get()
	return requestName
}

func (m *model) setRequestName(requestName string) {
	m.requestNameBinding.Set(requestName)
}
