package controller

import (
	"dumbky/internal/db"
	"dumbky/internal/events"
	"dumbky/internal/features/collectionsbrowser/model"
	"dumbky/internal/features/collectionsbrowser/view"
	"dumbky/internal/global"
	"dumbky/internal/log"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type controllerImpl struct {
	model model.Model
	view  view.View
}

func NewController() *controllerImpl {
	c := &controllerImpl{}

	c.model = model.NewModel()

	bindings := c.model.GetBindings()

	c.view = view.NewView(bindings.RequestsList, bindings.CollectionsList,
		view.ViewCallbacks{
			OnDeleteCollection: c.deleteCollection,
			OnDeleteRequest:    c.deleteRequest,
		})

	bindables := c.view.GetBindables()
	bindables.AddCollection.Bind(bindings.AddCollection)
	bindables.SelectedCollection.Bind(bindings.SelectedCollection)

	c.view.SetAddCollectionValidator(validators.ValidateCollectionName)
	c.view.SetAddCollectionHandler(func() {
		err := c.view.ValidateAddCollection()
		if err != nil {
			log.Error(err)
			return
		}
		collectionName := c.model.GetAddCollection()
		c.model.SetAddCollection("")
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

	c.view.SetRequestSelectedCallback(func(id int) {
		requestName := c.model.GetRequestNameById(id)
		events.GetBus[events.RequestSelected]().Publish(events.RequestSelected{RequestName: requestName})
	})

	c.view.SetCollectionSelectedCallback(func(id int) {
		name := c.model.GetCollectionNameById(id)
		c.refreshAndShowRequests(name)
	})

	c.view.SetBackHandler(func() { c.refreshAndShowCollections() })

	c.refreshAndShowCollections()
	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) GetSelectedCollection() string {
	return c.model.GetSelectedCollection()
}

func (c *controllerImpl) LazyRefreshAndShowRequests() {
	if c.view.ShowingRequests() {
		collectionName := c.model.GetSelectedCollection()
		c.refreshAndShowRequests(collectionName)
	}
}

func (c *controllerImpl) deleteRequest(name string) {
	collectionName := c.model.GetSelectedCollection()
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

func (c *controllerImpl) deleteCollection(name string) {
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

func (c *controllerImpl) refreshAndShowCollections() {
	c.model.SetSelectedCollection("")

	go func() {
		collectionNames := db.FetchCollectionNames()
		fyne.Do(func() {
			c.model.SetCollectionsList(collectionNames)
			c.view.ShowCollections()
		})
	}()
}

func (c *controllerImpl) refreshAndShowRequests(collectionName string) {
	c.model.SetSelectedCollection(collectionName)

	go func() {
		requestNames := db.FetchRequestNames(collectionName)
		fyne.Do(func() {
			c.model.SetRequestsList(requestNames)
			c.view.ShowRequests()
		})
	}()
}

func (c *controllerImpl) lazyRefreshAndShowCollections() {
	if c.view.ShowingCollections() {
		c.refreshAndShowCollections()
	}
}
