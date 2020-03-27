package client

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// URLBuilder helps build a URL for the Nirmata REST API
type URLBuilder interface {
	ToService(s Service) URLBuilder
	WithPath(path string) URLBuilder
	WithPaths(paths ...string) URLBuilder
	WithQuery(q Query) URLBuilder
	WithParameters(params map[string]string) URLBuilder
	WithMode(m OutputMode) URLBuilder
	SelectFields(fields []string) URLBuilder
	Build() string
}

type urlBldr struct {
	address string
	paths   []string
	fields  []string
	query   Query
	params  map[string]string
	mode    OutputMode
}

// NewURLBuilder creates a new URLBuilder
func NewURLBuilder(address string) URLBuilder {
	if !strings.HasSuffix(address, "/") {
		address = address + "/"
	}

	ub := &urlBldr{address: address, paths: make([]string, 0)}
	return ub
}

func (ub *urlBldr) ToService(s Service) URLBuilder {
	ub.paths = append(ub.paths, s.Name())
	ub.paths = append(ub.paths, "api")
	return ub
}

func (ub *urlBldr) WithPath(path string) URLBuilder {
	ub.paths = append(ub.paths, path)
	return ub
}

func (ub *urlBldr) WithPaths(paths ...string) URLBuilder {
	ub.paths = append(ub.paths, paths...)
	return ub
}

func (ub *urlBldr) WithQuery(q Query) URLBuilder {
	ub.query = q
	return ub
}

func (ub *urlBldr) WithParameters(params map[string]string) URLBuilder {
	ub.params = params
	return ub
}

func (ub *urlBldr) WithMode(m OutputMode) URLBuilder {
	ub.mode = m
	return ub
}

func (ub *urlBldr) SelectFields(fields []string) URLBuilder {
	ub.fields = fields
	return ub
}

func (ub *urlBldr) Build() string {
	path := path.Join(ub.paths...)
	s := fmt.Sprintf("%s%s", ub.address, path)

	hasQueryPart := false
	if ub.fields != nil {
		fields := strings.Join(ub.fields, ", ")
		fields = url.QueryEscape(fields)
		s = fmt.Sprintf("%s?fields=%s", s, fields)
		hasQueryPart = true
	}

	if ub.query != nil {
		qStr := url.QueryEscape(ub.query.String())
		if hasQueryPart {
			s = fmt.Sprintf("%s&query=%s", s, qStr)
		} else {
			s = fmt.Sprintf("%s?query=%s", s, qStr)
			hasQueryPart = true
		}
	}

	if ub.params != nil {
		pStr := convertParms(ub.params)
		if hasQueryPart {
			s = fmt.Sprintf("%s&%s", s, pStr)
		} else {
			s = fmt.Sprintf("%s?%s", s, pStr)
			hasQueryPart = true
		}
	}

	if ub.mode != OutputModeNone {
		modeString := ub.mode.String()
		if modeString != "" {
			if hasQueryPart {
				s = fmt.Sprintf("%s&mode=%s", s, modeString)
			} else {
				s = fmt.Sprintf("%s?mode=%s", s, modeString)
				hasQueryPart = true
			}
		}
	}

	return s
}

func convertParms(params map[string]string) string {
	results := ""
	for k, v := range params {
		part := url.QueryEscape(k) + "=" + url.QueryEscape(v)
		if results == "" {
			results = results + part
		} else {
			results = results + "&" + part
		}
	}

	return results
}
