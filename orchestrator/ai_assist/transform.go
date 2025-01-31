package ai_assist

import (
	"encoding/json"
	"errors"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Refusal string `json:"refusal"`
}

type Choice struct {
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletion struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"`
	Model     string    `json:"model"`
	Object    string    `json:"object"`
	Created   int64     `json:"created"`
	Choices   []Choice  `json:"choices"`
	Usage     Usage     `json:"usage"`
	CreatedAt time.Time `json:"-"`
}

func Transform(data []byte) (string, error) {
	var completion ChatCompletion
	err := json.Unmarshal(data, &completion)
	if err != nil {
		return "", err
	}

	completion.CreatedAt = time.Unix(completion.Created, 0)

	if completion.Choices == nil {
		return "", errors.New("Nil response from AI assistant")
	}

	return completion.Choices[0].Message.Content, nil
}
