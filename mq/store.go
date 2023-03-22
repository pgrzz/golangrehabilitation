package mq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

//每一个topic 会创建一个文件，当offset

// 一个header的offset
const FILEHEADER int64 = 8
const threshold int64 = 1073741824
const fileIndexPath string = "/Users/guanru.gr/go/src/myproject/golangrehabilitation/mq/resource/store.txt"
const fileIndexPre string = "/Users/guanru.gr/go/src/myproject/golangrehabilitation/mq/resource/"

// 打开文件时 文件头部为 8byte offset 然后后续偏移才是文件大小
type Store struct {
	//一组打开的文件句柄
	topicFiles sync.Map
	l          sync.Mutex
}

// filePath      use topic instead
// header  offset 8 messageLength 8
type topicStore struct {
	topic         string
	offset        int64 //    1024 1k * 1024 1m * 1024 1g
	messageLength int64
	version       int64 // after offset>thresshold(1G) 64 new topic will version+1 and offset unSubmitOffset clean to 0
}

type topicFile struct {
	f  *os.File
	ts topicStore
	mu sync.RWMutex
}

type ReaderFail struct {
	Topic string
}

// 偏移量不做清0 version 的版本是这样计算的 offset&threshold 取余操作
// 当offset+messageLength>threshold 时去新建一个覆盖当前的file
// 如果Store不做持久化操作那么断电后就没有可用性了
func NewStore2() *Store {
	p := &Store{}

	//创建文件index
	if _, err := os.Stat(fileIndexPath); os.IsNotExist(err) {
		// 文件不存在，创建文件
		file, err := os.Create(fileIndexPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}
	data, err := ioutil.ReadFile(fileIndexPath)
	if err != nil {
		panic(err)
	}
	var tt []topicStore
	json.Unmarshal(data, &tt)
	//open each file and load to  topicFiles
	for i := 0; i < len(tt); i++ {
		filePath := fileIndexPre + tt[i].topic + strconv.FormatInt(tt[i].version, 10)
		tempFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			f := &topicFile{f: tempFile, ts: tt[i]}
			p.topicFiles.Store(tt[i].topic, f)
		}
	}
	return p
}

// 当创建broker的时候进行注册
func (s *Store) Register(t topicStore) (bool, error) {
	filePath := fileIndexPre + t.topic + strconv.FormatInt(t.version, 10)
	tempFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		f := &topicFile{f: tempFile, ts: t}
		s.topicFiles.Store(t.topic, f)
		return true, nil
	}
	return false, err
}

// 这里为了方便使用json序列化，所有的文件结构存储在这里
func (s *Store) ReCreateTopicfile(t topicStore) (*topicFile, error) {
	//load and replace
	if value, ok := s.topicFiles.LoadAndDelete(t.topic); ok {
		tf := value.(topicFile)
		tf.mu.Lock()
		defer tf.mu.Unlock()
		defer tf.f.Close()
		oldVersion := tf.ts.version
		newVersion := tf.ts.offset & threshold
		filePath := fileIndexPre + t.topic + strconv.FormatInt(newVersion, 10)
		tempFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.version = newVersion
			f := topicFile{f: tempFile, ts: t, mu: sync.RWMutex{}}
			s.topicFiles.Store(t.topic, f)
			return &f, nil
		}

		//更新主表 定义
		//直接将内存中的值flush到磁盘
		s.flush()

		//异步删除对应的块
		oldFilePath := fileIndexPre + t.topic + strconv.FormatInt(oldVersion, 10)
		go deleteTopicfile(oldFilePath)
	}
	return nil, nil

}

// 在fileIndexPre 下扫描
func deleteTopicfile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Sprintf("delete file fail %s", filePath)
	}
}

// 定期刷盘操作更新主文件，和topic文件中的offset
// 同步内存数据到磁盘
// 这里有一个比较大问题，当topic很多的话这种偷懒的方式会造成明显的数据不一致同时保持meta最新的成本很高。
// 合适的方式是对更新部分做增量更新操作，
func (s *Store) flush() {

	s.l.Lock()
	defer s.l.Unlock()

	var tss []topicStore

	s.topicFiles.Range(func(key, value interface{}) bool {
		// 处理 key 和 value
		topicFile := value.(topicFile)
		tss = append(tss, topicFile.ts)
		return true
	})

	data, err := json.Marshal(tss)

	if err != nil {
		//log
		fmt.Sprintf("flush file  json fail %s", err.Error())
	}

	err = ioutil.WriteFile(fileIndexPath, data, 0644)

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func (s *Store) Shutdown() {

	s.topicFiles.Range(func(key, value interface{}) bool {
		// 处理 key 和 value
		topic := key.(string)
		topicFile := value.(topicFile)
		filePath := fileIndexPre + topicFile.ts.topic + strconv.FormatInt(topicFile.ts.version, 10)
		fmt.Sprintf("shutdown file topic,filePath %s,%s", topic, filePath)
		defer topicFile.f.Close()
		return true
	})

}

func (e ReaderFail) Error() string {
	return fmt.Sprintf("could not find broker for topic %s", e.Topic)
}

// offset 在初始化的时候就做好偏移, offset
func (s *Store) ReadTopic(topic string) ([]byte, error) {

	value, ok := s.topicFiles.Load(topic)

	if ok {
		tf := value.(topicFile)

		tf.mu.RLock()
		defer tf.mu.RUnlock()

		length := tf.ts.messageLength
		offset := tf.ts.offset + length
		data := make([]byte, length)
		_, err := tf.f.ReadAt(data, offset)
		if err != nil {
			return nil, err
		}
		//读取成功更新offset的值
		offset += length
		tf.ts.offset = offset
		return data, nil
	}
	return nil, ReaderFail{Topic: topic}

}

// 写入消息并返回偏移量
func (s *Store) WriteMessage(topic string, message []byte) (int64, error) {

	value, ok := s.topicFiles.Load(topic)
	if !ok {
		return 0, ReaderFail{Topic: topic}
	}
	tf := value.(topicFile)

	readOffset := tf.ts.loadOffset()
	//if offset more than threshold create new
	if readOffset+int64(len(message)) > threshold {
		f, error := s.ReCreateTopicfile(tf.ts)
		if error == nil {
			tf = *f
		}
	}

	tf.mu.Lock()
	defer tf.mu.Unlock()
	//repete seek
	offset, err := tf.f.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}
	_, err = tf.f.Write(message)
	if err != nil {
		return 0, err
	}
	// 返回偏移量
	return offset, nil
}

// return new value
func (s *topicStore) addOffset(value int64) int64 {
	return atomic.AddInt64(&s.offset, value) // 原子加 1

}

// return new value
func (s *topicStore) loadOffset() int64 {
	return atomic.LoadInt64(&s.offset)
}
