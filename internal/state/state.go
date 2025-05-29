package state

import "dumbky/internal/features/request"

type ExchangeState struct {
	Method  string               `json:"method"`
	URL     string               `json:"url"`
	UseSSL  bool                 `json:"ssl"`
	Request request.RequestState `json:"request"`
}
