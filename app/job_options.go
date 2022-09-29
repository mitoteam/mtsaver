package app

// Runtime options for job
var JobRuntimeOptions struct {
	ForceFull        bool
	ForceDiff        bool
	SettingsFilename string
}

func init() {
	JobRuntimeOptions.SettingsFilename = DefaultSettingsFilename
}
