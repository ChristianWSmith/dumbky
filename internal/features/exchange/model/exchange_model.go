package model

import (
	"dumbky/internal/constants"

	"fyne.io/fyne/v2/data/binding"
)

type modelImpl struct {
	methodBinding binding.String
	urlBinding    binding.String
	useSSLBinding binding.Bool
}

type Model interface {
	GetBindings() (method, url binding.String, ssl binding.Bool)
	GetMethod() string
	GetURL() string
	GetUseSSL() bool
	SetMethod(method string)
	SetURL(url string)
	SetUseSSL(useSSL bool)
	SetMethodListener(handler func())
}

func NewModel() Model {
	methodBinding := binding.NewString()
	methodBinding.Set(constants.HTTP_METHOD_DEFAULT)
	return &modelImpl{
		methodBinding: methodBinding,
		urlBinding:    binding.NewString(),
		useSSLBinding: binding.NewBool(),
	}
}

func (m *modelImpl) GetBindings() (method, url binding.String, ssl binding.Bool) {
	return m.methodBinding, m.urlBinding, m.useSSLBinding
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
