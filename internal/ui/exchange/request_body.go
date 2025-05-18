package exchange

import (
	"dumbky/internal/ui/common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type RequestBodyView struct {
	UI *fyne.Container
	BodyKeyValue common.KeyValueEditorView
	BodyType binding.String
	BodyRaw binding.String
}

func ComposeRequestBodyView() RequestBodyView {
	bodyTypeBind := binding.NewString()
	bodyRawBind := binding.NewString()

	bodyKeyValueEditorView := common.ComposeKeyValueEditorView()
	bodyTypeSelect := widget.NewSelect([]string{"FORM", "RAW", "NONE"}, nil)
	bodyRawEntry := widget.NewMultiLineEntry()
	bodyRawEntry.TextStyle.Monospace = true

	bodyContentStack := container.NewStack(bodyKeyValueEditorView.UI, bodyRawEntry)

	bodyTypeSelect.Bind(bodyTypeBind)
	bodyRawEntry.Bind(bodyRawBind)

	bodyTypeSelect.OnChanged = func(val string) {
		if val == "FORM" {
			bodyKeyValueEditorView.UI.Show()
			bodyRawEntry.Hide()
			bodyContentStack.Refresh()
		} else if val == "RAW" {
			bodyKeyValueEditorView.UI.Hide()
			bodyRawEntry.Show()
			bodyContentStack.Refresh()
		} else if val == "NONE" {
			bodyKeyValueEditorView.UI.Hide()
			bodyRawEntry.Hide()
			bodyContentStack.Refresh()
		}
	}
	bodyTypeSelect.SetSelectedIndex(0)

	ui := container.NewBorder(bodyTypeSelect, nil, nil, nil, bodyContentStack)

	return RequestBodyView{
		ui,
		bodyKeyValueEditorView,
		bodyTypeBind,
		bodyRawBind,
	}
}