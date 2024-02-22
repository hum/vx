# vx
Explain videos in plain text with LLMs.

## Usage
### As a CLI
The CLI allows specifying an alternative API to use (only tested with [Perplexity](https://perplexity.ai)), as well as a custom prompt or a model. For all options use:
```bash
> vx --help
```

Vx supports both streamed and non-streamed responses on the CLI. Use the `--stream` flag to stream to STDOUT.
```bash
> vx --url "url" --stream
```
### As a package
```go
package main

import (
  "fmt"

  "github.com/hum/vx"
  "github.com/sashabaranov/go-openai"
)

func main() {
  oapi := openai.NewClient("token")

  r, err := vx.GetVideoExplanationRequests(client, vx.VideoExplanationOpts{
    Url: "url",
    Prompt: "Give me 5 bullet points from this text: ",
    ChunkSize: 5 * time.Minute,
  })
  if err != nil {
    panic(err)
  }

  response, err := client.CreateChatCompletion(context.TODO(), r)
  if err != nil {
    panic(err)
  }
  fmt.Println(response.Choices[0].Message.Content)
}
```
