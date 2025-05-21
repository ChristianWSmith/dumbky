package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"dumbky/internal/ui/common"
	"dumbky/internal/ui/validators"

	"fyne.io/fyne/v2/container"
)

type RequestView struct {
	UI      *container.AppTabs
	Params  common.KeyValueEditorView
	Headers common.KeyValueEditorView
	Body    RequestBodyView
}

type RequestState struct {
	Params  common.KeyValueEditorState `json:"params"`
	Headers common.KeyValueEditorState `json:"headers"`
	Body    RequestBodyState           `json:"body"`
}

func (rv RequestView) ToState() (RequestState, error) {
	params, paramsErr := rv.Params.ToState()
	if paramsErr != nil {
		log.Error("", paramsErr.Error())
		return RequestState{}, paramsErr
	}
	headers, headersErr := rv.Headers.ToState()
	if headersErr != nil {
		log.Error("", headersErr.Error())
		return RequestState{}, headersErr
	}
	body, bodyErr := rv.Body.ToState()
	if bodyErr != nil {
		log.Error("", bodyErr.Error())
		return RequestState{}, bodyErr
	}
	return RequestState{
		Params:  params,
		Headers: headers,
		Body:    body,
	}, nil
}

func (rv RequestView) LoadState(requestState RequestState) error {
	paramsErr := rv.Params.LoadState(requestState.Params)
	if paramsErr != nil {
		log.Error("", paramsErr.Error())
		return paramsErr
	}
	headersErr := rv.Headers.LoadState(requestState.Headers)
	if headersErr != nil {
		log.Error("", headersErr.Error())
		return headersErr
	}
	bodyErr := rv.Body.LoadState(requestState.Body)
	if bodyErr != nil {
		log.Error("", bodyErr.Error())
		return bodyErr
	}
	return nil
}

func ComposeRequestView() RequestView {
	paramsView := common.ComposeKeyValueEditorView(validators.ValidateURLParamKey, validators.ValidateURLParamValue)
	headersView := common.ComposeKeyValueEditorView(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyView := ComposeRequestBodyView()

	paramsTab := container.NewTabItem(constants.UI_LABEL_PARAMETERS, paramsView.UI)
	headersTab := container.NewTabItem(constants.UI_LABEL_HEADERS, headersView.UI)
	bodyTab := container.NewTabItem(constants.UI_LABEL_BODY, bodyView.UI)

	ui := container.NewAppTabs(paramsTab, headersTab, bodyTab)

	return RequestView{
		ui,
		paramsView,
		headersView,
		bodyView,
	}
}
