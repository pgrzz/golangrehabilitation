package mq

import (
	"fmt"
	"sync"
)

type Broker struct {
	store *Store
	//cache 	topic:chan
	// per topic-buffer	在golang中使用chan来代替java中的队列用消息传递来代替并发的bufferPoll
	//真是一个非常不错的替代
	//同时channel天然的并发特性可以使的异步刷盘的情况下吞吐达到一个很大的量
	//可以定义多个 go rountine 并发消费 发送到tcp缓冲区，使得的broker性能达到最大值
	//channel的最大长度为一个int
	topicBuffers sync.Map
	//topic config
	tcs sync.Map
}

type topicConfig struct {
	topic       string
	ackSyncType int //1 代表 发送到channel 也就是buffer就可以ack， 2代表需要持久化后才能返回
	parallelnum int //default 1,代表会有多少个go rountine 来处理这个topic对应的channel
	//这个修改起来太麻烦了=。=，虽然在go来说可以搞 但是和rocketMq那些一样底层的buffer
	//需要用 cas自旋加上三个状态 1、停止接受 自旋阻塞 持续读 2、停止接受 停止读  修改引用关系resizeChannel（）3、可以读 可以写
	//所以保持和其他的mq一样重启broker生效
	bufferSize int //只是表示的数量，注意这里的实际size是 int * lens * messageBodySize

}

func newBroker(s *Store) *Broker {
	br := Broker{store: s}

	br.tcs.Range(func(key, value interface{}) bool {
		// 处理 key 和 value
		topic := key.(string)
		tc := value.(topicConfig)
		//根据tc分配channel
		buffer := make(chan Message, tc.bufferSize)
		br.topicBuffers.Store(topic, buffer)

		//根据并性度分配channel,这样同一个topic下是不能够保证顺序消费的,顺序消费需要设置parallelnum为1
		for i := 0; i < tc.parallelnum; i++ {
			go br.asyncWriteMsg(buffer)
		}

		return true
	})

	return &br
}

func (br *Broker) mockNameServerInitTopicConfig() {

	t1 := &topicConfig{topic: "a1", ackSyncType: 1, bufferSize: 1024}

	br.tcs.Store(t1.topic, t1)

}

// 收到producer发来的消息存储到store	//在这里可以做一层cache，
// 提供两种模式 1、一种写到broker时就认为成功可以返回ack，	也就是异步落盘
//
//	2、一种是等到message写入到store才返回ack信息
//
// 从nameServer加载topic的配置信息来决定采取哪种方式
func (b *Broker) receivedMessage(msg Message) (bool, error) {

	value, ok := b.topicBuffers.Load(msg.Topic)

	if !ok {
		return false, BrokerNotFoundError{Topic: msg.Topic}
	}
	buffer := value.(chan Message)
	value, ok = b.tcs.Load(msg.Topic)
	if !ok {
		return false, BrokerNotFoundError{Topic: msg.Topic}
	}
	config := value.(topicConfig)

	//打点 如果buffer 超过了 50%的积压
	go func() {
		if len(buffer) > config.bufferSize/2 {
			fmt.Sprintf("buffer obj wait to consume nums over 50 percent topic: %s", config.topic)
		}
	}()

	if config.ackSyncType == 1 {
		buffer <- msg
		return true, nil
	} else if config.ackSyncType == 2 {
		//同步直接写消息,由于file层已经有了,在异步消息的时候先查询一次如果topic的类型为ackSyncType==2
		cw := &crc32Writer{}
		mm := msg.message(cw)
		wb := &WriteBuffer{}
		mm.writeTo(wb)
		offset, err := b.store.WriteMessage(msg.Topic, wb.b[:])
		if err != nil {
			return false, err
		}
		msg.Offset = offset
		return true, nil
	}
	return false, nil
}

//这里会维护一个已经得到ack 的offset集合， 然后集合从小到大排序, 找到db当前的offset值，
//通过对当前的offset+1 找到最小开始点，如果没有的话就循环等待直到集合中的len>1 && db.offset+1 ==min(Msgs.offset)
//,然后看有多长的连续区间,可以在db一次性更新这么长的ack
//每一个发出去的msg都会加入到一个 wait to ack的map中  key就是  msg的唯一id（offset），
//value是计时器，会有一个计时器更新value的存活double rtt，
//如果超过3次 rtt 更新，就说明丢包了，需要根据offset去db中读取数据，

func (br *Broker) pollMsg(topic string) {

}

func (br *Broker) asyncWriteMsg(ch chan Message) {
	msg := <-ch

	value, ok := br.tcs.Load(msg.Topic)
	if !ok {
		fmt.Sprintf("load config fail when asyncWriteMsg topic: %s", msg.Topic)
	}
	config := value.(topicConfig)
	if config.ackSyncType == 1 {
		// flush memory
		cw := &crc32Writer{}
		mm := msg.message(cw)
		wb := &WriteBuffer{}
		mm.writeTo(wb)
		offset, err := br.store.WriteMessage(msg.Topic, wb.b[:])
		if err != nil {
			fmt.Sprintf("write msg fial when asyncWriteMsg topic: %s", msg.Topic)
		}
		msg.Offset = offset
	}
	//send to biz channel,tcp那边收到

}

type BrokerError struct {
	Topic string
}

func (e BrokerNotFoundError) BrokerError() string {
	return fmt.Sprintf("could not find broker for topic %s", e.Topic)
}

type TimeOut struct {
	Topic string
}

func (e TimeOut) Error() string {
	return fmt.Sprintf("channel full topic: %s", e.Topic)
}
