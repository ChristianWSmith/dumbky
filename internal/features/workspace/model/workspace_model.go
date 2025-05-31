package model

import (
	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/features/workspace/common"
	"dumbky/internal/log"
	"dumbky/internal/state"
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2/data/binding"
)

type modelImpl struct {
	documentDataMap             map[common.DocumentId]common.DocumentData
	requestName, collectionName binding.String
}

type Bindings struct {
	RequestName, CollectionName binding.String
}

type Model interface {
	DestroyDocumentData(id common.DocumentId)
	DocumentCount() int
	GetDocumentData(id common.DocumentId) common.DocumentData
	SetDocumentData(id common.DocumentId, data common.DocumentData)
	UpdateRequestName(id common.DocumentId, requestName string)
	SetRequestNameListener(handler func())
	GetRequestName() string
	SetRequestName(string)
	SetCollectionName(string)
	GetCollectionName() string
	GetBindings() Bindings
}

func NewModel() Model {
	return &modelImpl{
		documentDataMap: make(map[common.DocumentId]common.DocumentData),
		requestName:     binding.NewString(),
		collectionName:  binding.NewString(),
	}
}

func (m *modelImpl) DestroyDocumentData(id common.DocumentId) {
	delete(m.documentDataMap, id)
}

func (m *modelImpl) GetBindings() Bindings {
	return Bindings{
		RequestName:    m.requestName,
		CollectionName: m.collectionName,
	}
}

func (m *modelImpl) DocumentCount() int {
	return len(m.documentDataMap)
}

func (m *modelImpl) GetDocumentData(id common.DocumentId) common.DocumentData {
	return m.documentDataMap[id]
}

func (m *modelImpl) SetDocumentData(id common.DocumentId, data common.DocumentData) {
	m.documentDataMap[id] = data
}

func (m *modelImpl) UpdateRequestName(id common.DocumentId, requestName string) {
	documentData := m.documentDataMap[id]
	documentData.RequestName = requestName
	m.documentDataMap[id] = documentData
}

func (m *modelImpl) SetRequestNameListener(handler func()) {
	m.requestName.AddListener(binding.NewDataListener(handler))
}

func (m *modelImpl) GetRequestName() string {
	requestName, _ := m.requestName.Get()
	return requestName
}

func (m *modelImpl) SetRequestName(requestName string) {
	m.requestName.Set(requestName)
}

func (m *modelImpl) GetCollectionName() string {
	collectionName, _ := m.collectionName.Get()
	return collectionName
}

func (m *modelImpl) SetCollectionName(collectionName string) {
	m.collectionName.Set(collectionName)
}

func DocumentStateToRequest(documentState state.DocumentState) (db.Request, error) {
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

func RequestToDocumentState(request db.Request) (state.DocumentState, error) {
	log.Info(fmt.Sprintf("%v", request))
	document := state.DocumentState{}
	err := json.Unmarshal([]byte(request.Payload), &document)
	if err != nil {
		log.Error(err)
		return state.DocumentState{}, err
	}

	return document, nil
}
