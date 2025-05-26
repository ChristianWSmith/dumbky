package requestview

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/features/requestbodyview"
	"dumbky/internal/log"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type RequestView struct {
	UI          *fyne.Container
	QueryParams *keyvalueeditor.Controller
	PathParams  *keyvalueeditor.Controller
	Headers     *keyvalueeditor.Controller
	Body        requestbodyview.RequestBodyView
}

type RequestState struct {
	QueryParams keyvalueeditor.KeyValueEditorState `json:"queryParams"`
	PathParams  keyvalueeditor.KeyValueEditorState `json:"pathParams"`
	Headers     keyvalueeditor.KeyValueEditorState `json:"headers"`
	Body        requestbodyview.RequestBodyState   `json:"body"`
}

func (rv RequestView) ToState() (RequestState, error) {
	queryParams := rv.QueryParams.ToState()
	pathParams := rv.PathParams.ToState()
	headers := rv.Headers.ToState()
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
	pathParamsErr := rv.PathParams.LoadState(requestState.PathParams)
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
	queryParamsView := keyvalueeditor.NewController(validators.ValidateQueryParamKey, validators.ValidateQueryParamValue)
	pathParamsView := keyvalueeditor.NewController(validators.ValidatePathParamKey, validators.ValidatePathParamValue)
	headersView := keyvalueeditor.NewController(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyView := requestbodyview.ComposeRequestBodyView()

	queryParamsTab := container.NewTabItem(constants.UI_LABEL_QUERY_PARAMETERS, queryParamsView.GetUI())
	pathParamsTab := container.NewTabItem(constants.UI_LABEL_PATH_PARAMETERS, pathParamsView.GetUI())
	headersTab := container.NewTabItem(constants.UI_LABEL_HEADERS, headersView.GetUI())
	bodyTab := container.NewTabItem(constants.UI_LABEL_BODY, bodyView.UI)

	tabs := container.NewAppTabs(queryParamsTab, pathParamsTab, headersTab, bodyTab)
	ui := container.NewBorder(nil, nil, nil, nil, tabs)

	return RequestView{
		ui,
		queryParamsView,
		pathParamsView,
		headersView,
		bodyView,
	}
}
