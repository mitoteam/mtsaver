// Package app contains main mtsaver functionality.
package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mitoteam/mttools"
	"github.com/spf13/cobra"
)

var BuildVersion = "DEV"
var BuildCommit = "DEV"

var Global struct {
	AppName      string
	Version      string
	Commit       string
	BuiltWith    string
	SevenZipCmd  string
	SevenZipInfo string
}

func init() {
	Global.AppName = "mtsaver"
	Global.Version = BuildVersion
	Global.Commit = BuildCommit
	Global.BuiltWith = runtime.Version()
}

func SetupBeforeCommand(cmd *cobra.Command, args []string) error {
	if JobRuntimeOptions.NoConsole {
		if mttools.IsWindows() {
			mttools.HideConsole()
		} else {
			fmt.Println("--no-console option ignored under Linux")
		}
	}

	if Global.SevenZipCmd == "auto" {
		Global.SevenZipCmd = "" //reset to force autodetection
	}

	if len(Global.SevenZipCmd) > 0 {
		if !checkSevenZipCommand(Global.SevenZipCmd) {
			return errors.New("Can not run provided 7-Zip command: " + Global.SevenZipCmd)
		}
	} else {
		//try autodetect
		//try raw 7z command
		r := checkSevenZipCommand("7z")

		if !r {
			switch runtime.GOOS {
			case "windows":
				r = checkSevenZipCommand(os.Getenv("ProgramFiles") + "\\7-Zip\\7z.exe")

				if !r {
					r = checkSevenZipCommand(os.Getenv("ProgramFiles(x86)") + "\\7-Zip\\7z.exe")
				}
			case "linux":
				r = checkSevenZipCommand("/usr/lib/p7zip/7z")
				if !r {
					r = checkSevenZipCommand("/usr/bin/7z")
				}
				if !r {
					r = checkSevenZipCommand("/bin/7z")
				}
			}
		}

		if !r {
			return errors.New("Can not find 7-Zip. Please provide correct path with --7zip flag.")
		}
	}

	return nil
}

func checkSevenZipCommand(cmd string) bool {
	out, err := exec.Command(cmd).Output()

	if err != nil {
		return false
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		if scanner.Text() != "" {
			Global.SevenZipInfo = scanner.Text()
			Global.SevenZipCmd = cmd
			return true
		}
	}

	return false
}
