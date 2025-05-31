package view

import (
	"dumbky/internal/components"
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type viewImpl struct {
	ui                      *fyne.Container
	addCollectionEntry      *widget.Entry
	collectionsList         *widget.List
	requestsList            *widget.List
	collectionsContainer    *fyne.Container
	requestsContainer       *fyne.Container
	addButton               *widget.Button
	backButton              *widget.Button
	selectedCollectionLabel *components.FancyLabel
}

type Bindables struct {
	AddCollection, SelectedCollection features.StringBindable
}

type View interface {
	CanvasObject() fyne.CanvasObject
	GetBindables() Bindables
	SetBackHandler(handler func())
	SetRequestSelectedCallback(callback func(int))
	SetCollectionSelectedCallback(callback func(int))
	ShowRequests()
	ShowCollections()
	ShowingRequests() bool
	ShowingCollections() bool
	SetAddCollectionValidator(validator func(string) error)
	ValidateAddCollection() error
	SetAddCollectionHandler(handler func())
}

type ViewCallbacks struct {
	OnDeleteCollection, OnDeleteRequest func(string)
}

func NewView(requestsListBinding, collectionsListBinding binding.StringList, callbacks ViewCallbacks) View {
	addCollectionEntry := widget.NewEntry()

	requestsMenu := func(position fyne.Position, parent fyne.CanvasObject, requestName string) {
		pop := fyne.NewMenu("",
			fyne.NewMenuItem(constants.UI_LABEL_DELETE, func() {
				callbacks.OnDeleteRequest(requestName)
			}),
		)
		widget.ShowPopUpMenuAtRelativePosition(pop, global.Window.Canvas(), position, parent)
	}

	collectionsMenu := func(position fyne.Position, parent fyne.CanvasObject, collectionName string) {
		pop := fyne.NewMenu("",
			fyne.NewMenuItem(constants.UI_LABEL_DELETE, func() {
				callbacks.OnDeleteCollection(collectionName)
			}),
		)
		widget.ShowPopUpMenuAtRelativePosition(pop, global.Window.Canvas(), position, parent)
	}

	requestsList := widget.NewListWithData(
		requestsListBinding,
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			menuButton := widget.NewButtonWithIcon("", theme.MoreVerticalIcon(), nil)
			return container.NewBorder(nil, nil, nil, menuButton, label)
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			requestName, _ := item.(binding.String).Get()
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			menuButton := c.Objects[1].(*widget.Button)
			label.SetText(requestName)

			menuButton.OnTapped = func() {
				requestsMenu(menuButton.Position(), o, requestName)
			}
		},
	)

	collectionsList := widget.NewListWithData(
		collectionsListBinding,
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			menuBtn := widget.NewButtonWithIcon("", theme.MoreVerticalIcon(), nil)
			return container.NewBorder(nil, nil, nil, menuBtn, label)
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			collectionName, _ := item.(binding.String).Get()
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			menuBtn := c.Objects[1].(*widget.Button)
			label.SetText(collectionName)

			menuBtn.OnTapped = func() {
				collectionsMenu(menuBtn.Position(), o, collectionName)
			}
		},
	)

	backButton := widget.NewButtonWithIcon(constants.UI_LABEL_BACK, theme.NavigateBackIcon(), nil)
	selectedCollectionLabel := components.NewFancyLabel("")
	selectedCollectionLabel.SetTextStyle(fyne.TextStyle{
		Bold: true,
	})
	selectedCollectionLabel.SetBackgroundColor(theme.Color(theme.ColorNameFocus))
	selectedCollectionLabel.Refresh()
	requestsContainer := container.NewBorder(selectedCollectionLabel, backButton, nil, nil, requestsList)

	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)
	addCollectionContainer := container.NewBorder(nil, nil, nil, addButton, addCollectionEntry)
	collectionLabel := components.NewFancyLabel(constants.UI_LABEL_COLLECTIONS)
	collectionLabel.SetTextStyle(fyne.TextStyle{
		Bold: true,
	})
	collectionLabel.SetBackgroundColor(theme.Color(theme.ColorNameFocus))
	collectionLabel.Refresh()
	collectionsContainer := container.NewBorder(collectionLabel, addCollectionContainer, nil, nil, collectionsList)

	stack := container.NewStack(collectionsContainer, requestsContainer)

	ui := container.NewBorder(nil, nil, nil, nil, stack)

	return &viewImpl{
		ui:                      ui,
		addCollectionEntry:      addCollectionEntry,
		collectionsList:         collectionsList,
		requestsList:            requestsList,
		collectionsContainer:    collectionsContainer,
		requestsContainer:       requestsContainer,
		addButton:               addButton,
		backButton:              backButton,
		selectedCollectionLabel: selectedCollectionLabel,
	}
}

func (v *viewImpl) CanvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *viewImpl) GetBindables() Bindables {
	return Bindables{
		AddCollection:      v.addCollectionEntry,
		SelectedCollection: v.selectedCollectionLabel,
	}
}

func (v *viewImpl) SetBackHandler(handler func()) {
	v.backButton.OnTapped = handler
}

func (v *viewImpl) SetRequestSelectedCallback(callback func(int)) {
	v.requestsList.OnSelected = func(id widget.ListItemID) {
		v.requestsList.UnselectAll()
		callback(id)
	}
}

func (v *viewImpl) SetCollectionSelectedCallback(callback func(int)) {
	v.collectionsList.OnSelected = func(id widget.ListItemID) {
		v.collectionsList.UnselectAll()
		callback(id)
	}
}

func (v *viewImpl) ShowRequests() {
	v.requestsContainer.Show()
	v.collectionsContainer.Hide()
}

func (v *viewImpl) ShowCollections() {
	v.collectionsContainer.Show()
	v.requestsContainer.Hide()
}

func (v *viewImpl) ShowingRequests() bool {
	return !v.requestsContainer.Hidden
}

func (v *viewImpl) ShowingCollections() bool {
	return !v.collectionsContainer.Hidden
}

func (v *viewImpl) SetAddCollectionValidator(validator func(string) error) {
	v.addCollectionEntry.Validator = validator
}

func (v *viewImpl) ValidateAddCollection() error {
	return v.addCollectionEntry.Validate()
}

func (v *viewImpl) SetAddCollectionHandler(handler func()) {
	v.addButton.OnTapped = handler
}
