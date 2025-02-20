# mtsaver - differential backup archives management

[![Go Report Card](https://goreportcard.com/badge/github.com/mitoteam/mtsaver)](https://goreportcard.com/report/github.com/mitoteam/mtsaver)
![GitHub](https://img.shields.io/github/license/mitoteam/mtsaver)
![GitHub code size](https://img.shields.io/github/languages/code-size/mitoteam/mtsaver)
[![GitHub contributors](https://img.shields.io/github/contributors-anon/mitoteam/mtsaver)](https://github.com/mitoteam/mtsaver/graphs/contributors)
[![Build&Tests](https://github.com/mitoteam/mtsaver/actions/workflows/go.yml/badge.svg)](https://github.com/mitoteam/mtsaver/actions/workflows/go.yml)

[![GitHub Version](https://img.shields.io/github/v/release/mitoteam/mtsaver?logo=github)](https://github.com/mitoteam/mtsaver)
[![GitHub Release Date](https://img.shields.io/github/release-date/mitoteam/mtsaver)](https://github.com/mitoteam/mtsaver/releases)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/mitoteam/mtsaver)](https://github.com/mitoteam/mtsaver/commits)
[![GitHub downloads](https://img.shields.io/github/downloads/mitoteam/mtsaver/total)](https://github.com/mitoteam/mtsaver/releases)

Simple directory differential backups command-line utility. Based on 7-Zip archiver available for all major platforms.

**Basics of differential backups idea**:

1) Create full directory archive for first time backup is performed.
2) Create "_diff_" archive with only changed, added (or removed!) files every next time...
3) ...until certain criteria met (like "_total diff archives size exceeds full archive size_" or "_diff archive is more than 20% of full archive size_". Create new full archive at this point (so diffs will be small again).
4) Remove outdated full and diff archives (according to retention settings).

**Made to be simple yet powerful**:

* No dependencies. Only 7-Zip is required! So no need for special software to explore or unpack created backups.
* No database required: archiving and retention settings are in `.mtsaver.yml` file, other state info gathered on the fly from archives directory.
* No installation required. Designed to be used as-is. Distributed as single executable file.
* Both Windows and Linux, both x64 and x32 platforms support.
* Flexible full/diff archives retention settings.
* Ability to run some commands before or after backup creation.

Trying to be self-explanatory (`mtsaver help`). More docs and manuals are being added.

## Differential backup idea explained

Lets assume we have huge enough working directory (~12 Gb) we want to backup periodically (for example _daily_). And we want to have history of its changes (for example to be able to restore some file deleted forty two backups ago). Lets assume we want to keep history for _180 days_.

Day by day we work with very small amount of files: some new documents added, rarely some old ones are edited, even more rarely something being deleted.

Packing folder with 7-Zip creates 5Gb archive. Creating 180 archives (one per day) will give us something like this:

```text
backup_2022-01-01.7z 5Gb
backup_2022-01-02.7z 5Gb
backup_2022-01-03.7z 5Gb
...
backup_2022-03-14.7z 5Gb
backup_2022-03-15.7z 5.1Gb #something big was added in March, archives became bigger
backup_2022-03-16.7z 5.1Gb
backup_2022-03-17.7z 5.1Gb
...
backup_2022-06-30.7z 5.1Gb
---------------------------
total: 180 files, 910.5 Gb
```

Problems: 1) it takes 910Gb of storage 2) takes a lot of time to create each archive 3) requires some manual file management to delete outdated archives beyond 180 days window.

**mtsaver** is a solution to **pack only those new, changed or removed files** (thank you 7-Zip for anti-items support!) saving CPU-time and storage.
This is so called differential backups: we can have full archive with all files packed. Than we create differential archive containing only changed files.

**mtsaver** allows very flexible setups. For example: Create full archive once per month. Each day create differential archive. If differential archive becomes larger than 50% of full archive do not wait till month end and create new full archive immediately. Store maximum 6 full archives and delete older ones with all differential archives attached. But keep full archives at least 180 days even if there are more than 6 of them.

So for example above we will end with following files set:

```text
backup_2022-01-01_FULL.7z 5Gb
backup_2022-01-02_DIFF.7z 100Kb # changes since Jan-01 only!
backup_2022-01-03_DIFF.7z 150Kb # changes since Jan-01 only
...
backup_2022-01-30_DIFF.7z 32Mb  # changes since Jan-01 only
backup_2022-01-31_FULL.7z 5Gb   # 30 days passed, new full archive created
backup_2022-02-01_DIFF.7z 110Kb # changes since last full archive only
...
backup_2022-03-14_DIFF.7z 15Mb
backup_2022-03-15_DIFF.7z 2.8Gb # something big was added or a lot of documents
                                # was changed, differential backup became more
                                # than 50% of full backup
backup_2022-03-16_FULL.7z 5.1Gb # so next time new full archive created before 30 days
                                # window passed
backup_2022-03-17_DIFF.7z 140Kb # changes since last full archive only
...
backup_2022-06-30_DIFF.7z 51Mb  # changes since last full archive only
--------------------------------
total: 180 files, 37.5 Gb       # 6 planned full archives +1 unplanned in March (35.4Gb)
                                # 173 differential archives each of 12Mb average
```

Compare **910.5 Gb** of storage space taken and **37.5 Gb**. We saved more than **95%** of storage space and yet we have **full half-year day-by-day history** for directory!

Combine this with simple [syncthing](https://syncthing.net) setup for archives folder and you have differential cloud stored backups for free.

## Installation

### Using Scoop

Scoop is useful command-line installer and updater for Windows.

* Install `scoop`: https://github.com/ScoopInstaller/Scoop/wiki/Quick-Start

* Add bucket:

```sh
scoop bucket add mitoteam https://github.com/mitoteam/scoop-bucket
```

* Install:

```sh
scoop install mtsaver
```

Details: [scoop-bucket](https://github.com/mitoteam/scoop-bucket).

### Manual installation (Windows or Linux)

* Download latest release from [Releases](https://github.com/mitoteam/mtsaver/releases) page.
* Unpack with 7-zip.

## Usage

First you need to create a file with archiving settings: what to pack, where to pack, retention rules and so on. Just run `mtsaver init` in directory you want to backup. By default this file has `.mtsaver.yml` name. This will create `.mtsaver.yml` file with [all](app/job_settings.go) possible settings and explanations.

Open created `.mtsaver.yml` file and edit settings you need. You can remove untouched default settings from this file to keep things simple. Important settings are:

* **archives_path** path to store created archives.
* **keep_at_least** number of days to keep oldest full archive despite all other retention settings.

Default retention rules are: keep max 5 full archives, keep max 20 diff archives for each full archive, force full archive if previous diff size is 120% or more of latest full archive, do not store empty or unchanged diff archives, add _*.7z_ and _*.rar_ files to archives without compression (assuming they are already packed).

Run `mtsaver run` command in directory with `.mtsaver.yml` file to create new backup archive. First time it will be created as full archive. Next runs depending on conditions and settings either full or diff archives will be created and old ones will be removed.

You can use any scheduler (`cron` or _Windows Task Scheduler_) to run this command regularly to have your directory backups.

By default mtsaver creates file `_mtsaver.log` file in archives directory with archiving logs. It has explanations why full or diff archive was created. You can disable log file by setting `log_format:` option to _disable_ in `.mtsaver.yml` file (or use `--no-log` command-line argument).

## Help

Run `mtsaver help` for options and commands description.

## Inspired by

* https://www.cobiansoft.com/about.html
* https://nagimov.me/post/simple-differential-and-incremental-backups-using-7-zip/
* https://www.7-zip.org
