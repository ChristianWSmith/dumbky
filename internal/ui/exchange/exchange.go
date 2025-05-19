package exchange

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type ExchangeView struct {
	UI *fyne.Container
	HeaderView ExchangeHeaderView
	RequestView RequestView
	ResponseView ResponseView
}

func ComposeExchangeView() ExchangeView {
	headerView := ComposeExchangeHeaderView()
	requestView := ComposeRequestView()
	responseView := ComposeResponseView()

	headerView.SendButton.OnTapped = func () {
		// TODO: send the request
	}

	headerView.MethodSelect.OnChanged = func (val string) {
		// TODO: GET and HEAD don't support a body
	}

	requestResponseView := container.NewHSplit(requestView.UI, responseView.UI)
	ui := container.NewBorder(headerView.UI, nil, nil, nil, requestResponseView)

	return ExchangeView {
		ui,
		headerView,
		requestView,
		responseView,
	}
}