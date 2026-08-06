package lark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 处理器配置
type HandlerConfig struct {
	WorkerPoolSize    int           // 工作协程池大小
	QueueSize         int           // 消息队列大小
	MaxConcurrentAI   int           // AI 并发数限制
	MaxConcurrentHTTP int           // HTTP 并发数限制
	Timeout           time.Duration // 处理超时时间
}

// 消息任务
type MessageTask struct {
	Event     *larkim.P2MessageReceiveV1
	Context   context.Context
	CreatedAt time.Time
}

// 性能监控
type HandlerMetrics struct {
	processedTotal int64
	errorTotal     int64
	queuedTotal    int64
	droppedTotal   int64
	processingTime time.Duration
	mutex          sync.RWMutex
}

func NewHandlerMetrics() *HandlerMetrics {
	return &HandlerMetrics{}
}

func (m *HandlerMetrics) RecordSuccess(count int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.processedTotal += count
}

func (m *HandlerMetrics) RecordError(errorType string, count int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.errorTotal += count
}

func (m *HandlerMetrics) RecordQueued(count int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.queuedTotal += count
}

func (m *HandlerMetrics) RecordDropped(count int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.droppedTotal += count
}

func (m *HandlerMetrics) RecordProcessTime(duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.processingTime = duration
}

func (m *HandlerMetrics) GetStats() (processed, errors, queued, dropped int64, avgTime time.Duration) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.processedTotal, m.errorTotal, m.queuedTotal, m.droppedTotal, m.processingTime
}

// HTTP 客户端池
type HTTPClientPool struct {
	pool sync.Pool
}

func NewHTTPClientPool() *HTTPClientPool {
	return &HTTPClientPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &http.Client{
					Timeout: 30 * time.Second,
					Transport: &http.Transport{
						MaxIdleConns:        100,
						MaxIdleConnsPerHost: 10,
						IdleConnTimeout:     90 * time.Second,
						DisableKeepAlives:   false,
					},
				}
			},
		},
	}
}

func (p *HTTPClientPool) Get() *http.Client {
	return p.pool.Get().(*http.Client)
}

func (p *HTTPClientPool) Put(client *http.Client) {
	p.pool.Put(client)
}

// 消息处理器
type LarkHandler struct {
	Bot    *LarkBot
	config *HandlerConfig

	// AI客户端管理器
	aiManager *AIClientManager

	// 工作池
	messageQueue chan *MessageTask
	workerPool   sync.Pool

	// 限流器
	aiSemaphore   chan struct{}
	httpSemaphore chan struct{}

	// HTTP客户端池
	httpPool *HTTPClientPool

	// 插件管理器
	pluginManager *PluginManager

	// 监控
	metrics *HandlerMetrics

	// 优雅关闭
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewLarkHandler(config *HandlerConfig) *LarkHandler {
	if config.WorkerPoolSize == 0 {
		config.WorkerPoolSize = runtime.NumCPU() * 2
	}
	if config.QueueSize == 0 {
		config.QueueSize = 1000
	}
	if config.MaxConcurrentAI == 0 {
		config.MaxConcurrentAI = 10
	}
	if config.MaxConcurrentHTTP == 0 {
		config.MaxConcurrentHTTP = 20
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	handler := &LarkHandler{
		config:        config,
		messageQueue:  make(chan *MessageTask, config.QueueSize),
		aiSemaphore:   make(chan struct{}, config.MaxConcurrentAI),
		httpSemaphore: make(chan struct{}, config.MaxConcurrentHTTP),
		httpPool:      NewHTTPClientPool(),
		pluginManager: NewPluginManager(),
		aiManager:     NewAIClientManager(),
		metrics:       NewHandlerMetrics(),
		ctx:           ctx,
		cancel:        cancel,
	}

	// 初始化工作池
	handler.workerPool = sync.Pool{
		New: func() interface{} {
			return &MessageWorker{
				handler: handler,
			}
		},
	}

	// 启动工作协程
	handler.startWorkers()

	return handler
}

// 消息工作器
type MessageWorker struct {
	handler *LarkHandler
}

func (h *LarkHandler) startWorkers() {
	for i := 0; i < h.config.WorkerPoolSize; i++ {
		h.wg.Add(1)
		go h.worker()
	}
}

func (h *LarkHandler) worker() {
	defer h.wg.Done()

	for {
		select {
		case task := <-h.messageQueue:
			h.processTask(task)
		case <-h.ctx.Done():
			return
		}
	}
}

func (h *LarkHandler) processTask(task *MessageTask) {
	defer func() {
		if r := recover(); r != nil {
			h.metrics.RecordError("panic", 1)
			fmt.Printf("Worker panic recovered: %v\n", r)
		}
	}()

	// 设置处理超时
	ctx, cancel := context.WithTimeout(task.Context, h.config.Timeout)
	defer cancel()

	start := time.Now()

	// 处理消息
	err := h.handleMessage(ctx, task.Event)

	duration := time.Since(start)
	h.metrics.RecordProcessTime(duration)

	if err != nil {
		h.metrics.RecordError("process", 1)
		fmt.Printf("Message processing error: %v\n", err)
	} else {
		h.metrics.RecordSuccess(1)
	}
}

// 事件检查（使用并发安全的map）
var eventCache sync.Map

func (h *LarkHandler) EventCheck(eventId string) bool {
	_, loaded := eventCache.LoadOrStore(eventId, time.Now())
	if loaded {
		return false
	}

	// 定期清理过期事件
	go func() {
		time.Sleep(5 * time.Minute)
		eventCache.Delete(eventId)
	}()

	return true
}

func (h *LarkHandler) OnP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	// 事件去重检查 - 在接收消息时立即检查
	eventID := event.EventV2Base.Header.EventID
	if !h.EventCheck(eventID) {
		return nil // 重复事件，直接返回
	}

	// 快速返回，避免阻塞 WebSocket 连接
	task := &MessageTask{
		Event:     event,
		Context:   ctx,
		CreatedAt: time.Now(),
	}

	select {
	case h.messageQueue <- task:
		h.metrics.RecordQueued(1)
		return nil
	default:
		h.metrics.RecordDropped(1)
		return fmt.Errorf("message queue full, dropping message")
	}
}

func (h *LarkHandler) OnP2MessageReadV1(ctx context.Context, event *larkim.P2MessageReadV1) error {
	fmt.Printf("收到 message_read 事件, 消息已读\n")
	return nil
}

func (h *LarkHandler) handleMessage(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	// 解析消息内容
	var respContent map[string]string
	err := json.Unmarshal([]byte(*event.Event.Message.Content), &respContent)
	if err != nil || *event.Event.Message.MessageType != "text" {
		return h.SendMessage(event, larkim.MsgTypeText, "解析消息失败，请发送文本消息\n")
	}

	reqText := respContent["text"]

	// 根据消息类型分发处理
	if h.isJenkinsRequest(reqText) {
		return h.handleJenkinsRequest(ctx, event, reqText)
	} else {
		return h.handleAIRequest(ctx, event, reqText)
	}
}

func (h *LarkHandler) isJenkinsRequest(text string) bool {
	plugin := h.pluginManager.FindPlugin(text)
	return plugin != nil
}

func (h *LarkHandler) handleAIRequest(ctx context.Context, event *larkim.P2MessageReceiveV1, reqText string) error {
	// AI 并发限制
	select {
	case h.aiSemaphore <- struct{}{}:
		defer func() { <-h.aiSemaphore }()
	case <-ctx.Done():
		return ctx.Err()
	}

	// 创建卡片
	cardId := h.Bot.CreateNewCard(ctx)
	if cardId == "" {
		return fmt.Errorf("failed to create card")
	}

	// 发送卡片消息
	err := h.SendMessage(event, larkim.MsgTypeInteractive, cardId)
	if err != nil {
		return err
	}

	// 异步处理 AI 对话
	go h.processAIChat(context.Background(), cardId, reqText)

	return nil
}

// 注册AI客户端
func (h *LarkHandler) RegisterAIClient(name string, client AIClient) {
	h.aiManager.RegisterClient(name, client)
}

// 设置默认AI客户端
func (h *LarkHandler) SetDefaultAI(name string) bool {
	return h.aiManager.SetDefault(name)
}

// 获取AI客户端
func (h *LarkHandler) GetAIClient(name string) AIClient {
	if name == "" {
		return h.aiManager.GetDefault()
	}
	return h.aiManager.GetClient(name)
}

func (h *LarkHandler) processAIChat(ctx context.Context, cardId, reqText string) {
	// 使用AI客户端管理器
	aiClient := h.aiManager.GetDefault()
	if aiClient == nil {
		fmt.Println("没有可用的AI客户端")
		h.Bot.UpdateCardChat(ctx, cardId, []string{"错误：没有可用的AI客户端"})
		return
	}

	fmt.Printf("使用AI客户端: %s\n", aiClient.Name())

	// 调用AI客户端的流式对话方法
	responseChunks, err := aiClient.Chat(reqText)
	if err != nil {
		fmt.Printf("AI对话失败: %v\n", err)
		h.Bot.UpdateCardChat(ctx, cardId, []string{fmt.Sprintf("AI对话失败: %v", err)})
		return
	}
	responseChunks = ReplaceQwenNameIfMatched(responseChunks)
	// 更新卡片内容
	h.Bot.UpdateCardChat(ctx, cardId, responseChunks)
}

func (h *LarkHandler) handleJenkinsRequest(ctx context.Context, event *larkim.P2MessageReceiveV1, reqText string) error {
	// HTTP 并发限制
	select {
	case h.httpSemaphore <- struct{}{}:
		defer func() { <-h.httpSemaphore }()
	case <-ctx.Done():
		return ctx.Err()
	}
	// 创建卡片
	cardId := h.Bot.CreateNewCard(ctx)
	if cardId == "" {
		return fmt.Errorf("failed to create card")
	}

	// 发送卡片消息
	err := h.SendMessage(event, larkim.MsgTypeInteractive, cardId)
	if err != nil {
		return err
	}

	// 异步处理 Jenkins 构建
	go h.processJenkinsRequest(context.Background(), event, reqText, cardId)

	return nil
}

func (h *LarkHandler) processJenkinsRequest(ctx context.Context, event *larkim.P2MessageReceiveV1, reqText, cardId string) {
	msg := h.buildJenkinsJobs(ctx, reqText, *event.Event.Message.MessageId)
	// 更新卡片内容
	h.Bot.UpdateCardChat(ctx, cardId, msg)
}

func (h *LarkHandler) buildJenkinsJobs(ctx context.Context, reqText, messageId string) []string {
	// 查找匹配的插件
	plugin := h.pluginManager.FindPlugin(reqText)
	if plugin == nil {
		return []string{"没有找到匹配的插件"}
	}

	// 解析任务
	jobs := plugin.ParseJobs(reqText, messageId)
	if len(jobs) == 0 {
		return []string{"没有找到需要构建的任务"}
	}

	fmt.Printf("使用插件: %s, 解析到 %d 个任务\n", plugin.Name(), len(jobs))

	// 并发执行任务
	var wg sync.WaitGroup
	results := make(chan error, len(jobs))

	for _, job := range jobs {
		wg.Add(1)
		go func(j Job) {
			defer wg.Done()
			client := h.httpPool.Get()
			defer h.httpPool.Put(client)

			fmt.Printf("执行任务: %s\n", j.Name)
			err := plugin.ExecuteJob(ctx, j, client)
			results <- err
		}(job)
	}

	wg.Wait()
	close(results)

	// 收集结果
	var errors []error
	successCount := 0
	for err := range results {
		if err != nil {
			errors = append(errors, err)
			fmt.Printf("任务执行失败: %v\n", err)
		} else {
			successCount++
		}
	}
	msg := make([]string, 0, len(jobs))
	msg = append(msg, fmt.Sprintf("开始解析任务...\n"))
	for _, job := range jobs {
		msg = append(msg, fmt.Sprintf("解析到任务: %s \n", job.Name))
	}
	msg = append(msg, fmt.Sprintf("以上任务开始构建: 成功 %d/%d, 失败 %d 个\n", successCount, len(jobs), len(errors)))
	return msg
}

// 添加插件管理方法
func (h *LarkHandler) RegisterPlugin(plugin JobPlugin) {
	h.pluginManager.RegisterPlugin(plugin)
}

func (h *LarkHandler) GetPlugins() []JobPlugin {
	return h.pluginManager.GetAllPlugins()
}

func (h *LarkHandler) SendMessage(event *larkim.P2MessageReceiveV1, msgType, msg string) error {
	fmt.Printf("发送消息: %s, %s, %s\n", *event.Event.Message.ChatType, msgType, msg)
	if *event.Event.Message.ChatType == "p2p" {
		return h.Bot.SendP2PReqMessage(*event.Event.Message.ChatId, msgType, msg)
	} else {
		return h.Bot.SendReplyReqMessage(*event.Event.Message.MessageId, msgType, msg)
	}
}

// 启动监控报告器
func (h *LarkHandler) StartMetricsReporter() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			processed, errors, queued, dropped, avgTime := h.metrics.GetStats()
			fmt.Printf("Handler Stats - Processed: %d, Errors: %d, Queued: %d, Dropped: %d, AvgTime: %v\n",
				processed, errors, queued, dropped, avgTime)
		case <-h.ctx.Done():
			return
		}
	}
}

// 优雅关闭
func (h *LarkHandler) Shutdown() {
	fmt.Println("Shutting down handler...")
	h.cancel()
	h.wg.Wait()
	fmt.Println("Handler shutdown completed")
}

func ReplaceQwenNameIfMatched(chunks []string) []string {
	joined := strings.Join(chunks, "")
	keywords := []string{"通义千问", "阿里巴巴", "阿里", "通义", "qwen", "Qwen"}

	// 是否匹配关键词
	matched := false
	for _, kw := range keywords {
		if strings.Contains(joined, kw) {
			matched = true
			break
		}
	}

	if !matched {
		// 没有敏感词，原样返回
		return chunks
	}

	// 替换关键词
	replacer := strings.NewReplacer(
		"通义千问", "ShitGPT",
		"qwen", "ShitGPT",
		"Qwen", "ShitGPT",
		"阿里巴巴", "诗与剑",
		"阿里", "诗与剑",
		"通义实验室", "云梦录工作室",
		"Alibaba ", "shit",
	)

	replaced := replacer.Replace(joined)

	// 拆分回数组：可以按原始片段长度（不变），或者统一每N字分一段
	const chunkSize = 20
	var result []string
	for i := 0; i < len(replaced); i += chunkSize {
		end := i + chunkSize
		if end > len(replaced) {
			end = len(replaced)
		}
		result = append(result, replaced[i:end])
	}

	return result
}
