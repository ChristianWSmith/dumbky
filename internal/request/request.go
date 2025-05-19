package request

import (
	"crypto/tls"
	"dumbky/internal/constants"
	"dumbky/internal/log"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RequestPayload struct {
	URL      string
	Method   string
	UseSSL   bool
	Headers  map[string]string
	Params   map[string]string
	BodyType string
	BodyRaw  string
	BodyForm map[string]string
}

type ResponsePayload struct {
	Status string
	Time   string
	Body   string
}

func resolveBody(requestPayload RequestPayload) (*strings.Reader, error) {
	if requestPayload.BodyType == constants.UI_BODY_TYPE_FORM {
		pairs := []string{}
		for key, value := range requestPayload.BodyForm {
			pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
		}
		body := strings.Join(pairs, "&")
		return strings.NewReader(body), nil
	} else if requestPayload.BodyType == constants.UI_BODY_TYPE_RAW {
		return strings.NewReader(requestPayload.BodyRaw), nil
	} else if requestPayload.BodyType == constants.UI_BODY_TYPE_NONE {
		return strings.NewReader(""), nil
	}
	return strings.NewReader(""), errors.New("invalid body type")
}

func resolveURL(url string, params map[string]string) string {
	if len(params) != 0 {
		paramList := []string{}
		for key, value := range params {
			paramList = append(paramList, fmt.Sprintf("%s=%s", key, value))
		}
		url = url + "?" + strings.Join(paramList, "&")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

func resolveHeaders(request http.Request, headers map[string]string) {
	for key, value := range headers {
		request.Header.Set(key, value)
	}
}

func SendRequest(requestPayload RequestPayload) (ResponsePayload, error) {
	var client *http.Client

	if requestPayload.UseSSL {
		client = &http.Client{}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	body, err := resolveBody(requestPayload)
	if err != nil {
		log.Error("Failed to resolve request body", err.Error())
		return ResponsePayload{}, err
	}

	url := resolveURL(requestPayload.URL, requestPayload.Params)

	request, err := http.NewRequest(requestPayload.Method, url, body)
	if err != nil {
		log.Error("Failed to create HTTP Request", err.Error())
		return ResponsePayload{}, err
	}
	resolveHeaders(*request, requestPayload.Headers)

	log.Info("Sending request")
	start := time.Now()
	response, err := client.Do(request)
	elapsed := time.Since(start)
	log.Info("Request sent")

	if err != nil {
		log.Warn("Failed to send HTTP Request", err.Error())
		return ResponsePayload{}, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Failed to read response body", err.Error())
		return ResponsePayload{}, err
	}

	return ResponsePayload{
		response.Status,
		elapsed.String(),
		string(responseBody),
	}, nil
}
