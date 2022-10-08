package app

// Runtime options for job
var JobRuntimeOptions struct {
	ForceFull        bool // run --force-full
	ForceDiff        bool // run --force-diff
	Solid            bool // run --solid
	SettingsFilename string
	DefaultsFrom     string
}

func init() {
	//default values
	JobRuntimeOptions.SettingsFilename = DefaultSettingsFilename
}
