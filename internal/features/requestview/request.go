package requestview

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/features/requestbody"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type RequestView struct {
	UI          *fyne.Container
	QueryParams *keyvalueeditor.Controller
	PathParams  *keyvalueeditor.Controller
	Headers     *keyvalueeditor.Controller
	Body        *requestbody.Controller
}

type RequestState struct {
	QueryParams keyvalueeditor.KeyValueEditorState `json:"queryParams"`
	PathParams  keyvalueeditor.KeyValueEditorState `json:"pathParams"`
	Headers     keyvalueeditor.KeyValueEditorState `json:"headers"`
	Body        requestbody.RequestBodyState       `json:"body"`
}

func (rv RequestView) ToState() (RequestState, error) {
	queryParams := rv.QueryParams.ToState()
	pathParams := rv.PathParams.ToState()
	headers := rv.Headers.ToState()
	body := rv.Body.ToState()
	return RequestState{
		QueryParams: queryParams,
		PathParams:  pathParams,
		Headers:     headers,
		Body:        body,
	}, nil
}

func (rv RequestView) LoadState(requestState RequestState) error {
	rv.QueryParams.LoadState(requestState.QueryParams)
	rv.PathParams.LoadState(requestState.PathParams)
	rv.Headers.LoadState(requestState.Headers)
	rv.Body.LoadState(requestState.Body)
	return nil
}

func ComposeRequestView() RequestView {
	queryParamsView := keyvalueeditor.NewController(validators.ValidateQueryParamKey, validators.ValidateQueryParamValue)
	pathParamsView := keyvalueeditor.NewController(validators.ValidatePathParamKey, validators.ValidatePathParamValue)
	headersView := keyvalueeditor.NewController(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyView := requestbody.NewController()

	queryParamsTab := container.NewTabItem(constants.UI_LABEL_QUERY_PARAMETERS, queryParamsView.GetUI())
	pathParamsTab := container.NewTabItem(constants.UI_LABEL_PATH_PARAMETERS, pathParamsView.GetUI())
	headersTab := container.NewTabItem(constants.UI_LABEL_HEADERS, headersView.GetUI())
	bodyTab := container.NewTabItem(constants.UI_LABEL_BODY, bodyView.GetUI())

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
