package view

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	ui         fyne.CanvasObject
	sendButton *widget.Button

	urlEntry     *widget.Entry
	sslCheck     *widget.Check
	methodSelect *widget.Select

	bodyRawEntry     *widget.Entry
	bodyTypeSelect   *widget.Select
	bodyContentStack *fyne.Container
}

type Bindables struct {
	Method, URL       features.StringBindable
	UseSSL            features.BoolBindable
	BodyType, BodyRaw features.StringBindable
}

type View interface {
	GetBindables() Bindables
	CanvasObject() fyne.CanvasObject
	SetUrlValidator(validator func(string) error)
	SetSendHandler(handler func())
	EnableSend()
	DisableSend()
	SetBodyRawVisible(visible bool)
	SetBodyTypeSelectEnabled(enabled bool)
	Validate() error
}

func NewView(queryParamsUI, pathParamsUI, headersUI, bodyFormUI, responseUI fyne.CanvasObject) View {

	methodSelect := widget.NewSelect(constants.HttpMethods(), nil)
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_URL)
	urlEntry.TextStyle.Monospace = true
	sslCheck := widget.NewCheck(constants.UI_LABEL_SSL, nil)
	sendButton := widget.NewButton(constants.UI_LABEL_SEND, nil)
	sendButton.Icon = sendButton.Theme().Icon(theme.IconNameMailSend)

	methodSelect.SetSelected(constants.HTTP_METHOD_DEFAULT)

	sslSend := container.NewHBox(sslCheck, sendButton)
	headerUI := container.NewBorder(nil, nil, methodSelect, sslSend, urlEntry)

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
	requestUI := container.NewBorder(nil, nil, nil, nil, tabs)

	requestResponseView := container.NewHSplit(requestUI, responseUI)
	ui := container.NewBorder(headerUI, nil, nil, nil, requestResponseView)
	return &viewImpl{
		ui:               ui,
		sendButton:       sendButton,
		urlEntry:         urlEntry,
		sslCheck:         sslCheck,
		methodSelect:     methodSelect,
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
		Method:   v.methodSelect,
		URL:      v.urlEntry,
		UseSSL:   v.sslCheck,
		BodyType: v.bodyTypeSelect,
		BodyRaw:  v.bodyRawEntry,
	}
}

func (v *viewImpl) SetUrlValidator(validator func(string) error) {
	v.urlEntry.Validator = validator
}

func (v *viewImpl) SetSendHandler(handler func()) {
	v.sendButton.OnTapped = handler

}

func (v *viewImpl) EnableSend() {
	v.sendButton.Enable()
}

func (v *viewImpl) DisableSend() {
	v.sendButton.Disable()
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
	err := v.urlEntry.Validate()
	if err != nil {
		return err
	}
	return v.bodyRawEntry.Validate()
}
