package error_play

import (
	"errors"
	"fmt"
)

// common HTTP status codes
var ErrNotFoundHTTPCode = errors.New("404")
var ErrUnauthorizedHTTPCode = errors.New("401")

// database errors
var ErrRecordNotFoundErr = errors.New("DB: record not found")
var ErrAffectedRecordsMismatchErr = errors.New("DB: affected records mismatch")

// HTTP client errors
var ErrResourceNotFoundErr = errors.New("HTTP client: resource not found")
var ErrResourceUnauthorizedErr = errors.New("HTTP client: unauthorized")

// application errors (the new feature)
var ErrUserNotFoundErr = fmt.Errorf("user not found: %w (%w)",
	ErrRecordNotFoundErr, ErrNotFoundHTTPCode)
var ErrOtherResourceUnauthorizedErr = fmt.Errorf("unauthorized call: %w (%w)",
	ErrResourceUnauthorizedErr, ErrUnauthorizedHTTPCode)

func handleError(err error) {
	if errors.Is(err, ErrNotFoundHTTPCode) {
		fmt.Println("Will return 404")
	} else if errors.Is(err, ErrUnauthorizedHTTPCode) {
		fmt.Println("Will return 401")
	} else {
		fmt.Println("Will return 500")
	}
	fmt.Println(err.Error())
}

func Play_Organize_HTTP_Error() {
	handleError(ErrUserNotFoundErr)
	handleError(ErrOtherResourceUnauthorizedErr)
}
