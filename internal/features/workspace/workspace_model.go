package workspace

import (
	"dumbky/internal/features/exchange"
)

type model struct {
	documentDataMap map[documentId]documentData
}

func newModel() *model {
	return &model{
		documentDataMap: make(map[documentId]documentData),
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

func (m *model) destroyDocumentData(id documentId) {
	delete(m.documentDataMap, id)
}

func (m *model) documentCount() int {
	return len(m.documentDataMap)
}

func (m *model) getDocumentData(id documentId) documentData {
	return m.documentDataMap[id]
}

func (m *model) setDocumentData(id documentId, data documentData) {
	m.documentDataMap[id] = data
}

func (m *model) updateRequestName(id documentId, requestName string) {
	documentData := m.documentDataMap[id]
	documentData.RequestName = requestName
	m.documentDataMap[id] = documentData
}
