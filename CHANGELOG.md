# CHANGELOG

## 1.1.1 (July 6th, 2021)

Change:

- Rebuilt using latest version of goApiLib, to fix possible issue with connections via a proxy

## 1.1.0 (2nd September, 2020)

Changes:

- Fixed issue where tasks were not being cancelled
- Removed zone argument
- Removed userId and password override args, as API Key was mandatory anyway
- Removed unnecessary ___dryrun___ arg, as all it did was essentially count rows in the input file, or return 1 for when a single task record is provided
- Replaced references to calls with tasks
- Improved CLI output for SUCCESS and FAIL messages
- General tidy-up, removed orphaned code
- Tidied up docs, included all command line args

## 1.0.0 (22nd July 2019)

### Initial Release
