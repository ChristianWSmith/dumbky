package workspace

import (
	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/log"
	"dumbky/internal/state"
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2/data/binding"
)

type model struct {
	documentDataMap    map[documentId]documentData
	requestNameBinding binding.String
}

func newModel() *model {
	return &model{
		documentDataMap:    make(map[documentId]documentData),
		requestNameBinding: binding.NewString(),
	}
}

type DocumentState struct {
	CollectionName string              `json:"collection_name"`
	RequestName    string              `json:"request_name"`
	ExchangeState  state.ExchangeState `json:"exchange"`
}

type documentData struct {
	collectionName string
	requestName    string
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
	documentData.requestName = requestName
	m.documentDataMap[id] = documentData
}

func documentStateToRequest(documentState DocumentState) (db.Request, error) {
	if documentState.RequestName == "" {
		documentState.RequestName = constants.UI_PLACEHOLDER_UNTITLED
	}

	jsonData, err := json.Marshal(documentState)
	if err != nil {
		log.Error(err)
		return db.Request{}, err
	}

	jsonString := string(jsonData)

	return db.Request{
		CollectionName: documentState.CollectionName,
		Name:           documentState.RequestName,
		Payload:        jsonString,
	}, nil
}

func requestToDocumentState(request db.Request) (DocumentState, error) {
	log.Info(fmt.Sprintf("%v", request))
	document := DocumentState{}
	err := json.Unmarshal([]byte(request.Payload), &document)
	if err != nil {
		log.Error(err)
		return DocumentState{}, err
	}

	return document, nil
}

func (m *model) setRequestNameListener(handler func()) {
	m.requestNameBinding.AddListener(binding.NewDataListener(handler))
}

func (m *model) getRequestName() string {
	requestName, _ := m.requestNameBinding.Get()
	return requestName
}

func (m *model) setRequestName(requestName string) {
	m.requestNameBinding.Set(requestName)
}
