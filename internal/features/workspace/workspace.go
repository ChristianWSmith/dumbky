package workspace

import (
	"dumbky/internal/features"
	"dumbky/internal/features/workspace/controller"
)

type WorkspaceController interface {
	features.Controller
}

func New() WorkspaceController {
	return controller.NewController()
}
