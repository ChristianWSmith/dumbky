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

type Controller struct {
	model *model
	view  *view
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	c := &Controller{}

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

func (c *Controller) GetUI() *fyne.Container {
	return c.view.getUI()
}

func (c *Controller) GetSelectedCollection() string {
	return c.model.getSelectedCollection()
}

func (c *Controller) GetSelectedRequest() string {
	return c.model.getSelectedRequest()
}

func (c *Controller) SetSelectedRequest(requestName string) {
	c.model.setSelectedRequest(requestName)
}

func (c *Controller) SetSelectedRequestListener(handler func()) {
	c.model.setSelectedRequestListener(handler)
}

func (c *Controller) deleteRequest(name string) {
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

func (c *Controller) deleteCollection(name string) {
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

func (c *Controller) refreshAndShowCollections() {
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

func (c *Controller) refreshAndShowRequests(collectionName string) {
	c.model.setSelectedCollection(collectionName)

	go func() {
		requestNames := db.FetchRequestNames(collectionName)
		fyne.Do(func() {
			c.model.setRequestsList(requestNames)
			c.view.showRequests()
		})
	}()
}

func (c *Controller) lazyRefreshAndShowCollections() {
	if c.view.showingCollections() {
		c.refreshAndShowCollections()
	}
}

func (c *Controller) LazyRefreshAndShowRequests() {
	if c.view.showingRequests() {
		collectionName := c.model.getSelectedCollection()
		c.refreshAndShowRequests(collectionName)
	}
}
