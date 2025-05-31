package controller

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/response/model"
	"dumbky/internal/features/response/view"
	"dumbky/internal/httputils"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2"
)

type controllerImpl struct {
	model model.Model
	view  view.View
}

func NewController() *controllerImpl {
	model := model.NewModel()
	view := view.NewView()

	c := &controllerImpl{
		model: model,
		view:  view,
	}

	c.bindAll()

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) SetLoading(loading bool) {
	c.view.SetLoading(loading)
	if loading {
		c.model.SetResponse(constants.UI_LOADING_RESPONSE_STATUS,
			constants.UI_LOADING_RESPONSE_TIME,
			constants.UI_LOADING_RESPONSE_BODY)
	}
}

func (c *controllerImpl) Set(responsePayload httputils.ResponsePayload) {
	c.model.SetResponse(
		responsePayload.Status,
		responsePayload.Time,
		utils.SmartFormat(responsePayload.Body))
}

func (c *controllerImpl) bindAll() {
	bindings := c.model.GetBindings()
	bindables := c.view.GetBindables()
	bindables.Body.Bind(bindings.Body)
	bindables.Time.Bind(bindings.Time)
	bindables.Status.Bind(bindings.Status)
}
