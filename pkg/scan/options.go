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
	PrimaryTag  string
	Client      *http.Client
}
