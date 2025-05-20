package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/ui/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ExchangeHeaderView struct {
	UI            *fyne.Container
	SendButton    *widget.Button
	MethodSelect  *widget.Select
	URLEntry      *widget.Entry
	MethodBinding binding.String
	URLBinding    binding.String
	UseSSLBinding binding.Bool
}

func ComposeExchangeHeaderView() ExchangeHeaderView {
	methodBind := binding.NewString()
	urlBind := binding.NewString()
	sslBind := binding.NewBool()

	methodSelect := widget.NewSelect(constants.HttpMethods(), nil)
	methodSelect.SetSelectedIndex(0)
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_URL)
	urlEntry.TextStyle.Monospace = true
	sslCheck := widget.NewCheck(constants.UI_LABEL_SSL, nil)
	sendButton := widget.NewButton(constants.UI_LABEL_SEND, nil)
	sendButton.Icon = sendButton.Theme().Icon(theme.IconNameMailSend)

	methodSelect.Bind(methodBind)
	urlEntry.Bind(urlBind)
	sslCheck.Bind(sslBind)

	urlEntry.Validator = validators.ValidateURL

	sslSend := container.NewHBox(sslCheck, sendButton)
	ui := container.NewBorder(nil, nil, methodSelect, sslSend, urlEntry)

	return ExchangeHeaderView{
		ui,
		sendButton,
		methodSelect,
		urlEntry,
		methodBind,
		urlBind,
		sslBind,
	}
}
