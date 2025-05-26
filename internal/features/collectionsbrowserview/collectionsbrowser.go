package collectionsbrowserview

import (
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
	"dumbky/internal/validators"
)

type CollectionsBrowserView struct {
	UI                        *fyne.Container
	SelectedCollectionBinding binding.String
	SelectedRequestBinding    binding.String

	addCollectionBinding binding.String
	addCollectionEntry   *widget.Entry

	collectionsBinding binding.StringList
	collectionsList    *widget.List

	requestsBinding binding.StringList
	requestsList    *widget.List

	collectionsView *fyne.Container
	requestsView    *fyne.Container
}

func ComposeCollectionsBrowserView() CollectionsBrowserView {
	// Bindings
	selectedRequest := binding.NewString()
	selectedCollection := binding.NewString()
	addCollectionBind := binding.NewString()

	collectionsBind := binding.NewStringList()
	requestsBind := binding.NewStringList()

	cbv := CollectionsBrowserView{
		SelectedRequestBinding:    selectedRequest,
		SelectedCollectionBinding: selectedCollection,
		addCollectionBinding:      addCollectionBind,
		collectionsBinding:        collectionsBind,
		requestsBinding:           requestsBind,
	}

	// Collections List
	cbv.collectionsList = widget.NewListWithData(
		collectionsBind,
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			menuBtn := widget.NewButtonWithIcon("", nil, nil)
			menuBtn.Icon = menuBtn.Theme().Icon(theme.IconNameMoreVertical)
			return container.NewBorder(nil, nil, nil, menuBtn, label)
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			name, _ := item.(binding.String).Get()
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			menuBtn := c.Objects[1].(*widget.Button)
			label.SetText(name)

			menuBtn.OnTapped = func() {
				pop := fyne.NewMenu("",
					fyne.NewMenuItem("Delete", func() { cbv.deleteCollection(name) }),
				)
				widget.ShowPopUpMenuAtRelativePosition(pop, global.Window.Canvas(), menuBtn.Position(), o)
			}
		},
	)
	cbv.collectionsList.OnSelected = func(id widget.ListItemID) {
		cbv.collectionsList.UnselectAll()
		name, _ := cbv.collectionsBinding.GetValue(id)
		cbv.ShowRequests(name)
	}

	// Add Collection
	cbv.addCollectionEntry = widget.NewEntry()
	cbv.addCollectionEntry.Bind(addCollectionBind)
	cbv.addCollectionEntry.Validator = validators.ValidateCollectionName
	addBtn := widget.NewButtonWithIcon("", nil, nil)
	addBtn.Icon = addBtn.Theme().Icon(theme.IconNameContentAdd)
	addBtn.OnTapped = func() {
		err := cbv.addCollectionEntry.Validate()
		if err != nil {
			log.Error(err)
			return
		}
		name, err := addCollectionBind.Get()
		if err != nil {
			log.Error(err)
		}
		err = addCollectionBind.Set("")
		if err != nil {
			log.Error(err)
		}
		if name == "" {
			return
		}
		go func() {
			err := db.CreateCollection(name)
			if err != nil {
				dialog.ShowError(err, global.Window)
				return
			}
			fyne.Do(func() { cbv.ShowCollections() })
		}()
	}
	addCollectionView := container.NewBorder(nil, nil, nil, addBtn, cbv.addCollectionEntry)

	// Collections Container
	collectionLabel := widget.NewLabel(constants.UI_LABEL_COLLECTIONS)
	cbv.collectionsView = container.NewBorder(collectionLabel, addCollectionView, nil, nil, cbv.collectionsList)

	// Requests List
	cbv.requestsList = widget.NewListWithData(
		requestsBind,
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			menuBtn := widget.NewButtonWithIcon("", nil, nil)
			menuBtn.Icon = menuBtn.Theme().Icon(theme.IconNameMoreVertical)
			return container.NewBorder(nil, nil, nil, menuBtn, label)
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			name, _ := item.(binding.String).Get()
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			menuBtn := c.Objects[1].(*widget.Button)
			label.SetText(name)

			menuBtn.OnTapped = func() {
				pop := fyne.NewMenu("",
					fyne.NewMenuItem("Delete", func() { cbv.deleteRequest(name) }),
				)
				widget.ShowPopUpMenuAtRelativePosition(pop, global.Window.Canvas(), menuBtn.Position(), o)
			}
		},
	)
	cbv.requestsList.OnSelected = func(id widget.ListItemID) {
		cbv.requestsList.UnselectAll()
		name, _ := cbv.requestsBinding.GetValue(id)
		cbv.SelectedRequestBinding.Set(name)
	}

	// Requests Container
	backBtn := widget.NewButtonWithIcon("Back", nil, nil)
	backBtn.Icon = backBtn.Theme().Icon(theme.IconNameNavigateBack)
	backBtn.OnTapped = func() { cbv.ShowCollections() }
	selectedColLabel := widget.NewLabel("")
	selectedColLabel.Bind(selectedCollection)
	cbv.requestsView = container.NewBorder(selectedColLabel, backBtn, nil, nil, cbv.requestsList)

	// Root Stack
	stack := container.NewStack(cbv.collectionsView, cbv.requestsView)
	cbv.UI = container.NewBorder(nil, nil, nil, nil, stack)

	// Initialize view
	cbv.ShowCollections()
	return cbv
}

func (cbv CollectionsBrowserView) deleteRequest(name string) {
	collectionName, err := cbv.SelectedCollectionBinding.Get()
	if err != nil {
		log.Error(err)
		return
	}
	go func() {
		err = db.DeleteRequest(collectionName, name)
		if err != nil {
			log.Error(err)
			return
		}
		fyne.Do(func() {
			cbv.RefreshRequests()
		})
	}()
}

func (cbv CollectionsBrowserView) deleteCollection(name string) {
	go func() {
		err := db.DeleteCollection(name)
		if err != nil {
			log.Error(err)
		}
		fyne.Do(func() {
			cbv.RefreshCollections()
		})
	}()
}

func (cbv *CollectionsBrowserView) ShowCollections() {
	cbv.SelectedCollectionBinding.Set("")
	cbv.SelectedRequestBinding.Set("")
	// update binding list
	go func() {
		names := db.FetchCollectionNames()
		fyne.Do(func() {
			cbv.collectionsBinding.Set(names)

			cbv.requestsView.Hide()
			cbv.collectionsView.Show()
		})
	}()
}

func (cbv *CollectionsBrowserView) ShowRequests(collection string) {
	cbv.SelectedCollectionBinding.Set(collection)
	// update request names
	go func() {
		names := db.FetchRequestNames(collection)

		fyne.Do(func() {
			cbv.requestsBinding.Set(names)

			cbv.collectionsView.Hide()
			cbv.requestsView.Show()
		})
	}()
}

func (cbv *CollectionsBrowserView) RefreshCollections() {
	if !cbv.collectionsView.Hidden {
		cbv.ShowCollections()
	}
}

func (cbv *CollectionsBrowserView) RefreshRequests() error {
	if !cbv.requestsView.Hidden {
		collection, err := cbv.SelectedCollectionBinding.Get()
		if err != nil {
			return err
		}
		cbv.ShowRequests(collection)
	}
	return nil
}
