package arbor

import "time"

// StartTime returns a stable process start time used for "uptime" reporting.
func StartTime() time.Time {
	return processStart
}

var processStart = time.Now()
