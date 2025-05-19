package ui

import (
	"dumbky/internal/ui/exchange"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type RootView struct{
	UI *fyne.Container
	ExchangeView exchange.ExchangeView
}

func ComposeRootView() RootView {
	exchangeView := exchange.ComposeExchangeView()

	ui := container.NewBorder(nil, nil, nil, nil, exchangeView.UI)

	return RootView{
		ui,
		exchangeView,
	}
}