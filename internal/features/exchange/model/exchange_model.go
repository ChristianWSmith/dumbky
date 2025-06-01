package model

import (
	"dumbky/internal/constants"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2/data/binding"
)

type modelImpl struct {
	methodBinding   binding.String
	urlBinding      binding.String
	useSSLBinding   binding.Bool
	bodyTypeBinding binding.String
	bodyRawBinding  binding.String
	responseStatus  binding.String
	responseTime    binding.String
	responseBody    binding.String
}

type Bindings struct {
	Method, URL                                binding.String
	UseSSL                                     binding.Bool
	BodyType, BodyRaw                          binding.String
	ResponseStatus, ResponseTime, ResponseBody binding.String
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
	GetBodyType() string
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
		methodBinding:   methodBinding,
		urlBinding:      binding.NewString(),
		useSSLBinding:   binding.NewBool(),
		bodyTypeBinding: bodyTypeBind,
		bodyRawBinding:  binding.NewString(),
		responseStatus:  binding.NewString(),
		responseTime:    binding.NewString(),
		responseBody:    binding.NewString(),
	}
}

func (m *modelImpl) GetBindings() Bindings {
	return Bindings{
		Method:         m.methodBinding,
		URL:            m.urlBinding,
		UseSSL:         m.useSSLBinding,
		BodyType:       m.bodyTypeBinding,
		BodyRaw:        m.bodyRawBinding,
		ResponseStatus: m.responseStatus,
		ResponseTime:   m.responseTime,
		ResponseBody:   m.responseBody,
	}
}

func (m *modelImpl) GetMethod() string {
	method, _ := m.methodBinding.Get()
	return method
}

func (m *modelImpl) GetURL() string {
	url, _ := m.urlBinding.Get()
	return url
}

func (m *modelImpl) GetUseSSL() bool {
	useSSL, _ := m.useSSLBinding.Get()
	return useSSL
}

func (m *modelImpl) SetMethod(method string) {
	m.methodBinding.Set(method)
}

func (m *modelImpl) SetURL(url string) {
	m.urlBinding.Set(url)
}

func (m *modelImpl) SetUseSSL(useSSL bool) {
	m.useSSLBinding.Set(useSSL)
}

func (m *modelImpl) SetMethodListener(handler func()) {
	m.methodBinding.AddListener(binding.NewDataListener(handler))
}

func (m *modelImpl) SetBodyTypeListener(handler func()) {
	m.bodyTypeBinding.AddListener(binding.NewDataListener(handler))
}

func (m *modelImpl) GetBodyType() string {
	bodyType, _ := m.bodyTypeBinding.Get()
	return bodyType
}

func (m *modelImpl) GetBodyRaw() string {
	bodyRaw, _ := m.bodyRawBinding.Get()
	return bodyRaw
}

func (m *modelImpl) SetBodyType(bodyType string) {
	if !utils.ElementInSlice(constants.UIBodyTypes(), bodyType) {
		bodyType = constants.UI_BODY_TYPE_DEFAULT
	}
	m.bodyTypeBinding.Set(bodyType)
}

func (m *modelImpl) SetBodyRaw(bodyRaw string) {
	m.bodyRawBinding.Set(bodyRaw)
}

func (m *modelImpl) SetResponse(status, time, body string) {
	m.responseStatus.Set(status)
	m.responseTime.Set(time)
	m.responseBody.Set(body)
}
