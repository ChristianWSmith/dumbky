package response

import (
	"dumbky/internal/features"
	"dumbky/internal/features/template/controller"
)

type TemplateController interface {
	features.Controller
}

func New() TemplateController {
	return controller.NewController()
}
