package response

import (
	"dumbky/internal/features"
	"dumbky/internal/features/response/controller"
	"dumbky/internal/requesthelper"
)

type ResponseController interface {
	features.Controller
	SetLoading(loading bool)
	Set(responsePayload requesthelper.ResponsePayload)
}

func New() ResponseController {
	return controller.NewController()
}
