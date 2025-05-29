package collectionsbrowser

import (
	"dumbky/internal/db"
	"dumbky/internal/features"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type controller struct {
	model *model
	view  *view
}

type CollectionsBrowserController interface {
	features.Controller
	GetSelectedCollection() string
	GetSelectedRequest() string
	SetSelectedRequestCallback(callback func())
	LazyRefreshAndShowRequests()
}

var _ CollectionsBrowserController = (*controller)(nil)

func New() CollectionsBrowserController {
	return newController()
}

func newController() *controller {
	c := &controller{}

	c.model = newModel()
	c.view = newView(c.model.requestsListBinding, c.model.collectionsListBinding, c.deleteRequest, c.deleteCollection)

	c.view.addCollectionEntry.Bind(c.model.addCollectionBinding)
	c.view.selectedCollectionLabel.Bind(c.model.selectedCollectionBinding)

	c.view.setAddCollectionValidator(validators.ValidateCollectionName)
	c.view.setAddCollectionHandler(func() {
		err := c.view.validateAddCollection()
		if err != nil {
			log.Error(err)
			return
		}
		collectionName := c.model.getAddCollection()
		c.model.setAddCollection("")
		if collectionName == "" {
			return
		}
		go func() {
			err := db.CreateCollection(collectionName)
			if err != nil {
				dialog.ShowError(err, global.Window)
				return
			}
			fyne.Do(func() { c.refreshAndShowCollections() })
		}()
	})

	c.view.setRequestSelectedCallback(func(id int) {
		name := c.model.getRequestNameById(id)
		c.model.setSelectedRequest(name)
	})

	c.view.setCollectionSelectedCallback(func(id int) {
		name := c.model.getCollectionNameById(id)
		c.refreshAndShowRequests(name)
	})

	c.view.setBackHandler(func() { c.refreshAndShowCollections() })

	c.refreshAndShowCollections()
	return c
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) GetSelectedCollection() string {
	return c.model.getSelectedCollection()
}

func (c *controller) GetSelectedRequest() string {
	return c.model.getSelectedRequest()
}

func (c *controller) SetSelectedRequestCallback(callback func()) {
	c.model.setSelectedRequestCallback(callback)
}

func (c *controller) LazyRefreshAndShowRequests() {
	if c.view.showingRequests() {
		collectionName := c.model.getSelectedCollection()
		c.refreshAndShowRequests(collectionName)
	}
}

func (c *controller) deleteRequest(name string) {
	collectionName := c.model.getSelectedCollection()
	go func() {
		err := db.DeleteRequest(collectionName, name)
		if err != nil {
			log.Error(err)
			return
		}
		fyne.Do(func() {
			c.LazyRefreshAndShowRequests()
		})
	}()
}

func (c *controller) deleteCollection(name string) {
	go func() {
		err := db.DeleteCollection(name)
		if err != nil {
			log.Error(err)
		}
		fyne.Do(func() {
			c.lazyRefreshAndShowCollections()
		})
	}()
}

func (c *controller) refreshAndShowCollections() {
	c.model.setSelectedCollection("")
	c.model.setSelectedRequest("")

	go func() {
		collectionNames := db.FetchCollectionNames()
		fyne.Do(func() {
			c.model.setCollectionsList(collectionNames)
			c.view.showCollections()
		})
	}()
}

func (c *controller) refreshAndShowRequests(collectionName string) {
	c.model.setSelectedCollection(collectionName)

	go func() {
		requestNames := db.FetchRequestNames(collectionName)
		fyne.Do(func() {
			c.model.setRequestsList(requestNames)
			c.view.showRequests()
		})
	}()
}

func (c *controller) lazyRefreshAndShowCollections() {
	if c.view.showingCollections() {
		c.refreshAndShowCollections()
	}
}
