package openaiutil

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func Chart(token, prompt, diff, model, baseURL string) (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(token),
		option.WithBaseURL(baseURL),
	)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(diff),
			openai.SystemMessage(prompt),
		},
		Model: model,
	})
	if err != nil {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func ChartWithStream(token, prompt, diff, model, baseURL string) {
	ctx := context.Background()

	client := openai.NewClient(
		option.WithAPIKey(token),
		option.WithBaseURL(baseURL),
	)

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(diff),
			openai.SystemMessage(prompt),
		},
		Seed:  openai.Int(0),
		Model: model,
	})
	defer stream.Close()

	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if content, ok := acc.JustFinishedContent(); ok {
			println("Content stream finished:", content)
		}

		// if using tool calls
		if tool, ok := acc.JustFinishedToolCall(); ok {
			println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			println("Refusal stream finished:", refusal)
		}

		// it's best to use chunks after handling JustFinished events
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Println(chunk.Choices[0].Delta.Content)
			// printWithTypewriterEffect(chunk.Choices[0].Delta.Content, &lastPos)
		}
	}

	if stream.Err() != nil {
		// panic(stream.Err())
		println("失败：", stream.Err())
	}

	// After the stream is finished, acc can be used like a ChatCompletion
	_ = acc.Choices[0].Message.Content
}
