package requesthelper

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

type RequestConfig struct {
	URL         string
	Method      string
	UseSSL      bool
	Headers     map[string]string
	QueryParams map[string]string
	PathParams  map[string]string
	BodyType    string
	BodyRaw     string
	BodyForm    map[string]string
}

type ResponsePayload struct {
	Status string
	Time   string
	Body   string
}

func resolveBody(requestPayload RequestConfig) (*strings.Reader, error) {
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

func resolveURL(requestPayload RequestConfig) string {
	url := requestPayload.URL
	for key, value := range requestPayload.PathParams {
		url = strings.ReplaceAll(url, fmt.Sprintf(":%s:", key), value)
	}
	if len(requestPayload.QueryParams) != 0 {
		paramList := []string{}
		for key, value := range requestPayload.QueryParams {
			paramList = append(paramList, fmt.Sprintf("%s=%s", key, value))
		}
		url = url + "?" + strings.Join(paramList, "&")
	}
	if !strings.HasPrefix(url, constants.PROTOCOL_HTTP) && !strings.HasPrefix(url, constants.PROTOCOL_HTTPS) {
		if requestPayload.UseSSL {
			url = constants.PROTOCOL_HTTPS + url
		}
		url = constants.PROTOCOL_HTTP + url
	}
	return url
}

func resolveHeaders(request http.Request, headers map[string]string) {
	for key, value := range headers {
		request.Header.Set(key, value)
	}
}

func SendRequest(requestConfig RequestConfig) (ResponsePayload, error) {
	var client *http.Client

	if requestConfig.UseSSL {
		client = &http.Client{}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	body, err := resolveBody(requestConfig)
	if err != nil {
		log.Error(err)
		return ResponsePayload{}, err
	}

	url := resolveURL(requestConfig)

	request, err := http.NewRequest(requestConfig.Method, url, body)
	if err != nil {
		log.Error(err)
		return ResponsePayload{}, err
	}
	resolveHeaders(*request, requestConfig.Headers)

	log.Info("Sending request")
	start := time.Now()
	response, err := client.Do(request)
	elapsed := time.Since(start)

	if err != nil {
		log.Warn(err)
		return ResponsePayload{}, err
	} else {
		log.Info("Request sent")
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return ResponsePayload{}, err
	}

	return ResponsePayload{
		response.Status,
		elapsed.String(),
		string(responseBody),
	}, nil
}
