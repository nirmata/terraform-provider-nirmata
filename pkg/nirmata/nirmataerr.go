package nirmata

import (
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	client "github.com/nirmata/go-client/pkg/client"
)

// Returns true if the error matches all these conditions:
//  * err is of type client.Error
//  * Error.Code() matches code
//  * Error.Message() contains message
func isNirmataErr(err error, code string, message string) bool {
	var nirmataErr client.Error
	if errors.As(err, &nirmataErr) {
		return nirmataErr.Code() == code && strings.Contains(nirmataErr.Message(), message)
	}
	return false
}

// Returns true if the error matches all these conditions:
//  * err is of type client.RequestFailure
//  * RequestFailure.StatusCode() matches status code
// that sometimes only respond with status codes.
func isNirmataErrRequestFailureStatusCode(err error, statusCode int) bool {
	var nirmataErr client.RequestFailure
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
			nirmataErr, ok := err.(client.Error)
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
			clientErr, ok := err.(client.Error)
			if ok {
				for _, code := range codes {
					if clientErr.Code() == code {
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
