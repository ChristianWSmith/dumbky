package workspace

import (
	"dumbky/internal/features"
	"dumbky/internal/features/workspace/controller"
	"dumbky/internal/state"
)

type WorkspaceController interface {
	features.Controller
	SetAddHandler(handler func())
	SetSaveHandler(handler func())
	OpenTab(document state.DocumentState)
	SaveTab(callback func()) error
	LoadTab(collectionName, requestName string)
}

func New() WorkspaceController {
	return controller.NewController()
}
