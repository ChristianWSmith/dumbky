package constants

const UI_PLACEHOLDER_URL string = "https://www.example.com/"
const UI_PLACEHOLDER_BODY_TYPE_RAW string = "{}"
const UI_PLACEHOLDER_RESPONSE_STATUS string = "<response status>"
const UI_PLACEHOLDER_RESPONSE_TIME string = "<response time>"
const UI_PLACEHOLDER_RESPONSE_BODY string = "<response body>"
const UI_PLACEHOLDER_KEY string = "<key>"
const UI_PLACEHOLDER_VALUE string = "<value>"

const UI_LABEL_SSL string = "SSL"
const UI_LABEL_SEND string = "SEND"
const UI_LABEL_PARAMETERS string = "PARAMETERS"
const UI_LABEL_HEADERS string = "HEADERS"
const UI_LABEL_BODY string = "BODY"
const UI_LABEL_ADD string = "+"
const UI_LABEL_REMOVE string = "-"

const UI_BODY_TYPE_RAW string = "RAW"
const UI_BODY_TYPE_FORM string = "FORM"
const UI_BODY_TYPE_NONE string = "NONE"

func UIBodyTypes() []string {
	return []string{
		UI_BODY_TYPE_RAW,
		UI_BODY_TYPE_FORM,
		UI_BODY_TYPE_NONE,
	}
}