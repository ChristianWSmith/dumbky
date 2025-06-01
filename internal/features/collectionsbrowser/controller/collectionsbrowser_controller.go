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
		events.Publish(events.RequestSelected{RequestName: requestName, IsSelected: true})
	})

	c.view.SetCollectionSelectedCallback(func(id int) {
		collectionName := c.model.GetCollectionNameById(id)
		c.model.SetSelectedCollection(collectionName)
		events.Publish(events.CollectionSelected{CollectionName: collectionName, IsSelected: true})
		c.refreshAndShowRequests()
	})

	c.view.SetBackHandler(func() { c.refreshAndShowCollections() })

	events.Subscribe(func(event events.RequestSaved) {
		c.lazyRefreshAndShowRequests()
	})

	c.refreshAndShowCollections()
	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) lazyRefreshAndShowRequests() {
	if c.view.ShowingRequests() {
		c.refreshAndShowRequests()
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
			c.lazyRefreshAndShowRequests()
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
	events.Publish(events.CollectionSelected{IsSelected: false})

	go func() {
		collectionNames := db.FetchCollectionNames()
		fyne.Do(func() {
			c.model.SetCollectionsList(collectionNames)
			c.view.ShowCollections()
		})
	}()
}

func (c *controllerImpl) refreshAndShowRequests() {
	collectionName := c.model.GetSelectedCollection()
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
