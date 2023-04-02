package cubeOperatorKit

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

/*
用于管理Kafka客户端，处理与ML通信相关事务
*/
const kafkaAddress = "127.0.0.1:30901"
const kafkaTopic = "testTopic"

func InitKafkaProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner //写到随机分区中，默认设置8个分区
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{kafkaAddress}, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// ProduceObject 向消息队列发送一个对象
func ProduceObject(client sarama.SyncProducer, key string, value string) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = kafkaTopic
	msg.Key = sarama.StringEncoder(key)
	msg.Value = sarama.StringEncoder(value)
	_, _, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("向ML Kafka发送信息失败, ", err)
		return
	}
}

func InitKafkaConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	client, err := sarama.NewConsumer([]string{kafkaAddress}, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// ConsumerStartListening ！有阻塞 让消费者开始监听消息队列，理论阻塞不会结束
func ConsumerStartListening(consumer sarama.Consumer, handler func(key string, value string)) error {
	var wg sync.WaitGroup
	partitionList, err := consumer.Partitions(kafkaTopic) //获得该topic所有的分区
	if err != nil {
		log.Println("获取Kafka partition失败:, ", err)
		return err
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(kafkaTopic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("为分区%d创建消费者失败: %s\n\n", partition, err)
			return err
		}
		wg.Add(1)
		go func(sarama.PartitionConsumer) { //为每个分区开一个go协程去取值
			for msg := range pc.Messages() { //阻塞直到有值发送过来，然后再继续等待
				handler(string(msg.Key), string(msg.Value))
			}
			defer pc.AsyncClose()
			wg.Done()
		}(pc)
	}
	wg.Wait()
	return nil
}
