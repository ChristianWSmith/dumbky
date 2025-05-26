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
	Body                    *requestbody.Controller
}

var _ features.Controller = (*Controller)(nil)

func NewController() *Controller {
	queryParamsKeyValueCtrl := keyvalueeditor.NewController(validators.ValidateQueryParamKey, validators.ValidateQueryParamValue)
	pathParamsKeyValueCtrl := keyvalueeditor.NewController(validators.ValidatePathParamKey, validators.ValidatePathParamValue)
	headersKeyValueCtrl := keyvalueeditor.NewController(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	bodyKeyValueCtrl := requestbody.NewController()

	return &Controller{
		model:                   newModel(),
		view:                    newView(queryParamsKeyValueCtrl.GetUI(), pathParamsKeyValueCtrl.GetUI(), headersKeyValueCtrl.GetUI(), bodyKeyValueCtrl.GetUI()),
		queryParamsKeyValueCtrl: queryParamsKeyValueCtrl,
		pathParamsKeyValueCtrl:  pathParamsKeyValueCtrl,
		headersKeyValueCtrl:     headersKeyValueCtrl,
		Body:                    bodyKeyValueCtrl,
	}
}

func (c *Controller) GetUI() *fyne.Container {
	return c.view.getUI()
}

func (c *Controller) GetQueryParamsMap() map[string]string {
	return c.queryParamsKeyValueCtrl.GetMap()
}

func (c *Controller) GetPathParamsMap() map[string]string {
	return c.pathParamsKeyValueCtrl.GetMap()
}

func (c *Controller) GetHeadersMap() map[string]string {
	return c.headersKeyValueCtrl.GetMap()
}

func (c *Controller) ValidateQueryParams() error {
	return c.queryParamsKeyValueCtrl.Validate()
}

func (c *Controller) ValidatePathParams() error {
	return c.pathParamsKeyValueCtrl.Validate()
}

func (c *Controller) ValidateHeaders() error {
	return c.headersKeyValueCtrl.Validate()
}

func (c *Controller) GetRequestBodyFormMap() map[string]string {
	return c.Body.GetBodyFormMap()
}

func (c *Controller) ValidateRequestBodyForm() error {
	return c.Body.ValidateBodyForm()
}

func (c *Controller) GetRequestBodyRaw() string {
	return c.Body.GetBodyRaw()
}

func (c *Controller) ValidateRequestBodyRaw() error {
	return c.Body.ValidateBodyRaw()
}

func (c *Controller) EnableRequestBodyTypeSelect() {
	c.Body.EnableBodyTypeSelect()
}

func (c *Controller) DisableRequestBodyTypeSelect() {
	c.Body.DisableBodyTypeSelect()
}

func (c *Controller) SetRequestBodyType(bodyType string) {
	c.Body.SetBodyType(bodyType)
}

func (c *Controller) SetRequestBodyRaw(bodyRaw string) {
	c.Body.SetBodyRaw(bodyRaw)
}

func (c *Controller) GetRequestBodyType() string {
	return c.Body.GetBodyType()
}

func (c *Controller) ToState() (RequestState, error) {
	queryParams := c.queryParamsKeyValueCtrl.ToState()
	pathParams := c.pathParamsKeyValueCtrl.ToState()
	headers := c.headersKeyValueCtrl.ToState()
	body := c.Body.ToState()
	return RequestState{
		QueryParams: queryParams,
		PathParams:  pathParams,
		Headers:     headers,
		Body:        body,
	}, nil
}

func (c *Controller) LoadState(requestState RequestState) error {
	c.queryParamsKeyValueCtrl.LoadState(requestState.QueryParams)
	c.pathParamsKeyValueCtrl.LoadState(requestState.PathParams)
	c.headersKeyValueCtrl.LoadState(requestState.Headers)
	c.Body.LoadState(requestState.Body)
	return nil
}
