// 已修改：添加心跳检测功能 (heartbeat)
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// PushConfig 推送配置结构
type PushConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// ConfigManager 配置管理器
type ConfigManager struct {
	Route             string `json:"route"`
	HeartbeatURL      string `json:"heartbeat_url"`
	HeartbeatInterval int    `json:"heartbeat_interval"`
	Configs           map[string]PushConfig
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configFile string) (*ConfigManager, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 首先解析为通用map来提取route和其他配置
	var rawConfig map[string]interface{}
	err = json.Unmarshal(data, &rawConfig)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 提取route配置
	route := ""
	if routeValue, ok := rawConfig["route"].(string); ok {
		route = routeValue
	}

	// 提取心跳检测配置
	heartbeatURL := ""
	if urlValue, ok := rawConfig["heartbeat_url"].(string); ok {
		heartbeatURL = urlValue
	}

	heartbeatInterval := 0
	if intervalValue, ok := rawConfig["heartbeat_interval"].(float64); ok {
		heartbeatInterval = int(intervalValue)
	}

	// 提取推送配置（排除route、heartbeat_url、heartbeat_interval字段）
	configs := make(map[string]PushConfig)
	for key, value := range rawConfig {
		if key != "route" && key != "heartbeat_url" && key != "heartbeat_interval" {
			// 将interface{}转换为PushConfig
			valueBytes, err := json.Marshal(value)
			if err != nil {
				continue
			}
			var pushConfig PushConfig
			err = json.Unmarshal(valueBytes, &pushConfig)
			if err != nil {
				continue
			}
			configs[key] = pushConfig
		}
	}

	return &ConfigManager{
		Route:             route,
		HeartbeatURL:      heartbeatURL,
		HeartbeatInterval: heartbeatInterval,
		Configs:           configs,
	}, nil
}

// GetConfig 获取指定名称的配置
func (cm *ConfigManager) GetConfig(name string) (PushConfig, bool) {
	config, exists := cm.Configs[name]
	return config, exists
}

// GetAllConfigNames 获取所有配置名称
func (cm *ConfigManager) GetAllConfigNames() []string {
	names := make([]string, 0, len(cm.Configs))
	for name := range cm.Configs {
		names = append(names, name)
	}
	return names
}
