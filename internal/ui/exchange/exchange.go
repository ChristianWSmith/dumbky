package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type ExchangeView struct {
	UI           *fyne.Container
	HeaderView   ExchangeHeaderView
	RequestView  RequestView
	ResponseView ResponseView
}

func headerViewSendButtonOnTappedWorker(headerView ExchangeHeaderView, requestView RequestView, responseView ResponseView) {
	headerView.URLBinding.Set("123")
}

func headerViewMethodSelectOnChangedWorker(val string, requestView RequestView) {
	if val == constants.HTTP_METHOD_GET || val == constants.HTTP_METHOD_HEAD {
		fyne.Do(func() {
			err := requestView.Body.BodyTypeBinding.Set(constants.UI_BODY_TYPE_NONE)
			if err != nil {
				log.Error("Failed to set BodyTypeBinding to %s in headerViewMethodSelectOnChanged", constants.UI_BODY_TYPE_NONE)
			}
			requestView.Body.BodyTypeSelect.Disable()
		})
	} else {
		fyne.Do(func() {
			requestView.Body.BodyTypeSelect.Enable()
		})
	}

}

func ComposeExchangeView() ExchangeView {
	headerView := ComposeExchangeHeaderView()
	requestView := ComposeRequestView()
	responseView := ComposeResponseView()

	headerView.SendButton.OnTapped = func() {
		go headerViewSendButtonOnTappedWorker(headerView, requestView, responseView)
	}

	headerView.MethodSelect.OnChanged = func(val string) {
		go headerViewMethodSelectOnChangedWorker(val, requestView)
	}

	headerView.MethodSelect.SetSelectedIndex(0)
	requestView.Body.BodyTypeSelect.SetSelectedIndex(0)

	requestResponseView := container.NewHSplit(requestView.UI, responseView.UI)
	ui := container.NewBorder(headerView.UI, nil, nil, nil, requestResponseView)

	return ExchangeView{
		ui,
		headerView,
		requestView,
		responseView,
	}
}
