package process

import (
	"os"
)

func GetHost() string {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return h
}

func GetPid() int64 {
	return int64(os.Getpid())
}
