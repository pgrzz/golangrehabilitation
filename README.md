# golangrehabilitation
go相关的内容放在这里
其中mq 是用golang 实现的一个mini版本的mq，store 层基于文件+offset，的方式，在broker中消息直接从网络中到channel中就可以ack，
消费channel+store channel 同时尝试拉取消息，（这样减少了到磁盘的io，如果有消费者能力则尽最大努力投递），否则 store 层消费消息并且持久化，
可以理解持久化机制为内存+磁盘的方式。
在这之上也可以很容易实现容错机制（master-slave,额外增加一个slave消费者channel，如果要强一致那就channel接受到ack 使用两阶段提交需要额外增加一个文件结构中的version字段，如果接受最终一致性则可以通过循环重试即可）。
