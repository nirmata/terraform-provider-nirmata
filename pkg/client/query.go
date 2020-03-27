package client

import (
	"fmt"
	"strings"
)

// Query provides a way to filter which objects are matched
type Query interface {
	FieldEqualsValue(name, val string) Query
	String() string
}

// QueryOperator provides boolean and arithmetic operators for field-value pairs
type QueryOperator int

const (
	// QueryOperatorEquals ...
	QueryOperatorEquals QueryOperator = iota + 1
)

// NewQuery creates a new Query
func NewQuery() Query {
	return &query{queryParts: make([]queryPart, 0)}
}

type queryPart struct {
	fieldName string
	op        QueryOperator
	value     string
}

type query struct {
	queryParts []queryPart
}

func (q *query) FieldEqualsValue(name, val string) Query {
	q.queryParts = append(q.queryParts, queryPart{name, QueryOperatorEquals, val})
	return q
}

func (q *query) String() string {
	count := len(q.queryParts)
	if count == 0 {
		return ""
	}

	var sb strings.Builder
	for i, qp := range q.queryParts {
		if i == 0 {
			sb.WriteString("{")
		}

		s := fmt.Sprintf("\"%s\" : \"%s\"", qp.fieldName, qp.value)
		sb.WriteString(s)

		if i == (count - 1) {
			sb.WriteString("}")
		} else {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
