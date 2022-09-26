package mtsaver

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Job struct {
	Path         string
	ArchivesPath string
	ArchiveName  string
	FullSuffix   string
	Diffsuffix   string
	DateFormat   string
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
	name := filepath.Base(full_path)

	var job = &Job{
		Path:         full_path,
		DateFormat:   "2006-01-02_15-04-05",
		ArchiveName:  name,
		ArchivesPath: filepath.Dir(full_path) + string(filepath.Separator) + name + "_ARCHIVE",
		FullSuffix:   "FULL",
		Diffsuffix:   "DIFF",
	}

	return job, nil
}

func (job *Job) Run() error {
	fmt.Println("ArchivesPath: " + job.ArchivesPath)
	fmt.Println("GetFullArchiveName: " + job.GetFullArchiveName())

	if err := os.MkdirAll(job.ArchivesPath, 0777); err != nil {
		return err
	}

	var seven_zip_arguments = []string{
		"a",                      //add
		job.GetFullArchiveName(), //arch name
		job.Path,                 //folder
		"-ssw",                   //Compress files open for writing
		//"-xr!.git",               //exclude git
		//"-xr!*.exe",              //exclude exe
	}

	cmd := exec.Command(Global.SevenZipCmd, seven_zip_arguments...)
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))

	return nil
}

func (job *Job) GetFullArchiveName() string {
	return job.ArchivesPath + string(filepath.Separator) + job.ArchiveName + "_" + time.Now().Format(job.DateFormat) + "_" + job.FullSuffix + ".7z"
}
