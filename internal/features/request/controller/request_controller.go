package controller

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/features/request/model"
	"dumbky/internal/features/request/view"
	"dumbky/internal/log"
	"dumbky/internal/state"
	"dumbky/internal/utils"
	"errors"

	"fyne.io/fyne/v2"
)

type controllerImpl struct {
	model model.Model
	view  view.View

	queryParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController
	pathParamsKeyValueCtrl  keyvalueeditor.KeyValueEditorController
	headersKeyValueCtrl     keyvalueeditor.KeyValueEditorController
	bodyFormKeyValueCtrl    keyvalueeditor.KeyValueEditorController
}

func NewController(queryParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	pathParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	headersKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	bodyFormKeyValueCtrl keyvalueeditor.KeyValueEditorController) *controllerImpl {

	c := &controllerImpl{
		model:                   model.NewModel(),
		view:                    view.NewView(queryParamsKeyValueCtrl.CanvasObject(), pathParamsKeyValueCtrl.CanvasObject(), headersKeyValueCtrl.CanvasObject(), bodyFormKeyValueCtrl.CanvasObject()),
		queryParamsKeyValueCtrl: queryParamsKeyValueCtrl,
		pathParamsKeyValueCtrl:  pathParamsKeyValueCtrl,
		headersKeyValueCtrl:     headersKeyValueCtrl,
		bodyFormKeyValueCtrl:    bodyFormKeyValueCtrl,
	}

	c.bindAll()
	c.model.SetBodyTypeListener(c.showBodyType)

	return c
}

func (c *controllerImpl) CanvasObject() fyne.CanvasObject {
	return c.view.CanvasObject()
}

func (c *controllerImpl) GetQueryParamsMap() map[string]string {
	return c.queryParamsKeyValueCtrl.Get()
}

func (c *controllerImpl) GetPathParamsMap() map[string]string {
	return c.pathParamsKeyValueCtrl.Get()
}

func (c *controllerImpl) GetHeadersMap() map[string]string {
	return c.headersKeyValueCtrl.Get()
}

func (c *controllerImpl) ToState() state.RequestState {
	return state.RequestState{
		QueryParams: c.queryParamsKeyValueCtrl.ToState(),
		PathParams:  c.pathParamsKeyValueCtrl.ToState(),
		Headers:     c.headersKeyValueCtrl.ToState(),
		BodyForm:    c.bodyFormKeyValueCtrl.ToState(),
		BodyType:    c.model.GetBodyType(),
		BodyRaw:     c.model.GetBodyRaw(),
	}
}

func (c *controllerImpl) LoadState(requestState state.RequestState) {
	c.queryParamsKeyValueCtrl.LoadState(requestState.QueryParams)
	c.pathParamsKeyValueCtrl.LoadState(requestState.PathParams)
	c.headersKeyValueCtrl.LoadState(requestState.Headers)
	c.bodyFormKeyValueCtrl.LoadState(requestState.BodyForm)
	c.model.SetBodyType(requestState.BodyType)
	c.model.SetBodyRaw(requestState.BodyRaw)
}

func (c *controllerImpl) Validate() error {
	err := c.queryParamsKeyValueCtrl.Validate()
	if err != nil {
		return err
	}
	err = c.pathParamsKeyValueCtrl.Validate()
	if err != nil {
		return err
	}
	err = c.headersKeyValueCtrl.Validate()
	if err != nil {
		return err
	}
	err = c.bodyFormKeyValueCtrl.Validate()
	if err != nil {
		return err
	}
	return c.view.Validate()
}

func (c *controllerImpl) GetBodyFormMap() map[string]string {
	return c.bodyFormKeyValueCtrl.Get()
}
func (c *controllerImpl) SetBodyTypeSelectEnabled(enabled bool) {
	if enabled {
		c.view.SetBodyTypeSelectEnabled(true)
	} else {
		c.view.SetBodyTypeSelectEnabled(false)
		c.model.SetBodyType(constants.UI_BODY_TYPE_NONE)
	}
}

func (c *controllerImpl) GetBodyType() string {
	return c.model.GetBodyType()
}

func (c *controllerImpl) GetBodyRaw() string {
	return c.model.GetBodyRaw()
}

func (c *controllerImpl) FormatBodyRaw() {
	c.model.SetBodyRaw(utils.SmartFormat(c.model.GetBodyRaw()))
}

func (c *controllerImpl) bindAll() {
	bindings := c.model.GetBindings()
	bindables := c.view.GetBindables()
	bindables.BodyRaw.Bind(bindings.BodyRaw)
	bindables.BodyType.Bind(bindings.BodyType)

}

func (c *controllerImpl) showBodyType() {
	bodyType := c.model.GetBodyType()
	if bodyType == constants.UI_BODY_TYPE_FORM {
		c.bodyFormKeyValueCtrl.SetVisible(true)
		c.view.SetBodyRawVisible(false)
	} else if bodyType == constants.UI_BODY_TYPE_RAW {
		c.bodyFormKeyValueCtrl.SetVisible(false)
		c.view.SetBodyRawVisible(true)
	} else if bodyType == constants.UI_BODY_TYPE_NONE {
		c.bodyFormKeyValueCtrl.SetVisible(false)
		c.view.SetBodyRawVisible(false)
	} else {
		log.Error(errors.New("invalid body type"))
	}
}
