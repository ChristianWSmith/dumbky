package response

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/requesthelper"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2"
)

type controller struct {
	model *model
	view  *view
}

type ResponseController interface {
	features.Controller
	SetLoading(loading bool)
	Set(responsePayload requesthelper.ResponsePayload)
}

var _ ResponseController = (*controller)(nil)

func New() ResponseController {
	return newController()
}

func newController() *controller {
	model := newModel()
	view := newView()

	c := &controller{
		model: model,
		view:  view,
	}

	c.view.status.Bind(c.model.status)
	c.view.time.Bind(c.model.time)
	c.view.body.Bind(c.model.body)

	return c
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) SetLoading(loading bool) {
	c.view.setLoading(loading)
	if loading {
		c.model.setResponse(constants.UI_LOADING_RESPONSE_STATUS,
			constants.UI_LOADING_RESPONSE_TIME,
			constants.UI_LOADING_RESPONSE_BODY)
	}
}

func (c *controller) Set(responsePayload requesthelper.ResponsePayload) {
	c.model.setResponse(
		responsePayload.Status,
		responsePayload.Time,
		utils.SmartFormat(responsePayload.Body))
}
