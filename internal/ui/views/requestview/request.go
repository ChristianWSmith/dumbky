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
	UI          *container.AppTabs
	QueryParams keyvalueeditorview.KeyValueEditorView
	PathParams  keyvalueeditorview.KeyValueEditorView
	Headers     keyvalueeditorview.KeyValueEditorView
	Body        requestbodyview.RequestBodyView
}

type RequestState struct {
	QueryParams keyvalueeditorview.KeyValueEditorState `json:"queryParams"`
	PathParams  keyvalueeditorview.KeyValueEditorState `json:"pathParams"`
	Headers     keyvalueeditorview.KeyValueEditorState `json:"headers"`
	Body        requestbodyview.RequestBodyState       `json:"body"`
}

func (rv RequestView) ToState() (RequestState, error) {
	queryParams, queryParamsErr := rv.QueryParams.ToState()
	if queryParamsErr != nil {
		log.Error(queryParamsErr)
		return RequestState{}, queryParamsErr
	}
	pathParams, pathParamsErr := rv.PathParams.ToState()
	if pathParamsErr != nil {
		log.Error(pathParamsErr)
		return RequestState{}, pathParamsErr
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
		QueryParams: queryParams,
		PathParams:  pathParams,
		Headers:     headers,
		Body:        body,
	}, nil
}

func (rv RequestView) LoadState(requestState RequestState) error {
	queryParamsErr := rv.QueryParams.LoadState(requestState.QueryParams)
	if queryParamsErr != nil {
		log.Error(queryParamsErr)
		return queryParamsErr
	}
	pathParamsErr := rv.PathParams.LoadState(requestState.QueryParams)
	if pathParamsErr != nil {
		log.Error(pathParamsErr)
		return pathParamsErr
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
	queryParamsView := keyvalueeditorview.ComposeKeyValueEditorView(validators.ValidateQueryParamKey, validators.ValidateQueryParamValue)
	pathParamsView := keyvalueeditorview.ComposeKeyValueEditorView(validators.ValidatePathParamKey, validators.ValidatePathParamValue)
	headersView := keyvalueeditorview.ComposeKeyValueEditorView(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyView := requestbodyview.ComposeRequestBodyView()

	queryParamsTab := container.NewTabItem(constants.UI_LABEL_QUERY_PARAMETERS, queryParamsView.UI)
	pathParamsTab := container.NewTabItem(constants.UI_LABEL_PATH_PARAMETERS, pathParamsView.UI)
	headersTab := container.NewTabItem(constants.UI_LABEL_HEADERS, headersView.UI)
	bodyTab := container.NewTabItem(constants.UI_LABEL_BODY, bodyView.UI)

	ui := container.NewAppTabs(queryParamsTab, pathParamsTab, headersTab, bodyTab)

	return RequestView{
		ui,
		queryParamsView,
		pathParamsView,
		headersView,
		bodyView,
	}
}
