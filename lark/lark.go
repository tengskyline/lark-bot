package lark

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	dispatcher "github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkcardkit "github.com/larksuite/oapi-sdk-go/v3/service/cardkit/v1"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/tengskyline/lark-bot/conf"
)

type ILarkHandler interface {
	// OnP2MessageReceiveV1 收到用户消息
	OnP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error
	OnP2MessageReadV1(ctx context.Context, event *larkim.P2MessageReadV1) error
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
	eventDispatcher = eventDispatcher.OnP2MessageReadV1(l.larkHandler.OnP2MessageReadV1)
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
	} else if msgType == larkim.MsgTypeInteractive {
		content, _ := json.Marshal(map[string]any{
			"type": "card",
			"data": map[string]string{
				"card_id": msg,
			},
		})
		body = bodyBuild.Content(string(content)).Build()
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
	} else if msgType == larkim.MsgTypeInteractive {
		content, _ := json.Marshal(map[string]any{
			"type": "card",
			"data": map[string]string{
				"card_id": msg,
			},
		})
		body = bodybuild.Content(string(content)).Build()
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

var cardDataTemplate = `{"schema":"2.0","header":{"title":{"content":"%s","tag":"plain_text"}},"config":{"streaming_mode":true,"summary":{"content":""}},"body":{"elements":[{"tag":"markdown","content":"%s","element_id":"markdown_1"}]}}`

func (c *LarkBot) CreateNewCard(ctx context.Context) string {
	// create card
	cardData := fmt.Sprintf(cardDataTemplate, "思考中...", "稍等，让我想一想...")
	req := larkcardkit.NewCreateCardReqBuilder().
		Body(larkcardkit.NewCreateCardReqBodyBuilder().
			Type(`card_json`).
			Data(cardData).
			Build()).
		Build()
	resp, err := c.client.Cardkit.V1.Card.Create(ctx, req)
	if err != nil {
		fmt.Printf("failed to create card:%s\n", err.Error())
		return ""
	}
	if !resp.Success() {
		fmt.Printf("failed to create card:%s\n", resp.CodeError.String())
		return ""
	}
	return *resp.Data.CardId
}
func (c *LarkBot) UpdateCardChat(ctx context.Context, CardId string, answerCh []string) {
	answer := ""
	seq := 1
	for _, chunk := range answerCh {
		seq += 1
		answer += chunk
		// update card content streaming
		updateReq := larkcardkit.NewContentCardElementReqBuilder().
			CardId(CardId).
			ElementId(`markdown_1`).
			Body(larkcardkit.NewContentCardElementReqBodyBuilder().
				Uuid(uuid.New().String()).
				Content(answer).
				Sequence(seq).
				Build()).
			Build()
		updateResp, err := c.client.Cardkit.V1.CardElement.Content(ctx, updateReq)
		if err != nil {
			fmt.Printf("failed to  update card:%s\n", err.Error())
			return
		}
		if !updateResp.Success() {
			fmt.Printf("failed to  update card:%s\n", err.Error())
			return
		}
	}
}

func (c *LarkBot) SendQACard(ctx context.Context, receiveIdType string, receiveId string, answerCh []string) {
	// create card
	cardData := fmt.Sprintf(cardDataTemplate, "思考中...", "稍等，让我想一想...")
	req := larkcardkit.NewCreateCardReqBuilder().
		Body(larkcardkit.NewCreateCardReqBodyBuilder().
			Type(`card_json`).
			Data(cardData).
			Build()).
		Build()
	resp, err := c.client.Cardkit.V1.Card.Create(ctx, req)
	if err != nil {
		fmt.Printf("failed to create card:%s\n", err.Error())
		return
	}
	if !resp.Success() {
		fmt.Printf("failed to create card:%s\n", resp.CodeError.String())
		return
	}
	content, err := json.Marshal(map[string]any{
		"type": "card",
		"data": map[string]string{
			"card_id": *resp.Data.CardId,
		},
	})
	if err != nil {
		fmt.Printf("failed to create card:%s\n", err.Error())
		return
	}
	// send card to user or group
	res, err := c.client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType("interactive").
			ReceiveId(receiveId).
			Content(string(content)).
			Build()).
		Build())
	if err != nil {
		fmt.Printf("failed to create message:%s\n", err.Error())
		return
	}
	if !res.Success() {
		fmt.Printf("failed to create message:%s\n", err.Error())
		return
	}

	answer := ""
	seq := 1
	for _, chunk := range answerCh {
		seq += 1
		answer += chunk
		// update card content streaming
		updateReq := larkcardkit.NewContentCardElementReqBuilder().
			CardId(*resp.Data.CardId).
			ElementId(`markdown_1`).
			Body(larkcardkit.NewContentCardElementReqBodyBuilder().
				Uuid(uuid.New().String()).
				Content(answer).
				Sequence(seq).
				Build()).
			Build()
		updateResp, err := c.client.Cardkit.V1.CardElement.Content(ctx, updateReq)
		if err != nil {
			fmt.Printf("failed to  update card:%s\n", err.Error())
			return
		}
		if !updateResp.Success() {
			fmt.Printf("failed to  update card:%s\n", err.Error())
			return
		}
	}
	fmt.Printf("start processing QA %s/n", *res.Data.MessageId)
}

type Message struct {
	Text string `json:"text"`
}
