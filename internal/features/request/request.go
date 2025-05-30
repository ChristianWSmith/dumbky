package request

import (
	"dumbky/internal/features"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/features/request/controller"
	"dumbky/internal/state"
	"dumbky/internal/validators"
)

type RequestController interface {
	features.Controller
	GetQueryParamsMap() map[string]string
	GetPathParamsMap() map[string]string
	GetHeadersMap() map[string]string
	ToState() state.RequestState
	LoadState(requestState state.RequestState)
	Validate() error
	SetBodyTypeSelectEnabled(enabled bool)
	FormatBodyRaw()
	GetBodyFormMap() map[string]string
	GetBodyType() string
	GetBodyRaw() string
}

func New() RequestController {
	queryParamsKeyValueCtrl := keyvalueeditor.New(validators.ValidateQueryParamKey, validators.ValidateQueryParamValue)
	pathParamsKeyValueCtrl := keyvalueeditor.New(validators.ValidatePathParamKey, validators.ValidatePathParamValue)
	headersKeyValueCtrl := keyvalueeditor.New(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyFormKeyValueCtrl := keyvalueeditor.New(validators.ValidateFormBodyKey, validators.ValidateFormBodyValue)
	return controller.NewController(queryParamsKeyValueCtrl, pathParamsKeyValueCtrl, headersKeyValueCtrl, bodyFormKeyValueCtrl)
}
