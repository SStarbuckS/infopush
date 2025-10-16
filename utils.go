package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// timestamp 返回格式化的时间戳
func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

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
	ts := timestamp()
	fmt.Printf("[%s] %s - %s返回响应: %s\n", ts, configName, platform, responseStr)

	if strings.Contains(responseStr, successPattern) {
		return "Success", nil
	}

	// 返回错误信息（日志记录由 main.go 统一处理）
	return "", fmt.Errorf("%s", responseStr)
}

// logStartupTime 记录程序启动时间到日志文件
func logStartupTime() {
	// 确保 data 目录存在
	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Printf("创建data目录失败: %v\n", err)
		return
	}

	logFile, err := os.OpenFile("data/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("写入启动日志失败: %v\n", err)
		return
	}
	defer logFile.Close()

	ts := timestamp()
	logEntry := fmt.Sprintf("\n========================================\n本次启动时间: %s\n========================================\n", ts)
	if _, err := logFile.WriteString(logEntry); err != nil {
		fmt.Printf("写入启动日志失败: %v\n", err)
	}
}

// writeErrorLog 将错误信息写入日志文件
func writeErrorLog(timestamp, configName, platform, responseStr string, params map[string]string) {
	// 确保 data 目录存在
	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Printf("[%s] 创建data目录失败: %v\n", timestamp, err)
		return
	}

	logFile, err := os.OpenFile("data/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[%s] 写入错误日志失败: %v\n", timestamp, err)
		return
	}
	defer logFile.Close()

	// 构造请求参数字符串
	paramsStr := ""
	if msg, ok := params["msg"]; ok && msg != "" {
		paramsStr += fmt.Sprintf(" msg=%s", msg)
	}
	if title, ok := params["title"]; ok && title != "" {
		paramsStr += fmt.Sprintf(" title=%s", title)
	}

	// 记录完整信息：配置名、请求参数、平台、API响应
	logEntry := fmt.Sprintf("[%s] %s - %s返回响应: %s | 请求参数:%s\n",
		timestamp, configName, platform, responseStr, paramsStr)
	if _, err := logFile.WriteString(logEntry); err != nil {
		fmt.Printf("[%s] 写入错误日志失败: %v\n", timestamp, err)
	}
}
