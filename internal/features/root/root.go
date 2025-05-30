package root

import (
	"dumbky/internal/features"
	"dumbky/internal/features/dashboard"
	"dumbky/internal/features/root/controller"
)

type RootController interface {
	features.Controller
}

func New() RootController {
	dashboardController := dashboard.New()
	return controller.NewController(dashboardController)
}
