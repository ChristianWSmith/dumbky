package collectionsbrowserview

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"dumbky/internal/db"
	"dumbky/internal/global"
	"dumbky/internal/log"
)

type CollectionsBrowserView struct {
	UI                        *fyne.Container
	SelectedCollectionBinding binding.String
	SelectedRequestBinding    binding.String
	collectionNames           []string
	collectionsList           *widget.List
	requestsList              *widget.List
	stack                     *fyne.Container
	requestNames              []string
	currentCollection         string
	collectionsView           *fyne.Container
	requestsView              *fyne.Container
}

func ComposeCollectionsBrowserView() CollectionsBrowserView {
	selectedRequestBinding := binding.NewString()
	selectedCollectionBind := binding.NewString()
	view := CollectionsBrowserView{
		SelectedRequestBinding:    selectedRequestBinding,
		SelectedCollectionBinding: selectedCollectionBind,
		collectionNames:           db.FetchCollectionNames(),
	}

	// COLLECTIONS LIST
	view.collectionsList = widget.NewList(
		func() int { return len(view.collectionNames) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			menuButton := widget.NewButtonWithIcon("", nil, nil)
			menuButton.Icon = menuButton.Theme().Icon(theme.IconNameMoreVertical)
			return container.NewBorder(nil, nil, nil, menuButton, label)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			menuBtn := c.Objects[1].(*widget.Button)
			name := view.collectionNames[i]
			label.SetText(name)

			menuBtn.OnTapped = func() {
				pop := fyne.NewMenu("",
					fyne.NewMenuItem("Delete", func() {
						log.Info(fmt.Sprintf("TODO: Delete %s", name))
					}),
					fyne.NewMenuItem("Rename", func() {
						log.Info(fmt.Sprintf("TODO: Rename %s", name))
					}),
				)
				widget.ShowPopUpMenuAtRelativePosition(pop, global.Window.Canvas(), menuBtn.Position(), o)
			}
		},
	)
	view.collectionsList.OnSelected = func(id widget.ListItemID) {
		view.collectionsList.UnselectAll()
		view.showRequests(view.collectionNames[id])
		err := view.SelectedCollectionBinding.Set(view.collectionNames[id])
		log.Info(view.collectionNames[id])
		if err != nil {
			log.Error(err)
		}
	}

	// REQUESTS LIST (initially empty)
	view.requestsList = widget.NewList(
		func() int { return len(view.requestNames) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			menuButton := widget.NewButtonWithIcon("", nil, nil)
			menuButton.Icon = menuButton.Theme().Icon(theme.IconNameMoreVertical)
			return container.NewBorder(nil, nil, nil, menuButton, label)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			menuBtn := c.Objects[1].(*widget.Button)
			name := view.requestNames[i]
			label.SetText(name)

			menuBtn.OnTapped = func() {
				pop := fyne.NewMenu("",
					fyne.NewMenuItem("Delete", func() {
						log.Info(fmt.Sprintf("TODO: Delete %s", name))
					}),
					fyne.NewMenuItem("Rename", func() {
						log.Info(fmt.Sprintf("TODO: Rename %s", name))
					}),
				)
				widget.ShowPopUpMenuAtRelativePosition(pop, global.Window.Canvas(), menuBtn.Position(), o)
			}
		},
	)
	view.requestsList.OnSelected = func(id widget.ListItemID) {
		view.requestsList.UnselectAll()
		err := view.SelectedRequestBinding.Set(view.requestNames[id])
		log.Info(view.requestNames[id])
		if err != nil {
			log.Error(err)
		}
	}

	// BACK BUTTON
	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), nil)
	view.requestsView = container.NewBorder(backButton, nil, nil, nil, view.requestsList)

	// STACKED VIEWS

	// ADD COLLECTION BUTTON
	addBtn := widget.NewButton("Add Collection", func() {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Collection name")
		dialog.ShowForm("New Collection", "Add", "Cancel", []*widget.FormItem{
			widget.NewFormItem("Name", entry),
		}, func(ok bool) {
			if ok && entry.Text != "" {
				err := db.CreateCollection(entry.Text)
				if err != nil {
					dialog.ShowError(err, global.Window)
					return
				}
				view.refresh()
			}
		}, global.Window)
	})

	view.collectionsView = container.NewBorder(nil, addBtn, nil, nil, view.collectionsList)
	view.requestsView = container.NewBorder(backButton, nil, nil, nil, view.requestsList)
	view.requestsView.Hide()

	backButton.OnTapped = func() {
		view.requestsView.Hide()
		view.collectionsView.Show()
	}

	view.stack = container.NewStack(view.collectionsView, view.requestsView)
	view.UI = container.NewBorder(nil, nil, nil, nil, view.stack)
	return view
}

func (v *CollectionsBrowserView) refresh() {
	v.collectionNames = db.FetchCollectionNames()
	v.collectionsList.Refresh()
}

func (v *CollectionsBrowserView) showRequests(collection string) {
	v.currentCollection = collection
	v.requestNames = db.FetchRequestNames(collection)
	v.requestsList.Refresh()

	v.collectionsView.Hide()
	v.requestsView.Show()
}
