package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type LarkBotConfig struct {
	AppId             string `yaml:"AppId"`
	AppSecret         string `yaml:"AppSecret"`
	LogLevel          int    `yaml:"LogLevel"`
	VerificationToken string `yaml:"VerificationToken"`
	EncryptKey        string `yaml:"EncryptKey"`
}

// 全局配置对象
var GlobalConfig *LarkBotConfig = &LarkBotConfig{}

// 初始化配置
func ConfigInit(configPath string) error {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("[config] 读取配置文件失败: %v", err)
	}

	err = yaml.Unmarshal(content, GlobalConfig)
	if err != nil {
		return fmt.Errorf("[config] YAML 解析失败: %v", err)
	}

	fmt.Printf("✅ 配置加载成功：AppId = %s\n", GlobalConfig.AppId)
	return nil
}
