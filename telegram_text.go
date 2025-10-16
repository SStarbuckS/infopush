package main

import (
	"encoding/json"
	"fmt"
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

// SendTelegramText 发送Telegram文本消息 - 统一接口
func SendTelegramText(configName string, configData map[string]interface{}, params map[string]string) (string, error) {
	// 转换配置
	config, err := convertToTelegramTextConfig(configData)
	if err != nil {
		return "", err
	}

	// 获取消息内容
	message := params["msg"]

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
	response, err := httpRequest("POST", url, jsonData, 30*time.Second)
	if err != nil {
		return "", err
	}

	responseStr := string(response)
	return handleAPIResponse(configName, "Telegram文本", responseStr, `"ok":true`)
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
