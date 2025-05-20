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
	defer fyne.Do(func() {
		responseView.SetLoading(false)
	})

	fyne.Do(func() {
		headerView.SendButton.Disable()
		responseView.SetLoading(true)

		statusErr := responseView.StatusBinding.Set("")
		if statusErr != nil {
			log.Error("Failed to set StatusBinding", statusErr.Error())
		}
		timeErr := responseView.TimeBinding.Set("")
		if timeErr != nil {
			log.Error("Failed to set TimeBinding", timeErr.Error())
		}
		bodyErr := responseView.BodyBinding.Set("")
		if bodyErr != nil {
			log.Error("Failed to set BodyBinding", bodyErr.Error())
		}
	})

	defer fyne.Do(func() {
		headerView.SendButton.Enable()
	})

	url, urlErr := headerView.URLBinding.Get()
	if urlErr != nil {
		log.Error("Failed to Get URLBinding", urlErr.Error())
		return
	}
	urlErr = headerView.URLEntry.Validate()
	if urlErr != nil {
		log.Warn("Failed to Validate URLEntry", urlErr.Error())
		return
	}
	method, methodErr := headerView.MethodBinding.Get()
	if methodErr != nil {
		log.Error("Failed to Get MethodBinding", methodErr.Error())
		return
	}
	useSSL, useSSLErr := headerView.UseSSLBinding.Get()
	if useSSLErr != nil {
		log.Error("Failed to Get UseSSLBinding", useSSLErr.Error())
		return
	}
	headers := requestView.Headers.GetMap()
	headersErr := requestView.Headers.Validate()
	if headersErr != nil {
		log.Error("Failed to Get Headers Map", headersErr.Error())
		return
	}
	params := requestView.Params.GetMap()
	paramsErr := requestView.Params.Validate()
	if paramsErr != nil {
		log.Error("Failed to Get Params Map", paramsErr.Error())
		return
	}

	bodyType, bodyTypeErr := requestView.Body.BodyTypeBinding.Get()
	if bodyTypeErr != nil {
		log.Error("Failed to Get BodyTypeBinding", bodyTypeErr.Error())
		return
	}
	bodyRaw, bodyRawErr := requestView.Body.BodyRawBinding.Get()
	if bodyRawErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Error("Failed to Get BodyRawBinding", bodyRawErr.Error())
		return
	}
	bodyRawErr = requestView.Body.BodyRawEntry.Validate()
	if bodyRawErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Warn("Failed to Validate BodyRawEntry", bodyRawErr.Error())
		return
	}

	bodyForm := requestView.Body.BodyKeyValueEditor.GetMap()
	bodyFormErr := requestView.Body.BodyKeyValueEditor.Validate()
	if bodyFormErr != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn("Failed to Validate Body Map", bodyFormErr.Error())
		return
	}

	requestPayload := request.RequestPayload{
		URL:      url,
		Method:   method,
		UseSSL:   useSSL,
		Headers:  headers,
		Params:   params,
		BodyType: bodyType,
		BodyRaw:  bodyRaw,
		BodyForm: bodyForm,
	}

	responsePayload, err := request.SendRequest(requestPayload)
	if err != nil {
		log.Error("Failed to SendRequest", err.Error())
		return
	}

	fyne.Do(func() {
		statusErr := responseView.StatusBinding.Set(responsePayload.Status)
		if statusErr != nil {
			log.Error("Failed to set StatusBinding", statusErr.Error())
		}
		timeErr := responseView.TimeBinding.Set(responsePayload.Time)
		if timeErr != nil {
			log.Error("Failed to set TimeBinding", timeErr.Error())
		}
		bodyErr := responseView.BodyBinding.Set(responsePayload.Body)
		if bodyErr != nil {
			log.Error("Failed to set BodyBinding", bodyErr.Error())
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

	headerViewMethodSelectOnChangedOld := headerView.MethodSelect.OnChanged
	headerView.MethodSelect.OnChanged = func(val string) {
		headerViewMethodSelectOnChangedOld(val)
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
