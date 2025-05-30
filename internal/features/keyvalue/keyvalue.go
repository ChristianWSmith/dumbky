package keyvalue

import (
	"dumbky/internal/features"
	"dumbky/internal/features/keyvalue/controller"
	"dumbky/internal/state"
)

type KeyValueController interface {
	features.Controller
	Get() (key, value string)
	Validate() error
	IsEnabled() bool
	SetDestroyHandler(handler func())
	LoadState(state state.KeyValueState)
	ToState() state.KeyValueState
}

func New(keyValidator, valueValidator func(val string) error) KeyValueController {
	return controller.NewController(keyValidator, valueValidator)
}
