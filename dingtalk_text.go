package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// DingTalkTextConfig 钉钉机器人文本消息配置
type DingTalkTextConfig struct {
	AccessToken string
	APIBaseURL  string
}

// dingTalkTextMessage 钉钉文本消息结构
type dingTalkTextMessage struct {
	Content string `json:"content"`
}

// dingTalkTextRequest 发送文本消息的请求结构
type dingTalkTextRequest struct {
	MsgType string              `json:"msgtype"`
	Text    dingTalkTextMessage `json:"text"`
}

// SendDingTalkText 发送钉钉文本消息 - 统一接口
func SendDingTalkText(configName string, configData map[string]interface{}, params map[string]string) (string, error) {
	// 转换配置
	config, err := convertToDingTalkTextConfig(configData)
	if err != nil {
		return "", err
	}

	// 获取消息内容
	message := params["msg"]
	if message == "" {
		return "", fmt.Errorf("缺少消息内容")
	}

	// 构造完整的Webhook URL
	url := fmt.Sprintf("%s?access_token=%s", config.APIBaseURL, config.AccessToken)

	// 构造请求数据
	requestData := dingTalkTextRequest{
		MsgType: "text",
		Text: dingTalkTextMessage{
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
	return handleAPIResponse(configName, "钉钉文本", responseStr, `"errcode":0`)
}

// convertToDingTalkTextConfig 将通用配置转换为钉钉文本配置
func convertToDingTalkTextConfig(config map[string]interface{}) (DingTalkTextConfig, error) {
	// 使用类型断言提取配置值
	accessToken, _ := config["AccessToken"].(string)
	apiBaseURL, _ := config["APIBaseURL"].(string)

	if accessToken == "" || apiBaseURL == "" {
		return DingTalkTextConfig{}, fmt.Errorf("缺少必要的钉钉配置参数")
	}

	return DingTalkTextConfig{
		AccessToken: accessToken,
		APIBaseURL:  apiBaseURL,
	}, nil
}
