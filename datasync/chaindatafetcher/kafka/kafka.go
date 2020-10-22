// Copyright 2020 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
)

var logger = log.NewModuleLogger(log.ChainDataFetcher)

const (
	MsgIdxTotalSegments = iota
	MsgIdxSegmentIdx
)

const (
	KeyTotalSegments = "totalSegments"
	KeySegmentIdx    = "segmentIdx"
)

type IKey interface {
	Key() string
}

// Kafka connects to the brokers in an existing kafka cluster.
type Kafka struct {
	config   *KafkaConfig
	producer sarama.SyncProducer
	admin    sarama.ClusterAdmin
}

func NewKafka(conf *KafkaConfig) (*Kafka, error) {
	producer, err := sarama.NewSyncProducer(conf.Brokers, conf.SaramaConfig)
	if err != nil {
		logger.Error("Failed to create a new producer", "brokers", conf.Brokers)
		return nil, err
	}

	admin, err := sarama.NewClusterAdmin(conf.Brokers, conf.SaramaConfig)
	if err != nil {
		logger.Error("Failed to create a new cluster admin", "brokers", conf.Brokers)
		return nil, err
	}

	return &Kafka{
		config:   conf,
		producer: producer,
		admin:    admin,
	}, nil
}

func (k *Kafka) Close() {
	k.producer.Close()
	k.admin.Close()
}

func (k *Kafka) getTopicName(event string) string {
	return k.config.GetTopicName(event)
}

func (k *Kafka) CreateTopic(topic string) error {
	return k.admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     k.config.Partitions,
		ReplicationFactor: k.config.Replicas,
	}, false)
}

func (k *Kafka) DeleteTopic(topic string) error {
	return k.admin.DeleteTopic(topic)
}

func (k *Kafka) ListTopics() (map[string]sarama.TopicDetail, error) {
	return k.admin.ListTopics()
}

func (k *Kafka) split(data []byte) ([][]byte, int) {
	size := k.config.SegmentSizeBytes
	var segments [][]byte
	for len(data) > size {
		segments = append(segments, data[:size])
		data = data[size:]
	}
	segments = append(segments, data)
	return segments, len(segments)
}

func (k *Kafka) makeProducerMessage(topic, key string, segment []byte, segmentIdx, totalSegments uint64) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte(KeyTotalSegments),
				Value: common.Int64ToByteBigEndian(totalSegments),
			},
			{
				Key:   []byte(KeySegmentIdx),
				Value: common.Int64ToByteBigEndian(segmentIdx),
			},
		},
		Value: sarama.ByteEncoder(segment),
	}
}

func (k *Kafka) Publish(topic string, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	key := ""
	if v, ok := data.(IKey); ok {
		key = v.Key()
	}
	segments, totalSegments := k.split(dataBytes)
	for idx, segment := range segments {
		msg := k.makeProducerMessage(topic, key, segment, uint64(idx), uint64(totalSegments))
		_, _, err = k.producer.SendMessage(msg)
		if err != nil {
			logger.Error("sending kafka message is failed", "err", err, "segmentIdx", idx, "key", key)
			return err
		}
	}

	return err
}
