package request

import (
	"dumbky/internal/features"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/features/requestbody"
	"dumbky/internal/validators"

	"fyne.io/fyne/v2"
)

type controller struct {
	model *model
	view  *view

	queryParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController
	pathParamsKeyValueCtrl  keyvalueeditor.KeyValueEditorController
	headersKeyValueCtrl     keyvalueeditor.KeyValueEditorController
	requestBodyCtrl         requestbody.RequestBodyController
}

type RequestController interface {
	features.Controller
	GetQueryParamsMap() map[string]string
	GetPathParamsMap() map[string]string
	GetHeadersMap() map[string]string
	ToState() RequestState
	LoadState(requestState RequestState)
	Validate() error
	GetRequestBodyFormMap() map[string]string
	GetRequestBodyRaw() string
	GetRequestBodyType() string
	SetBodyTypeSelectEnabled(enabled bool)
	FormatBodyRaw()
}

var _ RequestController = (*controller)(nil)

func New() RequestController {
	queryParamsKeyValueCtrl := keyvalueeditor.New(validators.ValidateQueryParamKey, validators.ValidateQueryParamValue)
	pathParamsKeyValueCtrl := keyvalueeditor.New(validators.ValidatePathParamKey, validators.ValidatePathParamValue)
	headersKeyValueCtrl := keyvalueeditor.New(validators.ValidateHeaderKey, validators.ValidateHeaderValue)
	requestBodyCtrl := requestbody.New()
	return newController(queryParamsKeyValueCtrl, pathParamsKeyValueCtrl, headersKeyValueCtrl, requestBodyCtrl)
}

func newController(queryParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	pathParamsKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	headersKeyValueCtrl keyvalueeditor.KeyValueEditorController,
	requestBodyCtrl requestbody.RequestBodyController) *controller {

	return &controller{
		model:                   newModel(),
		view:                    newView(queryParamsKeyValueCtrl.CanvasObject(), pathParamsKeyValueCtrl.CanvasObject(), headersKeyValueCtrl.CanvasObject(), requestBodyCtrl.CanvasObject()),
		queryParamsKeyValueCtrl: queryParamsKeyValueCtrl,
		pathParamsKeyValueCtrl:  pathParamsKeyValueCtrl,
		headersKeyValueCtrl:     headersKeyValueCtrl,
		requestBodyCtrl:         requestBodyCtrl,
	}
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

func (c *controller) LoadState(requestState RequestState) {
	c.queryParamsKeyValueCtrl.LoadState(requestState.QueryParams)
	c.pathParamsKeyValueCtrl.LoadState(requestState.PathParams)
	c.headersKeyValueCtrl.LoadState(requestState.Headers)
	c.requestBodyCtrl.LoadState(requestState.Body)
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
	return c.requestBodyCtrl.Validate()
}

func (c *controller) GetRequestBodyFormMap() map[string]string {
	return c.requestBodyCtrl.GetBodyFormMap()
}

func (c *controller) GetRequestBodyRaw() string {
	return c.requestBodyCtrl.GetBodyRaw()
}

func (c *controller) GetRequestBodyType() string {
	return c.requestBodyCtrl.GetBodyType()
}

func (c *controller) SetBodyTypeSelectEnabled(enabled bool) {
	c.requestBodyCtrl.SetBodyTypeSelectEnabled(enabled)
}

func (c *controller) FormatBodyRaw() {
	c.requestBodyCtrl.FormatBodyRaw()
}
