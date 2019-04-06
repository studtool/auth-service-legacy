package config

import (
	"fmt"
)

func panicNotSet(name string) {
	panic(fmt.Sprintf("config: %s is required", name))
}

func panicInvalidFormat(name string, pattern string) {
	panic(fmt.Sprintf("config: %s format error; pattern - %s", name, pattern))
}
