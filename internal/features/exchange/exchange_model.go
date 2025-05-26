package exchange

import (
	"dumbky/internal/features/exchangeheader"
	"dumbky/internal/features/request"
)

type model struct {
}

func newModel() *model {
	return &model{}
}

type ExchangeState struct {
	Header  exchangeheader.ExchangeHeaderState `json:"header"`
	Request request.RequestState               `json:"request"`
}
