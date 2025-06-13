package lark

import (
	"context"
	"encoding/json"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/patrickmn/go-cache"
	"github.com/tengskyline/lark-bot/qwencli"
	"time"
)

var eventCache = cache.New(5*time.Minute, 10*time.Minute)

type LarkHandler struct {
	Bot     *LarkBot
	QClinet *qwencli.QwenClient
}

func NewLarkHandler() *LarkHandler {
	return &LarkHandler{}
}

func (e *LarkHandler) EventCheck(eventId string) bool {
	if _, found := eventCache.Get(eventId); found {
		return false
	}
	eventCache.Set(eventId, true, cache.DefaultExpiration)
	return true
}
func (e *LarkHandler) OnP2MessageReadV1(ctx context.Context, event *larkim.P2MessageReadV1) error {
	// 你可以留空，或者写日志
	fmt.Printf("收到 message_read 事件, 消息已读\n")
	return nil
}
func (e *LarkHandler) OnP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	// 处理消息 event，这里简单打印消息的内容
	fmt.Printf("[OnP2MessageReceiveV1 access], data: %s\n", larkcore.Prettify(event))
	/**
	 * 解析用户发送的消息。
	 * Parse the message sent by the user.
	 */
	eventID := event.EventV2Base.Header.EventID
	fmt.Printf("[OnP2MessageReceiveV1 access], data: %v\n", eventID)
	if !e.EventCheck(eventID) {
		return nil
	}
	/**
	 * 检查消息类型是否为文本
	 * Check if the message type is text
	 */
	var respContent map[string]string
	err := json.Unmarshal([]byte(*event.Event.Message.Content), &respContent)
	if err != nil || *event.Event.Message.MessageType != "text" {
		return e.SendMessage(event, larkim.MsgTypeText, "解析消息失败，请发送文本消息\n")
	}
	reqText := respContent["text"]
	chunks := make([]string, 0, 0)
	cardId := e.Bot.CreateNewCard(context.TODO())
	e.SendMessage(event, larkim.MsgTypeInteractive, cardId)
	e.QClinet.Chat(reqText, func(ctx context.Context, chunk []byte) error {
		if len(chunk) == 0 {
			return nil
		}
		chunks = append(chunks, string(chunk))
		return nil
	})
	fmt.Printf("[OnP2MessageReceiveV1 access], reqText: %v\n", reqText)
	fmt.Printf("[OnP2MessageReceiveV1 access], chat: %+v\n", chunks)
	e.Bot.UpdateCardChat(context.TODO(), cardId, chunks)
	return nil
}

func (e *LarkHandler) SendMessage(event *larkim.P2MessageReceiveV1, msgType, msg string) error {
	if *event.Event.Message.ChatType == "p2p" {
		return e.Bot.SendP2PReqMessage(*event.Event.Message.ChatId, msgType, msg)
	} else {
		return e.Bot.SendReplyReqMessage(*event.Event.Message.MessageId, msgType, msg)
	}
}
