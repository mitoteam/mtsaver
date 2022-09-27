package mtsaver

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type JobArchiveItem struct {
	Name   string
	Path   string
	IsFull bool
	Size   int64
	Time   time.Time
}

type JobArchiveFullItem struct {
	Item                 *JobArchiveItem
	DiffItemList         []*JobArchiveItem
	TotalDiffSizePercent int
}

type JobArchive struct {
	RawList      []JobArchiveItem
	FullItemList []JobArchiveFullItem
}

func (job *Job) ScanArchive() {
	files_list, err := os.ReadDir(job.Settings.ArchivesPath)
	if err != nil {
		log.Fatalln(err)
	}

	job.Archive = JobArchive{
		RawList:      make([]JobArchiveItem, 0, len(files_list)),
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
			job.Archive.RawList = append(
				job.Archive.RawList,
				JobArchiveItem{
					Name:   value.Name(),
					Path:   job.Settings.ArchivesPath + string(os.PathSeparator) + value.Name(),
					IsFull: strings.HasSuffix(value.Name(), full_suffix),
					Size:   info.Size(),
					Time:   info.ModTime(),
				},
			)
		}
	}

	//sort by time
	sort.Slice(job.Archive.RawList, func(i, j int) bool {
		return job.Archive.RawList[i].Time.Before(job.Archive.RawList[j].Time)
	})

	// build FULL -> DIFF[] tree
	var current_full_item *JobArchiveFullItem = nil

	for index := range job.Archive.RawList {
		var ai_pointer = &job.Archive.RawList[index]

		if ai_pointer.IsFull {
			full_item := JobArchiveFullItem{
				Item:         ai_pointer,
				DiffItemList: make([]*JobArchiveItem, 0),
			}

			job.Archive.FullItemList = append(job.Archive.FullItemList, full_item)

			//take address of just appended element
			current_full_item = &job.Archive.FullItemList[len(job.Archive.FullItemList)-1]
		} else {
			if current_full_item != nil { //skip diffs without parent
				current_full_item.DiffItemList = append(current_full_item.DiffItemList, ai_pointer)
			}
		}
	}

	//calculate diff sizes
	for index := range job.Archive.FullItemList {
		full_item := &job.Archive.FullItemList[index]

		if full_item.Item.Size == 0 {
			continue
		}

		var total_diff_size int64 = 0

		for index := range full_item.DiffItemList {
			total_diff_size += full_item.DiffItemList[index].Size
		}

		full_item.TotalDiffSizePercent = int(total_diff_size * 100 / full_item.Item.Size)
	}
}

func (ja *JobArchive) Dump() {
	fmt.Println("------ RAW LIST ------- ")
	for _, raw_item := range ja.RawList {
		fmt.Println("RAW: "+raw_item.Name, raw_item.Size)
	}

	fmt.Println("------------- ")
	for index := range ja.FullItemList {
		full_item := &ja.FullItemList[index]

		fmt.Printf("FULL: %s, size: %d, diff_size: %d%%\n", full_item.Item.Name, full_item.Item.Size, full_item.TotalDiffSizePercent)

		for index := range full_item.DiffItemList {
			diff_item := full_item.DiffItemList[index]
			fmt.Println("    DIFF: "+diff_item.Name, diff_item.Size)
		}
	}

	os.Exit(0)
}

func (afi *JobArchiveFullItem) Unlink() {
	//delete diffs
	for _, v := range afi.DiffItemList {
		if err := os.Remove(v.Path); err != nil {
			log.Fatalf("Error deleting file %s: %s", v.Path, err)
		}
	}

	//delete itself
	if err := os.Remove(afi.Item.Path); err != nil {
		log.Fatalf("Error deleting file %s: %s", afi.Item.Path, err)
	}
}
