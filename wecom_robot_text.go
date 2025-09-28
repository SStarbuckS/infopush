package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
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

// wecomRobotTextHttpRequest 发送HTTP请求
func wecomRobotTextHttpRequest(url string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
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
	response, err := wecomRobotTextHttpRequest(url, jsonData)
	if err != nil {
		return "", err
	}

	responseStr := string(response)
	fmt.Printf("[%s] %s - 企业微信群机器人文本返回响应: %s\n", time.Now().Format("2006-01-02 15:04:05.000"), configName, responseStr)

	// 检查响应状态
	if strings.Contains(responseStr, `"errcode":0`) {
		return "Success", nil
	} else {
		return "Error: " + responseStr, nil
	}
}
