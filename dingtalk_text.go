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

// dingTalkTextHttpRequest 发送HTTP请求
func dingTalkTextHttpRequest(url string, data []byte) ([]byte, error) {
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
	response, err := dingTalkTextHttpRequest(url, jsonData)
	if err != nil {
		return "", err
	}

	responseStr := string(response)
	fmt.Printf("[%s] %s - 钉钉文本返回响应: %s\n", time.Now().Format("2006-01-02 15:04:05.000"), configName, responseStr)

	// 直接检查响应字符串
	if strings.Contains(responseStr, `"errcode":0`) {
		return "Success", nil
	} else {
		return "Error: " + responseStr, nil
	}
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
