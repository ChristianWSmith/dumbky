package keyvalueeditor

import "dumbky/internal/features/keyvalue"

type model struct {
}

func newModel() *model {
	return &model{}
}

type KeyValueEditorState struct {
	KeyValueStates []keyvalue.KeyValueState `json:"keyValueStates"`
}
