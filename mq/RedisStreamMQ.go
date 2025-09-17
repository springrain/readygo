package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"readygo/cache"
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
	// GetStart 获取消费的起始位置
	GetStart(ctx context.Context) string
	// OnMessage 生产者发送消息
	SendMessage(ctx context.Context, messageObject T) (MessageID, error)
	// OnMessage 消费者处理消息
	OnMessage(ctx context.Context, messageID MessageID, messageObject T) (bool, error)
}

// MessageProducerConsumer 默认的消息队列实现
type MessageProducerConsumer[T any] struct {
	QueueName    string
	GroupName    string
	ConsumerName string
	Start        string
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

func (messageProducerConsumer *MessageProducerConsumer[T]) GetStart(ctx context.Context) string {
	return messageProducerConsumer.Start
}
func (messageProducerConsumer *MessageProducerConsumer[T]) SendMessage(ctx context.Context, messageObject T) (MessageID, error) {
	return sendMessage(ctx, messageProducerConsumer.QueueName, messageObject)
}

// CreateStreamConsumerGroup  创建 redis stream consumer group
// start 有 "0",从开始位置消费; $从最近的消息消费
func CreateStreamConsumerGroup(ctx context.Context, streamName, groupName, start string) error {
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
	messageID := MessageID{}
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

	if queueName == "" || groupName == "" || consumerName == "" {
		return errors.New("queueName or groupName or consumerName is empty")
	}

	//先创建组
	errGroup := CreateStreamConsumerGroup(ctx, queueName, groupName, "0")
	if errGroup != nil {
		return errGroup
	}

	for {
		// 使用 XREADGROUP 以阻塞方式读取消息
		// >：获取​​从未被该消费者组内任何消费者领取过​​的"全新"消息.这是最常用的模式.
		// 0：​​重新获取​​那些已经被领取但还躺在 PEL 中"未签收"的消息.常用于故障恢复和重试.
		// 两者都会获取到未确认的消息,但 > 是向前看(新消息),0是回头看(未完成的消息).

		streams, errResult := cache.RedisCMDContext(ctx, "xreadgroup", "group", groupName, consumerName, "count", 10, "streams", queueName, ">")

		if errResult != nil {
			if errResult == redis.Nil { // 超时,继续轮询
				continue
			}
			time.Sleep(1 * time.Second) // 出错后稍作等待
			continue
		}

		fmt.Println(streams)
		//map[geo:[[1758014730551-0 [name test]]]]
		streamMap := streams.(map[interface{}]interface{})
		for k, v := range streamMap {
			streamName := k.(string) //获取队列名称
			fmt.Println(streamName)
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
					GroupName:    groupName,
					QueueName:    queueName,
					ConsumerName: consumerName,
				}
				msgObjectBytes := values[1].(string)
				messageObj := new(T)
				err := json.Unmarshal([]byte(msgObjectBytes), messageObj)
				if err != nil {
					continue
				}
				ok, err := messageProducerConsumer.OnMessage(ctx, messageID, *messageObj)
				if err != nil {
					continue
				}
				if !ok {
					continue
				}
				result, errResult := cache.RedisCMDContext(ctx, "xack", queueName, groupName, consumerName, msgId)
				if errResult != nil {
					continue
				}
				if result.(int) == 1 { //success

				}

			}
		}

	}
}
