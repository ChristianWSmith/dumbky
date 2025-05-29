package request

import (
	"dumbky/internal/constants"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	ui fyne.CanvasObject

	bodyRawEntry     *widget.Entry
	bodyTypeSelect   *widget.Select
	bodyContentStack *fyne.Container
}

func newView(queryParamsUI, pathParamsUI, headersUI, bodyFormUI fyne.CanvasObject) *view {

	queryParamsTab := container.NewTabItem(constants.UI_LABEL_QUERY_PARAMETERS, queryParamsUI)
	pathParamsTab := container.NewTabItem(constants.UI_LABEL_PATH_PARAMETERS, pathParamsUI)
	headersTab := container.NewTabItem(constants.UI_LABEL_HEADERS, headersUI)

	bodyTypeSelect := widget.NewSelect(constants.UIBodyTypes(), nil)
	bodyRawEntry := widget.NewMultiLineEntry()
	bodyRawEntry.TextStyle.Monospace = true
	bodyRawEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_BODY_TYPE_RAW)

	bodyContentStack := container.NewStack(bodyFormUI, bodyRawEntry)

	bodyTypeSelect.SetSelected(constants.UI_BODY_TYPE_DEFAULT)

	bodyRawEntry.Validator = validators.ValidateRawBodyContent

	bodyUI := container.NewBorder(bodyTypeSelect, nil, nil, nil, bodyContentStack)

	bodyTab := container.NewTabItem(constants.UI_LABEL_BODY, bodyUI)

	tabs := container.NewAppTabs(queryParamsTab, pathParamsTab, headersTab, bodyTab)
	ui := container.NewBorder(nil, nil, nil, nil, tabs)

	return &view{
		ui:               ui,
		bodyRawEntry:     bodyRawEntry,
		bodyTypeSelect:   bodyTypeSelect,
		bodyContentStack: bodyContentStack,
	}
}

func (v *view) canvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *view) showBodyRaw() {
	v.bodyRawEntry.Show()
	v.bodyContentStack.Refresh()
}

func (v *view) hideBodyRaw() {
	v.bodyRawEntry.Hide()
	v.bodyContentStack.Refresh()
}

func (v *view) enabledBodyTypeSelect() {
	v.bodyTypeSelect.Enable()
}

func (v *view) disabledBodyTypeSelect() {
	v.bodyTypeSelect.Disable()
}

func (v *view) validateBodyRaw() error {
	return v.bodyRawEntry.Validate()
}
