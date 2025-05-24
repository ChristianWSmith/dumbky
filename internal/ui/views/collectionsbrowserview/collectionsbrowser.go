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

	collectionNames []string
	requestNames    []string
	collectionsList *widget.List
	requestsList    *widget.List
	collectionsView *fyne.Container
	requestsView    *fyne.Container
}

func ComposeCollectionsBrowserView() CollectionsBrowserView {
	selectedRequestBinding := binding.NewString()
	selectedCollectionBind := binding.NewString()
	cbv := CollectionsBrowserView{
		SelectedRequestBinding:    selectedRequestBinding,
		SelectedCollectionBinding: selectedCollectionBind,
		collectionNames:           db.FetchCollectionNames(),
	}

	// COLLECTIONS LIST
	cbv.collectionsList = widget.NewList(
		func() int { return len(cbv.collectionNames) },
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
			name := cbv.collectionNames[i]
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
	cbv.collectionsList.OnSelected = func(id widget.ListItemID) {
		cbv.collectionsList.UnselectAll()
		cbv.ShowRequests(cbv.collectionNames[id])
		err := cbv.SelectedCollectionBinding.Set(cbv.collectionNames[id])
		log.Info(cbv.collectionNames[id])
		if err != nil {
			log.Error(err)
		}
	}

	// REQUESTS LIST (initially empty)
	cbv.requestsList = widget.NewList(
		func() int { return len(cbv.requestNames) },
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
			name := cbv.requestNames[i]
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
	cbv.requestsList.OnSelected = func(id widget.ListItemID) {
		cbv.requestsList.UnselectAll()
		err := cbv.SelectedRequestBinding.Set(cbv.requestNames[id])
		log.Info(cbv.requestNames[id])
		if err != nil {
			log.Error(err)
		}
	}

	// BACK BUTTON
	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), nil)
	cbv.requestsView = container.NewBorder(backButton, nil, nil, nil, cbv.requestsList)

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
				cbv.ShowCollections()
			}
		}, global.Window)
	})

	cbv.collectionsView = container.NewBorder(nil, addBtn, nil, nil, cbv.collectionsList)
	cbv.requestsView = container.NewBorder(nil, backButton, nil, nil, cbv.requestsList)
	cbv.requestsView.Hide()

	backButton.OnTapped = func() {
		cbv.requestsView.Hide()
		cbv.collectionsView.Show()
	}

	stack := container.NewStack(cbv.collectionsView, cbv.requestsView)
	cbv.UI = container.NewBorder(nil, nil, nil, nil, stack)
	return cbv
}

func (v *CollectionsBrowserView) ShowCollections() {
	v.collectionNames = db.FetchCollectionNames()
	v.collectionsList.Refresh()

	v.requestsView.Hide()
	v.collectionsView.Show()
}

func (v *CollectionsBrowserView) ShowRequests(collection string) {
	v.requestNames = db.FetchRequestNames(collection)
	v.requestsList.Refresh()

	v.collectionsView.Hide()
	v.requestsView.Show()
}
