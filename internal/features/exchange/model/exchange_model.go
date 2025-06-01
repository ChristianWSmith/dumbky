package model

import (
	"dumbky/internal/constants"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2/data/binding"
)

type modelImpl struct {
	method          binding.String
	url             binding.String
	useSSL          binding.Bool
	requestBodyType binding.String
	requestBodyRaw  binding.String
	responseStatus  binding.String
	responseTime    binding.String
	responseBody    binding.String
}

type Bindings struct {
	Method, URL, BodyType, BodyRaw, ResponseStatus, ResponseTime, ResponseBody binding.String
	UseSSL                                                                     binding.Bool
}

type Model interface {
	GetBindings() Bindings
	GetMethod() string
	GetURL() string
	GetUseSSL() bool
	SetMethod(method string)
	SetURL(url string)
	SetUseSSL(useSSL bool)
	SetMethodListener(handler func())
	SetBodyTypeListener(handler func())
	GetRequestBodyType() string
	GetBodyRaw() string
	SetBodyType(bodyType string)
	SetBodyRaw(bodyRaw string)
	SetResponse(responseStatus, responseTime, responseBody string)
}

func NewModel() Model {
	methodBinding := binding.NewString()
	methodBinding.Set(constants.HTTP_METHOD_DEFAULT)
	bodyTypeBind := binding.NewString()
	bodyTypeBind.Set(constants.UI_BODY_TYPE_DEFAULT)
	return &modelImpl{
		method:          methodBinding,
		url:             binding.NewString(),
		useSSL:          binding.NewBool(),
		requestBodyType: bodyTypeBind,
		requestBodyRaw:  binding.NewString(),
		responseStatus:  binding.NewString(),
		responseTime:    binding.NewString(),
		responseBody:    binding.NewString(),
	}
}

func (m *modelImpl) GetBindings() Bindings {
	return Bindings{
		Method:         m.method,
		URL:            m.url,
		UseSSL:         m.useSSL,
		BodyType:       m.requestBodyType,
		BodyRaw:        m.requestBodyRaw,
		ResponseStatus: m.responseStatus,
		ResponseTime:   m.responseTime,
		ResponseBody:   m.responseBody,
	}
}

func (m *modelImpl) GetMethod() string {
	method, _ := m.method.Get()
	return method
}

func (m *modelImpl) GetURL() string {
	url, _ := m.url.Get()
	return url
}

func (m *modelImpl) GetUseSSL() bool {
	useSSL, _ := m.useSSL.Get()
	return useSSL
}

func (m *modelImpl) SetMethod(method string) {
	m.method.Set(method)
}

func (m *modelImpl) SetURL(url string) {
	m.url.Set(url)
}

func (m *modelImpl) SetUseSSL(useSSL bool) {
	m.useSSL.Set(useSSL)
}

func (m *modelImpl) SetMethodListener(handler func()) {
	m.method.AddListener(binding.NewDataListener(handler))
}

func (m *modelImpl) SetBodyTypeListener(handler func()) {
	m.requestBodyType.AddListener(binding.NewDataListener(handler))
}

func (m *modelImpl) GetRequestBodyType() string {
	bodyType, _ := m.requestBodyType.Get()
	return bodyType
}

func (m *modelImpl) GetBodyRaw() string {
	bodyRaw, _ := m.requestBodyRaw.Get()
	return bodyRaw
}

func (m *modelImpl) SetBodyType(bodyType string) {
	if !utils.ElementInSlice(constants.UIBodyTypes(), bodyType) {
		bodyType = constants.UI_BODY_TYPE_DEFAULT
	}
	m.requestBodyType.Set(bodyType)
}

func (m *modelImpl) SetBodyRaw(bodyRaw string) {
	m.requestBodyRaw.Set(bodyRaw)
}

func (m *modelImpl) SetResponse(status, time, body string) {
	m.responseStatus.Set(status)
	m.responseTime.Set(time)
	m.responseBody.Set(body)
}
