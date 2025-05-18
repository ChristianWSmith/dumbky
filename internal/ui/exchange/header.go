package exchange

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ExchangeHeaderView struct {
	UI *fyne.Container
	SendButton *widget.Button
	Method binding.String
	URL binding.String
	UseSSL binding.Bool
}

func ComposeExchangeHeaderView() ExchangeHeaderView {
	methodBind := binding.NewString()
	urlBind := binding.NewString()
	sslBind := binding.NewBool()

	methodSelect := widget.NewSelect([]string{"GET", "PUT", "POST", "DELETE"}, nil)
	methodSelect.SetSelectedIndex(0)
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://www.example.com/")
	sslCheck := widget.NewCheck("SSL", nil)
	sendButton := widget.NewButton("SEND", nil)

	methodSelect.Bind(methodBind)
	urlEntry.Bind(urlBind)
	sslCheck.Bind(sslBind)

	sslSend := container.NewHBox(sslCheck, sendButton)
	ui := container.NewBorder(nil, nil, methodSelect, sslSend, urlEntry)

	return ExchangeHeaderView {
		ui,
		sendButton,
		methodBind,
		urlBind,
		sslBind,
	}
}