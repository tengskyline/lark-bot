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

	// 然后检测任务类型关键字，只有同时匹配版本和任务类型才触发
	for taskType, keywords := range jp.config.TaskTypes {
		if taskType == "build" {
			jobs = append(jobs, jp.parseBuildJobs(text, messageId)...)
		} else if taskType == "hgame" {
			jobs = append(jobs, jp.parseHgameBuildJobs(text, messageId)...)
		} else {
			for _, keyword := range keywords {
				for versionKey, versionValue := range jp.config.Versions {
					keyJob := versionKey + keyword
					if strings.Contains(text, keyJob) {
						fmt.Println("解析到 real job", keyJob)
						// 同时匹配到版本和任务类型，才触发相应任务
						switch taskType {
						case "gm":
							jobs = append(jobs, jp.parseGMJobs(versionValue, messageId)...)
						case "package":
							jobs = append(jobs, jp.parsePackageJobs(versionValue, jp.getUpdateType(text), messageId)...)
						}
					}

				}

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

	// 设置其他头部
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
				tag = strings.ToUpper(tagKey) // 确保 tag 始终为大写
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
				Type:     "basic",
				Username: jp.Username,
				Password: jp.Token, // 使用 Password 字段而不是 Token
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

// 解析HGAME任务：构建 op=update，清档 op=clear
func (jp *Plugin) parseHgameBuildJobs(text, messageId string) []lark.Job {
	jobs := make([]lark.Job, 0)

	var tag, op string
	switch {
	case strings.Contains(text, "外网") || strings.Contains(text, "阿里") || strings.Contains(text, "ali"):
		tag = "ali"
	case strings.Contains(text, "集群"):
		tag = "cluster"
	case strings.Contains(text, "时间"):
		tag = "time"
	default:
		return jobs
	}
	switch {
	case strings.Contains(text, "清档"):
		op = "clear"
	case strings.Contains(text, "构建"):
		op = "update"
	default:
		return jobs
	}

	jobURL := fmt.Sprintf("%s/job/build_hgame_server/buildWithParameters?token=hgame&messageId=%s&op=%s&tag=%s",
		jp.BaseURL, messageId, op, tag)

	job := lark.Job{
		ID:     fmt.Sprintf("build_hgame_server_%s_%s", op, tag),
		Name:   fmt.Sprintf("HGAME任务 - %s - %s", op, tag),
		URL:    jobURL,
		Method: "POST",
		Auth: &lark.JobAuth{
			Type:     "basic",
			Username: jp.Username,
			Password: jp.Token,
		},
		Metadata: map[string]string{
			"type": "hgame",
			"op":   op,
		},
	}
	jobs = append(jobs, job)

	return jobs
}

// 解析后台任务
func (jp *Plugin) parseGMJobs(version, messageId string) []lark.Job {
	jobs := make([]lark.Job, 0)

	jobURL := fmt.Sprintf("%s/job/build_online_gm_patch/buildWithParameters?token=hgame&branch=new&messageId=%s&version=%s",
		jp.BaseURL, messageId, version)

	job := lark.Job{
		ID:     fmt.Sprintf("gm_%s", version),
		Name:   fmt.Sprintf("后台任务 - %s", version),
		URL:    jobURL,
		Method: "POST",
		Auth: &lark.JobAuth{
			Type:     "basic",
			Username: jp.Username,
			Password: jp.Token, // 使用 Password 字段而不是 Token
		},
		Metadata: map[string]string{
			"type":    "gm",
			"version": version,
		},
	}
	jobs = append(jobs, job)

	return jobs
}
func (jp *Plugin) getUpdateType(text string) string {
	// 使用配置化的更新类型关键字
	updateType := "all"
	for updateTypeKey, keywords := range jp.config.UpdateTypes {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				updateType = updateTypeKey
				break
			}
		}

	}
	return updateType
}

// 解析打包任务

func (jp *Plugin) parsePackageJobs(version, updateType, messageId string) []lark.Job {
	jobs := make([]lark.Job, 0)

	jobURL := fmt.Sprintf("%s/job/build_online_server_patch/buildWithParameters?token=hgame&branch=new&messageId=%s&version=%s&updateType=%s",
		jp.BaseURL, messageId, version, updateType)

	job := lark.Job{
		ID:     fmt.Sprintf("package_%s_%s", version, updateType),
		Name:   fmt.Sprintf("打包任务 - %s - %s", version, updateType),
		URL:    jobURL,
		Method: "POST",
		Auth: &lark.JobAuth{
			Type:     "basic",
			Username: jp.Username,
			Password: jp.Token, // 使用 Password 字段而不是 Token
		},
		Metadata: map[string]string{
			"type":       "package",
			"version":    version,
			"updateType": updateType,
		},
	}
	jobs = append(jobs, job)
	return jobs
}
