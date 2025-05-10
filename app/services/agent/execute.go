package agent_service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"strings"

	cfg "hintword.com/api/app/configs"
)

// loadSystemPrompt reads the content of the prompt file based on the type
func loadSystemPrompt(promptType string) (string, error) {
	projectRoot, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get project root: %v", err)
	}
	promptFile := fmt.Sprintf("%s/prompt/%s.txt", projectRoot, promptType)
	content, err := os.ReadFile(promptFile)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file: %v", err)
	}
	return string(content), nil
}

func Execute(promptType string, note string) (response Response, err error) {
	systemPrompt, err := loadSystemPrompt(promptType)
	if err != nil {
		return
	}

	payload := strings.NewReader(fmt.Sprintf(`{
		"model": "openai/gpt-4.1",
		"temperature": 1,
		"top_p": 1,
		"messages": [
			{
				"role": "system",
				"content": %q
			},
			{
				"role": "user",
				"content": %q
			}
		]
	}`, systemPrompt, note))

	client := &http.Client{}
	req, err := http.NewRequest("POST", cfg.GetConfig().OpenApiUrl, payload)

	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cfg.GetConfig().OpenApiKey))

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response)
	return
}
