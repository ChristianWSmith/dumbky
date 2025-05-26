package request

import (
	"dumbky/internal/features/keyvalueeditor"
	"dumbky/internal/features/requestbody"
)

type model struct {
}

func newModel() *model {
	return &model{}
}

type RequestState struct {
	QueryParams keyvalueeditor.KeyValueEditorState `json:"queryParams"`
	PathParams  keyvalueeditor.KeyValueEditorState `json:"pathParams"`
	Headers     keyvalueeditor.KeyValueEditorState `json:"headers"`
	Body        requestbody.RequestBodyState       `json:"body"`
}
