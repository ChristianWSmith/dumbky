package httputils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mattn/go-shellwords"
)

func CurlToRequest(curlCmd string) (*http.Request, error) {
	args, err := shellwords.Parse(curlCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cURL command: %v", err)
	}

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	state := ""
	var body bytes.Buffer

	for _, arg := range args {
		switch arg {
		case "curl":
		case "-X", "--request":
			state = "method"
		case "-H", "--header":
			state = "header"
		case "-d", "--data", "--data-ascii", "--data-raw":
			state = "data"
		case "-u", "--user":
			state = "user"
		case "-I", "--head":
			req.Method = http.MethodHead
		case "-b", "--cookie":
			state = "cookie"
		default:
			switch state {
			case "method":
				req.Method = arg
				state = ""
			case "header":
				fields := parseField(arg)
				req.Header.Add(fields[0], fields[1])
				state = ""
			case "data":
				if req.Method == http.MethodGet || req.Method == http.MethodHead {
					req.Method = http.MethodPost
				}
				if _, ok := req.Header["Content-Type"]; !ok {
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				}
				if body.Len() > 0 {
					body.WriteString("&")
				}
				body.WriteString(arg)
				state = ""
			case "user":
				auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(arg))
				req.Header.Set("Authorization", auth)
				state = ""
			case "cookie":
				req.Header.Set("Cookie", arg)
				state = ""
			default:
				if isURL(arg) {
					parsedURL, err := url.ParseRequestURI(arg)
					if err != nil {
						return nil, fmt.Errorf("failed to parse URL: %v", err)
					}
					req.URL = parsedURL
				}
			}
		}
	}

	if body.Len() > 0 {
		req.Body = io.NopCloser(&body)
	}

	return req, nil
}

func RequestToCurl(req *http.Request) (string, error) {
	var curlCmd strings.Builder

	curlCmd.WriteString("curl")
	curlCmd.WriteString(" -X " + req.Method)

	for key, values := range req.Header {
		for _, value := range values {
			curlCmd.WriteString(fmt.Sprintf(" -H '%s: %s'", key, value))
		}
	}

	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read request body: %v", err)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) > 0 {
			body := string(bodyBytes)
			curlCmd.WriteString(fmt.Sprintf(" -d '%s'", body))
		}
	}

	curlCmd.WriteString(" '" + req.URL.String() + "'")

	return curlCmd.String(), nil
}

func isURL(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func parseField(arg string) []string {
	index := strings.Index(arg, ":")
	if index == -1 {
		return []string{arg, ""}
	}
	return []string{arg[0:index], strings.TrimSpace(arg[index+1:])}
}
