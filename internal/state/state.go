package state

type DocumentState struct {
	CollectionName string        `json:"collection_name"`
	RequestName    string        `json:"request_name"`
	ExchangeState  ExchangeState `json:"exchange"`
}

type ExchangeState struct {
	Method      string              `json:"method"`
	URL         string              `json:"url"`
	UseSSL      bool                `json:"ssl"`
	QueryParams KeyValueEditorState `json:"queryParams"`
	PathParams  KeyValueEditorState `json:"pathParams"`
	Headers     KeyValueEditorState `json:"headers"`
	BodyType    string              `json:"bodyType"`
	BodyForm    KeyValueEditorState `json:"bodyForm"`
	BodyRaw     string              `json:"bodyRaw"`
	// TODO: save response details in exchange state
}

type KeyValueEditorState struct {
	KeyValueStates []KeyValueState `json:"keyValueStates"`
}

type KeyValueState struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}
