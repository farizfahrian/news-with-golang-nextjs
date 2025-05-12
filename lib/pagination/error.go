package pagination

import "errors"

var (
	ErrorMaxPage     = errors.New("choosen page is greater than max page")
	ErrorPage        = errors.New("page must greater than 0")
	ErrorPageEmpty   = errors.New("page cannot be empty")
	ErrorPageInvalid = errors.New("page is invalid, must be number")
)
