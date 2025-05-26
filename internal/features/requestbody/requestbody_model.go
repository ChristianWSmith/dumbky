package requestbody

import (
	"dumbky/internal/constants"
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2/data/binding"
)

type model struct {
	bodyTypeBinding binding.String
	bodyRawBinding  binding.String
}

type RequestBodyState struct {
	BodyType string                             `json:"bodyType"`
	BodyForm keyvalueeditor.KeyValueEditorState `json:"bodyForm"`
	BodyRaw  string                             `json:"bodyRaw"`
}

func newModel() *model {
	bodyTypeBind := binding.NewString()
	bodyTypeBind.Set(constants.UI_BODY_TYPE_DEFAULT)
	return &model{
		bodyTypeBinding: bodyTypeBind,
		bodyRawBinding:  binding.NewString(),
	}
}

func (m *model) setBodyTypeListener(handler func()) {
	m.bodyTypeBinding.AddListener(binding.NewDataListener(handler))
}

func (m *model) getBodyType() string {
	bodyType, _ := m.bodyTypeBinding.Get()
	return bodyType
}

func (m *model) getBodyRaw() string {
	bodyRaw, _ := m.bodyRawBinding.Get()
	return bodyRaw
}

func (m *model) setBodyType(bodyType string) {
	if !utils.ElementInSlice(constants.UIBodyTypes(), bodyType) {
		bodyType = constants.UI_BODY_TYPE_DEFAULT
	}
	m.bodyTypeBinding.Set(bodyType)
}

func (m *model) setBodyRaw(bodyRaw string) {
	m.bodyRawBinding.Set(bodyRaw)
}
