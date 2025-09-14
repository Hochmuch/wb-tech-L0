package cache

import "time"

type Config struct {
	Addr string
	TTL  time.Duration
}
