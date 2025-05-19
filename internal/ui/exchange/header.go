package exchange

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ExchangeHeaderView struct {
	UI *fyne.Container
	SendButton *widget.Button
	MethodSelect *widget.Select
	Method binding.String
	URL binding.String
	UseSSL binding.Bool
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

	methodSelect.Bind(methodBind)
	urlEntry.Bind(urlBind)
	sslCheck.Bind(sslBind)

	sslSend := container.NewHBox(sslCheck, sendButton)
	ui := container.NewBorder(nil, nil, methodSelect, sslSend, urlEntry)

	return ExchangeHeaderView {
		ui,
		sendButton,
		methodSelect,
		methodBind,
		urlBind,
		sslBind,
	}
}