package workspaceheader

import (
	"fyne.io/fyne/v2/data/binding"
)

type model struct {
	titleBinding binding.String
}

func newModel() *model {
	return &model{
		titleBinding: binding.NewString(),
	}
}

func (m *model) setTitleListener(handler func()) {
	m.titleBinding.AddListener(binding.NewDataListener(handler))
}

func (m *model) getTitle() string {
	title, _ := m.titleBinding.Get()
	return title
}

func (m *model) setTitle(title string) {
	m.titleBinding.Set(title)
}
