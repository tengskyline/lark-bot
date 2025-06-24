package conf

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// VipeExampleUsage 展示 Viper 的各种用法
func ViperExampleUsage() {
	fmt.Println("=== Viper 功能演示 ===")

	// 1. 直接读取配置值
	fmt.Printf("AppId: %s\n", viper.GetString("AppId"))
	fmt.Printf("LogLevel: %d\n", viper.GetInt("LogLevel"))
	fmt.Printf("WorkerPoolSize: %d\n", viper.GetInt("Concurrency.WorkerPoolSize"))

	// 2. 使用默认值
	serverPort := viper.GetInt("Server.Port")
	if serverPort == 0 {
		serverPort = 8080 // 如果没有配置，使用默认值
	}
	fmt.Printf("Server Port: %d\n", serverPort)

	// 3. 检查配置项是否存在
	if viper.IsSet("Jenkins.BaseURL") {
		fmt.Printf("Jenkins URL: %s\n", viper.GetString("Jenkins.BaseURL"))
	} else {
		fmt.Println("Jenkins 配置未设置")
	}

	// 4. 获取子配置
	concurrencyConfig := viper.Sub("Concurrency")
	if concurrencyConfig != nil {
		fmt.Println("并发配置:")
		fmt.Printf("  WorkerPoolSize: %d\n", concurrencyConfig.GetInt("WorkerPoolSize"))
		fmt.Printf("  QueueSize: %d\n", concurrencyConfig.GetInt("QueueSize"))
	}

	// 5. 获取所有配置键
	fmt.Println("所有配置键:")
	for _, key := range viper.AllKeys() {
		fmt.Printf("  %s\n", key)
	}

	// 6. 动态设置配置值
	viper.Set("Runtime.StartTime", "2024-01-01T00:00:00Z")
	fmt.Printf("动态设置的运行时间: %s\n", viper.GetString("Runtime.StartTime"))
}

// WatchConfigChanges 监听配置文件变化
func WatchConfigChanges() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("🔄 配置文件发生变化: %s\n", e.Name)

		// 重新解析配置到全局变量
		if err := viper.Unmarshal(GlobalConfig); err != nil {
			fmt.Printf("❌ 配置重新解析失败: %v\n", err)
		} else {
			fmt.Println("✅ 配置已自动重新加载")
		}
	})
}

// ValidateConfig 验证必要的配置项
func ValidateConfig() error {
	requiredKeys := []string{
		"AppId",
		"AppSecret",
	}

	for _, key := range requiredKeys {
		if !viper.IsSet(key) || viper.GetString(key) == "" {
			return fmt.Errorf("必需的配置项 %s 未设置", key)
		}
	}

	// 验证数值范围
	if viper.GetInt("Concurrency.WorkerPoolSize") <= 0 {
		return fmt.Errorf("WorkerPoolSize 必须大于 0")
	}

	if viper.GetInt("Concurrency.QueueSize") <= 0 {
		return fmt.Errorf("QueueSize 必须大于 0")
	}

	fmt.Println("✅ 配置验证通过")
	return nil
}

// GetConfigSummary 获取配置摘要（隐藏敏感信息）
func GetConfigSummary() map[string]interface{} {
	summary := make(map[string]interface{})

	// 安全的配置项
	safeKeys := []string{
		"LogLevel",
		"Concurrency.WorkerPoolSize",
		"Concurrency.QueueSize",
		"Concurrency.MaxConcurrentAI",
		"Concurrency.MaxConcurrentHTTP",
		"Concurrency.TimeoutSeconds",
		"Server.Port",
		"Server.Host",
		"Monitoring.MetricsEnabled",
	}

	for _, key := range safeKeys {
		if viper.IsSet(key) {
			summary[key] = viper.Get(key)
		}
	}

	// 敏感信息只显示是否配置
	sensitiveKeys := []string{
		"AppId",
		"AppSecret",
		"QwenKey",
		"VerificationToken",
		"EncryptKey",
	}

	for _, key := range sensitiveKeys {
		if viper.IsSet(key) && viper.GetString(key) != "" {
			summary[key] = "已配置"
		} else {
			summary[key] = "未配置"
		}
	}

	return summary
}
