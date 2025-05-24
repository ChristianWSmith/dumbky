package collectionsbrowserview

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"dumbky/internal/constants"
	"dumbky/internal/db"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/ui/validators"
)

type CollectionsBrowserView struct {
	UI                        *fyne.Container
	SelectedCollectionBinding binding.String
	SelectedRequestBinding    binding.String

	addCollectionBinding binding.String
	addCollectionEntry   *widget.Entry
	collectionNames      []string
	requestNames         []string
	collectionsList      *widget.List
	requestsList         *widget.List
	collectionsView      *fyne.Container
	requestsView         *fyne.Container
}

func ComposeCollectionsBrowserView() CollectionsBrowserView {
	selectedRequestBinding := binding.NewString()
	selectedCollectionBind := binding.NewString()
	cbv := CollectionsBrowserView{
		SelectedRequestBinding:    selectedRequestBinding,
		SelectedCollectionBinding: selectedCollectionBind,
		collectionNames:           []string{},
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
	}

	// ADD COLLECTION BUTTON
	addCollectionBind := binding.NewString()
	cbv.addCollectionBinding = addCollectionBind
	cbv.addCollectionEntry = widget.NewEntry()
	cbv.addCollectionEntry.Bind(addCollectionBind)
	cbv.addCollectionEntry.Validator = validators.ValidateCollectionName

	addCollectionButton := widget.NewButtonWithIcon("", nil, nil)
	addCollectionButton.Icon = addCollectionButton.Theme().Icon(theme.IconNameContentAdd)
	addCollectionButton.OnTapped = func() {
		err := cbv.addCollectionEntry.Validate()
		if err != nil {
			log.Error(err)
			return
		}
		collectionName, err := cbv.addCollectionBinding.Get()
		if err != nil {
			log.Error(err)
			return
		}
		err = cbv.addCollectionBinding.Set("")
		if err != nil {
			log.Error(err)
			return
		}
		go func() {
			err = db.CreateCollection(collectionName)
			if err != nil {
				dialog.ShowError(err, global.Window)
				return
			}
			fyne.Do(func() {
				cbv.ShowCollections()
			})
		}()
	}
	addCollectionView := container.NewBorder(nil, nil, nil, addCollectionButton, cbv.addCollectionEntry)

	// COLLECTIONS LABEL
	collectionLabel := widget.NewLabel(constants.UI_LABEL_COLLECTIONS)

	cbv.collectionsView = container.NewBorder(collectionLabel, addCollectionView, nil, nil, cbv.collectionsList)

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
		log.Debug(fmt.Sprintf("selected request %s", cbv.requestNames[id]))
		if err != nil {
			log.Error(err)
		}
	}

	// BACK BUTTON
	backButton := widget.NewButtonWithIcon("Back", nil, nil)
	backButton.Icon = backButton.Theme().Icon(theme.IconNameNavigateBack)
	backButton.OnTapped = func() {
		cbv.ShowCollections()
	}

	// SELECTED COLLECTION LABEL
	selectedCollectionLabel := widget.NewLabel("")
	selectedCollectionLabel.Bind(selectedCollectionBind)

	cbv.requestsView = container.NewBorder(selectedCollectionLabel, backButton, nil, nil, cbv.requestsList)

	// FINAL UI

	stack := container.NewStack(cbv.collectionsView, cbv.requestsView)
	cbv.UI = container.NewBorder(nil, nil, nil, nil, stack)
	cbv.ShowCollections()
	return cbv
}

func (cbv *CollectionsBrowserView) ShowCollections() {
	err := cbv.SelectedCollectionBinding.Set("")
	log.Debug("selected collection reset")
	if err != nil {
		log.Error(err)
	}
	err = cbv.SelectedRequestBinding.Set("")
	log.Debug("selected request reset")
	if err != nil {
		log.Error(err)
	}
	cbv.collectionNames = db.FetchCollectionNames()
	cbv.collectionsList.Refresh()

	cbv.requestsView.Hide()
	cbv.collectionsView.Show()
}

func (cbv *CollectionsBrowserView) ShowRequests(collection string) {
	err := cbv.SelectedCollectionBinding.Set(collection)

	log.Debug(fmt.Sprintf("selected collection %s", collection))
	if err != nil {
		log.Error(err)
	}
	cbv.requestNames = db.FetchRequestNames(collection)
	cbv.requestsList.Refresh()

	cbv.collectionsView.Hide()
	cbv.requestsView.Show()
}

func (cbv *CollectionsBrowserView) RefreshRequests() error {
	collectionName, err := cbv.SelectedCollectionBinding.Get()
	if err != nil {
		log.Error(err)
		return err
	}
	cbv.ShowRequests(collectionName)
	return nil
}
