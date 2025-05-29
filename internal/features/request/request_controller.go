package request

import (
	"dumbky/internal/features"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/features/requestbody"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
)

type Controller struct {
	model *model
	view  *view

	queryParamsKeyValueCtrl *keyvalueeditor.Controller
	pathParamsKeyValueCtrl  *keyvalueeditor.Controller
	headersKeyValueCtrl     *keyvalueeditor.Controller
	requestBodyCtrl         *requestbody.Controller
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	queryParamsKeyValueCtrl := keyvalueeditor.NewController(validators.ValidateQueryParamKey, validators.ValidateQueryParamValue)
	pathParamsKeyValueCtrl := keyvalueeditor.NewController(validators.ValidatePathParamKey, validators.ValidatePathParamValue)
	headersKeyValueCtrl := keyvalueeditor.NewController(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyKeyValueCtrl := requestbody.NewController()

	return &Controller{
		model:                   newModel(),
		view:                    newView(queryParamsKeyValueCtrl.CanvasObject(), pathParamsKeyValueCtrl.CanvasObject(), headersKeyValueCtrl.CanvasObject(), bodyKeyValueCtrl.CanvasObject()),
		queryParamsKeyValueCtrl: queryParamsKeyValueCtrl,
		pathParamsKeyValueCtrl:  pathParamsKeyValueCtrl,
		headersKeyValueCtrl:     headersKeyValueCtrl,
		requestBodyCtrl:         bodyKeyValueCtrl,
	}
}

func (c *Controller) CanvasObject() fyne.CanvasObject {
	return c.view.canvasObject()
}

func (c *Controller) GetQueryParamsMap() map[string]string {
	return c.queryParamsKeyValueCtrl.Get()
}

func (c *Controller) GetPathParamsMap() map[string]string {
	return c.pathParamsKeyValueCtrl.Get()
}

func (c *Controller) GetHeadersMap() map[string]string {
	return c.headersKeyValueCtrl.Get()
}

func (c *Controller) ToState() RequestState {
	queryParams := c.queryParamsKeyValueCtrl.ToState()
	pathParams := c.pathParamsKeyValueCtrl.ToState()
	headers := c.headersKeyValueCtrl.ToState()
	body := c.requestBodyCtrl.ToState()
	return RequestState{
		QueryParams: queryParams,
		PathParams:  pathParams,
		Headers:     headers,
		Body:        body,
	}
}

func (c *Controller) LoadState(requestState RequestState) {
	c.queryParamsKeyValueCtrl.LoadState(requestState.QueryParams)
	c.pathParamsKeyValueCtrl.LoadState(requestState.PathParams)
	c.headersKeyValueCtrl.LoadState(requestState.Headers)
	c.requestBodyCtrl.LoadState(requestState.Body)
}

func (c *Controller) Validate() error {
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
	return c.requestBodyCtrl.Validate()
}

func (c *Controller) GetRequestBodyFormMap() map[string]string {
	return c.requestBodyCtrl.GetBodyFormMap()
}

func (c *Controller) GetRequestBodyRaw() string {
	return c.requestBodyCtrl.GetBodyRaw()
}

func (c *Controller) GetRequestBodyType() string {
	return c.requestBodyCtrl.GetBodyType()
}

func (c *Controller) SetBodyTypeSelectEnabled(enabled bool) {
	c.requestBodyCtrl.SetBodyTypeSelectEnabled(enabled)
}

func (c *Controller) FormatBodyRaw() {
	c.requestBodyCtrl.FormatBodyRaw()
}
