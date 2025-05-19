package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/ui/common"
	"dumbky/internal/ui/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type RequestBodyView struct {
	UI                 *fyne.Container
	BodyKeyValueEditor common.KeyValueEditorView
	BodyTypeSelect     *widget.Select
	BodyRawEntry       *widget.Entry
	BodyTypeBinding    binding.String
	BodyRawBinding     binding.String
}

func bodyTypeSelectOnChangedWorker(val string, bodyRawEntry *widget.Entry, bodyKeyValueEditorView common.KeyValueEditorView, bodyContentStack *fyne.Container) {
	if val == constants.UI_BODY_TYPE_FORM {
		fyne.Do(func() {
			bodyRawEntry.Hide()
			bodyKeyValueEditorView.UI.Show()
			bodyContentStack.Refresh()
		})
	} else if val == constants.UI_BODY_TYPE_RAW {
		fyne.Do(func() {
			bodyKeyValueEditorView.UI.Hide()
			bodyRawEntry.Show()
			bodyContentStack.Refresh()
		})
	} else if val == constants.UI_BODY_TYPE_NONE {
		fyne.Do(func() {
			bodyKeyValueEditorView.UI.Hide()
			bodyRawEntry.Hide()
			bodyContentStack.Refresh()
		})
	}
}

func ComposeRequestBodyView() RequestBodyView {
	bodyTypeBind := binding.NewString()
	bodyRawBind := binding.NewString()

	bodyKeyValueEditorView := common.ComposeKeyValueEditorView(validators.ValidateFormBodyKey, validators.ValidateFormBodyValue)
	bodyTypeSelect := widget.NewSelect(constants.UIBodyTypes(), nil)
	bodyRawEntry := widget.NewMultiLineEntry()
	bodyRawEntry.TextStyle.Monospace = true
	bodyRawEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_BODY_TYPE_RAW)

	bodyContentStack := container.NewStack(bodyKeyValueEditorView.UI, bodyRawEntry)

	bodyTypeSelect.Bind(bodyTypeBind)
	bodyRawEntry.Bind(bodyRawBind)

	bodyRawEntry.Validator = validators.ValidateRawBodyContent

	bodyTypeSelectOnChangedOld := bodyTypeSelect.OnChanged
	bodyTypeSelect.OnChanged = func(val string) {
		bodyTypeSelectOnChangedOld(val)
		go bodyTypeSelectOnChangedWorker(val, bodyRawEntry, bodyKeyValueEditorView, bodyContentStack)
	}
	bodyTypeSelect.SetSelectedIndex(0)

	ui := container.NewBorder(bodyTypeSelect, nil, nil, nil, bodyContentStack)

	return RequestBodyView{
		ui,
		bodyKeyValueEditorView,
		bodyTypeSelect,
		bodyRawEntry,
		bodyTypeBind,
		bodyRawBind,
	}
}
