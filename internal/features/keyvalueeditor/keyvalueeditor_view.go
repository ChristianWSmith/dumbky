package keyvalueeditor

import (
	"fyne.io/fyne/v2"
)

type view struct {
	ui *fyne.Container
}

func newView() *view {
	return &view{
		ui: nil,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}

