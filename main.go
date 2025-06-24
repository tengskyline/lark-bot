package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tengskyline/lark-bot/ai"
	"github.com/tengskyline/lark-bot/conf"
	"github.com/tengskyline/lark-bot/lark"
	"github.com/tengskyline/lark-bot/plugins/jenkins"
)

func main() {
	configFile := flag.String("c", "conf/config.yaml", "default conf/config.yaml")
	flag.Parse()

	fmt.Printf("启动 Lark Bot (高并发模式)...\n")
	fmt.Printf("配置文件: %s\n", *configFile)

	err := conf.ConfigInit(*configFile)
	if err != nil {
		fmt.Printf("配置初始化失败: %v\n", err)
		return
	}
	fmt.Printf("配置加载成功: %+v\n", conf.GlobalConfig.AppId)

	// 创建高并发处理器
	fmt.Println("启用高并发处理模式")
	handlerConfig := &lark.HandlerConfig{
		WorkerPoolSize:    conf.GlobalConfig.Concurrency.WorkerPoolSize,
		QueueSize:         conf.GlobalConfig.Concurrency.QueueSize,
		MaxConcurrentAI:   conf.GlobalConfig.Concurrency.MaxConcurrentAI,
		MaxConcurrentHTTP: conf.GlobalConfig.Concurrency.MaxConcurrentHTTP,
		Timeout:           time.Duration(conf.GlobalConfig.Concurrency.TimeoutSeconds) * time.Second,
	}
	handler := lark.NewLarkHandler(handlerConfig)

	// 注册插件
	if conf.GlobalConfig.Jenkins.BaseURL != "" {
		jenkinsPlugin := jenkins.NewWithConfig(
			conf.GlobalConfig.Jenkins.BaseURL,
			conf.GlobalConfig.Jenkins.Username,
			conf.GlobalConfig.Jenkins.Token,
		)
		handler.RegisterPlugin(jenkinsPlugin)
		fmt.Printf("已注册插件: %s\n", jenkinsPlugin.Name())
	} else {
		fmt.Println("Jenkins配置未设置，跳过Jenkins插件注册")
	}

	// 注册AI客户端
	setupAIClients(handler)

	fmt.Printf("工作协程池大小: %d\n", handlerConfig.WorkerPoolSize)
	fmt.Printf("消息队列大小: %d\n", handlerConfig.QueueSize)
	fmt.Printf("AI 并发限制: %d\n", handlerConfig.MaxConcurrentAI)
	fmt.Printf("HTTP 并发限制: %d\n", handlerConfig.MaxConcurrentHTTP)
	fmt.Printf("处理超时时间: %v\n", handlerConfig.Timeout)

	app := lark.NewLark(handler, conf.GlobalConfig)

	// 设置依赖
	handler.Bot = app

	// 启动监控
	go handler.StartMetricsReporter()
	fmt.Println("性能监控已启动")

	// 设置优雅关闭
	setupGracefulShutdown(handler)

	fmt.Println("启动飞书事件监听...")
	err = app.Start()
	if err != nil {
		fmt.Printf("应用启动失败: %v\n", err)
	}
}

// setupGracefulShutdown 设置优雅关闭处理
func setupGracefulShutdown(handler *lark.LarkHandler) {
	// 创建信号通道
	c := make(chan os.Signal, 1)

	// 监听系统信号
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// 启动优雅关闭协程
	go func() {
		sig := <-c
		fmt.Printf("\n收到关闭信号: %v，正在优雅关闭...\n", sig)

		// 执行优雅关闭
		handler.Shutdown()

		fmt.Println("优雅关闭完成")
		fmt.Println("再见！")
		os.Exit(0)
	}()

	fmt.Println("优雅关闭处理器已启动")
}

// 设置AI客户端
func setupAIClients(handler *lark.LarkHandler) {
	factory := ai.NewClientFactory()
	aiConfig := conf.GlobalConfig.AI

	// 验证配置
	if err := factory.ValidateConfig(aiConfig); err != nil {
		fmt.Printf("AI配置验证失败: %v\n", err)
		return
	}

	// 使用工厂创建AI客户端
	client, err := factory.CreateClient(aiConfig)
	if err != nil {
		fmt.Printf("创建AI客户端失败: %v\n", err)
		return
	}

	handler.RegisterAIClient(aiConfig.Provider, client)
	handler.SetDefaultAI(aiConfig.Provider)
	fmt.Printf("已注册AI客户端: %s\n", client.Name())
}
