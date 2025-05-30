package dashboard

import (
	"dumbky/internal/features"
	"dumbky/internal/features/collectionsbrowser"
	"dumbky/internal/features/dashboard/controller"
	"dumbky/internal/features/workspace"
)

type DashboardController interface {
	features.Controller
}

func New() DashboardController {
	collectionsBrowserCtrl := collectionsbrowser.New()
	workspaceCtrl := workspace.New()
	return controller.NewController(collectionsBrowserCtrl, workspaceCtrl)
}
