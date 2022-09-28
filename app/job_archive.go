package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type JobArchiveFile struct {
	Name   string
	Path   string
	IsFull bool
	Size   int64
	Time   time.Time
}

type JobArchiveFullItem struct {
	File                 *JobArchiveFile
	DiffItemList         []*JobArchiveDiffItem
	TotalDiffSizePercent int
}

type JobArchiveDiffItem struct {
	File            *JobArchiveFile
	DiffSizePercent int
}

type JobArchive struct {
	FilesList    []JobArchiveFile
	FullItemList []JobArchiveFullItem
}

func (job *Job) ScanArchive() {
	files_list, err := os.ReadDir(job.Settings.ArchivesPath)
	if err != nil {
		log.Fatalln(err)
	}

	job.Archive = JobArchive{
		FilesList:    make([]JobArchiveFile, 0, len(files_list)),
		FullItemList: make([]JobArchiveFullItem, 0),
	}

	//prepare regexp and suffix
	re := regexp.MustCompile(
		"^" + regexp.QuoteMeta(job.Settings.ArchiveName) + "_(.*)_(" +
			regexp.QuoteMeta(job.Settings.FullSuffix) + "|" + regexp.QuoteMeta(job.Settings.DiffSuffix) +
			")\\.7z$",
	)
	full_suffix := "_" + job.Settings.FullSuffix + ".7z"

	//scan list
	for _, value := range files_list {
		//skip directories
		if value.IsDir() {
			continue
		}

		//check if this is our file (by name)
		if !re.MatchString(value.Name()) {
			continue
		}

		info, err := value.Info()
		if err == nil {
			job.Archive.FilesList = append(
				job.Archive.FilesList,
				JobArchiveFile{
					Name:   value.Name(),
					Path:   filepath.Join(job.Settings.ArchivesPath, value.Name()),
					IsFull: strings.HasSuffix(value.Name(), full_suffix),
					Size:   info.Size(),
					Time:   info.ModTime(),
				},
			)
		}
	}

	//sort by time
	sort.Slice(job.Archive.FilesList, func(i, j int) bool {
		return job.Archive.FilesList[i].Time.Before(job.Archive.FilesList[j].Time)
	})

	// build FULL -> DIFF[] tree
	var current_full_item *JobArchiveFullItem = nil

	for index := range job.Archive.FilesList {
		var ai_pointer = &job.Archive.FilesList[index]

		if ai_pointer.IsFull {
			full_item := JobArchiveFullItem{
				File:         ai_pointer,
				DiffItemList: make([]*JobArchiveDiffItem, 0),
			}

			job.Archive.FullItemList = append(job.Archive.FullItemList, full_item)

			//take address of just appended element
			current_full_item = &job.Archive.FullItemList[len(job.Archive.FullItemList)-1]
		} else {
			if current_full_item != nil { //skip diffs without parent
				current_full_item.DiffItemList = append(
					current_full_item.DiffItemList,
					&JobArchiveDiffItem{
						File: ai_pointer,
					},
				)
			}
		}
	}

	//calculate diff sizes
	for index := range job.Archive.FullItemList {
		full_item := &job.Archive.FullItemList[index]

		if full_item.File.Size == 0 {
			continue
		}

		var total_diff_size int64 = 0

		for index := range full_item.DiffItemList {
			diff_item := full_item.DiffItemList[index]
			diff_item.DiffSizePercent = int(diff_item.File.Size * 100 / full_item.File.Size)

			total_diff_size += diff_item.File.Size
		}

		full_item.TotalDiffSizePercent = int(total_diff_size * 100 / full_item.File.Size)
	}
}

func (ja *JobArchive) Dump(die bool) {
	fmt.Println("------ RAW LIST ------- ")
	for _, raw_file := range ja.FilesList {
		fmt.Println("RAW: "+raw_file.Name, raw_file.Size)
	}

	fmt.Println("------------- ")
	for index := range ja.FullItemList {
		full_item := &ja.FullItemList[index]

		fmt.Printf("FULL: %s, size: %d, diff_size: %d%%\n", full_item.File.Name, full_item.File.Size, full_item.TotalDiffSizePercent)

		for _, diff_item := range full_item.DiffItemList {
			fmt.Printf("    DIFF: %s, size %d %d%%\n", diff_item.File.Name, diff_item.File.Size, diff_item.DiffSizePercent)
		}
	}

	if die {
		os.Exit(0)
	}
}

func (afi *JobArchiveFullItem) Unlink() {
	//delete diffs
	for _, diff_item := range afi.DiffItemList {
		if err := os.Remove(diff_item.File.Path); err != nil {
			log.Fatalf("Error deleting file %s: %s", diff_item.File.Path, err)
		}
	}

	//delete itself
	if err := os.Remove(afi.File.Path); err != nil {
		log.Fatalf("Error deleting file %s: %s", afi.File.Path, err)
	}
}
