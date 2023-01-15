package scan

import (
	"net/http"
	"time"
)

type Options struct {
	Timeout     time.Duration
	Retries     int
	Rate        int
	Debug       bool
	ModelsRoute string
	Client      *http.Client
}
