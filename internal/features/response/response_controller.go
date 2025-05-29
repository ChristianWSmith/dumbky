package response

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/requesthelper"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model *model
	view  *view
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	model := newModel()
	view := newView()

	c := &Controller{
		model: model,
		view:  view,
	}

	c.view.status.Bind(c.model.status)
	c.view.time.Bind(c.model.time)
	c.view.body.Bind(c.model.body)

	return c
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) SetLoading(loading bool) {
	c.view.setLoading(loading)
	if loading {
		c.model.setResponse(constants.UI_LOADING_RESPONSE_STATUS,
			constants.UI_LOADING_RESPONSE_TIME,
			constants.UI_LOADING_RESPONSE_BODY)
	}
}

func (c *Controller) Set(responsePayload requesthelper.ResponsePayload) {
	c.model.setResponse(
		responsePayload.Status,
		responsePayload.Time,
		utils.SmartFormat(responsePayload.Body))
}
