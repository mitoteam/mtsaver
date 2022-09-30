package app

import (
	"fmt"
	"log"
	"mtsaver/mttools"
	"os"
	"reflect"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const DefaultSettingsFilename = ".mtsaver.yml"

// Setting for archived folder
type JobSettings struct {
	ArchivesPath string `yaml:"archives_path" yaml_comment:"full path to directory to create arhives in"`
	ArchiveName  string `yaml:"archive_name" yaml_comment:"base archive name (appended with timestamp and suffix)"`
	FullSuffix   string `yaml:"full_suffix" yaml_comment:"suffix for full archives"`
	DiffSuffix   string `yaml:"diff_suffix" yaml_comment:"suffix for differential archives"`
	DateFormat   string `yaml:"date_format" yaml_comment:"archive filename timestamp format, don't touch, Go's time formatting is crazyness https://mttm.ml/go-time-format"`

	//Int value from 0 = do not compress to 9 = max compression, longest time
	CompressionLevel int `yaml:"compression_level" yaml_comment:"7-Zip compression level from 0 to 9"`

	//List of patterns to exclude from archive
	Exclude []string `yaml_comment:"list of patterns to exclude from archiving"`

	//List of patterns to be added to archive without compression (works for FULL backups only)
	SkipCompression []string `yaml:"skip_compression" yaml_comment:"list of patterns to fast include to archive without compression"`

	//Run cleanup procedure before or after archive creation (default: after)
	Cleanup string `yaml_comment:"do old archives cleanup before or after archiving"`

	MaxFullCount int `yaml:"max_full_count" yaml_comment:"maximum full archives to keep"`

	//Maximum number of diff archives to have after full backup
	MaxDiffCount int `yaml:"max_diff_count" yaml_comment:"maximum differentiam archives to create before creating new full archive"`

	//Maximum latest diff size IN PERCENTS to force full backup
	MaxDiffSizePercent int `yaml:"max_diff_size_percent" yaml_comment:"last differential archive size in percents to force new full archive next run, 0 = not set"`

	//Maximum total diffs size IN PERCENTS to force full backup
	MaxTotalDiffSizePercent int `yaml:"max_total_diff_size_percent" yaml_comment:"total size of differential archives since latest full archive in percents to force new full archive next run, 0 = not set"`
}

func (js *JobSettings) LoadFromFile(path string) {
	//try to load only if it exists
	if !mttools.IsFileExists(path) {
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
	node_yaml := &yaml.Node{}

	if err := node_yaml.Encode(js); err != nil {
		log.Fatalln(err)
	}

	//setting header comment
	node_yaml.HeadComment = Global.AppName + " directory settings file" +
		"\n# " + strings.ReplaceAll(comment, "\n", "\n# ") +
		"\n# Created on: " + time.Now().Format(time.RFC3339) +
		"\n#\n\n"

	//adding comments
	r := reflect.TypeOf(js).Elem()

	for _, option_yaml := range node_yaml.Content {
		option_yaml.HeadComment = settingsOptionYamlComment(r, option_yaml.Value)
	}

	// unmarshalling to raw yaml
	file_yaml, err := yaml.Marshal(node_yaml)

	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Print(string(file_yaml)); os.Exit(0)

	if err := os.WriteFile(path, file_yaml, 0644); err != nil {
		return err
	}

	return nil
}

func settingsOptionYamlComment(r reflect.Type, yaml_field string) string {
	for i := 0; i < r.NumField(); i++ {
		tag := r.Field(i).Tag.Get("yaml")

		if tag == "" {
			tag = strings.ToLower(r.Field(i).Name)
		} else {
			tag = strings.TrimSpace(strings.Split(tag, ",")[0])
		}

		//fmt.Println(r.Field(i).Name, tag)
		if tag == yaml_field {
			return r.Field(i).Tag.Get("yaml_comment")
		}
	}

	return ""
}

func (js *JobSettings) Print() {
	yaml, err := yaml.Marshal(js)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(yaml))
}
