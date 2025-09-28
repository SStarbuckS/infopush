package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TelegramTextConfig Telegram Bot文本消息配置
type TelegramTextConfig struct {
	Token      string
	ChatID     string
	APIBaseURL string
}

// telegramTextRequest 发送文本消息的请求结构
type telegramTextRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// telegramTextHttpRequest 发送HTTP请求
func telegramTextHttpRequest(url string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// SendTelegramText 发送Telegram文本消息 - 统一接口
func SendTelegramText(configName string, configData map[string]interface{}, params map[string]string) (string, error) {
	// 转换配置
	config, err := convertToTelegramTextConfig(configData)
	if err != nil {
		return "", err
	}

	// 获取消息内容
	message := params["msg"]
	if message == "" {
		return "", fmt.Errorf("缺少消息内容")
	}

	// 构造API URL
	url := fmt.Sprintf("%s/bot%s/sendMessage", config.APIBaseURL, config.Token)

	// 构造请求数据
	requestData := telegramTextRequest{
		ChatID: config.ChatID,
		Text:   message,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// 发送请求
	response, err := telegramTextHttpRequest(url, jsonData)
	if err != nil {
		return "", err
	}

	responseStr := string(response)
	fmt.Printf("[%s] %s - Telegram文本返回响应: %s\n", time.Now().Format("2006-01-02 15:04:05.000"), configName, responseStr)

	// 直接检查响应字符串
	if strings.Contains(responseStr, `"ok":true`) {
		return "Success", nil
	} else {
		return "Error: " + responseStr, nil
	}
}

// convertToTelegramTextConfig 将通用配置转换为Telegram文本配置
func convertToTelegramTextConfig(config map[string]interface{}) (TelegramTextConfig, error) {
	// 使用类型断言提取配置值
	token, _ := config["Token"].(string)
	chatID, _ := config["ChatID"].(string)
	apiBaseURL, _ := config["APIBaseURL"].(string)

	if token == "" || chatID == "" || apiBaseURL == "" {
		return TelegramTextConfig{}, fmt.Errorf("缺少必要的Telegram配置参数")
	}

	return TelegramTextConfig{
		Token:      token,
		ChatID:     chatID,
		APIBaseURL: apiBaseURL,
	}, nil
}
