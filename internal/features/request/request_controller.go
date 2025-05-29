package request

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/log"
	"dumbky/internal/utils"
	"dumbky/internal/validators"
	"errors"

	"fyne.io/fyne/v2"
)

type controller struct {
	model *model
	view  *view

	queryParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController
	pathParamsKeyValueCtrl  keyvalueeditor.KeyValueEditorController
	headersKeyValueCtrl     keyvalueeditor.KeyValueEditorController
	bodyFormKeyValueCtrl    keyvalueeditor.KeyValueEditorController
}

type RequestController interface {
	features.Controller
	GetQueryParamsMap() map[string]string
	GetPathParamsMap() map[string]string
	GetHeadersMap() map[string]string
	ToState() RequestState
	LoadState(requestState RequestState)
	Validate() error
	SetBodyTypeSelectEnabled(enabled bool)
	FormatBodyRaw()
	GetBodyFormMap() map[string]string
	GetBodyType() string
	GetBodyRaw() string
}

var _ RequestController = (*controller)(nil)

func New() RequestController {
	queryParamsKeyValueCtrl := keyvalueeditor.New(validators.ValidateQueryParamKey, validators.ValidateQueryParamValue)
	pathParamsKeyValueCtrl := keyvalueeditor.New(validators.ValidatePathParamKey, validators.ValidatePathParamValue)
	headersKeyValueCtrl := keyvalueeditor.New(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyFormKeyValueCtrl := keyvalueeditor.New(validators.ValidateFormBodyKey, validators.ValidateFormBodyValue)
	return newController(queryParamsKeyValueCtrl, pathParamsKeyValueCtrl, headersKeyValueCtrl, bodyFormKeyValueCtrl)
}

func newController(queryParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	pathParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	headersKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	bodyFormKeyValueCtrl keyvalueeditor.KeyValueEditorController) *controller {

	c := &controller{
		model:                   newModel(),
		view:                    newView(queryParamsKeyValueCtrl.CanvasObject(), pathParamsKeyValueCtrl.CanvasObject(), headersKeyValueCtrl.CanvasObject(), bodyFormKeyValueCtrl.CanvasObject()),
		queryParamsKeyValueCtrl: queryParamsKeyValueCtrl,
		pathParamsKeyValueCtrl:  pathParamsKeyValueCtrl,
		headersKeyValueCtrl:     headersKeyValueCtrl,
		bodyFormKeyValueCtrl:    bodyFormKeyValueCtrl,
	}

	c.view.bodyTypeSelect.Bind(c.model.bodyTypeBinding)
	c.view.bodyRawEntry.Bind(c.model.bodyRawBinding)

	c.model.setBodyTypeListener(c.showBodyType)

	return c
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) GetQueryParamsMap() map[string]string {
	return c.queryParamsKeyValueCtrl.Get()
}

func (c *controller) GetPathParamsMap() map[string]string {
	return c.pathParamsKeyValueCtrl.Get()
}

func (c *controller) GetHeadersMap() map[string]string {
	return c.headersKeyValueCtrl.Get()
}

func (c *controller) ToState() RequestState {
	return RequestState{
		QueryParams: c.queryParamsKeyValueCtrl.ToState(),
		PathParams:  c.pathParamsKeyValueCtrl.ToState(),
		Headers:     c.headersKeyValueCtrl.ToState(),
		BodyForm:    c.bodyFormKeyValueCtrl.ToState(),
		BodyType:    c.model.getBodyType(),
		BodyRaw:     c.model.getBodyRaw(),
	}
}

func (c *controller) LoadState(requestState RequestState) {
	c.queryParamsKeyValueCtrl.LoadState(requestState.QueryParams)
	c.pathParamsKeyValueCtrl.LoadState(requestState.PathParams)
	c.headersKeyValueCtrl.LoadState(requestState.Headers)
	c.bodyFormKeyValueCtrl.LoadState(requestState.BodyForm)
	c.model.setBodyType(requestState.BodyType)
	c.model.setBodyRaw(requestState.BodyRaw)
}

func (c *controller) Validate() error {
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
	return c.view.validateBodyRaw()
}

func (c *controller) GetBodyFormMap() map[string]string {
	return c.bodyFormKeyValueCtrl.Get()
}
func (c *controller) SetBodyTypeSelectEnabled(enabled bool) {
	if enabled {
		c.view.enabledBodyTypeSelect()
	} else {
		c.view.disabledBodyTypeSelect()
		c.model.setBodyType(constants.UI_BODY_TYPE_NONE)
	}
}

func (c *controller) GetBodyType() string {
	return c.model.getBodyType()
}

func (c *controller) GetBodyRaw() string {
	return c.model.getBodyRaw()
}

func (c *controller) FormatBodyRaw() {
	c.model.setBodyRaw(utils.SmartFormat(c.model.getBodyRaw()))
}

func (c *controller) showBodyType() {
	bodyType := c.model.getBodyType()
	if bodyType == constants.UI_BODY_TYPE_FORM {
		c.bodyFormKeyValueCtrl.SetVisible(true)
		c.view.hideBodyRaw()
	} else if bodyType == constants.UI_BODY_TYPE_RAW {
		c.bodyFormKeyValueCtrl.SetVisible(false)
		c.view.showBodyRaw()
	} else if bodyType == constants.UI_BODY_TYPE_NONE {
		c.bodyFormKeyValueCtrl.SetVisible(false)
		c.view.hideBodyRaw()
	} else {
		log.Error(errors.New("invalid body type"))
	}
}
