package jenkins

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tengskyline/lark-bot/conf"
	"github.com/tengskyline/lark-bot/lark"
)

// Jenkins 插件实现
type Plugin struct {
	BaseURL  string
	Username string
	Token    string
	config   *conf.JenkinsKeywords
}

func New() *Plugin {
	return &Plugin{
		config: &conf.GlobalConfig.Jenkins.Keywords,
	}
}

func NewWithConfig(baseURL, username, token string) *Plugin {
	plugin := &Plugin{
		BaseURL:  baseURL,
		Username: username,
		Token:    token,
		config:   &conf.GlobalConfig.Jenkins.Keywords,
	}
	return plugin
}

func (jp *Plugin) Name() string {
	return "Jenkins构建插件"
}

func (jp *Plugin) Match(text string) bool {
	// 使用配置化的触发关键字
	for _, keyword := range jp.config.TriggerKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func (jp *Plugin) ParseJobs(text, messageId string) []lark.Job {
	jobs := make([]lark.Job, 0)

	// 使用配置化的任务类型关键字
	for taskType, keywords := range jp.config.TaskTypes {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				switch taskType {
				case "build":
					jobs = append(jobs, jp.parseBuildJobs(text, messageId)...)
				case "gm":
					jobs = append(jobs, jp.parseGMJobs(text, messageId)...)
				case "package":
					jobs = append(jobs, jp.parsePackageJobs(text, messageId)...)
				}
				break
			}
		}
	}

	return jobs
}

func (jp *Plugin) ExecuteJob(ctx context.Context, job lark.Job, client *http.Client) error {
	req, err := http.NewRequestWithContext(ctx, job.Method, job.URL, strings.NewReader(job.Body))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置认证
	if job.Auth != nil {
		switch job.Auth.Type {
		case "basic":
			req.SetBasicAuth(job.Auth.Username, job.Auth.Password)
		case "token":
			req.Header.Set("Authorization", "Bearer "+job.Auth.Token)
		}
	}

	// 设置头部
	for key, value := range job.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("执行请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Jenkins任务失败 状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Jenkins任务执行成功: %s\n", job.Name)
	return nil
}

// 解析构建任务
func (jp *Plugin) parseBuildJobs(text, messageId string) []lark.Job {
	jobs := make([]lark.Job, 0)

	// 使用配置化的分支关键字
	branch := ""
	for branchKey, keywords := range jp.config.Branches {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				branch = branchKey
				break
			}
		}
		if branch != "" {
			break
		}
	}

	if branch == "" {
		return jobs
	}

	// 使用配置化的标签关键字
	tag := ""
	for tagKey, keywords := range jp.config.Tags {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				tag = tagKey
				break
			}
		}
		if tag != "" {
			break
		}
	}

	if tag != "" {
		jobURL := fmt.Sprintf("%s/job/develop_server_cluster_build/buildWithParameters?token=hgame&messageId=%s&branch=%s&tag=%s",
			jp.BaseURL, messageId, branch, tag)

		job := lark.Job{
			ID:     fmt.Sprintf("build_%s_%s", branch, tag),
			Name:   fmt.Sprintf("构建任务 - %s分支 - %s", branch, tag),
			URL:    jobURL,
			Method: "POST",
			Auth: &lark.JobAuth{
				Type:     "token",
				Username: jp.Username,
				Token:    jp.Token,
			},
			Metadata: map[string]string{
				"type":   "build",
				"branch": branch,
				"tag":    tag,
			},
		}
		jobs = append(jobs, job)
	}

	return jobs
}

// 解析后台任务
func (jp *Plugin) parseGMJobs(text, messageId string) []lark.Job {
	jobs := make([]lark.Job, 0)

	// 使用配置化的版本关键字映射
	versions := make([]string, 0)
	for key, value := range jp.config.Versions {
		if strings.Contains(text, key) {
			versions = append(versions, value)
		}
	}

	for _, version := range versions {
		jobURL := fmt.Sprintf("%s/job/build_online_gm_patch/buildWithParameters?token=hgame&branch=new&messageId=%s&version=%s",
			jp.BaseURL, messageId, version)

		job := lark.Job{
			ID:     fmt.Sprintf("gm_%s", version),
			Name:   fmt.Sprintf("后台任务 - %s", version),
			URL:    jobURL,
			Method: "POST",
			Auth: &lark.JobAuth{
				Type:     "token",
				Username: jp.Username,
				Token:    jp.Token,
			},
			Metadata: map[string]string{
				"type":    "gm",
				"version": version,
			},
		}
		jobs = append(jobs, job)
	}

	return jobs
}

// 解析打包任务
func (jp *Plugin) parsePackageJobs(text, messageId string) []lark.Job {
	jobs := make([]lark.Job, 0)

	// 使用配置化的版本关键字映射
	versions := make([]string, 0)
	for key, value := range jp.config.Versions {
		if strings.Contains(text, key) {
			versions = append(versions, value)
		}
	}

	// 使用配置化的更新类型关键字
	updateType := "all"
	for updateTypeKey, keywords := range jp.config.UpdateTypes {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				updateType = updateTypeKey
				break
			}
		}
		if updateType != "all" {
			break
		}
	}

	for _, version := range versions {
		jobURL := fmt.Sprintf("%s/job/build_online_server_patch/buildWithParameters?token=hgame&branch=new&messageId=%s&version=%s&updateType=%s",
			jp.BaseURL, messageId, version, updateType)

		job := lark.Job{
			ID:     fmt.Sprintf("package_%s_%s", version, updateType),
			Name:   fmt.Sprintf("打包任务 - %s - %s", version, updateType),
			URL:    jobURL,
			Method: "POST",
			Auth: &lark.JobAuth{
				Type:     "token",
				Username: jp.Username,
				Token:    jp.Token,
			},
			Metadata: map[string]string{
				"type":       "package",
				"version":    version,
				"updateType": updateType,
			},
		}
		jobs = append(jobs, job)
	}

	return jobs
}
