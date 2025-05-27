package collectionsbrowser

import (
	"dumbky/internal/constants"

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

func newView(requestsListBinding, collectionsListBinding binding.StringList, requestsMenu, collectionsMenu func(fyne.Position, fyne.CanvasObject, string)) *view {
	addCollectionEntry := widget.NewEntry()

	requestsList := widget.NewListWithData(
		requestsListBinding,
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			menuBtn := widget.NewButtonWithIcon("", theme.MoreVerticalIcon(), nil)
			return container.NewBorder(nil, nil, nil, menuBtn, label)
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			name, _ := item.(binding.String).Get()
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			menuBtn := c.Objects[1].(*widget.Button)
			label.SetText(name)

			menuBtn.OnTapped = func() {
				requestsMenu(menuBtn.Position(), o, name)
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
			name, _ := item.(binding.String).Get()
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			menuBtn := c.Objects[1].(*widget.Button)
			label.SetText(name)

			menuBtn.OnTapped = func() {
				collectionsMenu(menuBtn.Position(), o, name)
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

func (v *view) getUI() *fyne.Container {
	return v.ui
}
