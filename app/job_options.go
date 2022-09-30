package app

// Runtime options for job
var JobRuntimeOptions struct {
	ForceFull        bool
	ForceDiff        bool
	SettingsFilename string
	DefaultsFrom     string
}

func init() {
	//default values
	JobRuntimeOptions.SettingsFilename = DefaultSettingsFilename
}
