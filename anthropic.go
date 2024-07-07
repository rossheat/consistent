package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ReqMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ReqBody struct {
	Model     string   `json:"model"`
	MaxTokens int      `json:"max_tokens"`
	Messages  []ReqMsg `json:"messages"`
}

type RespBody struct {
	Content []RespBodyContent `json:"content"`
}

type RespBodyContent struct {
	Text string `json:"text"`
}

type RespContentText struct {
	Answer string `json:"answer"`
}

func SendMessage(m Model) (string, error) {

	client := &http.Client{}

	content := fmt.Sprintf(
		"You MUST produce a correctly formatted JSON response to the following yes/no question '%v'. "+
			"If the 'question' is not a valid question your response will be no. "+
			"You can ONLY answer 'yes' or 'no' and you MUST respond in the following JSON format: "+
			"{'answer': <'yes'/'no'>}",
		m.TextInput.Value(),
	)

	reqBody := ReqBody{
		Model:     config.model,
		MaxTokens: 1024,
		Messages: []ReqMsg{
			{Role: "user", Content: content},
		},
	}

	bs, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := "https://api.anthropic.com/v1/messages"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("x-api-key", config.key)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("response status code %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var respBody RespBody
	json.Unmarshal(body, &respBody)

	var respJSON RespContentText
	json.Unmarshal([]byte(respBody.Content[0].Text), &respJSON)

	return respJSON.Answer, nil
}

func (m Model) AskQuestion() tea.Msg {
	results := LLMResults{yes: 0, no: 0}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	errChan := make(chan error, config.instances)

	for i := 0; i < config.instances; i++ {
		wg.Add(1)
		time.Sleep(time.Millisecond * time.Duration(config.delay))
		go func() {
			defer wg.Done()
			answer, err := SendMessage(m)
			if err != nil {
				errChan <- err
				return
			}

			mutex.Lock()
			defer mutex.Unlock()
			if answer == "yes" {
				results.yes++
			} else if answer == "no" {
				results.no++
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return ErrMsg(fmt.Errorf("error when calling the API: %v", err))
		}
	}

	return results
}
