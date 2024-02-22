package vx

import "time"

// TextBucket is an atomic unit a piece of text being explained
type TextBucket struct {
	// Start of the text. This should be a relative time only. I.e. StartTime=00:00:00
	Start time.Time
	// End of the text. This should be a relative time only. I.e. EndTime=00:00:05
	End time.Time
	// The text of the bucket, which happens between start time and end time.
	Text string
}
