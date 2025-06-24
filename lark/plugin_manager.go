package lark

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// 插件接口定义
type JobPlugin interface {
	// 插件名称
	Name() string

	// 检查是否匹配此插件
	Match(text string) bool

	// 解析任务
	ParseJobs(text, messageId string) []Job

	// 执行任务
	ExecuteJob(ctx context.Context, job Job, client *http.Client) error
}

// 任务定义
type Job struct {
	ID       string            // 任务ID
	Name     string            // 任务名称
	URL      string            // 任务URL
	Method   string            // HTTP方法
	Headers  map[string]string // HTTP头
	Body     string            // 请求体
	Auth     *JobAuth          // 认证信息
	Metadata map[string]string // 元数据
}

// 认证信息
type JobAuth struct {
	Type     string // basic, token, etc.
	Username string
	Password string
	Token    string
}

// 插件管理器
type PluginManager struct {
	plugins []JobPlugin
	mutex   sync.RWMutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make([]JobPlugin, 0),
	}
}

// 注册插件
func (pm *PluginManager) RegisterPlugin(plugin JobPlugin) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.plugins = append(pm.plugins, plugin)
	fmt.Printf("📦 注册插件: %s\n", plugin.Name())
}

// 查找匹配的插件
func (pm *PluginManager) FindPlugin(text string) JobPlugin {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	for _, plugin := range pm.plugins {
		if plugin.Match(text) {
			return plugin
		}
	}
	return nil
}

// 获取所有插件
func (pm *PluginManager) GetAllPlugins() []JobPlugin {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	result := make([]JobPlugin, len(pm.plugins))
	copy(result, pm.plugins)
	return result
}

// 获取插件数量
func (pm *PluginManager) Count() int {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return len(pm.plugins)
}

// 列出所有插件名称
func (pm *PluginManager) ListPluginNames() []string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	names := make([]string, len(pm.plugins))
	for i, plugin := range pm.plugins {
		names[i] = plugin.Name()
	}
	return names
}
