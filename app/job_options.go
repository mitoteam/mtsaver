package app

// Runtime options for job
var JobRuntimeOptions struct {
	SettingsFilename string
	NoConsole        bool // global: --no-console

	ForceFull        bool   // run --force-full
	ForceDiff        bool   // run --force-diff
	Solid            bool   // run --solid
	Password         string // run/restore --password <string>
	EncryptFilenames bool   // run --encrypt-filenames
	NoLog            bool   // run --no-log

	DefaultsFrom string // init --defaults-from <string>
	Print        bool   // init --print

	To     bool // restore --to
	Latest bool // restore --latest
}

func init() {
	//default values
	JobRuntimeOptions.SettingsFilename = DefaultSettingsFilename
}
