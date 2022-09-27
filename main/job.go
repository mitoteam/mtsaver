package mtsaver

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type Job struct {
	Path     string
	Settings JobSettings
	Archive  JobArchive
}

func NewJob(path string) (*Job, error) {
	file_info, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if !file_info.IsDir() {
		return nil, errors.New("\"" + path + "\" directory does not exists")
	}

	full_path, _ := filepath.Abs(path)

	var job = &Job{
		Path: full_path,
	}

	job.LoadSettings()
	job.ScanArchive()

	job.Archive.Dump()
	os.Exit(0)

	return job, nil
}

func (job *Job) Run() error {
	if err := os.MkdirAll(job.Settings.ArchivesPath, 0777); err != nil {
		return err
	}

	var seven_zip_arguments = []string{
		"u",                      //add
		job.getArchiveName(true), //arch name
		"-ssw",                   //Compress files open for writing
		"-mx" + strconv.Itoa(job.Settings.CompressionLevel), //compression level
	}

	for _, pattern := range job.Settings.Exclude {
		seven_zip_arguments = append(seven_zip_arguments, "-xr!"+pattern)
	}

	// final argument - folder to pack
	seven_zip_arguments = append(
		seven_zip_arguments,
		job.Path+string(filepath.Separator)+"*",
	)

	//run command
	cmd := exec.Command(Global.SevenZipCmd, seven_zip_arguments...)
	//fmt.Println("CMD: " + cmd.String())
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))

	return nil
}

func (job *Job) getArchiveName(full bool) string {
	var suffix string

	if full {
		suffix = job.Settings.FullSuffix
	} else {
		suffix = job.Settings.DiffSuffix
	}

	return job.Settings.ArchivesPath + string(filepath.Separator) +
		job.Settings.ArchiveName + "_" + time.Now().Format(job.Settings.DateFormat) +
		"_" + suffix + ".7z"
}

func (job *Job) LoadSettings() {
	job.Settings = JobSettings{
		CompressionLevel: -1,
	}

	job.Settings.LoadFromDir(job.Path)
	var s = &job.Settings

	name := filepath.Base(job.Path)

	if s.DateFormat == "" {
		s.DateFormat = "2006-01-02_15-04-05"
	}

	if len(s.ArchiveName) == 0 {
		s.ArchiveName = name
	}

	if len(s.ArchivesPath) == 0 {
		s.ArchivesPath = filepath.Dir(job.Path) + string(filepath.Separator) + name + "_ARCHIVE"
	}

	if len(s.FullSuffix) == 0 {
		s.FullSuffix = "FULL"
	}

	if len(s.DiffSuffix) == 0 {
		s.DiffSuffix = "DIFF"
	}

	if s.CompressionLevel == -1 {
		s.CompressionLevel = 5
	}

	// Do settings checks
	if s.FullSuffix == s.DiffSuffix {
		log.Fatalln("Full suffix should differ from diff suffix")
	}
}
