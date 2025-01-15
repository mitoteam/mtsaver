package app

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/mitoteam/mttools"
)

type JobArchiveFile struct {
	Name    string    //filename only
	Path    string    //full path
	IsFull  bool      //full or diff archive
	Size    int64     //file size
	ModTime time.Time //modification time
	Time    time.Time //timestamp from archive name
	Age     int       //age in days
	Hash    string    //sha256 for archives
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
	FilesList    []JobArchiveFile     // All archives raw list
	FullItemList []JobArchiveFullItem // Full archives list with diffs listed in DiffItemList
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
		if err != nil {
			continue
		}

		archive_file := JobArchiveFile{
			Name:    value.Name(),
			Path:    filepath.Join(job.Settings.ArchivesPath, value.Name()),
			IsFull:  strings.HasSuffix(value.Name(), full_suffix),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}

		//try to parse timestamp
		var filename_time time.Time

		matches := re.FindAllStringSubmatch(archive_file.Name, 1)
		if len(matches) > 0 {
			if len(matches[0]) > 1 {
				filename_time, _ = time.Parse(job.Settings.DateFormat, matches[0][1])
			}
		}

		if filename_time.IsZero() {
			archive_file.Time = archive_file.ModTime
		} else {
			archive_file.Time = filename_time
		}

		archive_file.Age = int(math.Ceil(time.Since(archive_file.Time).Hours() / 24))

		//calculate hash for diffs
		if !job.Settings.KeepSameDiff && !archive_file.IsFull {
			var hash string
			hash, err = mttools.FileSha256(archive_file.Path)

			if err == nil {
				archive_file.Hash = hash
			} else {
				log.Fatal(err)
			}
		}

		job.Archive.FilesList = append(job.Archive.FilesList, archive_file)
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

	job.Log(
		"Archives scan done. Total archives: %d. Full archives: %d",
		len(job.Archive.FilesList), len(job.Archive.FullItemList),
	)

	if len(job.Archive.FullItemList) > 0 {
		firstFullArch := job.Archive.FullItemList[0]

		job.Log(
			"Oldest archive: %s, age (days): %d",
			firstFullArch.File.Name, firstFullArch.File.Age,
		)

		lastFullArch := job.Archive.FullItemList[len(job.Archive.FullItemList)-1]
		lastFullArchInfo := fmt.Sprintf(
			"Newest full archive: %s, age (days): %d, diffs count: %d, total diffs size: %d%%",
			lastFullArch.File.Name, lastFullArch.File.Age, len(lastFullArch.DiffItemList),
			lastFullArch.TotalDiffSizePercent,
		)

		if len(lastFullArch.DiffItemList) > 0 {
			lastDiffArch := lastFullArch.DiffItemList[len(lastFullArch.DiffItemList)-1]
			lastFullArchInfo += fmt.Sprintf(", last diff size: %d%%", lastDiffArch.DiffSizePercent)
		}

		job.Log("%s", lastFullArchInfo)
	}
}

func (ja *JobArchive) LastFile() *JobArchiveFile {
	if ja.FilesList == nil {
		return nil
	}

	if len(ja.FilesList) == 0 {
		return nil
	}

	return &ja.FilesList[len(ja.FilesList)-1]
}

func (ja *JobArchive) Dump(die bool) {
	fmt.Println("------ PLAIN ARCHIVES LIST -------")
	for _, raw_file := range ja.FilesList {
		fmt.Println(raw_file.Name, mttools.FormatFileSize(raw_file.Size))
	}

	fmt.Println("\n------ DIFFs TREE -------")
	for index := range ja.FullItemList {
		full_item := &ja.FullItemList[index]

		info_str := fmt.Sprintf("size: %s, age: %d", mttools.FormatFileSize(full_item.File.Size), full_item.File.Age)

		if len(full_item.DiffItemList) > 0 {
			info_str += fmt.Sprintf(", diffs: %d (size %d%%)", len(full_item.DiffItemList), full_item.TotalDiffSizePercent)
		}

		fmt.Printf("FULL: %s, %s\n", full_item.File.Name, info_str)

		for _, diff_item := range full_item.DiffItemList {
			fmt.Printf("    DIFF: %s, size %s = %d%%, age: %d\n", diff_item.File.Name, mttools.FormatFileSize(diff_item.File.Size), diff_item.DiffSizePercent, diff_item.File.Age)
			if len(diff_item.File.Hash) > 0 {
				fmt.Println("    " + diff_item.File.Hash)
			}
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
