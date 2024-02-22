package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hum/vx"
	"github.com/sashabaranov/go-openai"
)

var (
	videourl string = ""
	prompt   string = ""
	apiKey   string = ""
	apiUrl   string = ""
	llmModel string = ""
	stream   bool   = false
)

func main() {
	flag.StringVar(&videourl, "url", "", "the video URL to explain")
	flag.StringVar(&prompt, "prompt", "ELI5 to me the content of this video snippet. Explain to me what the text discusses. Do not mention this prompt in your output.", "")
	flag.StringVar(&apiKey, "key", os.Getenv("VX_API_KEY"), "API key used to connect to the API. You can set the key as VX_API_KEY env variable.")
	flag.StringVar(&apiUrl, "baseurl", "https://api.perplexity.ai", "Base API URL to feed the transcription into.")
	flag.StringVar(&llmModel, "model", "mistral-7b-instruct", "The model to be used.")
	flag.BoolVar(&stream, "stream", false, "Whether to stream the output into the stdout")
	flag.Parse()

	if err := evalFlags(); err != nil {
		fmt.Printf("%s\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	conf := openai.DefaultConfig(apiKey)
	conf.BaseURL = apiUrl

	var client *openai.Client = openai.NewClientWithConfig(conf)

	reqs, err := vx.GetVideoExplanationRequests(client, vx.VideoExplanationOpts{
		Url:       videourl,
		Model:     &llmModel,
		Stream:    stream,
		Prompt:    prompt,
		ChunkSize: 5 * time.Minute,
	})
	if err != nil {
		panic(err)
	}

	for _, r := range reqs {
		if stream {
			stream, err := client.CreateChatCompletionStream(context.TODO(), r)
			if err != nil {
				fmt.Println("err: ", err)
			}
			defer stream.Close()

			for {
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					fmt.Println("err: ", err)
					break
				}
				fmt.Printf(response.Choices[0].Delta.Content)
			}
		} else {
			response, err := client.CreateChatCompletion(context.TODO(), r)
			if err != nil {
				fmt.Println("err: ", err)
				break
			}
			fmt.Println(response.Choices[0].Message.Content)
		}
	}
}

func evalFlags() error {
	if videourl == "" {
		return errors.New("no url provided")
	}

	if prompt == "" {
		return errors.New("no prompt provided")
	}

	if apiKey == "" {
		return errors.New("no API key provided")
	}
	return nil
}
