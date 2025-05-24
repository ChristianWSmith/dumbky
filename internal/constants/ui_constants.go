package constants

const UI_PLACEHOLDER_URL string = "https://www.example.com/"
const UI_PLACEHOLDER_BODY_TYPE_RAW string = "{}"
const UI_PLACEHOLDER_RESPONSE_STATUS string = ""
const UI_PLACEHOLDER_RESPONSE_TIME string = ""
const UI_PLACEHOLDER_RESPONSE_BODY string = ""
const UI_PLACEHOLDER_KEY string = ""
const UI_PLACEHOLDER_VALUE string = ""
const UI_PLACEHOLDER_UNTITLED string = "untitled"

const UI_LOADING_RESPONSE_STATUS string = ""
const UI_LOADING_RESPONSE_TIME string = ""
const UI_LOADING_RESPONSE_BODY string = ""

const UI_LABEL_SSL string = "SSL"
const UI_LABEL_SEND string = "Send"
const UI_LABEL_QUERY_PARAMETERS string = "Query Params"
const UI_LABEL_PATH_PARAMETERS string = "Path Params"
const UI_LABEL_HEADERS string = "Headers"
const UI_LABEL_BODY string = "Body"
const UI_LABEL_KEY_VALUE_ADD string = ""
const UI_LABEL_KEY_VALUE_ENABLE string = ""
const UI_LABEL_COLLECTIONS string = "Collections"

const UI_BODY_TYPE_RAW string = "Raw"
const UI_BODY_TYPE_FORM string = "Form"
const UI_BODY_TYPE_NONE string = "None"
const UI_BODY_TYPE_DEFAULT string = UI_BODY_TYPE_NONE

const UI_DASHBOARD_SIDEBAR_OFFSET float64 = 0.2

func UIBodyTypes() []string {
	return []string{
		UI_BODY_TYPE_RAW,
		UI_BODY_TYPE_FORM,
		UI_BODY_TYPE_NONE,
	}
}
