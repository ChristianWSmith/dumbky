package collectionsbrowser

import (
	"dumbky/internal/features"
	"dumbky/internal/features/collectionsbrowser/controller"
)

type CollectionsBrowserController interface {
	features.Controller
	GetSelectedCollection() string
	LazyRefreshAndShowRequests()
}

func New() CollectionsBrowserController {
	return controller.NewController()
}
