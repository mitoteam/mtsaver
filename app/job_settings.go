package app

import (
	"github.com/mitoteam/mttools"
)

const DefaultSettingsFilename = ".mtsaver.yml"

// Setting for archived folder
type JobSettings struct {
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

	KeepEmptyDiff bool `yaml:"keep_empty_diff" yaml_comment:"false = delete diff archives if no files where added to it (empty archive), true = keep anyway"`

	KeepSameDiff bool `yaml:"keep_same_diff" yaml_comment:"false = delete diff archives if it has same sha256 hash as previous one (nothing new added), true = keep anyway"`

	//Maximum total diffs size IN PERCENTS to force full backup
	MaxTotalDiffSizePercent int `yaml:"max_total_diff_size_percent" yaml_comment:"Total size of differential archives since latest full archive in percents to force new full archive next run, 0 = not set"`
}

// creates new settings with default values
func NewJobSettings() JobSettings {
	return JobSettings{
		CompressionLevel: -1,
		MaxFullCount:     5,
		MaxDiffCount:     20,
		KeepEmptyDiff:    false,
		KeepSameDiff:     false,
	}
}

func (js *JobSettings) LoadFromFile(path string) error {
	return mttools.LoadYamlSettingFromFile(path, js)
}

func (js *JobSettings) SaveToFile(path string, comment string) error {
	return mttools.SaveYamlSettingToFile(path, Global.AppName+" directory settings file", js)
}

func (js *JobSettings) Print() {
	mttools.PrintYamlSettings(js)
}
