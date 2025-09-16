package mq

import (
	"context"
	"errors"
	"fmt"
	"readygo/cache"
	"readygo/util"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var emptyMessageID = MessageID{}

var messageByteDataKey = "messageByteData"

// MessageID 消息ID,用于隔离依赖和属性显示
type MessageID struct {
	ID string
}

// IMessageProducerConsumer 生产消费者的接口
type IMessageProducerConsumer interface {
	// GetQueueName 获取队列名称
	GetQueueName(ctx context.Context) string
	// GetGroupName 获取消息组名称
	GetGroupName(ctx context.Context) string
	// GetConsumerName 获取消费者名称
	GetConsumerName(ctx context.Context) string
	// GetMessageObject 获取消息对象,是指针对象
	GetMessageObject(ctx context.Context) interface{}
	// GetStart 获取消费的起始位置
	GetStart(ctx context.Context) string
	// OnMessage 生产者发送消息
	SendMessage(ctx context.Context, messageObject interface{}) (MessageID, error)
	// OnMessage 消费者处理消息
	OnMessage(ctx context.Context, messageID MessageID, messageObject interface{}) (bool, error)
}

// MessageProducerConsumer 默认的消息队列实现
type MessageProducerConsumer struct {
	QueueName     string
	GroupName     string
	ConsumerName  string
	MessageObject interface{}
	Start         string
}

func (messageProducerConsumer *MessageProducerConsumer) GetQueueName(ctx context.Context) string {
	return messageProducerConsumer.QueueName
}
func (messageProducerConsumer *MessageProducerConsumer) GetGroupName(ctx context.Context) string {
	return messageProducerConsumer.GroupName
}
func (messageProducerConsumer *MessageProducerConsumer) GetConsumerName(ctx context.Context) string {
	return messageProducerConsumer.ConsumerName
}
func (messageProducerConsumer *MessageProducerConsumer) GetMessageObject(ctx context.Context) interface{} {
	return messageProducerConsumer.MessageObject
}
func (messageProducerConsumer *MessageProducerConsumer) GetStart(ctx context.Context) string {
	return messageProducerConsumer.Start
}
func (messageProducerConsumer *MessageProducerConsumer) OnMessage(ctx context.Context, messageID MessageID, messageObject interface{}) (bool, error) {
	return false, nil
}
func (messageProducerConsumer *MessageProducerConsumer) SendMessage(ctx context.Context, messageObject interface{}) (MessageID, error) {
	if messageObject == nil {
		return emptyMessageID, errors.New("messageObject is nil")
	}
	bytedata, err := util.Marshal(messageObject)
	if err != nil {
		return emptyMessageID, err
	}
	data := make([]interface{}, 0, 2)
	data[0] = messageByteDataKey
	data[1] = bytedata

	return SendMessage(ctx, messageProducerConsumer.GetQueueName(ctx), data)
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
func SendMessage(ctx context.Context, streamName string, values []interface{}) (MessageID, error) {
	_, errResult := cache.RedisCMDContext(ctx, "xadd", streamName, "*", values)

	if errResult != nil {
		return emptyMessageID, errResult
	}

	return emptyMessageID, nil

}

// StartConsumer 启动一个消费者
func StartConsumer(ctx context.Context, messageProducerConsumer IMessageProducerConsumer) error {
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

		streams, errResult := cache.RedisCMDContext(ctx, "xreadgroup", "group", groupName, consumerName, "count", 10, "block", 5*time.Second, "streams", queueName, ">")

		if errResult != nil {
			if errResult == redis.Nil { // 超时,继续轮询
				continue
			}
			time.Sleep(1 * time.Second) // 出错后稍作等待
			continue
		}

		fmt.Println(streams)

		/*
			// 处理从所有 Stream 中读取到的消息
			for _, stream := range streams {
				for _, message := range stream.Messages {
					// 调用回调函数处理消息
					if err := handler(message); err != nil {
						// 可根据错误类型决定是否重试或放入死信队列
						continue
					}

					// 处理成功,发送 ACK 确认消息
					if err := rdb.XAck(ctx, streamName, groupName, message.ID).Err(); err != nil {

					} else {

					}
				}
			}
		*/

	}
}
