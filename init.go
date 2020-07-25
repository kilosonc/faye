package faye

import (
	"net/http"
)

const (
	BYTE int64 = 1
	KB   int64 = 1 << 10
	MB   int64 = 1 << 20
	GB   int64 = 1 << 30
)

var (
	BlockSize  int64       = 8 * MB
	Headers    http.Header = make(http.Header)
	RetryTimes             = 3
)
