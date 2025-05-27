package collectionsbrowser

import (
	"fyne.io/fyne/v2/data/binding"
)

type model struct {
	selectedCollectionBinding binding.String
	selectedRequestBinding    binding.String
	addCollectionBinding      binding.String
	collectionsListBinding    binding.StringList
	requestsListBinding       binding.StringList
}

func newModel() *model {
	return &model{
		selectedCollectionBinding: binding.NewString(),
		selectedRequestBinding:    binding.NewString(),
		addCollectionBinding:      binding.NewString(),
		collectionsListBinding:    binding.NewStringList(),
		requestsListBinding:       binding.NewStringList(),
	}
}

func (m *model) getSelectedCollection() string {
	selectedCollection, _ := m.selectedCollectionBinding.Get()
	return selectedCollection
}

func (m *model) getSelectedRequest() string {
	selectedRequest, _ := m.selectedRequestBinding.Get()
	return selectedRequest
}

func (m *model) setSelectedCollection(collectionName string) {
	m.selectedCollectionBinding.Set(collectionName)
}

func (m *model) setSelectedRequest(requestName string) {
	m.selectedRequestBinding.Set(requestName)
}

func (m *model) setSelectedRequestListener(handler func()) {
	m.selectedRequestBinding.AddListener(binding.NewDataListener(handler))
}

func (m *model) setRequestsList(requestsList []string) {
	m.requestsListBinding.Set(requestsList)
}

func (m *model) setCollectionsList(collectionsList []string) {
	m.collectionsListBinding.Set(collectionsList)
}

func (m *model) getCollectionNameById(id int) string {
	collectionName, _ := m.collectionsListBinding.GetValue(id)
	return collectionName
}

func (m *model) getRequestNameById(id int) string {
	requestName, _ := m.requestsListBinding.GetValue(id)
	return requestName
}

func (m *model) getAddCollection() string {
	collectionName, _ := m.addCollectionBinding.Get()
	return collectionName
}

func (m *model) setAddCollection(collectionName string) {
	m.addCollectionBinding.Set(collectionName)
}
