package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// httpRequest 通用HTTP请求函数
func httpRequest(method, url string, data []byte, timeout time.Duration) ([]byte, error) {
	var req *http.Request
	var err error

	if data != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(data))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// handleAPIResponse 通用API响应处理函数
func handleAPIResponse(configName, platform, responseStr, successPattern string) (string, error) {
	fmt.Printf("[%s] %s - %s返回响应: %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		configName, platform, responseStr)

	if strings.Contains(responseStr, successPattern) {
		return "Success", nil
	}
	return "Error: " + responseStr, nil
}
