package rft

import (
	"path/filepath"
	"reflect"
	"runtime"
)

func FuncName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return filepath.Base(name)
}

func StructName(s interface{}) string {
	name := reflect.TypeOf(s).String()
	if reflect.ValueOf(s).Type().Kind() == reflect.Ptr {
		return name[1:]
	}
	return name
}
