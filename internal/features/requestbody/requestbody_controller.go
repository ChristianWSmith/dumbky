package requestbody

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/log"
	"dumbky/internal/validators"
	"errors"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model *model
	view  *view

	bodyFormKeyValueCtrl *keyvalueeditor.Controller
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	bodyFormKeyValueCtrl := keyvalueeditor.NewController(validators.ValidateFormBodyKey, validators.ValidateFormBodyValue)

	c := &Controller{
		model:                newModel(),
		view:                 newView(bodyFormKeyValueCtrl.CanvasObject()),
		bodyFormKeyValueCtrl: bodyFormKeyValueCtrl,
	}

	c.view.bodyTypeSelect.Bind(c.model.bodyTypeBinding)
	c.view.bodyRawEntry.Bind(c.model.bodyRawBinding)

	c.model.setBodyTypeListener(c.showBodyType)

	return c
}

func (c *Controller) GetBodyFormMap() map[string]string {
	return c.bodyFormKeyValueCtrl.GetMap()
}

func (c *Controller) ValidateBodyForm() error {
	return c.bodyFormKeyValueCtrl.Validate()
}

func (c *Controller) showBodyType() {
	bodyType := c.model.getBodyType()
	if bodyType == constants.UI_BODY_TYPE_FORM {
		c.bodyFormKeyValueCtrl.Show()
		c.view.hideBodyRaw()
	} else if bodyType == constants.UI_BODY_TYPE_RAW {
		c.bodyFormKeyValueCtrl.Hide()
		c.view.showBodyRaw()
	} else if bodyType == constants.UI_BODY_TYPE_NONE {
		c.bodyFormKeyValueCtrl.Hide()
		c.view.hideBodyRaw()
	} else {
		log.Error(errors.New("invalid body type"))
	}
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) ToState() RequestBodyState {
	return RequestBodyState{
		BodyType: c.model.getBodyType(),
		BodyForm: c.bodyFormKeyValueCtrl.ToState(),
		BodyRaw:  c.model.getBodyRaw(),
	}
}

func (c *Controller) LoadState(requestBodyState RequestBodyState) {
	c.model.setBodyType(requestBodyState.BodyType)
	c.bodyFormKeyValueCtrl.LoadState(requestBodyState.BodyForm)
	c.model.setBodyRaw(requestBodyState.BodyRaw)
}

func (c *Controller) ValidateBodyRaw() error {
	return c.view.validateBodyRaw()
}

func (c *Controller) EnableBodyTypeSelect() {
	c.view.enabledBodyTypeSelect()
}

func (c *Controller) DisableBodyTypeSelect() {
	c.view.disabledBodyTypeSelect()
}

func (c *Controller) GetBodyType() string {
	return c.model.getBodyType()
}

func (c *Controller) GetBodyRaw() string {
	return c.model.getBodyRaw()
}

func (c *Controller) SetBodyRaw(bodyRaw string) {
	c.model.setBodyRaw(bodyRaw)
}

func (c *Controller) SetBodyType(bodyType string) {
	c.model.setBodyType(bodyType)
}
