package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// WecomRobotTextConfig 企业微信群机器人文本配置
type WecomRobotTextConfig struct {
	APIBaseURL string   `json:"APIBaseURL"`
	Keys       []string `json:"Keys"`
}

// wecomRobotTextRequest 企业微信群机器人文本请求结构
type wecomRobotTextRequest struct {
	MsgType string                    `json:"msgtype"`
	Text    wecomRobotTextMessageBody `json:"text"`
}

// wecomRobotTextMessageBody 文本消息结构
type wecomRobotTextMessageBody struct {
	Content string `json:"content"`
}

// convertToWecomRobotTextConfig 转换配置数据到企业微信群机器人文本配置
func convertToWecomRobotTextConfig(configData map[string]interface{}) (*WecomRobotTextConfig, error) {
	config := &WecomRobotTextConfig{}

	if apiBaseURL, ok := configData["APIBaseURL"].(string); ok {
		config.APIBaseURL = apiBaseURL
	} else {
		return nil, fmt.Errorf("缺少 APIBaseURL 配置")
	}

	if keysInterface, ok := configData["Keys"].([]interface{}); ok {
		config.Keys = make([]string, len(keysInterface))
		for i, keyInterface := range keysInterface {
			if key, ok := keyInterface.(string); ok {
				config.Keys[i] = key
			} else {
				return nil, fmt.Errorf("Keys 配置格式错误")
			}
		}
	} else {
		return nil, fmt.Errorf("缺少 Keys 配置")
	}

	if len(config.Keys) == 0 {
		return nil, fmt.Errorf("Keys 配置不能为空")
	}

	return config, nil
}

// SendWecomRobotText 发送企业微信群机器人文本消息
func SendWecomRobotText(configName string, configData map[string]interface{}, params map[string]string) (string, error) {
	config, err := convertToWecomRobotTextConfig(configData)
	if err != nil {
		return "", err
	}

	message := params["msg"]
	if message == "" {
		return "", fmt.Errorf("缺少消息内容")
	}

	// 随机选择一个机器人key
	rand.Seed(time.Now().UnixNano())
	selectedKey := config.Keys[rand.Intn(len(config.Keys))]

	// 构造完整 Webhook URL
	url := fmt.Sprintf("%s/cgi-bin/webhook/send?key=%s", config.APIBaseURL, selectedKey)

	// 构造请求数据
	requestData := wecomRobotTextRequest{
		MsgType: "text",
		Text: wecomRobotTextMessageBody{
			Content: message,
		},
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
	return handleAPIResponse(configName, "企业微信群机器人文本", responseStr, `"errcode":0`)
}
