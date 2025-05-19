package constants

const HTTP_METHOD_GET string = "GET"
const HTTP_METHOD_PUT string = "PUT"
const HTTP_METHOD_DELETE string = "DELETE"
const HTTP_METHOD_POST string = "POST"
const HTTP_METHOD_HEAD string = "HEAD"
const HTTP_METHOD_PATCH string = "PATCH"
const HTTP_METHOD_OPTIONS string = "OPTIONS"

func HttpMethods() []string {
	return []string{
		HTTP_METHOD_GET, 
		HTTP_METHOD_PUT, 
		HTTP_METHOD_DELETE, 
		HTTP_METHOD_POST, 
		HTTP_METHOD_HEAD, 
		HTTP_METHOD_PATCH, 
		HTTP_METHOD_OPTIONS,
	}
}