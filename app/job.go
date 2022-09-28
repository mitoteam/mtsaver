package app

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

	return job, nil
}

func (job *Job) Run() error {
	if err := os.MkdirAll(job.Settings.ArchivesPath, 0777); err != nil {
		return err
	}

	job.ScanArchive()
	//job.Archive.Dump()

	if job.Settings.Cleanup == "before" {
		job.cleanup()
	}

	if len(job.Archive.FullItemList) == 0 { //no full archives at all
		//create one unconditionally
		job.createNewArchive(true, "")
	} else {
		//check diffs for the last one
		need_full := false
		last_full_item := job.Archive.FullItemList[len(job.Archive.FullItemList)-1]

		//check max count
		if len(last_full_item.DiffItemList) >= job.Settings.MaxDiffCount {
			need_full = true
		}

		//check max total size (in percents!)
		if job.Settings.MaxTotalDiffSizePercent > 0 {
			if last_full_item.TotalDiffSizePercent >= job.Settings.MaxTotalDiffSizePercent {
				need_full = true
			}
		}

		//check last diff size (in percents!)
		if job.Settings.MaxDiffSizePercent > 0 && len(last_full_item.DiffItemList) > 0 {
			last_diff_item := last_full_item.DiffItemList[len(last_full_item.DiffItemList)-1]
			if last_diff_item.DiffSizePercent >= job.Settings.MaxDiffSizePercent {
				need_full = true
			}
		}

		job.createNewArchive(need_full, last_full_item.File.Path)
	}

	if job.Settings.Cleanup == "after" {
		job.ScanArchive() //new archives may have been created
		job.cleanup()
	}

	job.ScanArchive()
	//job.Archive.Dump(false)

	return nil
}

func (job *Job) getArchiveName(is_full bool) string {
	var suffix string

	if is_full {
		suffix = job.Settings.FullSuffix
	} else {
		suffix = job.Settings.DiffSuffix
	}

	return filepath.Join(
		job.Settings.ArchivesPath,
		job.Settings.ArchiveName+"_"+time.Now().Format(job.Settings.DateFormat)+"_"+suffix+".7z",
	)
}

func (job *Job) LoadSettings() {
	job.Settings = JobSettings{
		CompressionLevel: -1,
		MaxFullCount:     5,
		MaxDiffCount:     20,
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
		s.ArchivesPath = filepath.Join(filepath.Dir(job.Path), name+"_ARCHIVE")
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

	if s.Cleanup == "" {
		s.Cleanup = "after"
	} else if s.Cleanup != "before" && s.Cleanup != "after" {
		log.Fatalln("Valid  values for 'cleanup' option are 'before', 'after'")
	}

	if s.MaxFullCount < 1 {
		log.Fatalln("Minumum value for max_full_count is 1")
	}

	if s.MaxDiffCount < 0 {
		log.Fatalln("Minumum value for max_diff_count is 0")
	}
}

func (job *Job) createNewArchive(is_full bool, full_archive_path string) {
	var seven_zip_arguments = make([]string, 0, 10)

	if is_full {
		seven_zip_arguments = append(
			seven_zip_arguments, "a",
			job.getArchiveName(true), //new full arch name
		)
	} else {
		seven_zip_arguments = append(
			seven_zip_arguments, "u",
			full_archive_path,
			"-u-", // disable updates in the base archive
			"-up3q3r2x2y2z0w2!"+job.getArchiveName(false),
		)
	}

	//common arguments
	seven_zip_arguments = append(seven_zip_arguments,
		"-ssw", //Compress files open for writing
		"-mx"+strconv.Itoa(job.Settings.CompressionLevel), //compression level
	)

	for _, pattern := range job.Settings.Exclude {
		seven_zip_arguments = append(seven_zip_arguments, "-xr!"+pattern)
	}

	// final argument - folder to pack
	seven_zip_arguments = append(
		seven_zip_arguments,
		filepath.Join(job.Path, "*"),
	)

	//run command
	cmd := exec.Command(Global.SevenZipCmd, seven_zip_arguments...)
	//fmt.Println("CMD: " + cmd.String())
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
}

func (job *Job) cleanup() {
	//delete FULL items
	for len(job.Archive.FullItemList) > job.Settings.MaxFullCount {
		job.Archive.FullItemList[0].Unlink()

		job.Archive.FullItemList = job.Archive.FullItemList[1:]
	}
}
