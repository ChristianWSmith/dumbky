package exchange

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/request"

	"fyne.io/fyne/v2/data/binding"
)

type model struct {
	methodBinding binding.String
	urlBinding    binding.String
	useSSLBinding binding.Bool
}

func newModel() *model {
	methodBinding := binding.NewString()
	methodBinding.Set(constants.HTTP_METHOD_DEFAULT)
	return &model{
		methodBinding: methodBinding,
		urlBinding:    binding.NewString(),
		useSSLBinding: binding.NewBool(),
	}
}

type ExchangeState struct {
	Method  string               `json:"method"`
	URL     string               `json:"url"`
	UseSSL  bool                 `json:"ssl"`
	Request request.RequestState `json:"request"`
}

func (m *model) getMethod() string {
	method, _ := m.methodBinding.Get()
	return method
}

func (m *model) getURL() string {
	url, _ := m.urlBinding.Get()
	return url
}

func (m *model) getUseSSL() bool {
	useSSL, _ := m.useSSLBinding.Get()
	return useSSL
}

func (m *model) setMethod(method string) {
	m.methodBinding.Set(method)
}

func (m *model) setURL(url string) {
	m.urlBinding.Set(url)
}

func (m *model) setUseSSL(useSSL bool) {
	m.useSSLBinding.Set(useSSL)
}

func (m *model) setMethodListener(handler func()) {
	m.methodBinding.AddListener(binding.NewDataListener(handler))
}
