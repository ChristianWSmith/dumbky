package model

import (
	"dumbky/internal/constants"
	"dumbky/internal/utils"

	"fyne.io/fyne/v2/data/binding"
)

type modelImpl struct {
	bodyTypeBinding binding.String
	bodyRawBinding  binding.String
}

type Bindings struct {
	BodyType, BodyRaw binding.String
}

type Model interface {
	GetBindings() Bindings
	SetBodyTypeListener(handler func())
	GetBodyType() string
	GetBodyRaw() string
	SetBodyType(bodyType string)
	SetBodyRaw(bodyRaw string)
}

func NewModel() Model {
	bodyTypeBind := binding.NewString()
	bodyTypeBind.Set(constants.UI_BODY_TYPE_DEFAULT)
	return &modelImpl{
		bodyTypeBinding: bodyTypeBind,
		bodyRawBinding:  binding.NewString(),
	}
}

func (m *modelImpl) SetBodyTypeListener(handler func()) {
	m.bodyTypeBinding.AddListener(binding.NewDataListener(handler))
}

func (m *modelImpl) GetBindings() Bindings {
	return Bindings{
		BodyType: m.bodyTypeBinding,
		BodyRaw:  m.bodyRawBinding,
	}
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
