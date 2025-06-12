package lark

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	dispatcher "github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/tengskyline/lark-bot/conf"
)

type ILarkHandler interface {
	// OnP2MessageReceiveV1 收到用户消息
	OnP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error
}
type LarkBot struct {
	config          *conf.LarkBotConfig
	eventDispatcher *dispatcher.EventDispatcher
	client          *lark.Client
	larkHandler     ILarkHandler
}

func NewLark(handler ILarkHandler, cf *conf.LarkBotConfig) *LarkBot {
	larkApp := &LarkBot{}
	config := cf
	larkApp.config = config
	larkApp.larkHandler = handler
	return larkApp
}
func (l *LarkBot) Start() error {
	//处理Event事件
	if err := l.EventDispatcher(); err != nil {
		return fmt.Errorf("event dispatcher error: %v", err)
	}
	return nil
}

// EventDispatcher 处理消息回调
func (l *LarkBot) EventDispatcher() error {
	// 注册事件回调
	eventDispatcher := dispatcher.NewEventDispatcher(l.config.VerificationToken, l.config.EncryptKey)
	// 处理消息回传
	eventDispatcher = eventDispatcher.OnP2MessageReceiveV1(l.larkHandler.OnP2MessageReceiveV1)
	// 保存事件回调
	l.eventDispatcher = eventDispatcher

	// 创建Client
	wsclient := larkws.NewClient(l.config.AppId, l.config.AppSecret,
		larkws.WithEventHandler(l.eventDispatcher),
		larkws.WithLogLevel(larkcore.LogLevel(l.config.LogLevel)),
	)
	l.client = lark.NewClient(l.config.AppId, l.config.AppSecret)
	// 启动客户端
	return wsclient.Start(context.Background())
}

// SendTextMessage 发送文本消息
func (c *LarkBot) SendP2PReqMessage(receiveId string, msgType, msg string) error {
	/**
	 * 使用SDK调用发送消息接口。 Use SDK to call send message interface.
	 * https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/create
	 */
	var body *larkim.CreateMessageReqBody
	bodyBuild := larkim.NewCreateMessageReqBodyBuilder().MsgType(msgType).ReceiveId(receiveId)
	if msgType == larkim.MsgTypeText {
		content := larkim.NewTextMsgBuilder().TextLine(msg).Build()
		body = bodyBuild.Content(content).Build()
	} else if msgType == larkim.MsgTypeImage {
		image := `{"image_key":"` + msg + `"}`
		body = bodyBuild.Content(image).Build()
	}
	resp, err := c.client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId). // 消息接收者的 ID 类型，设置为会话ID。 ID type of the message receiver, set to chat ID.
		Body(body).Build())

	if err != nil || !resp.Success() {
		fmt.Println(err)
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil
	}
	return nil
}

func (c *LarkBot) SendReplyReqMessage(MessageId string, msgType, msg string) error {
	/**
	 * 使用SDK调用发送消息接口。 Use SDK to call send message interface.
	 * https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/message/create
	 */
	var body *larkim.ReplyMessageReqBody
	bodybuild := larkim.NewReplyMessageReqBodyBuilder().MsgType(msgType)
	if msgType == larkim.MsgTypeText {
		content := larkim.NewTextMsgBuilder().TextLine(msg).Build()
		body = bodybuild.Content(content).Build()
	} else if msgType == larkim.MsgTypeImage {
		image := `{"image_key":"` + msg + `"}`
		body = bodybuild.Content(image).Build()
	}
	resp, err := c.client.Im.Message.Reply(context.Background(), larkim.NewReplyMessageReqBuilder().
		MessageId(MessageId).
		Body(body).
		Build())
	if err != nil || !resp.Success() {
		fmt.Println(err)
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil
	}
	return nil
}
