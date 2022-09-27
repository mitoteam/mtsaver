package mtsaver

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const settings_filename = ".mtsaver.yml"

type JobSettings struct {
	ArchivesPath string
	ArchiveName  string
	FullSuffix   string
	DiffSuffix   string
	DateFormat   string

	CompressionLevel int `yaml:"compression_level"`
	Exclude          []string
	MaxFullCount     int    `yaml:"max_full_count"`
	MaxDiffCount     int    `yaml:"max_diff_count"`
	Cleanup          string `yaml:"cleanup"`
}

func (job_settings *JobSettings) LoadFromDir(dirPath string) {
	var filename = dirPath + string(filepath.Separator) + settings_filename

	//try to load only if it exists
	if _, err := os.Stat(filename); err != nil {
		return
	}

	yamlFile, err := os.ReadFile(filename)

	if err != nil {
		log.Fatalf("Error while reading %v file: %v", filename, err)
	}

	err = yaml.Unmarshal(yamlFile, job_settings)
	if err != nil {
		log.Fatalf("Error parsing yaml: %v", err)
	}
}
