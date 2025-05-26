package workspace

import "dumbky/internal/features/exchange"

type model struct {
}

func newModel() *model {
	return &model{}
}

type Document struct {
	CollectionName string                 `json:"collection_name"`
	Title          string                 `json:"title"`
	ExchangeState  exchange.ExchangeState `json:"exchange"`
}

type WorkspaceTab struct {
	CollectionName string
	Title          string
	ExchangeView   *exchange.Controller
}
