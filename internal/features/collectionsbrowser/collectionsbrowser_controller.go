package collectionsbrowser

import (
	"dumbky/internal/db"
	"dumbky/internal/features"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Controller struct {
	model *model
	view  *view
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	c := &Controller{}

	requestsMenu := func(position fyne.Position, parent fyne.CanvasObject, name string) {
		pop := fyne.NewMenu("",
			fyne.NewMenuItem("Delete", func() {
				c.deleteRequest(name)
			}),
		)
		widget.ShowPopUpMenuAtRelativePosition(pop, global.Window.Canvas(), position, parent)
	}

	collectionsMenu := func(position fyne.Position, parent fyne.CanvasObject, name string) {
		pop := fyne.NewMenu("",
			fyne.NewMenuItem("Delete", func() {
				c.deleteCollection(name)
			}),
		)
		widget.ShowPopUpMenuAtRelativePosition(pop, global.Window.Canvas(), position, parent)
	}

	c.model = newModel()
	c.view = newView(c.model.requestsListBinding, c.model.collectionsListBinding, requestsMenu, collectionsMenu)

	c.view.collectionsList.OnSelected = func(id widget.ListItemID) {
		c.view.collectionsList.UnselectAll()
		name, _ := c.model.collectionsListBinding.GetValue(id)
		c.ShowRequests(name)
	}

	c.view.addCollectionEntry.Validator = validators.ValidateCollectionName
	c.view.addButton.OnTapped = func() {
		err := c.view.addCollectionEntry.Validate()
		if err != nil {
			log.Error(err)
			return
		}
		name, err := c.model.addCollectionBinding.Get()
		if err != nil {
			log.Error(err)
		}
		err = c.model.addCollectionBinding.Set("")
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
			fyne.Do(func() { c.ShowCollections() })
		}()
	}

	c.view.requestsList.OnSelected = func(id widget.ListItemID) {
		c.view.requestsList.UnselectAll()
		name, _ := c.model.requestsListBinding.GetValue(id)
		c.model.selectedRequestBinding.Set(name)
	}

	c.view.backButton.OnTapped = func() { c.ShowCollections() }

	c.view.addCollectionEntry.Bind(c.model.addCollectionBinding)
	c.view.selectedCollectionLabel.Bind(c.model.selectedCollectionBinding)

	c.ShowCollections()
	return c
}

func (c *Controller) GetUI() *fyne.Container {
	return c.view.getUI()
}

func (c *Controller) GetSelectedCollection() string {
	selectedCollection, _ := c.model.selectedCollectionBinding.Get()
	return selectedCollection
}

func (c *Controller) GetSelectedRequest() string {
	selectedRequest, _ := c.model.selectedRequestBinding.Get()
	return selectedRequest
}

func (c *Controller) SetSelectedRequest(requestName string) {
	c.model.selectedRequestBinding.Set(requestName)
}

func (c *Controller) SetSelectedRequestListener(handler func()) {
	c.model.selectedRequestBinding.AddListener(binding.NewDataListener(handler))
}

func (c *Controller) deleteRequest(name string) {
	collectionName, err := c.model.selectedCollectionBinding.Get()
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
			c.RefreshRequests()
		})
	}()
}

func (c *Controller) deleteCollection(name string) {
	go func() {
		err := db.DeleteCollection(name)
		if err != nil {
			log.Error(err)
		}
		fyne.Do(func() {
			c.RefreshCollections()
		})
	}()
}

func (c *Controller) ShowCollections() {
	c.model.selectedCollectionBinding.Set("")
	c.model.selectedRequestBinding.Set("")
	// update binding list
	go func() {
		names := db.FetchCollectionNames()
		fyne.Do(func() {
			c.model.collectionsListBinding.Set(names)

			c.view.requestsContainer.Hide()
			c.view.collectionsContainer.Show()
		})
	}()
}

func (c *Controller) ShowRequests(collection string) {
	c.model.selectedCollectionBinding.Set(collection)
	// update request names
	go func() {
		names := db.FetchRequestNames(collection)

		fyne.Do(func() {
			c.model.requestsListBinding.Set(names)

			c.view.collectionsContainer.Hide()
			c.view.requestsContainer.Show()
		})
	}()
}

func (c *Controller) RefreshCollections() {
	if !c.view.collectionsContainer.Hidden {
		c.ShowCollections()
	}
}

func (c *Controller) RefreshRequests() error {
	if !c.view.requestsContainer.Hidden {
		collection, err := c.model.selectedCollectionBinding.Get()
		if err != nil {
			return err
		}
		c.ShowRequests(collection)
	}
	return nil
}
