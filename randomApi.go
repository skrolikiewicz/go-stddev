package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type randomApiRequest struct {
	JsronRPC string          `json:"jsonrpc"`
	Method   string          `json:"method"`
	Params   randomApiParams `json:"params"`
	ID       int             `json:"id"`
}

type randomApiParams struct {
	APIKey      string `json:"apiKey"`
	N           int    `json:"n"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Replacement bool   `json:"replacement"`
}

type randomApiResponse struct {
	Result *randomApiResult `json:"result"`
	Error  *randomApiError  `json:"error"`
	ID     int              `json:"id"`
}

type randomApiResult struct {
	Random randomApiResultData `json:"random"`
}

type randomApiResultData struct {
	Data []int `json:"data"`
}

type randomApiError struct {
	Message string `json:"message"`
}

func getRandomNumbers(ctx context.Context, length int) (randomNumbers []int, err error) {
	url := "https://api.random.org/json-rpc/2/invoke"
	apiKey := os.Getenv("RANDOM_ORG_API_KEY")
	params := randomApiParams{apiKey, length, 1, 1000, true}
	payload := randomApiRequest{"2.0", "generateIntegers", params, length}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.WithContext(ctx)

	httpClient := http.Client{
		Timeout: time.Second * 30,
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Unsuccessful status code returned from random.org API")
	}

	var result randomApiResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, errors.New(result.Error.Message)
	}

	return result.Result.Random.Data, nil
}
