package keyvalueeditor

import (
	"dumbky/internal/features"
	"dumbky/internal/features/keyvalueeditor/controller"
	"dumbky/internal/state"
)

type KeyValueEditorController interface {
	features.Controller
	ToState() state.KeyValueEditorState
	LoadState(keyValueEditorState state.KeyValueEditorState)
	SetVisible(visible bool)
	Validate() error
	Get() map[string]string
}

func New(keyValidator, valueValidator func(string) error) KeyValueEditorController {
	return controller.NewController(keyValidator, valueValidator)
}
