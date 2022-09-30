package app

// Runtime options for job
var JobRuntimeOptions struct {
	ForceFull        bool
	ForceDiff        bool
	ForceCleanup     bool
	SettingsFilename string
}

func init() {
	//default values
	JobRuntimeOptions.SettingsFilename = DefaultSettingsFilename
}
