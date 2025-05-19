package exchange

import (
	"dumbky/internal/constants"
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
	bodyTypeSelect := widget.NewSelect(constants.UIBodyTypes(), nil)
	bodyRawEntry := widget.NewMultiLineEntry()
	bodyRawEntry.TextStyle.Monospace = true
	bodyRawEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_BODY_TYPE_RAW)

	bodyContentStack := container.NewStack(bodyKeyValueEditorView.UI, bodyRawEntry)

	bodyTypeSelect.Bind(bodyTypeBind)
	bodyRawEntry.Bind(bodyRawBind)

	bodyTypeSelect.OnChanged = func(val string) {
		// TODO: background
		if val == constants.UI_BODY_TYPE_FORM {
			bodyRawEntry.Hide()
			bodyKeyValueEditorView.UI.Show()
			bodyContentStack.Refresh()
		} else if val == constants.UI_BODY_TYPE_RAW {
			bodyKeyValueEditorView.UI.Hide()
			bodyRawEntry.Show()
			bodyContentStack.Refresh()
		} else if val == constants.UI_BODY_TYPE_NONE {
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