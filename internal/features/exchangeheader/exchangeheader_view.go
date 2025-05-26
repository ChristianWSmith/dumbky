package exchangeheader

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	ui         *fyne.Container
	sendButton *widget.Button

	urlEntry     *widget.Entry
	sslCheck     *widget.Check
	methodSelect *widget.Select
}

func newView() *view {

	methodSelect := widget.NewSelect(constants.HttpMethods(), nil)
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_URL)
	urlEntry.TextStyle.Monospace = true
	sslCheck := widget.NewCheck(constants.UI_LABEL_SSL, nil)
	sendButton := widget.NewButton(constants.UI_LABEL_SEND, nil)
	sendButton.Icon = sendButton.Theme().Icon(theme.IconNameMailSend)

	methodSelect.SetSelected(constants.HTTP_METHOD_DEFAULT)

	sslSend := container.NewHBox(sslCheck, sendButton)
	ui := container.NewBorder(nil, nil, methodSelect, sslSend, urlEntry)
	return &view{
		ui:           ui,
		sendButton:   sendButton,
		urlEntry:     urlEntry,
		sslCheck:     sslCheck,
		methodSelect: methodSelect,
	}
}

func (v *view) getUI() *fyne.Container {
	return v.ui
}

func (v *view) setUrlValidator(validator func(string) error) {
	v.urlEntry.Validator = validator
}

func (v *view) setSendHandler(handler func()) {
	v.sendButton.OnTapped = handler

}

func (v *view) enableSend() {
	v.sendButton.Enable()
}

func (v *view) disableSend() {
	v.sendButton.Disable()
}

func (v *view) validateURL() error {
	return v.urlEntry.Validate()
}
