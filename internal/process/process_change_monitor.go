package process

import (
	"fmt"

	"github.com/ochinchina/filechangemonitor"
)

const fileMonitorWorkers = 10 // Number of worker goroutines for file monitoring

var fileChangeMonitor = filechangemonitor.NewFileChangeMonitor(fileMonitorWorkers)

// AddProgramChangeMonitor adds program change listener to monitor if the program binary.
func AddProgramChangeMonitor(path string, fileChangeCb func(path string, mode filechangemonitor.FileChangeMode)) {
	_ = fileChangeMonitor.AddMonitorFile(path,
		false,
		filechangemonitor.NewExactFileMatcher(path),
		filechangemonitor.NewFileChangeCallbackWrapper(fileChangeCb),
		filechangemonitor.NewFileMD5CompareInfo()) // Ignore monitor registration error
}

// AddConfigChangeMonitor adds program change listener to monitor if any of its configuration files is changed.
func AddConfigChangeMonitor(path string, filePattern string, fileChangeCb func(path string, mode filechangemonitor.FileChangeMode)) {
	fmt.Printf("filePattern=%s\n", filePattern)
	_ = fileChangeMonitor.AddMonitorFile(path,
		true,
		filechangemonitor.NewPatternFileMatcher(filePattern),
		filechangemonitor.NewFileChangeCallbackWrapper(fileChangeCb),
		filechangemonitor.NewFileMD5CompareInfo()) // Ignore monitor registration error
}
