package circleci

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type respStatus struct {
	Indiactor   string `json:"indicator"`
	Description string `json:"description"`
}

type response struct {
	Status respStatus `json:"status"`
}

func fetchStatus() (*response, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://status.circleci.com/api/v2/status.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("API has an error")
	}

	var result response
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}
