package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const (
	pingInterval = 30 * time.Second // 每30秒ping一次
	logFile      = "ping.log"
	configFile   = "config.json"
	webPort      = 8080
)

// 配置结构体
type Config struct {
	Interval int    `json:"interval"`
	Gateway  string `json:"gateway"`
}

// 状态结构体
type Status struct {
	Connected     bool   `json:"connected"`
	Gateway       string `json:"gateway"`
	ResponseTime  int64  `json:"responseTime"`
	Error         string `json:"error"`
}

// 全局变量
var (
	currentConfig Config
	ticker        *time.Ticker
)

func main() {
	// 加载配置
	loadConfig()

	// 打开日志文件
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer f.Close()

	// 设置日志输出到文件
	log.SetOutput(f)
	log.Println("启动网关ping服务...")

	// 定义网关地址
	gateway := currentConfig.Gateway
	if gateway == "" {
		gateway = getDefaultGateway()
		if gateway == "" {
			log.Println("无法获取默认网关，使用192.168.1.1作为默认值")
			gateway = "192.168.1.1"
			currentConfig.Gateway = gateway
			saveConfig()
		}
	}
	log.Printf("默认网关: %s\n", gateway)

	// 启动定时ping任务
	startTicker()

	// 立即执行一次ping
	pingGateway(gateway)

	// 启动HTTP服务器
	go startHTTPServer()

	// 处理信号，优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("收到退出信号，正在停止服务...")
	if ticker != nil {
		ticker.Stop()
	}
}

// 加载配置
func loadConfig() {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		// 配置文件不存在，使用默认值
		currentConfig = Config{
			Interval: 30,
			Gateway:  "",
		}
		saveConfig()
		return
	}

	err = json.Unmarshal(data, &currentConfig)
	if err != nil {
		// 配置文件损坏，使用默认值
		currentConfig = Config{
			Interval: 30,
			Gateway:  "",
		}
		saveConfig()
	}

	// 确保值有效
	if currentConfig.Interval <= 0 {
		currentConfig.Interval = 30
	}
}

// 保存配置
func saveConfig() {
	data, err := json.MarshalIndent(currentConfig, "", "  ")
	if err != nil {
		log.Printf("保存配置失败: %v\n", err)
		return
	}

	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		log.Printf("保存配置文件失败: %v\n", err)
	}
}

// 启动定时器
func startTicker() {
	if ticker != nil {
		ticker.Stop()
	}
	ticker = time.NewTicker(time.Duration(currentConfig.Interval) * time.Second)

	go func() {
		for range ticker.C {
			pingGateway(currentConfig.Gateway)
		}
	}()
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

// 启动HTTP服务器
func startHTTPServer() {
	// 静态文件服务
	http.Handle("/", http.FileServer(http.Dir("../www")))

	// API路由
	http.HandleFunc("/api/settings", handleSettings)
	http.HandleFunc("/api/status", handleStatus)
	http.HandleFunc("/api/logs", handleLogs)

	addr := fmt.Sprintf(":%d", webPort)
	log.Printf("HTTP服务器启动在 http://localhost%s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("HTTP服务器启动失败: %v\n", err)
	}
}

// 处理设置请求
func handleSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		// 获取设置
		json.NewEncoder(w).Encode(currentConfig)

	case "POST":
		// 保存设置
		var newConfig Config
		err := json.NewDecoder(r.Body).Decode(&newConfig)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "无效的请求数据",
			})
			return
		}

		// 验证设置
		if newConfig.Interval <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "间隔时间必须大于0",
			})
			return
		}

		if newConfig.Gateway == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "网关地址不能为空",
			})
			return
		}

		// 更新配置
		currentConfig = newConfig
		saveConfig()
		startTicker() // 重启定时器

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// 处理状态请求
func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 执行一次ping来获取当前状态
	start := time.Now()
	_, err := net.DialTimeout("tcp", currentConfig.Gateway+":80", 5*time.Second)
	elapsed := time.Since(start)

	status := Status{
		Gateway:      currentConfig.Gateway,
		ResponseTime: elapsed.Milliseconds(),
	}

	if err != nil {
		status.Connected = false
		status.Error = err.Error()
	} else {
		status.Connected = true
	}

	json.NewEncoder(w).Encode(status)
}

// 处理日志请求
func handleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 读取日志文件
	data, err := ioutil.ReadFile(logFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("无法读取日志文件"))
		return
	}

	w.Write(data)
}