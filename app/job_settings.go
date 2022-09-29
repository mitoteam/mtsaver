package app

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const DefaultSettingsFilename = ".mtsaver.yml"

// Setting for archived folder
type JobSettings struct {
	ArchivesPath string `yaml:"archives_path"`
	ArchiveName  string `yaml:"archive_name"`
	FullSuffix   string `yaml:"full_suffix"`
	DiffSuffix   string `yaml:"diff_suffix"`
	DateFormat   string `yaml:"date_format"`

	//Int value from 0 = do not compress to 9 = max compression, longest time
	CompressionLevel int `yaml:"compression_level"`

	//List of patterns to exclude from archive
	Exclude []string

	//List of patterns to be added to archive without compression (works for FULL backups only)
	SkipCompression []string `yaml:"skip_compression"`

	//Run cleanup procedure before or after archive creation (default: after)
	Cleanup string

	MaxFullCount int `yaml:"max_full_count"`

	//Maximum number of diff archives to have after full backup
	MaxDiffCount int `yaml:"max_diff_count"`

	//Maximum latest diff size IN PERCENTS to force full backup
	MaxDiffSizePercent int `yaml:"max_diff_size_percent"`

	//Maximum total diffs size IN PERCENTS to force full backup
	MaxTotalDiffSizePercent int `yaml:"max_total_diff_size_percent"`
}

func (js *JobSettings) LoadFromFile(path string) {
	//try to load only if it exists
	if !IsFileExists(path) {
		return
	}

	yamlFile, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("Error while reading %v file: %v", path, err)
	}

	err = yaml.Unmarshal(yamlFile, js)
	if err != nil {
		log.Fatalf("Error parsing yaml: %v", err)
	}
}

func (js *JobSettings) SaveToFile(path string, comment string) error {
	settings_yaml := &yaml.Node{}

	if err := settings_yaml.Encode(js); err != nil {
		log.Fatalln(err)
	}

	if len(settings_yaml.Content) > 0 {
		settings_yaml.Content[0].HeadComment = Global.AppName + " directory settings file" +
			"\n# " + strings.ReplaceAll(comment, "\n", "\n# ") +
			"\n# Created on: " + time.Now().Format(time.RFC3339) +
			"\n#\n\n"
	}

	file_yaml, err := yaml.Marshal(settings_yaml)

	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(path, file_yaml, 0644); err != nil {
		return err
	}

	return nil
}

func (js *JobSettings) Print() {
	yaml, err := yaml.Marshal(js)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(yaml))
}
