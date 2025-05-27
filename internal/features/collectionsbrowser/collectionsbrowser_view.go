package collectionsbrowser

import (
	"dumbky/internal/constants"
	"dumbky/internal/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type view struct {
	ui                      *fyne.Container
	addCollectionEntry      *widget.Entry
	collectionsList         *widget.List
	requestsList            *widget.List
	collectionsContainer    *fyne.Container
	requestsContainer       *fyne.Container
	addButton               *widget.Button
	backButton              *widget.Button
	selectedCollectionLabel *widget.Label
}

func newView(requestsListBinding, collectionsListBinding binding.StringList, onDeleteRequest, onDeleteCollection func(string)) *view {
	addCollectionEntry := widget.NewEntry()

	requestsMenu := func(position fyne.Position, parent fyne.CanvasObject, requestName string) {
		pop := fyne.NewMenu("",
			fyne.NewMenuItem(constants.UI_LABEL_DELETE, func() {
				onDeleteRequest(requestName)
			}),
		)
		widget.ShowPopUpMenuAtRelativePosition(pop, global.Window.Canvas(), position, parent)
	}

	collectionsMenu := func(position fyne.Position, parent fyne.CanvasObject, collectionName string) {
		pop := fyne.NewMenu("",
			fyne.NewMenuItem(constants.UI_LABEL_DELETE, func() {
				onDeleteCollection(collectionName)
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

	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), nil)
	selectedCollectionLabel := widget.NewLabel("")
	requestsContainer := container.NewBorder(selectedCollectionLabel, backButton, nil, nil, requestsList)

	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)
	addCollectionContainer := container.NewBorder(nil, nil, nil, addButton, addCollectionEntry)
	collectionLabel := widget.NewLabel(constants.UI_LABEL_COLLECTIONS)
	collectionsContainer := container.NewBorder(collectionLabel, addCollectionContainer, nil, nil, collectionsList)

	stack := container.NewStack(collectionsContainer, requestsContainer)

	ui := container.NewBorder(nil, nil, nil, nil, stack)

	return &view{
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

func (v *view) canvasObject() fyne.CanvasObject {
	return v.ui
}

func (v *view) setBackHandler(handler func()) {
	v.backButton.OnTapped = handler
}

func (v *view) setRequestSelectedCallback(callback func(int)) {
	v.requestsList.OnSelected = func(id widget.ListItemID) {
		v.requestsList.UnselectAll()
		callback(id)
	}
}

func (v *view) setCollectionSelectedCallback(callback func(int)) {
	v.collectionsList.OnSelected = func(id widget.ListItemID) {
		v.collectionsList.UnselectAll()
		callback(id)
	}
}

func (v *view) showRequests() {
	v.requestsContainer.Show()
	v.collectionsContainer.Hide()
}

func (v *view) showCollections() {
	v.collectionsContainer.Show()
	v.requestsContainer.Hide()
}

func (v *view) showingRequests() bool {
	return !v.requestsContainer.Hidden
}

func (v *view) showingCollections() bool {
	return !v.collectionsContainer.Hidden
}

func (v *view) setAddCollectionValidator(validator func(string) error) {
	v.addCollectionEntry.Validator = validator
}

func (v *view) validateAddCollection() error {
	return v.addCollectionEntry.Validate()
}

func (v *view) setAddCollectionHandler(handler func()) {
	v.addButton.OnTapped = handler
}
