package app

import (
	"fmt"
	"log"
	"log/slog"
	"os"
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

	logger  *slog.Logger
	logfile *os.File
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

	// make sure archives directory exists
	if !mttools.IsDirExists(job.Settings.ArchivesPath) {
		if err := os.MkdirAll(job.Settings.ArchivesPath, 0777); err != nil {
			return nil, err
		}

		job.Log("Archives directory created: %s", job.Settings.ArchivesPath)
	}

	//initialize logger
	if job.Settings.LogFormat == "text" || job.Settings.LogFormat == "json" {
		logFilename := filepath.Join(job.Settings.ArchivesPath, job.Settings.LogFilename)
		logExists := mttools.IsFileExists(logFilename)

		job.logfile, err = os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening log file %s: %v", logFilename, err)
		}

		if logExists {
			//add some space to distinguish runs
			job.logfile.WriteString("\n\n")
		}

		var logHandler slog.Handler

		if job.Settings.LogFormat == "text" {
			logHandler = slog.NewTextHandler(job.logfile, nil)
		} else if job.Settings.LogFormat == "json" {
			logHandler = slog.NewJSONHandler(job.logfile, nil)
		} else {
			log.Panicf("Unknown log format %s", job.Settings.LogFormat)
		}

		job.logger = slog.New(logHandler)
	}

	return job, nil
}

func (job *Job) Run() error {
	job.Log("Starting directory backup: %s", job.Path)

	if job.Settings.Cleanup == "before" {
		job.Cleanup()
	}

	job.ScanArchive()
	//job.Archive.Dump(false)

	//run commands before creating new archive
	if len(job.Settings.RunBefore) > 0 {
		job.Log("Executing commands from 'run_before' option")

		for _, command := range job.Settings.RunBefore {
			job.Log("Command: %s", command)

			if output, err := mttools.ExecCommandLine(command); err != nil {
				job.Log("Command error: %s", err.Error())
			} else {
				if job.Settings.LogCommandOutput {
					job.Log("Command output:")
					job.RawLog(output)
				} else {
					//screen only
					log.Println("Command output:")
					fmt.Println(output)
				}
			}
		}
	}

	if JobRuntimeOptions.ForceFull {
		job.Log("Full archive was forced")
		job.createArchive(true, "")
	} else if JobRuntimeOptions.ForceDiff {
		if len(job.Archive.FullItemList) == 0 {
			log.Fatalln("Can not force differential backup because no full backups found.")
		}

		job.Log("Diff archive was forced")
		job.createArchive(false, job.Archive.FullItemList[len(job.Archive.FullItemList)-1].File.Path)
	} else if len(job.Archive.FullItemList) == 0 {
		//no full archives at all, create one unconditionally
		job.Log("No full archives found. Creating one.")
		job.createArchive(true, "")
	} else {
		//check diffs for the last one full-arch
		need_full := false
		last_full_arch := job.Archive.FullItemList[len(job.Archive.FullItemList)-1]

		//check max count
		if len(last_full_arch.DiffItemList) >= job.Settings.MaxDiffCount {
			job.Log(
				"Diff archives count (%d) exceeds maximum (%d). Creating full archive.",
				len(last_full_arch.DiffItemList), job.Settings.MaxDiffCount,
			)
			need_full = true
		}

		//check max total size (in percents!)
		if !need_full && job.Settings.MaxTotalDiffSizePercent > 0 {
			if last_full_arch.TotalDiffSizePercent >= job.Settings.MaxTotalDiffSizePercent {
				job.Log(
					"Diff archives total size (%d%% of full archive) exceeds maximum (%d%%). Creating full archive.",
					last_full_arch.TotalDiffSizePercent, job.Settings.MaxTotalDiffSizePercent,
				)

				need_full = true
			}
		}

		//check last diff size (in percents!)
		if !need_full && job.Settings.MaxDiffSizePercent > 0 && len(last_full_arch.DiffItemList) > 0 {
			last_diff_item := last_full_arch.DiffItemList[len(last_full_arch.DiffItemList)-1]
			if last_diff_item.DiffSizePercent >= job.Settings.MaxDiffSizePercent {
				job.Log(
					"Last diff archive size (%d%% of full archive) exceeds maximum (%d%%). Creating full archive.",
					last_diff_item.DiffSizePercent, job.Settings.MaxDiffSizePercent,
				)

				need_full = true
			}
		}

		if !need_full {
			job.Log("Creating diff archive for %s", last_full_arch.File.Name)
		}

		job.createArchive(need_full, last_full_arch.File.Path)
	}

	if job.Settings.Cleanup == "after" {
		job.Cleanup()
	}

	if job.logfile != nil {
		job.logfile.Close()
	}

	return nil
}

func (job *Job) Dump() {
	if !mttools.IsDirExists(job.Settings.ArchivesPath) {
		fmt.Printf("%s directory does not exists\n", job.Settings.ArchivesPath)
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
	s.ApplyDefaultsAndCheck(job.Path)
}

func (job *Job) createArchive(is_full bool, full_archive_path string) {
	var common_arguments = []string{} //7-zip command (add or update), basic compression settings
	job_archive_filename := job.getArchiveName(is_full)
	var err error
	js := &job.Settings //convenience variable

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
	if js.Solid {
		common_arguments = append(common_arguments, "-ms=on")
	}

	//set password for archive
	if len(js.Password) > 0 {
		common_arguments = append(common_arguments, "-p"+js.Password)

		if js.EncryptFilenames {
			//mhe = encrypt headers
			common_arguments = append(common_arguments, "-mhe")
		}
	}

	//exclusions
	for _, pattern := range js.Exclude {
		common_arguments = append(common_arguments, "-xr!"+pattern)
	}

	//// RUN BASIC COMPRESSION
	var basic_arguments = make([]string, len(common_arguments))
	copy(basic_arguments, common_arguments)

	basic_arguments = append(basic_arguments,
		"-mx"+strconv.Itoa(js.CompressionLevel), //compression level
	)

	if is_full {
		// exclude skip_compression patterns
		for _, pattern := range js.SkipCompression {
			basic_arguments = append(basic_arguments, "-xr!"+pattern)
		}
	}

	// final argument - whole folder to pack
	basic_arguments = append(basic_arguments, filepath.Join(job.Path, "*"))

	// run command
	main_output := job.runSevenZip(basic_arguments)

	//// ADD ITEMS WITHOUT COMPRESSION - works only for full archives now
	if is_full && len(js.SkipCompression) > 0 {
		var skip_compression_arguments = make([]string, len(common_arguments))
		copy(skip_compression_arguments, common_arguments)

		skip_compression_arguments = append(skip_compression_arguments,
			"-m0=copy", //do not compress at all
		)

		// exclude skip_compression patterns
		for _, pattern := range js.SkipCompression {
			skip_compression_arguments = append(skip_compression_arguments, filepath.Join(job.Path, pattern))
		}

		job.runSevenZip(skip_compression_arguments)
	}

	var archType string
	if is_full {
		archType = "Full"
	} else {
		archType = "Diff"
	}

	job.Log("%s archive created: %s", archType, job_archive_filename)

	//check if empty diff was created
	if !is_full {
		is_empty := strings.Contains(main_output, "Add new data to archive: 0 files, 0 bytes")

		if is_empty {
			if !js.KeepEmptyDiff {
				job.Log("Empty diff archive detected (%s). Removing it.", filepath.Base(job_archive_filename))

				if err = os.Remove(job_archive_filename); err != nil {
					log.Fatalf("Error deleting file %s: %s", filepath.Base(job_archive_filename), err)
				}
			}
		} else {
			if !js.KeepSameDiff {
				if prev_archive := job.Archive.LastFile(); prev_archive != nil {
					if !prev_archive.IsFull {

						var last_hash string
						last_hash, err = mttools.FileSha256(job_archive_filename)

						if err == nil {
							if len(last_hash) > 0 && last_hash == prev_archive.Hash {
								job.Log("Diff archive with same sha256 created (%s). Removing it.", filepath.Base(job_archive_filename))

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

func (job *Job) runSevenZip(arguments []string) string {
	output, err := mttools.ExecCmdWaitAndPrint(Global.SevenZipCmd, arguments)

	if err != nil {
		job.Log("Error running 7-zip: %s", err.Error())
	}

	if job.Settings.LogCommandOutput {
		job.RawLog(output)
	}

	return output
}

func (job *Job) Cleanup() error {
	//always re-scan archives before cleaning up
	job.ScanArchive()
	//DBG: job.Archive.Dump(true)

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

// Adds message to job's log
func (job *Job) Log(format string, args ...any) {
	message := fmt.Sprintf(format, args...)

	//always print to screen
	log.Print(message)

	// if file logger is defined write to it as well
	if job.logger != nil {
		job.logger.Info(message)
	}
}

// Adds content to logs without attributing and formatting
func (job *Job) RawLog(content string) {
	//make sure it ends with EOL
	if content[len(content)-1] != '\n' {
		content += "\n"
	}

	//always print to screen
	fmt.Print(content)

	// add to log file only if asked
	if job.logfile != nil && job.Settings.LogCommandOutput {
		job.logfile.WriteString(content)
	}
}
