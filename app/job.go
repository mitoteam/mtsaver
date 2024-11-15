package app

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitoteam/mttools"
)

type Job struct {
	Path     string
	Settings JobSettings
	Archive  JobArchive
}

// Creates new Job. If first argument given - using it as path to directory. If absent - using current directory.
func NewJobFromArgs(args []string) (*Job, error) {
	var path string

	if len(args) > 0 {
		path = args[0]
	} else {
		path = "." //current directory
	}

	path, err := mttools.GetDirAbsolutePath(path)
	if err != nil {
		return nil, err
	}

	var job = &Job{
		Path: path,
	}

	job.LoadSettings()

	return job, nil
}

func (job *Job) Run() error {
	fmt.Printf("Starting directory backup: %s\n", job.Path)

	if err := os.MkdirAll(job.Settings.ArchivesPath, 0777); err != nil {
		return err
	}

	if job.Settings.Cleanup == "before" {
		job.Cleanup()
	}

	job.ScanArchive()
	//job.Archive.Dump(false)

	if JobRuntimeOptions.ForceFull {
		job.createArchive(true, "")
	} else if JobRuntimeOptions.ForceDiff {
		if len(job.Archive.FullItemList) == 0 {
			log.Fatalln("can not force differential backup because no full backups found.")
		}

		job.createArchive(false, job.Archive.FullItemList[len(job.Archive.FullItemList)-1].File.Path)
	} else if len(job.Archive.FullItemList) == 0 {
		//no full archives at all
		//create one unconditionally
		job.createArchive(true, "")
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

		job.createArchive(need_full, last_full_item.File.Path)
	}

	if job.Settings.Cleanup == "after" {
		job.Cleanup()
	}

	return nil
}

func (job *Job) Dump() {
	if !mttools.IsDirExists(job.Settings.ArchivesPath) {
		fmt.Printf("%s directory does not exists", job.Settings.ArchivesPath)
	}

	job.ScanArchive()
	job.Archive.Dump(false)
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

func (job *Job) SettingsFilename() (filename string) {
	if filepath.IsAbs(JobRuntimeOptions.SettingsFilename) {
		filename = JobRuntimeOptions.SettingsFilename
	} else {
		filename = filepath.Join(job.Path, JobRuntimeOptions.SettingsFilename)
	}

	return
}

func (job *Job) LoadSettings() {
	job.Settings = NewJobSettings()

	if mttools.IsFileExists(job.SettingsFilename()) {
		if err := job.Settings.LoadFromFile(job.SettingsFilename()); err != nil {
			log.Fatalln(err)
		}
	}

	var s = &job.Settings

	// set defaults if something is missing in file

	if s.DateFormat == "" {
		s.DateFormat = "2006-01-02_15-04-05"
	}

	name := filepath.Base(job.Path)

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
		log.Fatalln("Minimum value for max_full_count is 1")
	}

	if s.MaxDiffCount < 0 {
		log.Fatalln("Minimum value for max_diff_count is 0")
	}
}

func (job *Job) createArchive(is_full bool, full_archive_path string) {
	var common_arguments = []string{} //7-zip command (add or update), basic compression settings
	job_archive_filename := job.getArchiveName(is_full)
	var err error

	if is_full {
		common_arguments = append(common_arguments, "a", job_archive_filename)
	} else {
		// thanks: https://nagimov.me/post/simple-differential-and-incremental-backups-using-7-zip/

		common_arguments = append(common_arguments,
			"u",
			full_archive_path, //existing full archive
			"-u-",             // disable updates in the base archive
			"-up3q3r2x2y2z0w2!"+job_archive_filename,
		)
	}

	common_arguments = append(common_arguments,
		"-r0",       //recursion only for patterns with wildcard
		"-ssw",      //compress files open for writing
		"-bb1",      //show names of processed files
		"-bse1",     //error messages to stdout
		"-sccUTF-8", //console output encoding
	)

	//turn on solid mode for archives
	if JobRuntimeOptions.Solid || job.Settings.Solid {
		common_arguments = append(common_arguments, "-ms=on")
	}

	//set password for archive
	password := ""

	if len(JobRuntimeOptions.Password) > 0 {
		password = JobRuntimeOptions.Password
	} else if len(job.Settings.Password) > 0 {
		password = job.Settings.Password
	}

	if len(password) > 0 {
		common_arguments = append(common_arguments, "-p"+password)

		if JobRuntimeOptions.EncryptFilenames || job.Settings.EncryptFilenames {
			//mhe = encrypt headers
			common_arguments = append(common_arguments, "-mhe")
		}
	}

	//exclusions
	for _, pattern := range job.Settings.Exclude {
		common_arguments = append(common_arguments, "-xr!"+pattern)
	}

	//// RUN BASIC COMPRESSION
	var basic_arguments = make([]string, len(common_arguments))
	copy(basic_arguments, common_arguments)

	basic_arguments = append(basic_arguments,
		"-mx"+strconv.Itoa(job.Settings.CompressionLevel), //compression level
	)

	if is_full {
		// exclude skip_compression patterns
		for _, pattern := range job.Settings.SkipCompression {
			basic_arguments = append(basic_arguments, "-xr!"+pattern)
		}
	}

	// final argument - whole folder to pack
	basic_arguments = append(basic_arguments, filepath.Join(job.Path, "*"))

	// run command
	main_output := runSevenZip(basic_arguments)

	//// ADD ITEMS WITHOUT COMPRESSION - works only for full archives now
	if is_full && len(job.Settings.SkipCompression) > 0 {
		var skip_compression_arguments = make([]string, len(common_arguments))
		copy(skip_compression_arguments, common_arguments)

		skip_compression_arguments = append(skip_compression_arguments,
			"-m0=copy", //do not compress at all
		)

		// exclude skip_compression patterns
		for _, pattern := range job.Settings.SkipCompression {
			skip_compression_arguments = append(skip_compression_arguments, filepath.Join(job.Path, pattern))
		}

		runSevenZip(skip_compression_arguments)
	}

	//check if empty diff was created
	if !is_full {
		is_empty := strings.Contains(main_output, "Add new data to archive: 0 files, 0 bytes")

		if is_empty {
			if !job.Settings.KeepEmptyDiff {
				fmt.Printf("Empty diff archive created (%s). Removing it.", filepath.Base(job_archive_filename))

				if err = os.Remove(job_archive_filename); err != nil {
					log.Fatalf("Error deleting file %s: %s", filepath.Base(job_archive_filename), err)
				}
			}
		} else {
			if !job.Settings.KeepSameDiff {
				if prev_archive := job.Archive.LastFile(); prev_archive != nil {
					if !prev_archive.IsFull {
						//log.Printf("Prev hash: %s", prev_archive.Hash)

						var last_hash string
						last_hash, err = mttools.FileSha256(job_archive_filename)

						if err == nil {
							//log.Printf("Last hash: %s", last_hash)

							if len(last_hash) > 0 && last_hash == prev_archive.Hash {
								fmt.Printf("Diff archive with same sha256 created (%s). Removing it.", filepath.Base(job_archive_filename))

								if err = os.Remove(job_archive_filename); err != nil {
									log.Fatalf("Error deleting file %s: %s", filepath.Base(job_archive_filename), err)
								}

							}
						}
					}
				}
			}
		}
	}
}

func runSevenZip(arguments []string) string {
	cmd := exec.Command(Global.SevenZipCmd, arguments...)
	//fmt.Println("7z CMD: " + cmd.String())

	pipe, _ := cmd.StdoutPipe()

	cmd.Start()

	output := ""

	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		text := scanner.Text()
		output = output + text
		fmt.Println(text)
	}

	cmd.Wait()

	return output
}

func (job *Job) Cleanup() error {
	//do cleanup only if archives directory exists
	if !mttools.IsDirExists(job.Settings.ArchivesPath) {
		return nil
	}

	//always re-scan archives before cleaning up
	job.ScanArchive()
	//job.Archive.Dump(true)

	//delete FULL items
	out_of_window_count := len(job.Archive.FullItemList) - job.Settings.MaxFullCount

	for i := 0; i < out_of_window_count; i++ {
		if job.Settings.KeepAtLeast > 0 {
			if job.Archive.FullItemList[i].File.Age < job.Settings.KeepAtLeast {
				continue
			}
		}

		job.Archive.FullItemList[i].Unlink()
	}

	return nil
}
