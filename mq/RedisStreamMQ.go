package mq

import (
	"context"
	"encoding/json"
	"errors"
	"readygo/cache"
	"readygo/util"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// emptyMessageID 空的消息ID
var emptyMessageID = MessageID{}

// streamRawDataJSONKey  redis stream 中的值只能使用字符串,所以使用Json格式!!!!
const streamRawDataJSONKey = "stream_raw_data_json"

// MessageID 消息ID,用于隔离依赖和属性显示
type MessageID struct {
	ID           string
	QueueName    string
	GroupName    string
	ConsumerName string
}

// IMessageProducerConsumer 生产消费者的接口
type IMessageProducerConsumer[T any] interface {
	// GetQueueName 获取队列名称
	GetQueueName(ctx context.Context) string
	// GetGroupName 获取消息组名称
	GetGroupName(ctx context.Context) string
	// GetConsumerName 获取消费者名称
	GetConsumerName(ctx context.Context) string
	// GetCount 一次获取消息的总数,默认10
	GetCount(ctx context.Context) int
	// GetBlock 阻塞毫秒数,默认5000毫秒
	GetBlock(ctx context.Context) int
	// GetStart 获取消费的起始位置,默认 ">"
	GetStart(ctx context.Context) string
	// SendMessage 生产者发送消息
	SendMessage(ctx context.Context, messageObject T) (MessageID, error)
	// OnMessage 消费者处理消息
	OnMessage(ctx context.Context, messageID MessageID, messageObject T) (bool, error)
	// GetMinIdleTime XPENDING命令的min-idle-time毫秒数,避免处理最新的消息.默认300秒
	GetMinIdleTime(ctx context.Context) int
	//MaxRetryCount 最大的重试次数,默认20次
	GetMaxRetryCount(ctx context.Context) int
}

// MessageProducerConsumer 默认的消息队列实现
type MessageProducerConsumer[T any] struct {
	QueueName     string
	GroupName     string
	ConsumerName  string
	Count         int
	Block         int
	Start         string
	MinIdleTime   int
	MaxRetryCount int
}

func (messageProducerConsumer *MessageProducerConsumer[T]) GetQueueName(ctx context.Context) string {
	return messageProducerConsumer.QueueName
}
func (messageProducerConsumer *MessageProducerConsumer[T]) GetGroupName(ctx context.Context) string {
	return messageProducerConsumer.GroupName
}
func (messageProducerConsumer *MessageProducerConsumer[T]) GetConsumerName(ctx context.Context) string {
	return messageProducerConsumer.ConsumerName
}
func (messageProducerConsumer *MessageProducerConsumer[T]) GetCount(ctx context.Context) int {
	if messageProducerConsumer.Count == 0 {
		return 10
	}
	return messageProducerConsumer.Count
}
func (messageProducerConsumer *MessageProducerConsumer[T]) GetBlock(ctx context.Context) int {
	if messageProducerConsumer.Block == 0 {
		return 5000
	} else if messageProducerConsumer.Block < 0 {
		return 0
	}
	return messageProducerConsumer.Block
}
func (messageProducerConsumer *MessageProducerConsumer[T]) GetStart(ctx context.Context) string {
	if messageProducerConsumer.Start == "" {
		return ">"
	}
	return messageProducerConsumer.Start
}
func (messageProducerConsumer *MessageProducerConsumer[T]) SendMessage(ctx context.Context, messageObject T) (MessageID, error) {
	return sendMessage(ctx, messageProducerConsumer.QueueName, messageObject)
}

func (messageProducerConsumer *MessageProducerConsumer[T]) GetMinIdleTime(ctx context.Context) int {
	if messageProducerConsumer.MinIdleTime == 0 {
		return 300000
	}
	return messageProducerConsumer.MinIdleTime
}

func (messageProducerConsumer *MessageProducerConsumer[T]) GetMaxRetryCount(ctx context.Context) int {
	if messageProducerConsumer.MaxRetryCount == 0 {
		return 20
	}
	return messageProducerConsumer.MaxRetryCount
}

// createStreamConsumerGroup  创建 redis stream consumer group
// start 有 "0",从开始位置消费; $从最近的消息消费
func createStreamConsumerGroup(ctx context.Context, streamName, groupName, start string) error {
	if streamName == "" || groupName == "" {
		return errors.New("值不能为空")
	}
	// start 有 "0",从开始位置消费; $从最近的消息消费
	if start == "" {
		start = "0"
	}
	_, errResult := cache.RedisCMDContext(ctx, "xgroup", "create", streamName, groupName, start, "MKSTREAM")
	// 获值错误
	if errResult != nil {
		if strings.Contains(errResult.Error(), "already exists") { // 已经存在,不再创建
			return nil
		}
		util.FuncLogError(ctx, errResult)
		return errResult
	}

	return nil
}

// SendMessage  发送消息队列
func sendMessage[T any](ctx context.Context, queueName string, messageObject T) (MessageID, error) {
	jsonData, err := json.Marshal(messageObject)
	if err != nil {
		return emptyMessageID, err
	}

	args := make([]interface{}, 0)
	args = append(args, "xadd")
	args = append(args, queueName)
	args = append(args, "*")
	args = append(args, streamRawDataJSONKey)
	args = append(args, jsonData)

	//result, errResult := cache.RedisCMDContext(ctx, "xadd", streamName, "*", values)
	result, errResult := cache.RedisCMDContext(ctx, args...)
	if errResult != nil {
		return emptyMessageID, errResult
	}
	messageID := MessageID{QueueName: queueName}
	if msgID, ok := result.(string); ok {
		messageID.ID = msgID
	}
	return messageID, nil
}

// StartConsumer 启动一个消费者
func StartConsumer[T any](ctx context.Context, messageProducerConsumer IMessageProducerConsumer[T]) error {
	queueName := messageProducerConsumer.GetQueueName(ctx)
	groupName := messageProducerConsumer.GetGroupName(ctx)
	consumerName := messageProducerConsumer.GetConsumerName(ctx)
	count := messageProducerConsumer.GetCount(ctx)
	block := messageProducerConsumer.GetBlock(ctx)
	start := messageProducerConsumer.GetStart(ctx)

	if queueName == "" || groupName == "" || consumerName == "" {
		return errors.New("queueName or groupName or consumerName is empty")
	}

	//先创建组
	errGroup := createStreamConsumerGroup(ctx, queueName, groupName, "0")
	if errGroup != nil {
		util.FuncLogError(ctx, errGroup)
		return errGroup
	}

	for {
		// 使用 XREADGROUP 以阻塞方式读取消息
		// >：获取​​从未被该消费者组内任何消费者领取过​​的"全新"消息.这是最常用的模式.
		// 0：​​重新获取​​那些已经被领取但还躺在 PEL 中"未签收"的消息.常用于故障恢复和重试.
		// 两者都会获取到未确认的消息,但 > 是向前看(新消息),0是回头看(未完成的消息).

		streams, errResult := cache.RedisCMDContext(ctx, "xreadgroup", "group", groupName, consumerName, "count", count, "block", block, "streams", queueName, start)

		if errResult != nil {
			if errResult == redis.Nil { // 超时,继续轮询
				continue
			}
			util.FuncLogError(ctx, errResult)
			time.Sleep(1 * time.Second) // 出错后稍作等待
			continue
		}

		//map[geo:[[1758014730551-0 [name test]]]]
		streamMap := streams.(map[interface{}]interface{})
		for k, v := range streamMap {
			streamName := k.(string) //获取队列名称
			vs := v.([]interface{})
			for _, msgObject := range vs { //循环消息
				msgs := msgObject.([]interface{})
				if len(msgs) != 2 {
					continue
				}
				msgId := msgs[0].(string)
				values := msgs[1].([]interface{})
				if len(values) != 2 {
					continue
				}
				messageID := MessageID{
					ID:           msgId,
					QueueName:    streamName,
					GroupName:    groupName,
					ConsumerName: consumerName,
				}
				msgObjectBytes := values[1].(string)
				messageObj := new(T)
				err := json.Unmarshal([]byte(msgObjectBytes), messageObj)
				if err != nil {
					util.FuncLogError(ctx, err)
					continue
				}
				ok, err := messageProducerConsumer.OnMessage(ctx, messageID, *messageObj)
				if err != nil {
					util.FuncLogError(ctx, err)
					continue
				}
				if !ok {
					continue
				}
				result, errResult := cache.RedisCMDContext(ctx, "xack", queueName, groupName, consumerName, msgId)
				if errResult != nil {
					util.FuncLogError(ctx, errResult)
					continue
				}
				if result.(int) == 1 { //success

				}

			}
		}

	}
}

// RetryConsumer 重试消费者的XPENDING消息.minIdleTime是消息的最小空闲毫秒,只有空闲时间超过此值的消息才会被重试
// 使用 XPENDING,XCLAIM,XRANGE 然后调用OnMessage处理
func RetryConsumer[T any](ctx context.Context, messageProducerConsumer IMessageProducerConsumer[T]) {
	queueName := messageProducerConsumer.GetQueueName(ctx)
	groupName := messageProducerConsumer.GetGroupName(ctx)
	consumerName := messageProducerConsumer.GetConsumerName(ctx)
	count := messageProducerConsumer.GetCount(ctx)
	//block := messageProducerConsumer.GetBlock(ctx)
	minIdleTime := messageProducerConsumer.GetMinIdleTime(ctx)
	maxRetryCount := messageProducerConsumer.GetMaxRetryCount(ctx)
	//start := messageProducerConsumer.GetStart(ctx)
	for {
		// XPENDING geo deepseek_group IDLE 300000 - + 10 consumer1
		xpending, errResult := cache.RedisCMDContext(ctx, "xpending", queueName, groupName, "idle", minIdleTime, "-", "+", count, consumerName)
		if errResult != nil || xpending == nil {
			util.FuncLogError(ctx, errResult)
			continue
		}
		msgsSlice := xpending.([]interface{})
		if len(msgsSlice) == 0 { //没有消息,睡眠一会
			time.Sleep(time.Millisecond * time.Duration(minIdleTime))
			continue
		}

		for _, msgObject := range msgsSlice { //循环所有的消息
			msg := msgObject.([]interface{})
			if len(msg) != 4 {
				continue
			}
			msgId := msg[0].(string)
			consumerName := msg[1].(string)
			//idleTime := msg[2].(int64)
			deliveryCount := msg[3].(int64)
			if deliveryCount > int64(maxRetryCount) { //超过最大的投递次数,强制ACK消息
				_, err := cache.RedisCMDContext(ctx, "xack", queueName, groupName, consumerName, msgId)
				if err != nil {
					util.FuncLogError(ctx, err)
				}
				continue
			}
			// 使用 XCLAIM 迁移id给自己,用于XPENDING重新计算投递次数
			// XCLAIM geo deepseek_group consumer1 0  1758158851046-0 1758158851042-0 JUSTID
			_, err := cache.RedisCMDContext(ctx, "xclaim", queueName, groupName, consumerName, 0, msgId, "JUSTID")
			if err != nil {
				util.FuncLogError(ctx, err)
				continue
			}

			//XRANGE geo 1758018581240-0 1758018581240-0 ,不会增加 投递次数(delivery count),需要使用 XCLAIM 重新分配消息给自己
			// xreadgroup 参数id是开始并不包括!!!!,所以使用 XRANGE geo 1758018581240-0 1758018581240-0 读取到消息内容,然后调用OnMessage
			xrange, err := cache.RedisCMDContext(ctx, "xrange", queueName, msgId, msgId)
			if err != nil {
				util.FuncLogError(ctx, err)
				continue
			}

			// [[1758165638806-0 [abc 456]]]
			//fmt.Println(xrange)
			xs := xrange.([]interface{})
			if len(xs) != 1 {
				continue
			}
			msgObj := (xs[0]).([]interface{})
			if len(msgObj) != 2 || msgObj[0].(string) != msgId {
				continue
			}
			msgJson := (msgObj[1]).([]interface{})
			//if len(msgJson) != 2 {
			if len(msgJson) != 2 || msgJson[0].(string) != streamRawDataJSONKey {
				continue
			}

			messageID := MessageID{
				ID:           msgId,
				QueueName:    queueName,
				GroupName:    groupName,
				ConsumerName: consumerName,
			}
			msgObjectBytes := msgJson[1].(string)
			messageObj := new(T)
			err = json.Unmarshal([]byte(msgObjectBytes), messageObj)
			if err != nil {
				util.FuncLogError(ctx, err)
				continue
			}
			ok, err := messageProducerConsumer.OnMessage(ctx, messageID, *messageObj)
			if err != nil {
				util.FuncLogError(ctx, err)
				continue
			}
			if !ok {
				continue
			}
			result, errResult := cache.RedisCMDContext(ctx, "xack", queueName, groupName, consumerName, msgId)
			if errResult != nil {
				util.FuncLogError(ctx, errResult)
				continue
			}
			if result.(int) == 1 { //success

			}

		}

	}
}
