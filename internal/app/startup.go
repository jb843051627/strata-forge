package app

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func Probe(ctx context.Context, address string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+address+"/healthz", nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 2 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("health probe returned %s", response.Status)
	}
	return nil
}

func AddressURL(address string) string {
	if address == "" {
		return "http://127.0.0.1:8080"
	}
	return "http://" + address
}
