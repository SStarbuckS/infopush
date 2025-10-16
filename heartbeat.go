// 已修改：添加心跳检测功能 (heartbeat)
package main

import (
	"fmt"
	"time"
)

// HeartbeatService 心跳检测服务
type HeartbeatService struct {
	URL      string
	Interval int
}

// Start 启动心跳检测服务
func (h *HeartbeatService) Start() {
	// 如果URL为空，不执行任何逻辑
	if h.URL == "" {
		return
	}

	// 在独立的goroutine中运行心跳检测
	go func() {
		ticker := time.NewTicker(time.Duration(h.Interval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			h.sendHeartbeat()
		}
	}()
}

// sendHeartbeat 发送心跳请求
func (h *HeartbeatService) sendHeartbeat() {
	response, err := httpRequest("GET", h.URL, nil, 30*time.Second)
	if err != nil {
		fmt.Printf("[%s] 心跳检测失败: %v\n", timestamp(), err)
		return
	}

	fmt.Printf("[%s] 心跳检测响应: %s\n", timestamp(), string(response))
}
