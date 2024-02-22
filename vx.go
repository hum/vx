package vx

import (
	"time"

	"github.com/hum/vcat"
	"github.com/sashabaranov/go-openai"
)

// GetVideoExplanationStreamRequests retrieves the transcription of a given URL,
// then splits the transcript into smaller segments according to the provided opts.ChunkSize.
//
// It returns a slice of openai.ChatCompletionRequest objects that can be used to interactively stream responses from the API.
func GetVideoExplanationRequests(c *openai.Client, opts VideoExplanationOpts) ([]openai.ChatCompletionRequest, error) {
	transcript, err := vcat.GetTranscription(opts.Url, "en")
	if err != nil {
		return nil, err
	}

	chunks, err := chunkTranscriptByTimeDelta(transcript, opts.ChunkSize)
	if err != nil {
		return nil, err
	}

	requests := makeRequests(chunks, opts)
	return requests, nil
}

// Given a single TextBucket, this function constructs an individual request object for the API
// If no model is explicitly set via opts.Model, it defaults to 'mistral-7b-instruct'.
func NewChatCompletionRequest(chunk TextBucket, opts VideoExplanationOpts) *openai.ChatCompletionRequest {
	if opts.Model == nil {
		model := "mistral-7b-instruct"
		opts.Model = &model
	}

	return &openai.ChatCompletionRequest{
		Model:  *opts.Model,
		Stream: opts.Stream,
		Messages: []openai.ChatCompletionMessage{
			{
				Content: opts.Prompt + chunk.Text,
				Role:    openai.ChatMessageRoleUser,
			},
		},
	}
}

// makeRequest takes in a slice of TextBuckets and returns API requests based on the provided opts.
func makeRequests(chunks []TextBucket, opts VideoExplanationOpts) []openai.ChatCompletionRequest {
	var result = make([]openai.ChatCompletionRequest, 0, len(chunks))

	for _, chunk := range chunks {
		request := NewChatCompletionRequest(chunk, opts)
		result = append(result, *request)
	}
	return result
}

// chunkTranscriptByTimeDelta takes in a whole video transcript and splits it into buckets that are the length of the provided duration parameter d.
func chunkTranscriptByTimeDelta(transcript *vcat.Transcript, d time.Duration) ([]TextBucket, error) {
	var (
		result = make([]TextBucket, 0)
		tmp    = ""

		currLatestTime, _ = time.Parse(time.TimeOnly, transcript.Text[0].Start)
	)

	for _, text := range transcript.Text {
		var (
			startTime, _ = time.Parse(time.TimeOnly, text.Start)
			endTime, _   = time.Parse(time.TimeOnly, text.End)
		)

		if startTime.After(currLatestTime.Add(d)) {
			result = append(result, TextBucket{Start: currLatestTime, End: endTime, Text: tmp})
			currLatestTime = startTime
			tmp = ""
			continue
		}
		tmp += " " + text.Text
	}

	// Append any leftover data if it did not fit within the duration
	var (
		startTime, _ = time.Parse(time.TimeOnly, transcript.Text[len(transcript.Text)-1].Start)
		endTime, _   = time.Parse(time.TimeOnly, transcript.Text[len(transcript.Text)-1].End)
		text         = transcript.Text[len(transcript.Text)-1].Text
	)
	if tmp != "" {
		result = append(result, TextBucket{Start: startTime, End: endTime, Text: text})
	}
	return result, nil
}
