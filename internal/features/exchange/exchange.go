package exchange

import (
	"dumbky/internal/features"
	"dumbky/internal/features/exchange/controller"
	"dumbky/internal/features/request"
	"dumbky/internal/features/response"
	"dumbky/internal/state"
)

type ExchangeController interface {
	features.Controller
	ToState() state.ExchangeState
	LoadState(exchangeState state.ExchangeState)
	Validate() error
}

func New() ExchangeController {
	requestCtrl := request.New()
	responseCtrl := response.New()
	return controller.NewController(requestCtrl, responseCtrl)
}
