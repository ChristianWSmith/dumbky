package view

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"

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
}

type Bindables struct {
	Method, URL features.StringBindable
	UseSSL      features.BoolBindable
}

type View interface {
	GetBindables() Bindables
	CanvasObject() fyne.CanvasObject
	SetUrlValidator(validator func(string) error)
	SetSendHandler(handler func())
	EnableSend()
	DisableSend()
	ValidateURL() error
}

func NewView(requestUI, responseUI fyne.CanvasObject) View {

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

	requestResponseView := container.NewHSplit(requestUI, responseUI)
	ui := container.NewBorder(headerUI, nil, nil, nil, requestResponseView)
	return &viewImpl{
		ui:           ui,
		sendButton:   sendButton,
		urlEntry:     urlEntry,
		sslCheck:     sslCheck,
		methodSelect: methodSelect,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{
		Method: v.methodSelect,
		URL:    v.urlEntry,
		UseSSL: v.sslCheck,
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

func (v *viewImpl) ValidateURL() error {
	return v.urlEntry.Validate()
}
