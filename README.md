# mtsaver - differential backup archives management

[![Go Report Card](https://goreportcard.com/badge/github.com/mitoteam/mtsaver)](https://goreportcard.com/report/github.com/mitoteam/mtsaver)
![GitHub](https://img.shields.io/github/license/mitoteam/mtsaver)

[![GitHub Version](https://img.shields.io/github/v/release/mitoteam/mtsaver?logo=github)](https://github.com/mitoteam/mtsaver)
[![GitHub Release Date](https://img.shields.io/github/release-date/mitoteam/mtsaver)](https://github.com/mitoteam/mtsaver/releases)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/mitoteam/mtsaver)
[![GitHub contributors](https://img.shields.io/github/contributors-anon/mitoteam/mtsaver)](https://github.com/mitoteam/mtsaver/graphs/contributors)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/mitoteam/mtsaver)](https://github.com/mitoteam/mtsaver/commits)

Simple directory differential backup command line solution. Based on 7-Zip archiver available for all major platforms.

Trying to be self-explainatory (`mtsaver help`). More docs and manuals are being added.

## Basic idea

Lets assume we have huge enought working directory (~12 Gb) we want to backup periodicly (for example _daily_). And we want to have history of its changes (for example to be able to restore some file deleted forty two backups ago). Lets assume we want to keep history for _180 days_.

Day by day we work with very small amount of files: some new documents added, rarely some old ones are edited, even more rarely something being deleted.

Packing folder with 7-Zip creates 5Gb arhive. Creating 180 archives (one per day) will give us something like this:

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

Problems: 1) it takes 910Gb of storage 2) takes a lot of time to create each archive 3) requires some manual file management to delete outdated arhives beyond 180 days window

**mtsaver** is a solution to **pack only those new, changed or removed files** (thank you 7-Zip for anti-items support!) saving time and storage.
This is so called differential backups: we can have full archive with all files packed. Than we create differential archive containing only changed files.

**mtsaver** allows very flexible setups. For example: Create full archive once per month. Each day create differential archive. If differential archive becames larger than 50% of full archive do not wait till month end and create new full archive immediately. Store maximum 6 full archives and delete older ones with all differential archives attached. But keep full archives at least 180 days even if there are more than 6 of them.

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

Compare **910.5 Gb** of storage space taken and **37.5 Gb**. We saved more than 95% of storage space and yet we have full half-year history for directory!

## Usage

Build from sources (`make`) or unpack one of pre-compiled binaries.

Run `mtsaver help` for help.

## Inspired by

- https://www.cobiansoft.com/about.html
- https://nagimov.me/post/simple-differential-and-incremental-backups-using-7-zip/
- https://www.7-zip.org
