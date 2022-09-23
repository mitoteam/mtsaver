package mtsaver

import "runtime"

var (
	AppName   = "mtsaver"
	Version   = "1.0.0-alpha"
	BuiltWith string
)

func init() {
	BuiltWith = runtime.Version()
}
