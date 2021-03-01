package utils

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	conf "tesla/config"
	"time"
)

var (
	Producer sarama.SyncProducer
)

func init() {
	var err error
	config := sarama.NewConfig()
	// request.timeout.ms
	config.Producer.Timeout = time.Second * 5
	// message.max.bytes
	config.Producer.MaxMessageBytes = 1024 * 1024
	// request.required.acks
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_6_0_0

	if err := config.Validate(); err != nil {
		panic(fmt.Errorf("invalid configuration, error: %v", err))
	}
	kafkaUrlString := conf.AppConfig.KafkaUrl
	kafkaUrl := strings.Split(kafkaUrlString, ",")
	Producer, err = sarama.NewSyncProducer(kafkaUrl, config)
	if err != nil {
		panic(err)
	}
}
