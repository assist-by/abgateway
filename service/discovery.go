package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	lib "github.com/assist-by/libStruct"
)

// 서비스 주소 가져오는 API
func GetServiceAddress(serviceName, serviceDiscoveryURL string) (string, error) {
	url :=
		fmt.Sprintf("%s/services/%s", serviceDiscoveryURL, serviceName)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error getting service info: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var service lib.Service
	err = json.Unmarshal(body, &service)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling service data: %v", err)
	}

	return service.Address, nil
}
