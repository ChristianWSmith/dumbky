package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"dumbky/internal/request"

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
	url, urlErr := headerView.URLBinding.Get()
	if urlErr != nil {
		log.Error("Failed to Get URLBinding")
		return
	}
	urlErr = headerView.URLEntry.Validate()
	if urlErr != nil {
		log.Error("Failed to Validate URLEntry")
		return
	}
	method, methodErr := headerView.MethodBinding.Get()
	if methodErr != nil {
		log.Error("Failed to Get MethodBinding")
		return
	}
	useSSL, useSSLErr := headerView.UseSSLBinding.Get()
	if useSSLErr != nil {
		log.Error("Failed to Get UseSSLBinding")
		return
	}
	headers := requestView.Headers.GetMap()
	headersErr := requestView.Headers.Validate()
	if headersErr != nil {
		log.Error("Failed to Get Headers Map")
		return
	}
	params := requestView.Params.GetMap()
	paramsErr := requestView.Params.Validate()
	if paramsErr != nil {
		log.Error("Failed to Get Params Map")
		return
	}

	bodyType, bodyTypeErr := requestView.Body.BodyTypeBinding.Get()
	bodyRaw, bodyRawErr := requestView.Body.BodyRawBinding.Get()
	if bodyRawErr != nil && (bodyType == constants.UI_BODY_TYPE_RAW || bodyTypeErr != nil) {
		log.Error("Failed to Get BodyRawBinding")
		return
	}
	bodyRawErr = requestView.Body.BodyRawEntry.Validate()
	if bodyRawErr != nil && (bodyType == constants.UI_BODY_TYPE_RAW || bodyTypeErr != nil) {
		log.Error("Failed to Validate BodyRawEntry")
		return
	}

	bodyForm := requestView.Body.BodyKeyValueEditor.GetMap()
	bodyFormErr := requestView.Body.BodyKeyValueEditor.Validate()
	if bodyFormErr != nil && (bodyType == constants.UI_BODY_TYPE_FORM || bodyTypeErr != nil) {
		log.Error("Failed to Validate Body Map")
		return
	}

	requestPayload := request.RequestPayload{
		URL:      url,
		Method:   method,
		UseSSL:   useSSL,
		Headers:  headers,
		Params:   params,
		BodyRaw:  bodyRaw,
		BodyForm: bodyForm,
	}

	responsePayload, err := request.SendRequest(requestPayload)
	if err != nil {
		log.Error("Failed to SendRequest")
		return
	}

	fyne.Do(func() {
		statusErr := responseView.StatusBinding.Set(responsePayload.Status)
		if statusErr != nil {
			log.Error("Failed to set StatusBinding")
		}
		timeErr := responseView.TimeBinding.Set(responsePayload.Time)
		if timeErr != nil {
			log.Error("Failed to set TimeBinding")
		}
		bodyErr := responseView.BodyBinding.Set(responsePayload.Body)
		if bodyErr != nil {
			log.Error("Failed to set BodyBinding")
		}
	})
}

func headerViewMethodSelectOnChangedWorker(val string, requestView RequestView) {
	if val == constants.HTTP_METHOD_GET || val == constants.HTTP_METHOD_HEAD {
		fyne.Do(func() {
			requestView.Body.BodyTypeSelect.SetSelected(constants.UI_BODY_TYPE_NONE)
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

	headerView.MethodSelect.SetSelectedIndex(headerView.MethodSelect.SelectedIndex())
	requestView.Body.BodyTypeSelect.SetSelectedIndex(requestView.Body.BodyTypeSelect.SelectedIndex())

	requestResponseView := container.NewHSplit(requestView.UI, responseView.UI)
	ui := container.NewBorder(headerView.UI, nil, nil, nil, requestResponseView)

	return ExchangeView{
		ui,
		headerView,
		requestView,
		responseView,
	}
}
