package helper

import (
	"crypto/tls"
	"net/http"
)

func DoRequest(url string) (int, error) {
	res, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: transport}
	resp, err := httpClient.Do(res)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
