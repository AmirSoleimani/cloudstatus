package gcp

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// TODO use https://status.cloud.google.com/incidents.json?

func fetchStatus() (bool, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://status.cloud.google.com/")
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New("API has an error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if !strings.Contains(string(body), "All services available") {
		return false, nil
	}

	return true, nil
}
