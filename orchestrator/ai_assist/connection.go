package ai_assist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"orchestrator/pkg/logger"
	"os"
	"time"
)

func Connection(ctx context.Context, key string, description string) (serviceName string, err error) {
	prompt, err := ReadPrompt()
	if err != nil {
		logger.GetLoggerFromContext(ctx).Error(ctx, err.Error())
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "meta-llama/llama-3.1-70b-instruct:free", // модель нейронки
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf("%s %s", prompt, description),
			},
		},
	})
	if err != nil {
		logger.GetLoggerFromContext(ctx).Error(ctx, err.Error())
	}

	resp, err := SendRequest(ctx, requestBody, key)
	if err != nil {
		waitTime := 1 * time.Second
		for i := 0; i < 5; i++ {
			logger.GetLoggerFromContext(ctx).Info(ctx, "AI assistant doesnt answer. Retry to connect again...")
			time.Sleep(waitTime)
			resp, err = SendRequest(ctx, requestBody, key)
			if err == nil {
				break
			}
		}
		return "Service define fail", errors.New("Service define fail")
	}

	logger.GetLoggerFromContext(ctx).Info(ctx, "Response from AI took successfully")

	if resp == "Service define fail" {
		return resp, errors.New("Service define fail")
	}

	return resp, nil
}

func ReadPrompt() (string, error) {
	jsonFile, err := os.Open("ai_assist/prompt.json")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	return string(byteValue), nil
}

func SendRequest(ctx context.Context, requestBody []byte, key string) (response string, err error) {
	req, err := http.NewRequestWithContext(context.Background(),
		"POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		logger.GetLoggerFromContext(ctx).Error(ctx, err.Error())
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Request to AI created successfully")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))
	logger.GetLoggerFromContext(ctx).Info(ctx, "AI assist authorization completed successfully")

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.GetLoggerFromContext(ctx).Error(ctx, err.Error())
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Request sent successfully to https://openrouter.ai/api/v1/chat/completions")

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	aiResponse, err := Transform(responseBody)
	if err != nil {
		return "", err
	}

	return aiResponse, nil
}
