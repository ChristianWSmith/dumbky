package model

import (
	"fyne.io/fyne/v2/data/binding"
)

type modelImpl struct {
	selectedCollectionBinding binding.String
	addCollectionBinding      binding.String
	collectionsListBinding    binding.StringList
	requestsListBinding       binding.StringList
}

type Bindings struct {
	SelectedCollection, AddCollection binding.String
	CollectionsList, RequestsList     binding.StringList
}

type Model interface {
	GetBindings() Bindings
	GetSelectedCollection() string
	SetSelectedCollection(collectionName string)
	SetRequestsList(requestsList []string)
	SetCollectionsList(collectionsList []string)
	GetCollectionNameById(id int) string
	GetRequestNameById(id int) string
	GetAddCollection() string
	SetAddCollection(collectionName string)
}

func NewModel() Model {
	return &modelImpl{
		selectedCollectionBinding: binding.NewString(),
		addCollectionBinding:      binding.NewString(),
		collectionsListBinding:    binding.NewStringList(),
		requestsListBinding:       binding.NewStringList(),
	}
}

func (m *modelImpl) GetBindings() Bindings {
	return Bindings{
		SelectedCollection: m.selectedCollectionBinding,
		AddCollection:      m.addCollectionBinding,
		CollectionsList:    m.collectionsListBinding,
		RequestsList:       m.requestsListBinding,
	}
}

func (m *modelImpl) GetSelectedCollection() string {
	selectedCollection, _ := m.selectedCollectionBinding.Get()
	return selectedCollection
}

func (m *modelImpl) SetSelectedCollection(collectionName string) {
	m.selectedCollectionBinding.Set(collectionName)
}

func (m *modelImpl) SetRequestsList(requestsList []string) {
	m.requestsListBinding.Set(requestsList)
}

func (m *modelImpl) SetCollectionsList(collectionsList []string) {
	m.collectionsListBinding.Set(collectionsList)
}

func (m *modelImpl) GetCollectionNameById(id int) string {
	collectionName, _ := m.collectionsListBinding.GetValue(id)
	return collectionName
}

func (m *modelImpl) GetRequestNameById(id int) string {
	requestName, _ := m.requestsListBinding.GetValue(id)
	return requestName
}

func (m *modelImpl) GetAddCollection() string {
	collectionName, _ := m.addCollectionBinding.Get()
	return collectionName
}

func (m *modelImpl) SetAddCollection(collectionName string) {
	m.addCollectionBinding.Set(collectionName)
}
