package kafka

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"
	"unicode/utf8"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
)

type DlqMessage struct {
	Error           string            `json:"error"`
	SourceTopic     string            `json:"source_topic"`
	SourcePartition int32             `json:"source_partition"`
	SourceOffset    int64             `json:"source_offset"`
	Timestamp       time.Time         `json:"timestamp"`
	KeyBase64       string            `json:"key_base64,omitempty"`
	ValueBase64     string            `json:"value_base64,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
}

type DlqWriter struct {
	log   *zerolog.Logger
	topic string
	p     *ckafka.Producer
}

func NewDlqWriter(log *zerolog.Logger, brokers, topic string) (*DlqWriter, error) {
	p, err := ckafka.NewProducer(&ckafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"enable.idempotence": true,
	})
	if err != nil {
		return nil, err
	}

	w := &DlqWriter{log: log, topic: topic, p: p}
	return w, nil
}

func (w *DlqWriter) Start() {
	go func() {
		for e := range w.p.Events() {
			switch ev := e.(type) {
			case *ckafka.Message:
				if ev.TopicPartition.Error != nil {
					w.log.Error().Err(ev.TopicPartition.Error).Msg("DLQ delivery failed")
				}
			}
		}
	}()
}

func (w *DlqWriter) Close() {
	if w == nil || w.p == nil {
		return
	}
	w.p.Flush(3000)
	w.p.Close()
}

func (w *DlqWriter) Publish(ctx context.Context, source *ckafka.Message, cause error) error {
	if w == nil || w.p == nil {
		return nil
	}

	headers := make(map[string]string)
	for _, h := range source.Headers {
		v := string(h.Value)
		if !utf8.Valid(h.Value) {
			v = "b64:" + base64.StdEncoding.EncodeToString(h.Value)
		}
		headers[h.Key] = v
	}

	dlqMsg := DlqMessage{
		Error:           cause.Error(),
		SourceTopic:     *source.TopicPartition.Topic,
		SourcePartition: source.TopicPartition.Partition,
		SourceOffset:    int64(source.TopicPartition.Offset),
		Timestamp:       source.Timestamp,
		Headers:         headers,
	}

	if len(source.Key) > 0 {
		dlqMsg.KeyBase64 = base64.StdEncoding.EncodeToString(source.Key)
	}
	if len(source.Value) > 0 {
		dlqMsg.ValueBase64 = base64.StdEncoding.EncodeToString(source.Value)
	}

	payload, err := json.Marshal(dlqMsg)
	if err != nil {
		return err
	}

	pMsg := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &w.topic, Partition: ckafka.PartitionAny},
		Value:          payload,
		Key:            source.Key,
		Timestamp:      time.Now(),
	}

	if err := w.p.Produce(pMsg, nil); err != nil {
		return err
	}

	w.p.Flush(2000)
	w.log.Warn().Str("dlq_topic", w.topic).Str("source_topic", dlqMsg.SourceTopic).Int32("partition", dlqMsg.SourcePartition).Int64("offset", dlqMsg.SourceOffset).Msg("Message sent to DLQ")
	return nil
}
