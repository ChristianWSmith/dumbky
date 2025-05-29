package exchangeheader

import (
	"dumbky/internal/constants"
	"dumbky/internal/features"
	"dumbky/internal/utils"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
)

type controller struct {
	model *model
	view  *view
}

type ExchangeHeaderController interface {
	features.Controller
	GetMethod() string
	GetURL() string
	GetUseSSL() bool
	SetMethodListener(handler func())
	SetSendHandler(handler func())
	SetSendEnabled(enabled bool)
	ToState() ExchangeHeaderState
	LoadState(exchangeHeaderState ExchangeHeaderState)
	Validate() error
}

var _ ExchangeHeaderController = (*controller)(nil)

func NewController() ExchangeHeaderController {
	return newController()
}

func newController() *controller {
	model := newModel()
	view := newView()

	c := &controller{
		model: model,
		view:  view,
	}

	c.view.methodSelect.Bind(c.model.methodBinding)
	c.view.urlEntry.Bind(c.model.urlBinding)
	c.view.sslCheck.Bind(c.model.useSSLBinding)

	c.view.setUrlValidator(validators.ValidateURL)

	return c
}

func (c *controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *controller) GetMethod() string {
	return c.model.getMethod()
}

func (c *controller) GetURL() string {
	return c.model.getURL()
}

func (c *controller) GetUseSSL() bool {
	return c.model.getUseSSL()
}

func (c *controller) SetMethodListener(handler func()) {
	c.model.setMethodListener(handler)
}

func (c *controller) SetSendHandler(handler func()) {
	c.view.setSendHandler(handler)
}

func (c *controller) SetSendEnabled(enabled bool) {
	if enabled {
		c.view.enableSend()
	} else {
		c.view.disableSend()
	}
}

func (c *controller) ToState() ExchangeHeaderState {
	return ExchangeHeaderState{
		Method: c.model.getMethod(),
		URL:    c.model.getURL(),
		UseSSL: c.model.getUseSSL(),
	}
}

func (c *controller) LoadState(exchangeHeaderState ExchangeHeaderState) {
	method := exchangeHeaderState.Method
	if !utils.ElementInSlice(constants.HttpMethods(), method) {
		method = constants.HTTP_METHOD_DEFAULT
	}
	c.model.setMethod(method)
	c.model.setURL(exchangeHeaderState.URL)
	c.model.setUseSSL(exchangeHeaderState.UseSSL)
}

func (c *controller) Validate() error {
	return c.view.validateURL()
}
