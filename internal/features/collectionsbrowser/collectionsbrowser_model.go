package collectionsbrowser

import "fyne.io/fyne/v2/data/binding"

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
