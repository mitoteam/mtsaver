package app

import (
	"github.com/mitoteam/mttools"
)

const DefaultSettingsFilename = ".mtsaver.yml"

// Setting for archived folder
type JobSettings struct {
	ArchivesPath string `yaml:"archives_path" yaml_comment:"full path to directory to create archives in"`
	ArchiveName  string `yaml:"archive_name" yaml_comment:"base archive name (appended with timestamp and suffix)"`
	FullSuffix   string `yaml:"full_suffix" yaml_comment:"suffix for full archives"`
	DiffSuffix   string `yaml:"diff_suffix" yaml_comment:"suffix for differential archives"`
	DateFormat   string `yaml:"date_format" yaml_comment:"archive filename timestamp format. Don't touch it! Go's time formatting is craziness https://mttm.ml/go-time-format"`

	//Int value from 0 = do not compress to 9 = max compression, longest time
	CompressionLevel int `yaml:"compression_level" yaml_comment:"7-Zip compression level from 0 to 9"`

	//Create solid archives
	Solid bool `yaml_comment:"create solid 7-zip archives"`

	//List of patterns to exclude from archive
	Exclude []string `yaml_comment:"list of patterns to exclude from archiving"`

	//List of patterns to be added to archive without compression (works for FULL backups only)
	SkipCompression []string `yaml:"skip_compression" yaml_comment:"list of patterns to fast include to archive without compression"`

	//Run cleanup procedure before or after archive creation (default: after)
	Cleanup string `yaml_comment:"do old archives cleanup before or after archiving"`

	MaxFullCount int `yaml:"max_full_count" yaml_comment:"maximum full archives to keep"`
	KeepAtLeast  int `yaml:"keep_at_least" yaml_comment:"do not remove full archives if they younger than this count of days"`

	//Maximum number of diff archives to have after full backup
	MaxDiffCount int `yaml:"max_diff_count" yaml_comment:"maximum differential archives to create before creating new full archive"`

	//Maximum latest diff size IN PERCENTS to force full backup
	MaxDiffSizePercent int `yaml:"max_diff_size_percent" yaml_comment:"last differential archive size in percents to force new full archive next run, 0 = not set"`

	//Maximum total diffs size IN PERCENTS to force full backup
	MaxTotalDiffSizePercent int `yaml:"max_total_diff_size_percent" yaml_comment:"total size of differential archives since latest full archive in percents to force new full archive next run, 0 = not set"`
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
