package kafka

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"
	"unicode/utf8"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type DlqMessage struct {
	Error           string            `json:"error"`
	SourceTopic     string            `json:"source_topic"`
	SourcePartition int               `json:"source_partition"`
	SourceOffset    int64             `json:"source_offset"`
	Timestamp       time.Time         `json:"timestamp"`
	KeyBase64       string            `json:"key_base64,omitempty"`
	ValueBase64     string            `json:"value_base64,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
}

type DlqWriter struct {
	log   *zerolog.Logger
	topic string
	kw    *kafka.Writer
}

func NewDlqWriter(log *zerolog.Logger, brokers []string, topic string) (*DlqWriter, error) {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
	}
	return &DlqWriter{log: log, topic: topic, kw: w}, nil
}

func (w *DlqWriter) Close() {
	_ = w.kw.Close()
}

func (w *DlqWriter) Publish(ctx context.Context, source *kafka.Message, cause error) error {
	if w == nil || w.kw == nil {
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
		SourceTopic:     source.Topic,
		SourcePartition: source.Partition,
		SourceOffset:    source.Offset,
		Timestamp:       source.Time,
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

	err = w.kw.WriteMessages(ctx, kafka.Message{
		Topic: w.topic,
		Key:   source.Key,
		Value: payload,
		Time:  time.Now(),
	})
	if err != nil {
		return err
	}

	w.log.Warn().Str("dlq_topic", w.topic).Str("source_topic", dlqMsg.SourceTopic).Int("partition", dlqMsg.SourcePartition).Int64("offset", dlqMsg.SourceOffset).Msg("Message sent to DLQ")
	return nil
}
