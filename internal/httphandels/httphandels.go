package httphandels

import (
	"fmt"
	"io"
	"net/http"
)

func HTTPGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	var body []byte
	if err != nil {
		return body, fmt.Errorf("Get error: %v\n", err)
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Reader error")
		return body, fmt.Errorf("Reader error: %v\n", err)
	}
	return body, nil
}
