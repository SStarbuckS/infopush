// 已修改：添加心跳检测功能 (heartbeat)
package main

import (
	"fmt"
	"net/http"
	"strings"
)

// 全局配置管理器
var configManager *ConfigManager

// dynamicHandler 动态路由处理器
func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// 从URL路径中提取配置名称
	fullPath := strings.Trim(r.URL.Path, "/")

	// 处理全局路由前缀并提取配置路径
	var configPath string
	if configManager.Route != "" && configManager.Route != "/" {
		// 移除路由前缀
		routePrefix := strings.Trim(configManager.Route, "/") + "/"
		if strings.HasPrefix(fullPath+"/", routePrefix) {
			configPath = strings.TrimPrefix(fullPath, strings.TrimSuffix(routePrefix, "/"))
			configPath = strings.Trim(configPath, "/")
		}
	} else {
		configPath = fullPath
	}

	// 统一检查配置路径
	if configPath == "" {
		ts := timestamp()
		errorMsg := "目的地空无一物 - 缺少配置路径"

		// 写入错误日志
		writeErrorLog(ts, "", "unknown", errorMsg, nil)

		http.Error(w, "目的地空无一物", http.StatusBadRequest)
		fmt.Printf("[%s] %s\n", ts, errorMsg)
		return
	}

	// 获取配置
	config, exists := configManager.GetConfig(configPath)
	if !exists {
		ts := timestamp()
		errorMsg := fmt.Sprintf("这里是一片荒原 - 配置 '%s' 不存在", configPath)

		// 写入错误日志
		writeErrorLog(ts, configPath, "unknown", errorMsg, nil)

		http.Error(w, "这里是一片荒原", http.StatusNotFound)
		fmt.Printf("[%s] %s\n", ts, errorMsg)
		return
	}

	// 获取消息内容 - 缺少msg参数
	msg := r.FormValue("msg")
	if msg == "" {
		ts := timestamp()
		errorMsg := "Wel Come! - 缺少msg参数"

		// 写入错误日志
		writeErrorLog(ts, configPath, config.Type, errorMsg, nil)

		http.Error(w, "Wel Come!", http.StatusBadRequest)
		fmt.Printf("[%s] %s\n", ts, errorMsg)
		return
	}

	// 获取所有表单参数
	params := make(map[string]string)
	params["msg"] = msg
	params["title"] = r.FormValue("title")

	// 根据配置类型处理消息
	var result string
	var err error

	switch config.Type {
	case "dingtalk_text":
		result, err = SendDingTalkText(configPath, config.Config, params)
	case "telegram_text":
		result, err = SendTelegramText(configPath, config.Config, params)
	case "wecom_mpnews":
		result, err = SendWecomMPNews(configPath, config.Config, params)
	case "wecom_robot_text":
		result, err = SendWecomRobotText(configPath, config.Config, params)
	default:
		ts := timestamp()
		errorMsg := fmt.Sprintf("不支持的推送类型: %s", config.Type)

		// 写入错误日志
		writeErrorLog(ts, configPath, config.Type, errorMsg, params)

		http.Error(w, errorMsg, http.StatusBadRequest)
		fmt.Printf("[%s] %s - %s\n", ts, configPath, errorMsg)
		return
	}

	if err != nil {
		ts := timestamp()
		errorMsg := fmt.Sprintf("Error: %v", err)

		// 写入错误日志
		writeErrorLog(ts, configPath, config.Type, errorMsg, params)

		http.Error(w, errorMsg, http.StatusInternalServerError)
		fmt.Printf("[%s] %s - %s\n", ts, configPath, errorMsg)
		return
	}

	// 返回响应
	fmt.Printf("[%s] %s - %s\n", timestamp(), configPath, result)
	fmt.Fprint(w, result)
}

func main() {
	// 记录启动时间到日志文件
	logStartupTime()

	// 加载配置文件
	var err error
	configManager, err = NewConfigManager("data/config.json")
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}

	// 检查全局路由配置
	if configManager.Route == "" {
		fmt.Printf("错误: 配置文件中缺少 'route' 字段或值为空\n")
		fmt.Printf("请在 config.json 中设置全局路由，例如:\n")
		fmt.Printf("  \"route\": \"/\"          # 无前缀\n")
		fmt.Printf("  \"route\": \"/push\"      # 有前缀\n")
		return
	}

	// 启动心跳检测服务
	heartbeat := &HeartbeatService{
		URL:      configManager.HeartbeatURL,
		Interval: configManager.HeartbeatInterval,
	}
	heartbeat.Start()

	// 显示心跳检测状态信息
	if configManager.HeartbeatURL != "" {
		fmt.Printf("心跳检测已启动: %s (间隔: %d秒)\n", configManager.HeartbeatURL, configManager.HeartbeatInterval)
	} else {
		fmt.Println("心跳检测未配置")
	}

	// 注册动态路由
	http.HandleFunc("/", dynamicHandler)

	// 启动服务器
	fmt.Println("多配置消息推送服务启动中...")
	fmt.Println("服务地址: http://localhost:8080")

	// 显示路由前缀信息
	if configManager.Route != "" && configManager.Route != "/" {
		fmt.Printf("全局路由前缀: %s\n", configManager.Route)
	}

	fmt.Println("支持的配置路由:")
	for _, name := range configManager.GetAllConfigNames() {
		config, _ := configManager.GetConfig(name)
		routePath := name
		if configManager.Route != "" && configManager.Route != "/" {
			routePath = strings.Trim(configManager.Route, "/") + "/" + name
		}
		fmt.Printf("  - http://localhost:8080/%s/ (类型: %s)\n", routePath, config.Type)
	}
	fmt.Println("使用方法: POST/GET 请求，参数 msg=消息内容 [title=标题]")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}
