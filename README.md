# mtsaver - differential backup archives management (7-zip based)

[![Go Report Card](https://goreportcard.com/badge/github.com/mitoteam/mtsaver)](https://goreportcard.com/report/github.com/mitoteam/mtsaver)
![GitHub](https://img.shields.io/github/license/mitoteam/mtsaver)

[![GitHub Version](https://img.shields.io/github/v/release/mitoteam/mtsaver?logo=github)](https://github.com/mitoteam/mtsaver)
[![GitHub Release Date](https://img.shields.io/github/release-date/mitoteam/mtsaver)](https://github.com/mitoteam/mtsaver/releases)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/mitoteam/mtsaver)
[![GitHub contributors](https://img.shields.io/github/contributors-anon/mitoteam/mtsaver)](https://github.com/mitoteam/mtsaver/graphs/contributors)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/mitoteam/mtsaver)](https://github.com/mitoteam/mtsaver/commits)

7-Zip based simple directory backup command line solution.

Trying to be self-explainatory (`mtsaver help`). More docs will be added later.

## Basic idea

Lets assume we have huge enought working directory (~12 Gb) we want to backup periodicly (for example _daily_). And we want to have history of its changes (for example to be able to restore some file deleted twelve backups ago). Lets assume we want to keep history for _180 days_.

Packing folder with 7-Zip creates 5Gb arhive. Creating 180 archives (one per day) will 1) take 900Gb of storage 2) takes a lot of time to create each archive 3) requires some manual file management to delete outdated arhives beyond 180 days window.

Day by day we work with very small amount of files: some new documents added, rarely some old ones are edited, even more rarely something being deleted.

**mtsaver** is a solution to **pack only those new, changed or removed files** (thank you 7-Zip for anti-items support!) saving time and storage.

## Usage

Build from sources (`make`) or unpack one of pre-compiled binaries.

Run `mtsaver help` for help.

## Inspired by
- https://www.cobiansoft.com/about.html
- https://nagimov.me/post/simple-differential-and-incremental-backups-using-7-zip/
- https://www.7-zip.org
