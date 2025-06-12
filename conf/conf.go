package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type LarkBotConfig struct {
	// 应用APP ID
	AppId string
	// 应用APP Secret
	AppSecret string
	// 日志级别
	// LogLevelDebug LogLevel = 1 ;
	// LogLevelInfo  LogLevel = 2 ;
	// LogLevelWarn  LogLevel = 3 ;
	// LogLevelError LogLevel = 4 ;;
	LogLevel int
	// 用于加密事件或回调的请求内容，校验请求来源
	VerificationToken string
	// 用于加密事件或回调的请求内容，用于解密请求内容
	EncryptKey string
}

var GlobalConfig *LarkBotConfig

func ConfigInit(configPath string) error {
	content, _ := ioutil.ReadFile(configPath)
	err := yaml.Unmarshal(content, GlobalConfig)
	if err != nil {
		log.Println("[config] [yaml]", err.Error())
	}
	fmt.Printf("%+v\n", GlobalConfig)
	return nil
}
