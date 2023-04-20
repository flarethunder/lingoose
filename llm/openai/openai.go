package openai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/henomis/lingoose/chat"
	"github.com/sashabaranov/go-openai"
)

type Model string

const (
	GPT4               Model = openai.GPT4
	GPT3Dot5Turbo      Model = openai.GPT3Dot5Turbo
	GPT3TextDavinci003 Model = openai.GPT3TextDavinci003
	GPT3TextDavinci002 Model = openai.GPT3TextDavinci002
	GPT3TextCurie001   Model = openai.GPT3TextCurie001
	GPT3TextBabbage001 Model = openai.GPT3TextBabbage001
	GPT3TextAda001     Model = openai.GPT3TextAda001
	GPT3TextDavinci001 Model = openai.GPT3TextDavinci001
	GPT3Davinci        Model = openai.GPT3Davinci
	GPT3Curie          Model = openai.GPT3Curie
	GPT3Ada            Model = openai.GPT3Ada
	GPT3Babbage        Model = openai.GPT3Babbage
)

type OpenAI struct {
	openAIClient *openai.Client
	model        Model
	verbose      bool
}

func New(model Model, verbose bool) (*OpenAI, error) {

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	return &OpenAI{
		openAIClient: openai.NewClient(openAIKey),
		model:        model,
		verbose:      verbose,
	}, nil
}

func (o *OpenAI) Completion(prompt string) (string, error) {

	response, err := o.openAIClient.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model:  openai.GPT3TextDavinci003,
			Prompt: prompt,
		},
	)

	if err != nil {
		return "", fmt.Errorf("openai: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	output := strings.TrimSpace(response.Choices[0].Text)
	if o.verbose {
		fmt.Printf("---USER---\n%s\n", prompt)
		fmt.Printf("---AI---\n%s\n", output)
	}

	return output, nil
}

func (o *OpenAI) Chat(prompt *chat.Chat) (interface{}, error) {

	var messages []openai.ChatCompletionMessage
	promptMessages, err := prompt.ToMessages()
	if err != nil {
		return nil, err
	}

	for _, message := range promptMessages {
		if message.Type == chat.MessageTypeUser {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: message.Content,
			})
		} else if message.Type == chat.MessageTypeAssistant {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: message.Content,
			})
		} else if message.Type == chat.MessageTypeSystem {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: message.Content,
			})
		}
	}

	response, err := o.openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("openai: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	content := response.Choices[0].Message.Content

	message := &chat.Message{
		Type:    chat.MessageTypeAssistant,
		Content: content,
	}

	if o.verbose {
		for _, message := range promptMessages {
			if message.Type == chat.MessageTypeUser {
				fmt.Printf("---USER---\n%s\n", message.Content)
			} else if message.Type == chat.MessageTypeAssistant {
				fmt.Printf("---AI---\n%s\n", message.Content)
			} else if message.Type == chat.MessageTypeSystem {
				fmt.Printf("---SYSTEM---\n%s\n", message.Content)
			}
		}
		fmt.Printf("---AI---\n%s\n", content)
	}

	return message, nil
}