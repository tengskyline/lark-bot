package jenkins

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var developServer = make(map[string]string)
var OnlineServer = make(map[string]string)

func init() {
	{
		developServer["吴志阳"] = "Y"
		developServer["阿房"] = "H"
		developServer["潘浩哲"] = "W"
	}
	{
		OnlineServer["新马"] = "xinma"
		OnlineServer["xinma"] = "xinma"
		OnlineServer["东南亚"] = "xinma"
	}
	{
		OnlineServer["游族"] = "uuzu"
		OnlineServer["国内"] = "uuzu"
		OnlineServer["国服"] = "uuzu"
		OnlineServer["uuzu"] = "uuzu"
	}
	{
		OnlineServer["越南"] = "vn"
		OnlineServer["vn"] = "vn"
		OnlineServer["vn"] = "vn"
	}
	{
		OnlineServer["台湾"] = "tw"
		OnlineServer["tw"] = "tw"
		OnlineServer["港台"] = "tw"
		OnlineServer["TW"] = "tw"
	}
	{
		OnlineServer["韩国"] = "kr"
		OnlineServer["kr"] = "kr"
		OnlineServer["KR"] = "kr"
	}
}

type Jenkins struct {
	Username string
	Password string
}

func (j *Jenkins) JenkinsJob(messageId, reqText string) error {
	//if strings.Contains(reqText, "构建") {
	//	branch := ""
	//	if strings.Contains(reqText, "主干") || strings.Contains(reqText, "trunk") {
	//		branch = "trunk"
	//	} else if strings.Contains(reqText, "分支") || strings.Contains(reqText, "new") {
	//		branch = "new"
	//	} else {
	//		msg = "你是打包主干还是分支啊 分支只能打最新 想好再来"
	//		return msg
	//	}
	//	neiwang := "http://192.168.100.3:8080/job/develop_server_cluster_build/buildWithParameters?token=hgame&messageId=" + messageId + "&"
	//	tag := ""
	//	if strings.Contains(reqText, "吴志阳") {
	//		tag = "Y"
	//	} else if strings.Contains(reqText, "阿房") {
	//		tag = "H"
	//	} else if strings.Contains(reqText, "潘浩哲") {
	//		tag = "W"
	//	}
	//	if len(tag) == 0 {
	//		msg = "没有找到需要构建的对象 请检查后再试"
	//	} else {
	//		jenkins := neiwang + "branch=" + branch + "&tag=" + tag
	//		jenkinss = append(jenkinss, jenkins)
	//	}
	//}
	return nil
}
func (j *Jenkins) JenkinsBuild(jenkinsUrl string) error {
	fmt.Printf(jenkinsUrl + "\n")
	// 创建请求
	req, rerr := http.NewRequest("POST", jenkinsUrl, nil)
	if rerr != nil {
		return fmt.Errorf("执行任务失败 请检查重试")
	}
	// 设置 Basic Auth 认证
	req.SetBasicAuth(j.Username, j.Password)

	// 发送请求
	rclient := &http.Client{}
	resp, terr := rclient.Do(req)
	if terr != nil {
		return fmt.Errorf("执行任务失败 请检查重试")
	}
	defer resp.Body.Close()
	// 读取响应
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
	return nil
}
