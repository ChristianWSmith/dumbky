package requestview

import (
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"dumbky/internal/ui/validators"
	"dumbky/internal/ui/views/keyvalueeditorview"
	"dumbky/internal/ui/views/requestbodyview"

	"fyne.io/fyne/v2/container"
)

type RequestView struct {
	UI      *container.AppTabs
	Params  keyvalueeditorview.KeyValueEditorView
	Headers keyvalueeditorview.KeyValueEditorView
	Body    requestbodyview.RequestBodyView
}

type RequestState struct {
	Params  keyvalueeditorview.KeyValueEditorState `json:"params"`
	Headers keyvalueeditorview.KeyValueEditorState `json:"headers"`
	Body    requestbodyview.RequestBodyState       `json:"body"`
}

func (rv RequestView) ToState() (RequestState, error) {
	params, paramsErr := rv.Params.ToState()
	if paramsErr != nil {
		log.Error(paramsErr)
		return RequestState{}, paramsErr
	}
	headers, headersErr := rv.Headers.ToState()
	if headersErr != nil {
		log.Error(headersErr)
		return RequestState{}, headersErr
	}
	body, bodyErr := rv.Body.ToState()
	if bodyErr != nil {
		log.Error(bodyErr)
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
		log.Error(paramsErr)
		return paramsErr
	}
	headersErr := rv.Headers.LoadState(requestState.Headers)
	if headersErr != nil {
		log.Error(headersErr)
		return headersErr
	}
	bodyErr := rv.Body.LoadState(requestState.Body)
	if bodyErr != nil {
		log.Error(bodyErr)
		return bodyErr
	}
	return nil
}

func ComposeRequestView() RequestView {
	paramsView := keyvalueeditorview.ComposeKeyValueEditorView(validators.ValidateURLParamKey, validators.ValidateURLParamValue)
	headersView := keyvalueeditorview.ComposeKeyValueEditorView(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyView := requestbodyview.ComposeRequestBodyView()

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
