package mq

type Message struct {
	topic      string
	messageId  uint64
	messageKey string //use for hashIndex   key:topic+"#"+key,value:FileIndex	map[key][FileIndex]
	len        int
	body       []byte
}
