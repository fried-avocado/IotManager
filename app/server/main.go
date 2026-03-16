package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	pingInterval = 30 * time.Second // 每30秒ping一次
	logFile      = "ping.log"
)

func main() {
	// 打开日志文件
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer f.Close()

	// 设置日志输出到文件
	log.SetOutput(f)
	log.Println("启动网关ping服务...")

	// 定义网关地址（默认网关，可根据实际情况修改）
	gateway := getDefaultGateway()
	if gateway == "" {
		log.Println("无法获取默认网关，使用192.168.1.1作为默认值")
		gateway = "192.168.1.1"
	}
	log.Printf("默认网关: %s\n", gateway)

	// 启动定时ping任务
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	// 立即执行一次ping
	pingGateway(gateway)

	// 处理信号，优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			pingGateway(gateway)
		case <-sigChan:
			log.Println("收到退出信号，正在停止服务...")
			return
		}
	}
}

// pingGateway 执行ping操作并记录结果
func pingGateway(gateway string) {
	start := time.Now()
	_, err := net.DialTimeout("tcp", gateway+":80", 5*time.Second)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("❌ 网关 %s 不可达: %v (耗时: %v)\n", gateway, err, elapsed)
	} else {
		log.Printf("✅ 网关 %s 可达 (耗时: %v)\n", gateway, elapsed)
	}
}

// getDefaultGateway 获取默认网关
func getDefaultGateway() string {
	// 在实际生产环境中，这里应该调用系统API获取默认网关
	// 这里为了简化，返回空字符串，使用默认值
	return ""
}