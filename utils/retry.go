package utils

import (
	"time"
)

type RetryFunc func() error

func Retry(f RetryFunc, n int, d time.Duration) (err error) {
	t := time.NewTicker(d)
	for i := 0; i < n; i++ {
		if err = f(); err == nil {
			return nil
		}
		<-t.C
	}
	return err
}
