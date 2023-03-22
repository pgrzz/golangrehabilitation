package mq

import (
	"fmt"
	"sync"
)

// 提供broker注册，管理着topic-broker的映射关系
type NameServer struct {
	topicBroker sync.Map
}

type BrokerNotFoundError struct {
	Topic string
}

func (e BrokerNotFoundError) Error() string {
	return fmt.Sprintf("could not find broker for topic %s", e.Topic)
}

// 根据topic获取到对应的broker
func (n *NameServer) GetBrokerBytopic(topic string) (Broker, error) {

	value, ok := n.topicBroker.Load(topic)
	if ok {
		broker := value.(Broker)
		return broker, nil
	}
	return Broker{}, BrokerNotFoundError{Topic: topic}
}

//broker注册到注册中心

func (n *NameServer) register(topic string, b Broker) {
	n.topicBroker.Store(topic, b)
}
