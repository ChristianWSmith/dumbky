package response

import (
	"dumbky/internal/features"
	"dumbky/internal/features/response/controller"
	"dumbky/internal/httputils"
)

type ResponseController interface {
	features.Controller
	SetLoading(loading bool)
	Set(responsePayload httputils.ResponsePayload)
}

func New() ResponseController {
	return controller.NewController()
}
