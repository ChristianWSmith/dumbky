package root

import (
	"dumbky/internal/features"
	"dumbky/internal/features/collectionsbrowser"
	"dumbky/internal/features/root/controller"
	"dumbky/internal/features/workspace"
)

type RootController interface {
	features.Controller
}

func New() RootController {
	collectionsBrowserCtrl := collectionsbrowser.New()
	workspaceCtrl := workspace.New()
	return controller.NewController(collectionsBrowserCtrl, workspaceCtrl)
}
