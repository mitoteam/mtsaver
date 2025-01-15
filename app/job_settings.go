package app

import (
	"log"
	"path/filepath"

	"github.com/mitoteam/mttools"
)

const DefaultSettingsFilename = ".mtsaver.yml"

// Setting for archived folder
type JobSettings struct {
	LoadedFromFile bool `yaml:"-"` //ignored in yaml

	ArchivesPath string `yaml:"archives_path" yaml_comment:"Full path to directory to create archives in"`
	ArchiveName  string `yaml:"archive_name" yaml_comment:"Base archive name (appended with timestamp, suffix and .7z extension)"`
	FullSuffix   string `yaml:"full_suffix" yaml_comment:"Suffix for full archives"`
	DiffSuffix   string `yaml:"diff_suffix" yaml_comment:"Suffix for differential archives"`
	DateFormat   string `yaml:"date_format" yaml_comment:"Archive filename timestamp format. Don't touch it if you don't understand! Golang's time formatting is a bit crazy https://mttm.ml/go-time-format"`

	CompressionLevel int    `yaml:"compression_level" yaml_comment:"7-zip compression level from 0 to 9. Default: 5"`
	Password         string `yaml:"password" yaml_comment:"Set this to protect .7z file with password."`
	EncryptFilenames bool   `yaml:"encrypt_filenames" yaml_comment:"Encrypt filenames in .7z archive (used only when 'password' is set)."`

	//Create solid archives
	Solid bool `yaml_comment:"Create solid 7-zip archives"`

	//List of patterns to exclude from archive
	Exclude []string `yaml_comment:"List of patterns to exclude from archive"`

	//List of patterns to be added to archive without compression (works for FULL backups only)
	SkipCompression []string `yaml:"skip_compression" yaml_comment:"List of patterns for fast adding to archive without compression"`

	//Run cleanup procedure before or after archive creation (default: after)
	Cleanup string `yaml_comment:"Cleanup old archives before or after archiving (before|after, default: after)"`

	MaxFullCount int `yaml:"max_full_count" yaml_comment:"Maximum count of full archives to keep"`
	KeepAtLeast  int `yaml:"keep_at_least" yaml_comment:"Do not remove full archives if they younger than this count of days"`

	//Maximum number of diff archives to have after full backup
	MaxDiffCount int `yaml:"max_diff_count" yaml_comment:"Maximum count of differential archives to create before creating new full archive"`

	//Maximum latest diff size IN PERCENTS to force full backup
	MaxDiffSizePercent int `yaml:"max_diff_size_percent" yaml_comment:"Last differential archive size in percents to force new full archive at next run, 0 = not set"`

	//Maximum total diffs size IN PERCENTS to force full backup
	MaxTotalDiffSizePercent int `yaml:"max_total_diff_size_percent" yaml_comment:"Total size of differential archives since latest full archive in percents to force new full archive next run, 0 = not set"`

	KeepEmptyDiff bool `yaml:"keep_empty_diff" yaml_comment:"false = delete diff archives if no files where added to it (empty archive), true = keep anyway"`

	KeepSameDiff bool `yaml:"keep_same_diff" yaml_comment:"false = delete diff archives if it has same sha256 hash as previous one (nothing new added), true = keep anyway"`

	// Commands to run
	RunBefore []string `yaml:"run_before" yaml_comment:"List of commands to run before creating archive"`

	// Log file name
	LogFilename      string `yaml:"log_filename" yaml_comment:"Name of file to add log messages to."`
	LogFormat        string `yaml:"log_format" yaml_comment:"Log file format: text|json|no. Default: text. 'no' = disable logging."`
	LogCommandOutput bool   `yaml:"log_command_output" yaml_comment:"Add commands (from run_before) and 7-Zip output to log file."`
	LogMaxSize       int64  `yaml:"log_max_size" yaml_comment:"Log file size for it to be rotated. Default: 1Mb."`
}

// creates new settings with default values
func NewJobSettings() JobSettings {
	return JobSettings{
		LoadedFromFile:     false,
		CompressionLevel:   -1,
		MaxFullCount:       5,
		MaxDiffCount:       20,
		SkipCompression:    []string{"*.7z", "*.rar"},
		MaxDiffSizePercent: 120,
		KeepEmptyDiff:      false,
		KeepSameDiff:       false,
		LogFilename:        "_mtsaver.log",
		LogFormat:          "text",
		LogCommandOutput:   false,
		LogMaxSize:         1024 * 1024, // 1Mb
	}
}

func (js *JobSettings) LoadFromFile(path string) error {
	if err := mttools.LoadYamlSettingFromFile(path, js); err != nil {
		return err
	}

	js.LoadedFromFile = true

	return nil
}

func (js *JobSettings) SaveToFile(path string, comment string) error {
	return mttools.SaveYamlSettingToFile(path, Global.AppName+" directory settings file", js)
}

func (js *JobSettings) Print() {
	mttools.PrintYamlSettings(js)
}

func (js *JobSettings) ApplyDefaultsAndCheck(job_path string) {
	//// Set defaults for missing values
	if js.DateFormat == "" {
		js.DateFormat = "2006-01-02_15-04-05"
	}

	name := filepath.Base(job_path)

	if len(js.ArchiveName) == 0 {
		js.ArchiveName = name
	}

	if len(js.ArchivesPath) == 0 {
		js.ArchivesPath = filepath.Join(filepath.Dir(job_path), name+"_ARCHIVE")
	}

	if len(js.FullSuffix) == 0 {
		js.FullSuffix = "FULL"
	}

	if len(js.DiffSuffix) == 0 {
		js.DiffSuffix = "DIFF"
	}

	if js.CompressionLevel == -1 {
		js.CompressionLevel = 5
	}

	//--------------------------------------
	// Override values from runtime options
	//--------------------------------------

	//turn on solid mode for archives
	if JobRuntimeOptions.Solid {
		js.Solid = true
	}

	//password
	if len(JobRuntimeOptions.Password) > 0 {
		js.Password = JobRuntimeOptions.Password
	}

	//encrypt filenames
	if JobRuntimeOptions.EncryptFilenames {
		js.EncryptFilenames = true
	}

	//skip logging
	if JobRuntimeOptions.NoLog {
		js.LogFormat = "no"
	}

	//--------------------
	// Do settings checks
	//--------------------

	if js.FullSuffix == js.DiffSuffix {
		log.Fatalln("Full suffix should differ from diff suffix")
	}

	if js.Cleanup == "" {
		js.Cleanup = "after"
	} else if js.Cleanup != "before" && js.Cleanup != "after" {
		log.Fatalln("Valid  values for 'cleanup' option are 'before', 'after'")
	}

	if js.MaxFullCount < 1 {
		log.Fatalln("Minimum value for max_full_count is 1")
	}

	if js.MaxDiffCount < 0 {
		log.Fatalln("Minimum value for max_diff_count is 0")
	}

	if js.LogFormat != "no" && js.LogFormat != "text" && js.LogFormat != "json" {
		log.Fatalf("Wrong log format: %s\n", js.LogFormat)
	}

	if js.LogMaxSize < 10240 {
		log.Fatalln("Minimum value for log_size_max is 10240 (10kb)")
	}
}
