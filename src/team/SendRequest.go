package team

import (
	projectConfig "ApiChernoToRustat/config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"runtime"
)

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func (g *GraphQLRequest) SendPostRequest(graphqlRequest GraphQLRequest, config projectConfig.Config) (map[string]interface{}, error) {
	jsonData, marshalErr := json.Marshal(graphqlRequest)
	if marshalErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Ошибка marshal тела запроса: file: %s | line: %d | error: %v", file, line-2, marshalErr)
	}

	reader := bytes.NewReader(jsonData)

	url := fmt.Sprintf("https://%s/graphql", config.DirectusURL)
	req, createReqErr := http.NewRequest("POST", url, reader)
	if createReqErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Ошибка создания запроса с url: %s | file: %s | line: %d | error: %v", url, file, line-2, createReqErr)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}

	resp, reqErr := client.Do(req)
	if reqErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Ошибка отправки запроса по url: %s | file: %s | line: %d | error: %v", url, file, line-2, reqErr)
	}

	var result map[string]interface{}
	if decoderErr := json.NewDecoder(resp.Body).Decode(&result); decoderErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Ошибка при decode тела ответа по url:  %s | file: %s | line: %d | error: %v", url, file, line-1, decoderErr)
	}

	if closeErr := resp.Body.Close(); closeErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Ошибка при попытке закрыть тело ответа по url: %s | file: %s | line: %d | error: %v", url, file, line-1, closeErr)
	}

	return result, nil
}

func (g *GraphQLRequest) SendGetRequest(rustatApi string) (map[string]interface{}, error) {
	client := resty.New()
	resp, reqErr := client.R().Get(rustatApi)
	if reqErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Ошибка отправки get запроса по url: %s | file: %s | line: %d | error: %v", rustatApi, file, line-2, reqErr)
	}

	var result map[string]interface{}
	if unmarshalErr := json.Unmarshal(resp.Body(), &result); unmarshalErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Ошибка unmarshal в структуру из тела ответа по url: %s | file: %s | line: %d | error: %v", rustatApi, file, line-2, unmarshalErr)
	}

	return result, nil
}

// Я ИСПОЛЬЗУЮ ДВЕ РАЗНЫЕ БИБЛЫ ДЛЯ HTTP ЗАПРОСОВ Т.К:
// POST запросы работают только с базовой библой
// GET работает только с RESTY :). Я не даун, я просто искренне не понимаю
