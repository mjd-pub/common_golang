/**
* @Author: zhangjian@mioji.com
* @Date: 2019/9/9 下午6:05
 */
package mq

import (
	"log"
	"testing"
)

func TestProducer(t *testing.T) {
	product := NewProducerMQ("writer", "miaoji1109", "10.10.7.241", "5672", "dev")
	defer product.Close()
	err := product.Declare("publish_test", "zhangjian", "key")
	if err != nil {
		log.Println(err)
		return
	}
	err = product.Publish("this is a msg")
	if err != nil {
		log.Println(err)
		return
	}
}

func TestConsumer(t *testing.T) {
	mq := NewConsumerMQ("writer", "miaoji1109", "10.10.7.241", "5672", "dev")
	//注册消费者1
	mq.RegisterReceiver("publish_test", 0, Queue1)
	//注册消费者2
	mq.RegisterReceiver("publish_test", 0, Queue1)

	mq.Start()
}

//OnReceiver 消费消息 返回0代表正常（方便以后扩展处理）
func Queue1(msg []byte) int {
	//fmt.Println(string(msg))
	return 0
}
