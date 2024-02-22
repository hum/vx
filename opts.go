package vx

import "time"

// VideoExplanationOpts encapsulates all options to use for video explanation.
type VideoExplanationOpts struct {
	// The video URL to the content
	Url string
	// Name of the model used as an LLM
	Model *string
	// Whether to stream the request
	Stream bool
	// The default prompt to be appended before each TextBucket
	Prompt string
	// Time interval of the chunked text. If the text is an hour long and the duration is set to 30*time.Minute, there will be 2 buckets internally.
	ChunkSize time.Duration
}
