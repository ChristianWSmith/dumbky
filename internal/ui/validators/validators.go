package validators

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

func ValidateNil(val string) error {
	return nil
}

func ValidateHeaderKey(s string) error {
	if s == "" {
		return errors.New("header key cannot be empty")
	}
	for _, r := range s {
		if r <= 32 || r == 127 || strings.ContainsRune("()<>@,;:\\\"/[]?={} \t", r) {
			return errors.New("invalid character in header key")
		}
	}
	return nil
}

func ValidateHeaderValue(s string) error {
	for _, r := range s {
		if r < 32 && r != 9 {
			return errors.New("control character in header value")
		}
		if r == 127 {
			return errors.New("del character in header value")
		}
	}
	return nil
}

func ValidateURL(s string) error {
	if s == "" {
		return errors.New("url cannot be empty")
	}
	_, err := url.Parse(s)
	if err != nil {
		return err
	}
	return nil
}

func ValidateQueryParamKey(s string) error {
	if _, err := url.QueryUnescape(s); err != nil {
		return err
	}
	return nil
}

func ValidateQueryParamValue(s string) error {
	if _, err := url.QueryUnescape(s); err != nil {
		return err
	}
	return nil
}

func ValidatePathParamKey(s string) error {
	match, err := regexp.MatchString("^[a-zA-Z_]+$", s)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("invalid path param key")
	}
	return nil
}

func ValidateCollectionName(s string) error {
	match, err := regexp.MatchString("^[a-zA-Z_]+$", s)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("invalid collection name")
	}
	return nil
}

func ValidatePathParamValue(s string) error {
	if _, err := url.QueryUnescape(s); err != nil {
		return err
	}
	return nil
}

func ValidateFormBodyKey(s string) error {
	if _, err := url.ParseQuery(s + "=x"); err != nil {
		return err
	}
	return nil
}

func ValidateFormBodyValue(s string) error {
	if _, err := url.ParseQuery("x=" + s); err != nil {
		return err
	}
	return nil
}

func ValidateRawBodyContent(s string) error {
	if !utf8.ValidString(s) {
		return errors.New("invalid UTF-8 in raw body")
	}
	return nil
}
