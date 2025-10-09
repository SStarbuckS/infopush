package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// WecomMPNewsConfig 企业微信图文消息配置
type WecomMPNewsConfig struct {
	APIBaseURL   string
	CorpID       string
	CorpSecret   string
	AgentID      string
	ThumbMediaID string
	Author       string
}

// accessTokenResponse 获取访问令牌的响应结构
type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// article 图文消息文章结构
type article struct {
	Title        string `json:"title"`
	ThumbMediaID string `json:"thumb_media_id"`
	Author       string `json:"author"`
	Content      string `json:"content"`
	Digest       string `json:"digest"`
}

// mpNews 图文消息结构
type mpNews struct {
	Articles []article `json:"articles"`
}

// mpNewsRequest 发送图文消息的请求结构
type mpNewsRequest struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
	AgentID string `json:"agentid"`
	MPNews  mpNews `json:"mpnews"`
}

// getWecomAccessToken 获取企业微信访问令牌
func getWecomAccessToken(config WecomMPNewsConfig) (string, error) {
	url := fmt.Sprintf("%s/cgi-bin/gettoken?corpid=%s&corpsecret=%s",
		config.APIBaseURL, config.CorpID, config.CorpSecret)

	response, err := httpRequest("GET", url, nil, 30*time.Second)
	if err != nil {
		return "", err
	}

	var tokenResp accessTokenResponse
	err = json.Unmarshal(response, &tokenResp)
	if err != nil {
		return "", err
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("获取访问令牌失败: %s", tokenResp.ErrMsg)
	}

	return tokenResp.AccessToken, nil
}

// createWecomMPNewsData 构造图文消息数据
func createWecomMPNewsData(config WecomMPNewsConfig, title, content string) mpNewsRequest {
	// 将换行符转换为HTML换行
	htmlContent := strings.ReplaceAll(content, "\n", "<br>")
	htmlContent = strings.ReplaceAll(htmlContent, "\r\n", "<br>")
	htmlContent = strings.ReplaceAll(htmlContent, "\r", "<br>")

	articleData := article{
		Title:        title,
		ThumbMediaID: config.ThumbMediaID,
		Author:       config.Author,
		Content:      htmlContent,
		Digest:       content,
	}

	return mpNewsRequest{
		ToUser:  "@all",
		MsgType: "mpnews",
		AgentID: config.AgentID,
		MPNews: mpNews{
			Articles: []article{articleData},
		},
	}
}

// SendWecomMPNews 发送企业微信图文消息 - 统一接口
func SendWecomMPNews(configName string, configData map[string]interface{}, params map[string]string) (string, error) {
	// 转换配置
	config, err := convertToWecomMPNewsConfig(configData)
	if err != nil {
		return "", err
	}

	// 处理标题参数
	title := params["title"]
	if title == "" {
		if defaultTitle, ok := configData["DefaultTitle"].(string); ok {
			title = defaultTitle
		} else {
			title = "新提醒" // 默认标题
		}
	}

	// 获取消息内容
	content := params["msg"]

	// 获取访问令牌
	accessToken, err := getWecomAccessToken(config)
	if err != nil {
		return "", fmt.Errorf("获取访问令牌失败: %v", err)
	}

	// 构造消息数据
	msgData := createWecomMPNewsData(config, title, content)

	// 发送消息
	url := fmt.Sprintf("%s/cgi-bin/message/send?access_token=%s", config.APIBaseURL, accessToken)

	jsonData, err := json.Marshal(msgData)
	if err != nil {
		return "", err
	}

	response, err := httpRequest("POST", url, jsonData, 30*time.Second)
	if err != nil {
		return "", err
	}

	responseStr := string(response)
	return handleAPIResponse(configName, "企业微信图文", responseStr, `"errcode":0`)
}

// convertToWecomMPNewsConfig 将通用配置转换为企业微信配置
func convertToWecomMPNewsConfig(config map[string]interface{}) (WecomMPNewsConfig, error) {
	// 使用类型断言提取配置值
	apiBaseURL, _ := config["APIBaseURL"].(string)
	corpID, _ := config["CorpID"].(string)
	corpSecret, _ := config["CorpSecret"].(string)
	agentID, _ := config["AgentID"].(string)
	thumbMediaID, _ := config["ThumbMediaID"].(string)
	author, _ := config["Author"].(string)

	if apiBaseURL == "" || corpID == "" || corpSecret == "" || agentID == "" {
		return WecomMPNewsConfig{}, fmt.Errorf("缺少必要的企业微信配置参数")
	}

	return WecomMPNewsConfig{
		APIBaseURL:   apiBaseURL,
		CorpID:       corpID,
		CorpSecret:   corpSecret,
		AgentID:      agentID,
		ThumbMediaID: thumbMediaID,
		Author:       author,
	}, nil
}
