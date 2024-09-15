package kafka

import "github.com/segmentio/kafka-go"

// 서비스 등록 카프카 producer 생성
func NewWriter(broker, topic string) *kafka.Writer {
	return kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:     []string{broker},
			Topic:       topic,
			MaxAttempts: 5,
		})
}
