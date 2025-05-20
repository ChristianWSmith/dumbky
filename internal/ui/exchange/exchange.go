package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"dumbky/internal/request"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type ExchangeView struct {
	UI           *fyne.Container
	HeaderView   ExchangeHeaderView
	RequestView  RequestView
	ResponseView ResponseView
}

func headerViewSendButtonOnTappedWorker(requestView RequestView, responseView ResponseView, requestPayload request.RequestPayload, cleanupFunc func()) {
	defer fyne.Do(cleanupFunc)

	fyne.Do(func() {
		bodyType, bodyTypeGetErr := requestView.Body.BodyTypeBinding.Get()
		if bodyTypeGetErr != nil {
			log.Error("Failed to Get BodyTypeBinding during autoformat", bodyTypeGetErr.Error())
		}
		if bodyType != constants.UI_BODY_TYPE_RAW && bodyTypeGetErr == nil {
			log.Debug("Skipping autoformat for non-raw bodyType")
			return
		}
		bodyRaw, bodyRawGetErr := requestView.Body.BodyRawBinding.Get()
		if bodyRawGetErr != nil {
			log.Error("Failed to Get BodyRawBinding during autoformat", bodyRawGetErr.Error())
			return
		}
		bodyRawSetErr := requestView.Body.BodyRawBinding.Set(utils.SmartFormat(bodyRaw))
		if bodyRawSetErr != nil {
			log.Error("Failed to Set BodyRawBinding during autoformat", bodyRawSetErr.Error())
			return
		}
	})

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

		bodyErr := responseView.BodyBinding.Set(utils.SmartFormat(responsePayload.Body))
		if bodyErr != nil {
			log.Error("Failed to set BodyBinding", bodyErr.Error())
		}
	})
}

func headerViewSendButtonOnTapped(headerView ExchangeHeaderView, requestView RequestView, responseView ResponseView) (marshalled bool) {
	cleanupFunc := func() {
		responseView.SetLoading(false)
		headerView.SendButton.Enable()
	}

	defer func() {
		if !marshalled {
			cleanupFunc()
		}
	}()

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

	url, urlGetErr := headerView.URLBinding.Get()
	if urlGetErr != nil {
		log.Error("Failed to Get URLBinding", urlGetErr.Error())
		return false
	}
	urlValidateErr := headerView.URLEntry.Validate()
	if urlValidateErr != nil {
		log.Warn("Failed to Validate URLEntry", urlValidateErr.Error())
		return false
	}

	method, methodGetErr := headerView.MethodBinding.Get()
	if methodGetErr != nil {
		log.Error("Failed to Get MethodBinding", methodGetErr.Error())
		return false
	}

	useSSL, useSSLGetErr := headerView.UseSSLBinding.Get()
	if useSSLGetErr != nil {
		log.Error("Failed to Get UseSSLBinding", useSSLGetErr.Error())
		return false
	}

	headers, headersGetErr := requestView.Headers.GetMap()
	if headersGetErr != nil {
		log.Error("Failed to Get Headers Map", headersGetErr.Error())
		return false
	}
	headersValidatErr := requestView.Headers.Validate()
	if headersValidatErr != nil {
		log.Error("Failed to Validate Headers Map", headersValidatErr.Error())
		return false
	}

	params, paramsGetErr := requestView.Params.GetMap()
	if paramsGetErr != nil {
		log.Error("Failed to Get Params Map", paramsGetErr.Error())
		return false
	}
	paramsValidatErr := requestView.Params.Validate()
	if paramsValidatErr != nil {
		log.Error("Failed to Validate Params Map", paramsValidatErr.Error())
		return false
	}

	bodyType, bodyTypeGetErr := requestView.Body.BodyTypeBinding.Get()
	if bodyTypeGetErr != nil {
		log.Error("Failed to Get BodyTypeBinding", bodyTypeGetErr.Error())
		return false
	}

	bodyRaw, bodyRawGetErr := requestView.Body.BodyRawBinding.Get()
	if bodyRawGetErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Error("Failed to Get BodyRawBinding", bodyRawGetErr.Error())
		return false
	}
	bodyRawValidateErr := requestView.Body.BodyRawEntry.Validate()
	if bodyRawValidateErr != nil && bodyType == constants.UI_BODY_TYPE_RAW {
		log.Warn("Failed to Validate BodyRawEntry", bodyRawValidateErr.Error())
		return false
	}

	bodyForm, bodyFormGetErr := requestView.Body.BodyKeyValueEditor.GetMap()
	if bodyFormGetErr != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn("Failed to Get Body Map", bodyFormGetErr.Error())
		return false
	}
	bodyFormValidateErr := requestView.Body.BodyKeyValueEditor.Validate()
	if bodyFormValidateErr != nil && (bodyType == constants.UI_BODY_TYPE_FORM) {
		log.Warn("Failed to Validate Body Map", bodyFormValidateErr.Error())
		return false
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

	go headerViewSendButtonOnTappedWorker(requestView, responseView, requestPayload, cleanupFunc)
	return true
}

func headerViewMethodSelectOnChanged(val string, requestView RequestView) {
	if val == constants.HTTP_METHOD_GET || val == constants.HTTP_METHOD_HEAD {
		requestView.Body.BodyTypeSelect.SetSelected(constants.UI_BODY_TYPE_NONE)
		requestView.Body.BodyTypeSelect.Disable()
	} else {
		requestView.Body.BodyTypeSelect.Enable()
	}

}

func ComposeExchangeView() ExchangeView {
	headerView := ComposeExchangeHeaderView()
	requestView := ComposeRequestView()
	responseView := ComposeResponseView()

	headerView.SendButton.OnTapped = func() {
		marshalled := headerViewSendButtonOnTapped(headerView, requestView, responseView)
		if !marshalled {
			log.Warn("Failed to marshal request")
		}
	}

	headerViewMethodSelectOnChangedOld := headerView.MethodSelect.OnChanged
	headerView.MethodSelect.OnChanged = func(val string) {
		headerViewMethodSelectOnChangedOld(val)
		headerViewMethodSelectOnChanged(val, requestView)
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
