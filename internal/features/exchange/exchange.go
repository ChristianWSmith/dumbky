package exchange

import (
	"dumbky/internal/features"
	"dumbky/internal/features/exchange/controller"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/state"
	"dumbky/internal/validators"
)

type ExchangeController interface {
	features.Controller
	ToState() state.ExchangeState
	LoadState(exchangeState state.ExchangeState)
}

func New() ExchangeController {
	queryParamsKeyValueCtrl := keyvalueeditor.New(validators.ValidateQueryParamKey, validators.ValidateQueryParamValue)
	pathParamsKeyValueCtrl := keyvalueeditor.New(validators.ValidatePathParamKey, validators.ValidatePathParamValue)
	headersKeyValueCtrl := keyvalueeditor.New(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyFormKeyValueCtrl := keyvalueeditor.New(validators.ValidateFormBodyKey, validators.ValidateFormBodyValue)
	return controller.NewController(queryParamsKeyValueCtrl, pathParamsKeyValueCtrl, headersKeyValueCtrl, bodyFormKeyValueCtrl)
}
