package app

import "errors"

type Error error

var (
	ERR_TIMEOUT Error = errors.New("Timed out")
)
