package nirmata

import (
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	nirmataerr "github.com/nirmata/go-client/pkg/nirmataErr"
)

// Returns true if the error matches all these conditions:
//  * err is of type nirmataerr.Error
//  * Error.Code() matches code
//  * Error.Message() contains message
func isNirmataErr(err error, code string, message string) bool {
	var nirmataErr nirmataerr.Error
	if errors.As(err, &nirmataErr) {
		return nirmataErr.Code() == code && strings.Contains(nirmataErr.Message(), message)
	}
	return false
}

// Returns true if the error matches all these conditions:
//  * err is of type nirmataerr.RequestFailure
//  * RequestFailure.StatusCode() matches status code
// that sometimes only respond with status codes.
func isNirmataErrRequestFailureStatusCode(err error, statusCode int) bool {
	var nirmataErr nirmataerr.RequestFailure
	if errors.As(err, &nirmataErr) {
		return nirmataErr.StatusCode() == statusCode
	}
	return false
}

func retryOnNirmataCode(code string, f func() (interface{}, error)) (interface{}, error) {
	var resp interface{}
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = f()
		if err != nil {
			nirmataErr, ok := err.(nirmataerr.Error)
			if ok && nirmataErr.Code() == code {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return resp, err
}

// RetryOnNirmataCodes retries Nirmata error codes for one minute
// Note: This function will be moved out of the nirmata package in the future.
func RetryOnNirmataCodes(codes []string, f func() (interface{}, error)) (interface{}, error) {
	var resp interface{}
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = f()
		if err != nil {
			nirmataErr, ok := err.(nirmataerr.Error)
			if ok {
				for _, code := range codes {
					if nirmataErr.Code() == code {
						return resource.RetryableError(err)
					}
				}
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return resp, err
}
