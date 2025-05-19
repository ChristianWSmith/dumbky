package request

import "dumbky/internal/log"

type RequestPayload struct {
	URL      string
	Method   string
	UseSSL   bool
	Headers  map[string]string
	Params   map[string]string
	BodyRaw  string
	BodyForm map[string]string
}

type ResponsePayload struct {
	Status string
	Time   string
	Body   string
}

func SendRequest(requestPayload RequestPayload) (ResponsePayload, error) {
	log.Info("Sending request")
	status := "200 OK"
	time := "0 ms"
	body := "{}"
	return ResponsePayload{
		status,
		time,
		body,
	}, nil
}
