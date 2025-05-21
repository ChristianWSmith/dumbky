package rootview

import (
	"dumbky/internal/ui/views/exchangeview"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type RootView struct {
	UI           *fyne.Container
	ExchangeView exchangeview.ExchangeView
}

func ComposeRootView() RootView {
	exchangeView := exchangeview.ComposeExchangeView()

	ui := container.NewBorder(nil, nil, nil, nil, exchangeView.UI)

	return RootView{
		ui,
		exchangeView,
	}
}
