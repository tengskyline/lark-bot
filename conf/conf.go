package conf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// AI 客户端配置
type AIConfig struct {
	Provider string            `yaml:"Provider" mapstructure:"Provider"` // qwen, openai, claude, etc.
	APIKey   string            `yaml:"APIKey" mapstructure:"APIKey"`
	BaseURL  string            `yaml:"BaseURL" mapstructure:"BaseURL"`
	Model    string            `yaml:"Model" mapstructure:"Model"`
	Options  map[string]string `yaml:"Options" mapstructure:"Options"` // 额外配置选项
}

type ConcurrencyConfig struct {
	WorkerPoolSize    int `mapstructure:"WorkerPoolSize"`
	QueueSize         int `mapstructure:"QueueSize"`
	MaxConcurrentAI   int `mapstructure:"MaxConcurrentAI"`
	MaxConcurrentHTTP int `mapstructure:"MaxConcurrentHTTP"`
	TimeoutSeconds    int `mapstructure:"TimeoutSeconds"`
}

// Jenkins配置
type JenkinsConfig struct {
	BaseURL  string          `mapstructure:"BaseURL"`
	Username string          `mapstructure:"Username"`
	Token    string          `mapstructure:"Token"`
	Keywords JenkinsKeywords `mapstructure:"Keywords"`
}

// Jenkins关键字配置
type JenkinsKeywords struct {
	// 触发关键字
	TriggerKeywords []string `mapstructure:"TriggerKeywords"`

	// 任务类型关键字
	TaskTypes map[string][]string `mapstructure:"TaskTypes"`

	// 分支关键字
	Branches map[string][]string `mapstructure:"Branches"`

	// 标签关键字
	Tags map[string][]string `mapstructure:"Tags"`

	// 版本关键字映射
	Versions map[string]string `mapstructure:"Versions"`

	// 更新类型关键字
	UpdateTypes map[string][]string `mapstructure:"UpdateTypes"`
}

type LarkBotConfig struct {
	AppId             string            `mapstructure:"AppId"`
	AppSecret         string            `mapstructure:"AppSecret"`
	LogLevel          int               `mapstructure:"LogLevel"`
	VerificationToken string            `mapstructure:"VerificationToken"`
	EncryptKey        string            `mapstructure:"EncryptKey"`
	Concurrency       ConcurrencyConfig `mapstructure:"Concurrency"`
	AI                AIConfig          `mapstructure:"AI"`
	Jenkins           JenkinsConfig     `mapstructure:"Jenkins"`
}

// 全局配置对象
var GlobalConfig *LarkBotConfig = &LarkBotConfig{}

// 初始化配置
func ConfigInit(configPath string) error {
	// 设置配置文件路径和名称
	if configPath != "" {
		// 如果指定了完整路径，解析路径和文件名
		lastSlash := strings.LastIndex(configPath, "/")
		if lastSlash != -1 {
			viper.AddConfigPath(configPath[:lastSlash])
			fileName := configPath[lastSlash+1:]
			// 去掉扩展名
			if dotIndex := strings.LastIndex(fileName, "."); dotIndex != -1 {
				viper.SetConfigName(fileName[:dotIndex])
			} else {
				viper.SetConfigName(fileName)
			}
		} else {
			viper.SetConfigName(strings.TrimSuffix(configPath, ".yaml"))
			viper.AddConfigPath(".")
		}
	} else {
		// 默认配置
		viper.SetConfigName("config")
		viper.AddConfigPath("./conf")
		viper.AddConfigPath(".")
	}

	// 设置配置文件类型
	viper.SetConfigType("yaml")

	// 设置环境变量前缀
	viper.SetEnvPrefix("LARK_BOT")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("[config] 读取配置文件失败: %v", err)
	}

	// 将配置解析到结构体
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("[config] 配置解析失败: %v", err)
	}

	fmt.Printf("配置加载成功：AppId = %s\n", GlobalConfig.AppId)
	fmt.Printf("使用配置文件：%s\n", viper.ConfigFileUsed())
	return nil
}

// 设置默认配置值
func setDefaults() {
	// 基础配置默认值
	viper.SetDefault("LogLevel", 2)

	// 并发配置默认值
	viper.SetDefault("Concurrency.WorkerPoolSize", 10)
	viper.SetDefault("Concurrency.QueueSize", 1000)
	viper.SetDefault("Concurrency.MaxConcurrentAI", 5)
	viper.SetDefault("Concurrency.MaxConcurrentHTTP", 20)
	viper.SetDefault("Concurrency.TimeoutSeconds", 30)

	// AI配置默认值
	viper.SetDefault("AI.Provider", "qwen")
	viper.SetDefault("AI.Model", "qwen-turbo")

	// Jenkins配置默认值
	viper.SetDefault("Jenkins.Keywords.TriggerKeywords", []string{"打包", "后台", "构建"})
	viper.SetDefault("Jenkins.Keywords.TaskTypes", map[string][]string{
		"build":   {"构建"},
		"gm":      {"后台"},
		"package": {"打包"},
	})
	viper.SetDefault("Jenkins.Keywords.Branches", map[string][]string{
		"trunk": {"主干", "trunk"},
		"new":   {"分支", "new"},
	})
	viper.SetDefault("Jenkins.Keywords.Tags", map[string][]string{
		"Y": {"WZY"},
		"H": {"AF"},
		"W": {"PHZ"},
	})
	viper.SetDefault("Jenkins.Keywords.Versions", map[string]string{
		"新马":    "xinma",
		"xinma": "xinma",
		"东南亚":   "xinma",
		"越南":    "vn",
		"vn":    "vn",
		"台湾":    "tw",
		"tw":    "tw",
		"港台":    "tw",
		"韩国":    "kr",
		"kr":    "kr",
	})
	viper.SetDefault("Jenkins.Keywords.UpdateTypes", map[string][]string{
		"all":  {"全部", "all"},
		"code": {"代码", "code"},
		"xml":  {"配置", "xml"},
	})
}

// 获取配置值的便捷方法
func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

// 重新加载配置
func ReloadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("[config] 重新加载配置失败: %v", err)
	}

	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("[config] 配置重新解析失败: %v", err)
	}

	fmt.Println("配置重新加载成功")
	return nil
}
