package qwencli

import (
	"context"
	"fmt"
	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
)

type QwenClient struct {
	APIKey string
	Model  string
	Client *dashscopego.TongyiClient
}

func NewClient(apiKey, model string) *QwenClient {
	if model == "" {
		model = qwen.QwenTurbo
	}
	cli := &QwenClient{APIKey: apiKey, Model: model}
	cli.Client = dashscopego.NewTongyiClient(model, apiKey)
	return cli
}

func (c *QwenClient) Chat(msg string, streamCallbackFn qwen.StreamingFunc) string {
	content := qwen.TextContent{Text: msg}

	input := dashscopego.TextInput{
		Messages: []dashscopego.TextMessage{
			{Role: qwen.RoleUser, Content: &content},
		},
	}

	req := &dashscopego.TextRequest{
		Input:       input,            // 请求内容
		StreamingFn: streamCallbackFn, // 流式输出的回调函数, 默认为 nil, 表示不使用流式输出.
	}

	// 发送请求.
	ctx := context.TODO()
	resp, err := c.Client.CreateCompletion(ctx, req)
	if err != nil {
		panic(err)
	}

	/*
		获取结果.
		详细字段说明请查阅 'HTTP调用接口 -> 出参描述'.
		如果request中没有定义流式输出的回调函数 StreamingFn, 则使用此方法获取应答内容.
	*/
	fmt.Println(resp.Output.Choices[0].Message.Content.ToString())

	// 获取 RequestcID, Token 消耗， 结束标识等信息
	fmt.Println(resp.RequestID)
	fmt.Println(resp.Output.Choices[0].FinishReason)
	fmt.Println(resp.Usage.TotalTokens)
	return resp.Output.Choices[0].Message.Content.ToString()
}
