package workspace

import (
	"dumbky/internal/features/exchange"
)

type model struct {
	documentDataMap map[string]documentData
}

func newModel() *model {
	return &model{
		documentDataMap: make(map[string]documentData),
	}
}

type DocumentState struct {
	CollectionName string                 `json:"collection_name"`
	RequestName    string                 `json:"request_name"`
	ExchangeState  exchange.ExchangeState `json:"exchange"`
}

type documentData struct {
	CollectionName string
	RequestName    string
}
