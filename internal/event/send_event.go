package event

import (
	"encoding/json"
	"postit/model"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	topic         = "EVENT"
	kafkaProducer *kafka.Producer
)

func Init() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	})
	if err != nil {
		logger.Error(err)
	}

	kafkaProducer = producer
}

func SendPostEvent(username, status string, post *model.Post) {
	eventData := model.EventData{
		Data:   post,
		Status: status,
	}
	event := model.Event{
		Username:      username,
		CustomerEvent: model.EventPost,
		Data:          eventData,
	}
	sendEvent(event)
}

func SendLikeEvent(username, status string, like *model.Like) {
	eventData := model.EventData{
		Data:   like,
		Status: status,
	}
	event := model.Event{
		Username:      username,
		CustomerEvent: model.EventLike,
		Data:          eventData,
	}
	sendEvent(event)
}

func SendCommentEvent(username, status string, comment *model.Comment) {
	eventData := model.EventData{
		Data:   comment,
		Status: status,
	}
	event := model.Event{
		Username:      username,
		CustomerEvent: model.EventComment,
		Data:          eventData,
	}
	sendEvent(event)
}

func sendEvent(event model.Event) {
	stringEvent, err := json.Marshal(event)
	if err != nil {
		logger.Error(err)
	}

	err = kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(stringEvent)},
		nil,
	)

	if err != nil {
		logger.Error(err)
	} else {
		logger.Infof("Sent event off to kafka %s", event)
	}
}
