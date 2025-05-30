package view

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	ui fyne.CanvasObject

	bodyRawEntry     *widget.Entry
	bodyTypeSelect   *widget.Select
	bodyContentStack *fyne.Container
}

type Bindables struct {
	BodyType, BodyRaw features.StringBindable
}

type View interface {
	GetBindables() Bindables
	CanvasObject() fyne.CanvasObject
	SetBodyRawVisible(visible bool)
	SetBodyTypeSelectEnabled(enabled bool)
	Validate() error
}

func NewView(queryParamsUI, pathParamsUI, headersUI, bodyFormUI fyne.CanvasObject) View {

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

	return &viewImpl{
		ui:               ui,
		bodyRawEntry:     bodyRawEntry,
		bodyTypeSelect:   bodyTypeSelect,
		bodyContentStack: bodyContentStack,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{
		BodyType: v.bodyTypeSelect,
		BodyRaw:  v.bodyRawEntry,
	}
}

func (v *viewImpl) SetBodyRawVisible(visible bool) {
	if visible {

		v.bodyRawEntry.Show()
		v.bodyContentStack.Refresh()
	} else {

		v.bodyRawEntry.Hide()
		v.bodyContentStack.Refresh()
	}
}

func (v *viewImpl) SetBodyTypeSelectEnabled(enabled bool) {
	if enabled {
		v.bodyTypeSelect.Enable()
	} else {
		v.bodyTypeSelect.Disable()
	}
}

func (v *viewImpl) Validate() error {
	return v.bodyRawEntry.Validate()
}
