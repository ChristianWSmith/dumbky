package exchangeheader

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/utils"
	"dumbky/internal/validators"

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

	c.view.methodSelect.Bind(c.model.methodBinding)
	c.view.urlEntry.Bind(c.model.urlBinding)
	c.view.sslCheck.Bind(c.model.useSSLBinding)

	c.view.setUrlValidator(validators.ValidateURL)

	return c
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) GetMethod() string {
	return c.model.getMethod()
}

func (c *Controller) GetURL() string {
	return c.model.getURL()
}

func (c *Controller) GetUseSSL() bool {
	return c.model.getUseSSL()
}

func (c *Controller) SetMethodListener(handler func()) {
	c.model.setMethodListener(handler)
}

func (c *Controller) SetSendHandler(handler func()) {
	c.view.setSendHandler(handler)
}

func (c *Controller) SetSendEnabled(enabled bool) {
	if enabled {
		c.view.enableSend()
	} else {
		c.view.disableSend()
	}
}

func (c *Controller) ToState() ExchangeHeaderState {
	return ExchangeHeaderState{
		Method: c.model.getMethod(),
		URL:    c.model.getURL(),
		UseSSL: c.model.getUseSSL(),
	}
}

func (c *Controller) LoadState(exchangeHeaderState ExchangeHeaderState) {
	method := exchangeHeaderState.Method
	if !utils.ElementInSlice(constants.HttpMethods(), method) {
		method = constants.HTTP_METHOD_DEFAULT
	}
	c.model.setMethod(method)
	c.model.setURL(exchangeHeaderState.URL)
	c.model.setUseSSL(exchangeHeaderState.UseSSL)
}

func (c *Controller) Validate() error {
	return c.view.validateURL()
}
