package exchangeheaderview

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
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
	MethodBinding binding.String
	URLBinding    binding.String
	UseSSLBinding binding.Bool

	urlEntry *widget.Entry
}

type ExchangeHeaderState struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	UseSSL bool   `json:"ssl"`
}

func (ehv ExchangeHeaderView) ToState() (ExchangeHeaderState, error) {
	method, methodErr := ehv.MethodBinding.Get()
	if methodErr != nil {
		log.Error(methodErr)
		return ExchangeHeaderState{}, methodErr
	}
	url, urlErr := ehv.URLBinding.Get()
	if urlErr != nil {
		log.Error(urlErr)
		return ExchangeHeaderState{}, urlErr
	}
	useSSL, useSSLErr := ehv.UseSSLBinding.Get()
	if useSSLErr != nil {
		log.Error(useSSLErr)
		return ExchangeHeaderState{}, useSSLErr
	}
	return ExchangeHeaderState{
		Method: method,
		URL:    url,
		UseSSL: useSSL,
	}, nil
}

func (ehv ExchangeHeaderView) LoadState(exchangeHeaderState ExchangeHeaderState) error {
	methodErr := ehv.MethodBinding.Set(exchangeHeaderState.Method)
	if methodErr != nil {
		log.Error(methodErr)
		return methodErr
	}
	urlErr := ehv.URLBinding.Set(exchangeHeaderState.URL)
	if urlErr != nil {
		log.Error(urlErr)
		return urlErr
	}
	useSSLErr := ehv.UseSSLBinding.Set(exchangeHeaderState.UseSSL)
	if useSSLErr != nil {
		log.Error(useSSLErr)
		return useSSLErr
	}
	return nil
}

func (ehv ExchangeHeaderView) ValidateURL() error {
	return ehv.urlEntry.Validate()
}

func ComposeExchangeHeaderView() ExchangeHeaderView {
	methodBind := binding.NewString()
	urlBind := binding.NewString()
	sslBind := binding.NewBool()

	methodSelect := widget.NewSelect(constants.HttpMethods(), nil)
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder(constants.UI_PLACEHOLDER_URL)
	urlEntry.TextStyle.Monospace = true
	sslCheck := widget.NewCheck(constants.UI_LABEL_SSL, nil)
	sendButton := widget.NewButton(constants.UI_LABEL_SEND, nil)
	sendButton.Icon = sendButton.Theme().Icon(theme.IconNameMailSend)

	methodSelect.Bind(methodBind)
	urlEntry.Bind(urlBind)
	sslCheck.Bind(sslBind)

	methodBind.Set(constants.HTTP_METHOD_GET)
	urlEntry.Validator = validators.ValidateURL

	sslSend := container.NewHBox(sslCheck, sendButton)
	ui := container.NewBorder(nil, nil, methodSelect, sslSend, urlEntry)

	return ExchangeHeaderView{
		ui,
		sendButton,
		methodBind,
		urlBind,
		sslBind,
		urlEntry,
	}
}
