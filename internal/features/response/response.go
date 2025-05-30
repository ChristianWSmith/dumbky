package response

import (
	"dumbky/internal/features"
	"dumbky/internal/features/response/controller"
	"dumbky/internal/restclient"
)

type ResponseController interface {
	features.Controller
	SetLoading(loading bool)
	Set(responsePayload restclient.ResponsePayload)
}

func New() ResponseController {
	return controller.NewController()
}
