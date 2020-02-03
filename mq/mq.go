//mq 实现rabbitmq的producer和consumer
package mq

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"common_golang/utils"
	"log"
	"sync"
	"time"
)

// ReceiveMQ 用于管理和维护rabbitmq的对象
type ConsumerMQ struct {
	wg        sync.WaitGroup
	channel   *amqp.Channel
	receivers []Receiver

	user   string
	passwd string
	host   string
	port   string
	vhost  string
}

type ProducerMQ struct {
	channel  *amqp.Channel
	declare  int // 0 未初始化 1 已初始化
	queue    string
	exchange string
	key      string
}

// Receiver 观察者模式需要的接口
type Receiver struct {
	logStat     int              //0 默认记录日志 1 不记录日志
	queueName   string           // 获取接收者需要监听的队列
	receiveFunc func([]byte) int // 处理收到的消息 返回0代表正常
}

// New 建立RabbitMQ连接
func newConn(user string, passwd string, host string, port string, vhost string) *amqp.Channel {
	rabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, passwd, host, port)
	conn, err := amqp.DialConfig(rabbitUrl, amqp.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
		Vhost:     vhost,
	})
	if err != nil {
		panic(err)
	}
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return channel
}

// New 创建一个新的操作RabbitMQ的对象
func NewConsumerMQ(user string, passwd string, host string, port string, vhost string) *ConsumerMQ {
	channel := newConn(user, passwd, host, port, vhost)
	return &ConsumerMQ{
		channel: channel,
		user:    user,
		passwd:  passwd,
		host:    host,
		port:    port,
		vhost:   vhost,
	}
}

// Start 启动Rabbitmq的客户端
func (mq *ConsumerMQ) Start() {
	for {
		mq.run()
		// 一旦连接断开，那么需要隔一段时间去重连
		time.Sleep(10 * time.Second)
		log.Println("链接中断，正在发起重试")
		mq.retry()
		log.Println("中断重试连接成功")
	}
}

// RegisterReceiver 注册一个用于接收指定队列指定路由的数据接收者
func (mq *ConsumerMQ) RegisterReceiver(queueName string, logStat int, receiveFunc func([]byte) int) {
	mq.receivers = append(mq.receivers, Receiver{
		queueName:   queueName,
		receiveFunc: receiveFunc,
		logStat:     logStat,
	})
}

// run 开始获取连接并初始化相关操作
func (mq *ConsumerMQ) run() {

	if mq.channel == nil {
		log.Printf("rabbit刷新连接失败，将要重连: %s\n", mq.host)
		return
	}

	for _, receiver := range mq.receivers {
		mq.wg.Add(1)
		go mq.listen(receiver) // 每个接收者单独启动一个goroutine接收消息
	}

	mq.wg.Wait()

	log.Println("所有处理queue的任务都意外退出了")
	// 理论上mq.run()在程序的执行过程中是不会结束的
	// 一旦结束就说明所有的接收者都退出了，那么意味着程序与rabbitmq的连接断开
	// 那么则需要重新连接，这里尝试销毁当前连接
	mq.channel.Close()
}

// retry 重连
func (mq *ConsumerMQ) retry() {
	mq.channel = newConn(mq.user, mq.passwd, mq.host, mq.port, mq.vhost)
}

// Listen 监听指定路由发来的消息
func (mq *ConsumerMQ) listen(receiver Receiver) {
	defer mq.wg.Done()
	// 这里获取每个接收者需要监听的队列和路由
	queueName := receiver.queueName
	// 获取消费通道 确保rabbitmq会一个一个发消息
	mq.channel.Qos(1, 0, true)
	msgs, err := mq.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if nil != err {
		log.Printf("获取队列 %s 的消费通道失败: %s", queueName, err.Error())
		return
	}

	// 使用callback消费数据
	for msg := range msgs {
		// 当接收者消息处理失败的时候，
		// 比如网络问题导致的数据库连接失败，redis连接失败等等这种
		// 通过重试可以成功的操作，那么这个时候是需要重试的
		receiver.receiveFunc(msg.Body)
		// 确认收到本条消息, multiple必须为false
		msg.Ack(false)
		if receiver.logStat == 0 {
			printLog(1, string(msg.Body), receiver.queueName)
		}
	}
}

// NewProductMQ
func NewProducerMQ(user string, passwd string, host string, port string, vhost string) *ProducerMQ {
	channel := newConn(user, passwd, host, port, vhost)
	return &ProducerMQ{
		channel: channel,
	}
}

// Declare 初始化队列 交换机 队列
func (p *ProducerMQ) Declare(queue, exchange, key string) error {
	p.queue = queue
	p.exchange = exchange
	p.key = key
	// 确认队列是否存在 不存在就declare
	_, err := p.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	// 确认exchange是否存在
	err = p.channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}
	// 绑定exchange key和队列
	err = p.channel.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		return err
	}
	p.declare = 1
	return nil
}

// Publish 消息推送 推完消息后需业务端自己决定关闭的时机
func (p *ProducerMQ) Publish(msg string) error {
	if p.declare == 0 {
		return errors.New("Declare未执行")
	}
	err := p.channel.Publish(p.exchange, p.key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		return err
	}
	printLog(0, msg, p.queue)
	return nil
}

// Close 关闭连接
func (p *ProducerMQ) Close() {
	p.channel.Close()
}

// printLog 打印日志
func printLog(mode int, msg string, queue string) {
	var (
		oriIp, destIp string
		errorId                   int
	)
	msgMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(msg), &msgMap)
	if err != nil {
		errorId = -1
	}
	if mode == 0 {
		oriIp, _ = utils.GetLocalHostIp()
	} else {
		destIp, _ = utils.GetLocalHostIp()
	}
	qid, ok := msgMap["qid"].(string)
	if !ok {
		qid = ""
	}
	csuid, ok := msgMap["csuid"].(string)
	if !ok {
		csuid = ""
	}
	modeType, ok := msgMap["type"].(string)
	if !ok {
		modeType = ""
	}
	uid, ok := msgMap["uid"].(string)
	if !ok {
		uid = ""
	}
	role, ok := msgMap["role"].(string)
	if !ok {
		role = ""
	}
	ptid, ok := msgMap["ptid"].(string)
	if !ok {
		ptid = ""
	}
	acc, ok := msgMap["acc"].(string)
	if !ok {
		acc = ""
	}
	mqLog := map[string]interface{}{
		"qid":      qid,
		"type":     modeType,
		"role":     role,
		"uid":      uid,
		"csuid":    csuid,
		"acc":      acc,
		"ptid":     ptid,
		"queue":    queue,
		"msg_mode": mode,
		"t":        time.Now().Format("2006-01-02T15:04:05+08:00"),
		"key":      "MJ_MD_MQ_LOG",
		"ori_ip":   oriIp,
		"dest_ip":  destIp,
		"msg":      msg,
		"error_id": errorId,
	}
	logByte, _ := json.Marshal(mqLog)
	fmt.Printf("\n%s\n", string(logByte))
}
